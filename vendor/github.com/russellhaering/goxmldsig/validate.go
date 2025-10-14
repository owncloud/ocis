package dsig

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"

	"github.com/beevik/etree"
	"github.com/russellhaering/goxmldsig/etreeutils"
	"github.com/russellhaering/goxmldsig/types"
)

var uriRegexp = regexp.MustCompile("^#[a-zA-Z_][\\w.-]*$")
var whiteSpace = regexp.MustCompile("\\s+")

var (
	// ErrMissingSignature indicates that no enveloped signature was found referencing
	// the top level element passed for signature verification.
	ErrMissingSignature = errors.New("Missing signature referencing the top-level element")
	ErrInvalidSignature = errors.New("Invalid Signature")
)

type ValidationContext struct {
	CertificateStore X509CertificateStore
	IdAttribute      string
	Clock            *Clock
}

func NewDefaultValidationContext(certificateStore X509CertificateStore) *ValidationContext {
	return &ValidationContext{
		CertificateStore: certificateStore,
		IdAttribute:      DefaultIdAttr,
	}
}

// TODO(russell_h): More flexible namespace support. This might barely work.
func inNamespace(el *etree.Element, ns string) bool {
	for _, attr := range el.Attr {
		if attr.Value == ns {
			if attr.Space == "" && attr.Key == "xmlns" {
				return el.Space == ""
			} else if attr.Space == "xmlns" {
				return el.Space == attr.Key
			}
		}
	}

	return false
}

func childPath(space, tag string) string {
	if space == "" {
		return "./" + tag
	} else {
		return "./" + space + ":" + tag
	}
}

func mapPathToElement(tree, el *etree.Element) []int {
	for i, child := range tree.Child {
		if child == el {
			return []int{i}
		}
	}

	for i, child := range tree.Child {
		if childElement, ok := child.(*etree.Element); ok {
			childPath := mapPathToElement(childElement, el)
			if childPath != nil {
				return append([]int{i}, childPath...)
			}
		}
	}

	return nil
}

func removeElementAtPath(el *etree.Element, path []int) bool {
	if len(path) == 0 {
		return false
	}

	if len(el.Child) <= path[0] {
		return false
	}

	childElement, ok := el.Child[path[0]].(*etree.Element)
	if !ok {
		return false
	}

	if len(path) == 1 {
		el.RemoveChild(childElement)
		return true
	}

	return removeElementAtPath(childElement, path[1:])
}

// Transform returns a new element equivalent to the passed root el, but with
// the set of transformations described by the ref applied.
//
// The functionality of transform is currently very limited and purpose-specific.
func (ctx *ValidationContext) transform(
	el *etree.Element,
	sig *types.Signature,
	ref *types.Reference) (*etree.Element, Canonicalizer, error) {
	transforms := ref.Transforms.Transforms

	// map the path to the passed signature relative to the passed root, in
	// order to enable removal of the signature by an enveloped signature
	// transform
	signaturePath := mapPathToElement(el, sig.UnderlyingElement())

	// make a copy of the passed root
	el = el.Copy()

	var canonicalizer Canonicalizer

	for _, transform := range transforms {
		algo := transform.Algorithm

		switch AlgorithmID(algo) {
		case EnvelopedSignatureAltorithmId:
			if !removeElementAtPath(el, signaturePath) {
				return nil, nil, errors.New("Error applying canonicalization transform: Signature not found")
			}

		case CanonicalXML10ExclusiveAlgorithmId:
			var prefixList string
			if transform.InclusiveNamespaces != nil {
				prefixList = transform.InclusiveNamespaces.PrefixList
			}

			canonicalizer = MakeC14N10ExclusiveCanonicalizerWithPrefixList(prefixList)

		case CanonicalXML10ExclusiveWithCommentsAlgorithmId:
			var prefixList string
			if transform.InclusiveNamespaces != nil {
				prefixList = transform.InclusiveNamespaces.PrefixList
			}

			canonicalizer = MakeC14N10ExclusiveWithCommentsCanonicalizerWithPrefixList(prefixList)

		case CanonicalXML11AlgorithmId:
			canonicalizer = MakeC14N11Canonicalizer()

		case CanonicalXML11WithCommentsAlgorithmId:
			canonicalizer = MakeC14N11WithCommentsCanonicalizer()

		case CanonicalXML10RecAlgorithmId:
			canonicalizer = MakeC14N10RecCanonicalizer()

		case CanonicalXML10WithCommentsAlgorithmId:
			canonicalizer = MakeC14N10WithCommentsCanonicalizer()

		default:
			return nil, nil, errors.New("Unknown Transform Algorithm: " + algo)
		}
	}

	if canonicalizer == nil {
		canonicalizer = MakeNullCanonicalizer()
	}

	return el, canonicalizer, nil
}

// deprecated
func (ctx *ValidationContext) digest(el *etree.Element, digestAlgorithmId string, canonicalizer Canonicalizer) ([]byte, error) {
	data, err := canonicalizer.Canonicalize(el)
	if err != nil {
		return nil, err
	}

	digestAlgorithm, ok := digestAlgorithmsByIdentifier[digestAlgorithmId]
	if !ok {
		return nil, errors.New("Unknown digest algorithm: " + digestAlgorithmId)
	}

	hash := digestAlgorithm.New()
	_, err = hash.Write(data)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func (ctx *ValidationContext) getCanonicalSignedInfo(sig *types.Signature) ([]byte, error) {
	signatureElement := sig.UnderlyingElement()

	nsCtx, err := etreeutils.NSBuildParentContext(signatureElement)
	if err != nil {
		return nil, err
	}

	signedInfo, err := etreeutils.NSFindOneChildCtx(nsCtx, signatureElement, Namespace, SignedInfoTag)
	if err != nil {
		return nil, err
	}

	if signedInfo == nil {
		return nil, errors.New("Missing SignedInfo")
	}

	// Canonicalize the xml
	canonical, err := canonicalSerialize(signedInfo)
	if err != nil {
		return nil, err
	}

	return canonical, nil
}

// deprecated
func (ctx *ValidationContext) verifySignedInfo(sig *types.Signature, canonicalizer Canonicalizer, signatureMethodId string, cert *x509.Certificate, decodedSignature []byte) error {
	signatureElement := sig.UnderlyingElement()

	nsCtx, err := etreeutils.NSBuildParentContext(signatureElement)
	if err != nil {
		return err
	}

	signedInfo, err := etreeutils.NSFindOneChildCtx(nsCtx, signatureElement, Namespace, SignedInfoTag)
	if err != nil {
		return err
	}

	if signedInfo == nil {
		return errors.New("Missing SignedInfo")
	}

	// Canonicalize the xml
	canonical, err := canonicalSerialize(signedInfo)
	if err != nil {
		return err
	}

	algo, ok := x509SignatureAlgorithmByIdentifier[signatureMethodId]
	if !ok {
		return errors.New("Unknown signature method: " + signatureMethodId)
	}

	err = cert.CheckSignature(algo, canonical, decodedSignature)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *ValidationContext) validateSignature(el *etree.Element, sig *types.Signature, cert *x509.Certificate) (*etree.Element, error) {

	// Actually verify the 'SignedInfo' was signed by a trusted source
	signatureMethod := sig.SignedInfo.SignatureMethod.Algorithm

	canonicalSignedInfoBytes, err := ctx.getCanonicalSignedInfo(sig)
	if err != nil {
		return nil, errors.New("Could not obtain canonical signed info bytes")
	}

	if canonicalSignedInfoBytes == nil {
		return nil, errors.New("Missing SignedInfo")
	}

	algo, ok := x509SignatureAlgorithmByIdentifier[signatureMethod]
	if !ok {
		return nil, errors.New("Unknown signature method: " + signatureMethod)
	}

	if sig.SignatureValue == nil {
		return nil, errors.New("Signature could not be verified")
	}

	// Decode the 'SignatureValue' so we can compare against it
	decodedSignature, err := base64.StdEncoding.DecodeString(sig.SignatureValue.Data)
	if err != nil {
		return nil, errors.New("Could not decode signature")
	}

	err = cert.CheckSignature(algo, canonicalSignedInfoBytes, decodedSignature)
	if err != nil {
		return nil, err
	}

	// only use the verified canonicalSignedInfoBytes
	// unmarshal canonicalSignedInfoBytes into a new SignedInfo type
	// to obtain the reference
	signedInfo := &types.SignedInfo{}
	err = xml.Unmarshal(canonicalSignedInfoBytes, signedInfo)
	if err != nil {
		return nil, err
	}

	idAttrEl := el.SelectAttr(ctx.IdAttribute)
	idAttr := ""
	if idAttrEl != nil {
		idAttr = idAttrEl.Value
	}

	var ref *types.Reference

	// Find the first reference which references the top-level element
	for _, _ref := range signedInfo.References {
		if _ref.URI == "" || _ref.URI[1:] == idAttr {
			ref = &_ref
		}
	}

	// prevents null pointer deref
	if ref == nil {
		return nil, errors.New("Missing reference")
	}

	digestAlgorithmId := ref.DigestAlgo.Algorithm
	signedDigestValue, err := base64.StdEncoding.DecodeString(ref.DigestValue)
	if err != nil {
		return nil, err
	}

	// Perform all transformations listed in the 'SignedInfo'
	// Basically, this means removing the 'SignedInfo'
	transformed, canonicalizer, err := ctx.transform(el, sig, ref)
	if err != nil {
		return nil, err
	}

	referencedBytes, err := canonicalizer.Canonicalize(transformed)
	if err != nil {
		return nil, err
	}

	// use a known digest hashing algorithm
	hashAlgorithm, ok := digestAlgorithmsByIdentifier[digestAlgorithmId]
	if !ok {
		return nil, errors.New("Unknown digest algorithm: " + digestAlgorithmId)
	}

	hash := hashAlgorithm.New()
	_, err = hash.Write(referencedBytes)
	if err != nil {
		return nil, err
	}

	computedDigest := hash.Sum(nil)
	/* Digest the transformed XML and compare it to the 'DigestValue' from the 'SignedInfo'
	digest, err := ctx.digest(transformed, digestAlgorithm, canonicalizer)
	*/

	if !bytes.Equal(computedDigest, signedDigestValue) {
		return nil, errors.New("Signature could not be verified")
	}

	if !(len(computedDigest) >= 20) {
		return nil, errors.New("Computed digest is less than 20 something went wrong")
	}

	// now only the referencedBytes is verified,
	// unmarshal into new etree
	doc := etree.NewDocument()
	err = doc.ReadFromBytes(referencedBytes)
	if err != nil {
		return nil, err
	}

	return doc.Root(), nil
}

func contains(roots []*x509.Certificate, cert *x509.Certificate) bool {
	for _, root := range roots {
		if root.Equal(cert) {
			return true
		}
	}
	return false
}

// In most places, we use etree Elements, but while deserializing the Signature, we use
// encoding/xml unmarshal directly to convert to a convenient go struct. This presents a problem in some cases because
// when an xml element repeats under the parent, the last element will win and/or be appended. We need to assert that
// the Signature object matches the expected shape of a Signature object.
func validateShape(signatureEl *etree.Element) error {
	children := signatureEl.ChildElements()

	childCounts := map[string]int{}
	for _, child := range children {
		childCounts[child.Tag]++
	}

	validateCount := childCounts[SignedInfoTag] == 1 && childCounts[KeyInfoTag] <= 1 && childCounts[SignatureValueTag] == 1
	if !validateCount {
		return ErrInvalidSignature
	}
	return nil
}

// findSignature searches for a Signature element referencing the passed root element.
func (ctx *ValidationContext) findSignature(root *etree.Element) (*types.Signature, error) {
	idAttrEl := root.SelectAttr(ctx.IdAttribute)
	idAttr := ""
	if idAttrEl != nil {
		idAttr = idAttrEl.Value
	}

	var sig *types.Signature

	// Traverse the tree looking for a Signature element
	err := etreeutils.NSFindIterate(root, Namespace, SignatureTag, func(ctx etreeutils.NSContext, signatureEl *etree.Element) error {
		err := validateShape(signatureEl)
		if err != nil {
			return err
		}
		found := false
		err = etreeutils.NSFindChildrenIterateCtx(ctx, signatureEl, Namespace, SignedInfoTag,
			func(ctx etreeutils.NSContext, signedInfo *etree.Element) error {
				detachedSignedInfo, err := etreeutils.NSDetatch(ctx, signedInfo)
				if err != nil {
					return err
				}

				c14NMethod, err := etreeutils.NSFindOneChildCtx(ctx, detachedSignedInfo, Namespace, CanonicalizationMethodTag)
				if err != nil {
					return err
				}

				if c14NMethod == nil {
					return errors.New("missing CanonicalizationMethod on Signature")
				}

				c14NAlgorithm := c14NMethod.SelectAttrValue(AlgorithmAttr, "")

				var canonicalSignedInfo *etree.Element

				switch alg := AlgorithmID(c14NAlgorithm); alg {
				case CanonicalXML10ExclusiveAlgorithmId, CanonicalXML10ExclusiveWithCommentsAlgorithmId:
					err := etreeutils.TransformExcC14n(detachedSignedInfo, "", alg == CanonicalXML10ExclusiveWithCommentsAlgorithmId)
					if err != nil {
						return err
					}

					// NOTE: TransformExcC14n transforms the element in-place,
					// while canonicalPrep isn't meant to. Once we standardize
					// this behavior we can drop this, as well as the adding and
					// removing of elements below.
					canonicalSignedInfo = detachedSignedInfo

				case CanonicalXML11AlgorithmId, CanonicalXML10RecAlgorithmId:
					canonicalSignedInfo = canonicalPrep(detachedSignedInfo, true, false)

				case CanonicalXML11WithCommentsAlgorithmId, CanonicalXML10WithCommentsAlgorithmId:
					canonicalSignedInfo = canonicalPrep(detachedSignedInfo, true, true)

				default:
					return fmt.Errorf("invalid CanonicalizationMethod on Signature: %s", c14NAlgorithm)
				}

				signatureEl.InsertChildAt(signedInfo.Index(), canonicalSignedInfo)
				signatureEl.RemoveChild(signedInfo)

				found = true

				return etreeutils.ErrTraversalHalted
			})
		if err != nil {
			return err
		}

		if !found {
			return errors.New("Missing SignedInfo")
		}

		// Unmarshal the signature into a structured Signature type
		_sig := &types.Signature{}
		err = etreeutils.NSUnmarshalElement(ctx, signatureEl, _sig)
		if err != nil {
			return err
		}

		// Traverse references in the signature to determine whether it has at least
		// one reference to the top level element. If so, conclude the search.
		for _, ref := range _sig.SignedInfo.References {
			if ref.URI == "" || ref.URI[1:] == idAttr {
				sig = _sig
				return etreeutils.ErrTraversalHalted
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if sig == nil {
		return nil, ErrMissingSignature
	}

	return sig, nil
}

func (ctx *ValidationContext) verifyCertificate(sig *types.Signature) (*x509.Certificate, error) {
	now := ctx.Clock.Now()

	roots, err := ctx.CertificateStore.Certificates()
	if err != nil {
		return nil, err
	}

	var untrustedCert *x509.Certificate

	if sig.KeyInfo != nil {
		// If the Signature includes KeyInfo, extract the certificate from there
		if len(sig.KeyInfo.X509Data.X509Certificates) == 0 || sig.KeyInfo.X509Data.X509Certificates[0].Data == "" {
			return nil, errors.New("missing X509Certificate within KeyInfo")
		}

		certData, err := base64.StdEncoding.DecodeString(
			whiteSpace.ReplaceAllString(sig.KeyInfo.X509Data.X509Certificates[0].Data, ""))
		if err != nil {
			return nil, errors.New("Failed to parse certificate")
		}

		untrustedCert, err = x509.ParseCertificate(certData)
		if err != nil {
			return nil, err
		}
	} else {
		// If the Signature doesn't have KeyInfo, Use the root certificate if there is only one
		if len(roots) == 1 {
			untrustedCert = roots[0]
		} else {
			return nil, errors.New("Missing x509 Element")
		}
	}

	rootIdx := -1
	for i, root := range roots {
		if root.Equal(untrustedCert) {
			rootIdx = i
		}
	}

	if rootIdx == -1 {
		return nil, errors.New("Could not verify certificate against trusted certs")
	}
	var trustedCert *x509.Certificate

	trustedCert = roots[rootIdx]

	// Verify that the certificate is one we trust

	if now.Before(trustedCert.NotBefore) || now.After(trustedCert.NotAfter) {
		return nil, errors.New("Cert is not valid at this time")
	}

	return trustedCert, nil
}

// Validate verifies that the passed element contains a valid enveloped signature
// matching a currently-valid certificate in the context's CertificateStore.
func (ctx *ValidationContext) Validate(el *etree.Element) (*etree.Element, error) {
	// Make a copy of the element to avoid mutating the one we were passed.
	el = el.Copy()

	sig, err := ctx.findSignature(el)
	if err != nil {
		return nil, err
	}

	// function to get the trusted certificate
	cert, err := ctx.verifyCertificate(sig)
	if err != nil {
		return nil, err
	}

	return ctx.validateSignature(el, sig, cert)
}
