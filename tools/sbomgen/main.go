// sbomgen scans a project (Go modules and a pnpm lockfile) and writes an enriched SPDX 3.0.1 JSON-LD
// SBOM.
//
// The SBOM is assembled as an SPDX 2.3 document, converted to SPDX 3 and validated by
// github.com/spdx/tools-golang.
//
// What it records per dependency:
//   - name, version and purl, and the dependency graph (Go toolchain, pnpm-lock.yaml)
//   - a sha256 of the Go module zip / the integrity hash of the npm package from the lockfile
//   - the declared license (npm metadata) or the license found in the package's license files
//   - the copyright statements found in the package's license files
//   - the distribution source (supplier) and download location
//
// The project itself is one root package (license and copyright from REUSE.toml, version from
// -version, hash of the git source archive).
package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
	spdxjson "github.com/spdx/tools-golang/json"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
	"github.com/spdx/tools-golang/spdx/v3/v3_0"
	"gopkg.in/yaml.v3"
)

// toolVersion identifies this tool in the SBOM; bump it when the output changes.
const toolVersion = "0.1.0"

const (
	repoURL      = "https://github.com/owncloud/ocis"
	owner        = "ownCloud GmbH - A Kiteworks Company"
	ownerContact = "devops@owncloud.com"                // PackageSupplier contact in REUSE.toml
	goProxy      = "Go Module Proxy - proxy.golang.org" // no parentheses: tools-golang parses "(x)" as an email
	npmRegistry  = "npm Registry - registry.npmjs.org"
	noAssertion  = "NOASSERTION"
	toolName     = "sbomgen-" + toolVersion
)

func main() {
	var (
		out, out23, version, pnpmLock, reuseToml, npmCache, sourceDir string
		npmDev                                                        bool
		goModules                                                     stringList
	)
	flag.StringVar(&out, "out", "", "SPDX 3 JSON-LD output")
	flag.StringVar(&out23, "out23", "", "also write the intermediate SPDX 2.3 JSON here")
	flag.StringVar(&version, "version", "", "project version")
	flag.Var(&goModules, "go-module", "directory of a Go module to scan (repeatable)")
	flag.StringVar(&pnpmLock, "pnpm-lock", "", "pnpm-lock.yaml")
	flag.BoolVar(&npmDev, "npm-dev", false, "include npm devDependencies")
	flag.StringVar(&reuseToml, "reuse-toml", "", "REUSE.toml (root license and copyright)")
	flag.StringVar(&npmCache, "npm-cache", "", "file that caches per-package npm license and copyright facts")
	flag.StringVar(&sourceDir, "source-dir", ".", "git checkout of the project (its archive is hashed for the root package)")
	flag.Parse()
	for name, v := range map[string]string{"out": out, "version": version, "pnpm-lock": pnpmLock, "reuse-toml": reuseToml} {
		if v == "" {
			fatal("missing -%s", name)
		}
	}
	if len(goModules) == 0 {
		fatal("missing -go-module")
	}

	rootLicense, rootCopyright := readReuse(reuseToml)
	integrity := readPnpmLock(pnpmLock)
	spdxIDs := loadLicenseList()
	stats := map[string]int{}

	doc, goZips, npmRefs := scan(goModules, pnpmLock, npmDev, version, spdxIDs)
	for i := range npmRefs {
		npmRefs[i].Pinned = integrity[npmRefs[i].Name+"@"+npmRefs[i].Version]
	}
	npmFacts := loadNPMInfo(npmRefs, npmCache, spdxIDs)

	// ---- the root package ---------------------------------------------------------------------
	root := doc.Packages[0]
	root.PackageVersion = version
	root.PackageDownloadLocation = repoURL
	root.PackageHomePage = repoURL
	root.PackageCopyrightText = rootCopyright
	root.PackageLicenseDeclared = rootLicense
	root.PackageSupplier = &common.Supplier{SupplierType: "Organization", Supplier: owner + " (" + ownerContact + ")"}
	root.PackageExternalReferences = []*v2_3.PackageExternalReference{
		{Category: "PACKAGE-MANAGER", RefType: "purl", Locator: "pkg:github/owncloud/ocis@" + version}}
	if commit, sum, err := gitArchiveSHA256(sourceDir); err != nil {
		fmt.Fprintf(os.Stderr, "sbomgen: warning: no hash for the root package: %v\n", err)
		stats["no_hash"]++
	} else {
		root.PackageChecksums = []common.Checksum{{Algorithm: common.SHA256, Value: sum}}
		root.PackageSourceInfo = "git commit " + commit
		root.PackageComment = fmt.Sprintf("The checksum is the SHA-256 of `git archive --format=tar %s`.", commit)
	}

	// ---- dependencies --------------------------------------------------------------------------
	for _, p := range doc.Packages[1:] {
		var hash *common.Checksum
		switch ecosystem(p) {
		case "golang":
			hash = goHash(goZips, p.PackageName, p.PackageVersion)
			p.PackageSupplier = &common.Supplier{SupplierType: "Organization", Supplier: goProxy}
			p.PackageDownloadLocation = fmt.Sprintf("https://proxy.golang.org/%s/@v/%s.zip",
				goEscape(p.PackageName), goEscape(p.PackageVersion))
			if p.PackageLicenseDeclared != noAssertion {
				p.PackageLicenseComments = "Detected by classifying the license files of the module."
			}
		case "npm":
			hash = integrity[p.PackageName+"@"+p.PackageVersion]
			p.PackageSupplier = &common.Supplier{SupplierType: "Organization", Supplier: npmRegistry}
			base := p.PackageName[strings.LastIndex(p.PackageName, "/")+1:]
			p.PackageDownloadLocation = fmt.Sprintf("https://registry.npmjs.org/%s/-/%s-%s.tgz", p.PackageName, base, p.PackageVersion)
			facts := npmFacts[p.PackageName+"@"+p.PackageVersion]
			p.PackageLicenseDeclared = facts.License
			p.PackageCopyrightText = facts.Copyright
			switch facts.Source {
			case "license file", "readme":
				p.PackageLicenseComments = "Detected from the " + facts.Source + " in the package tarball; the package metadata declares no license."
			}
		}
		if hash != nil {
			p.PackageChecksums = []common.Checksum{*hash}
		} else {
			stats["no_hash"]++
		}
		if p.PackageLicenseDeclared != noAssertion && !validExpression(p.PackageLicenseDeclared, spdxIDs) {
			stats["license_invalid"]++
			p.PackageLicenseDeclared = noAssertion
		}
		if p.PackageLicenseDeclared == noAssertion {
			stats["no_license"]++
		}
		if p.PackageCopyrightText == noAssertion {
			stats["no_copyright"]++
		}
	}
	sort.SliceStable(doc.Relationships, func(i, j int) bool {
		return doc.Relationships[i].RefA.ElementRefID < doc.Relationships[j].RefA.ElementRefID
	})
	doc.CreationInfo.Creators = []common.Creator{
		{CreatorType: "Organization", Creator: owner + " (" + ownerContact + ")"},
		{CreatorType: "Tool", Creator: toolName}}

	if out23 != "" {
		w, err := os.Create(out23)
		check(err)
		check(spdxjson.Write(doc, w, spdxjson.Indent("  ")))
		check(w.Close())
	}
	// NOASSERTION would become a license edge to a NOASSERTION object in SPDX 3; leave those out instead
	for _, p := range doc.Packages {
		if p.PackageLicenseDeclared == noAssertion {
			p.PackageLicenseDeclared = ""
		}
		p.PackageLicenseConcluded = ""
	}

	// ---- convert to SPDX 3 and validate -------------------------------------------------------
	d3 := &v3_0.Document{}
	d3.Name = doc.DocumentName
	v3_0.From_v2_3(*doc, d3)
	// the 2.3 conversion leaves the (shared) creation info without a spec version
	d3.CreationInfo.SetSpecVersion(v3_0.Version)
	d3.ProfileConformances = []v3_0.ProfileIdentifierType{
		v3_0.ProfileIdentifierType_Core, v3_0.ProfileIdentifierType_Software, v3_0.ProfileIdentifierType_ExpandedLicensing}
	d3.RootElements[0].(*v3_0.SBOM).SbomTypes = []v3_0.SbomType{v3_0.SbomType_Source} // made from manifests and lockfiles
	dedupe(d3)
	if err := d3.Validate(true); err != nil {
		fmt.Fprintln(os.Stderr, "VALIDATION ERRORS (tools-golang):\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "tools-golang validation: OK")

	o, err := os.Create(out)
	check(err)
	check(d3.Write(o))
	check(o.Close())
	fmt.Fprintf(os.Stderr, "%d packages, %d relationships, stats=%v\n", len(doc.Packages), len(doc.Relationships), stats)
}

// ---- helpers --------------------------------------------------------------------------------------

func ecosystem(p *v2_3.Package) string {
	for _, r := range p.PackageExternalReferences {
		if r.RefType == "purl" {
			s := strings.TrimPrefix(r.Locator, "pkg:")
			return s[:strings.Index(s, "/")]
		}
	}
	return ""
}

func goEscape(s string) string { // module path/version case-encoding used by the Go proxy
	return regexp.MustCompile(`[A-Z]`).ReplaceAllStringFunc(s, func(m string) string { return "!" + strings.ToLower(m) })
}

// gitArchiveSHA256 hashes `git archive --format=tar HEAD` of the checkout in dir.
func gitArchiveSHA256(dir string) (commit, sum string, err error) {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", "", fmt.Errorf("git rev-parse HEAD in %s: %v", dir, err)
	}
	commit = strings.TrimSpace(string(out))
	cmd := exec.Command("git", "-C", dir, "archive", "--format=tar", commit)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	if err := cmd.Start(); err != nil {
		return "", "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, stdout); err != nil {
		return "", "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", "", fmt.Errorf("git archive: %v", err)
	}
	return commit, hex.EncodeToString(h.Sum(nil)), nil
}

func readReuse(path string) (license, copyright string) {
	b, err := os.ReadFile(path)
	check(err)
	var t struct {
		Annotations []struct {
			License   string `toml:"SPDX-License-Identifier"`
			Copyright string `toml:"SPDX-FileCopyrightText"`
		} `toml:"annotations"`
	}
	check(toml.Unmarshal(b, &t))
	if len(t.Annotations) == 0 || t.Annotations[0].License == "" || t.Annotations[0].Copyright == "" {
		fatal("%s: no [[annotations]] license/copyright", path)
	}
	return t.Annotations[0].License, t.Annotations[0].Copyright
}

func goHash(zips map[string]string, name, version string) *common.Checksum {
	zip, ok := zips[name+"@"+version]
	if !ok {
		return nil
	}
	f, err := os.Open(zip)
	check(err)
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	check(err)
	return &common.Checksum{Algorithm: common.SHA256, Value: hex.EncodeToString(h.Sum(nil))}
}

func readPnpmLock(path string) map[string]*common.Checksum {
	b, err := os.ReadFile(path)
	check(err)
	var lock struct {
		Packages map[string]struct {
			Resolution struct{ Integrity string } `yaml:"resolution"`
		} `yaml:"packages"`
	}
	check(yaml.Unmarshal(b, &lock))
	algos := map[string]common.ChecksumAlgorithm{"sha1": common.SHA1, "sha256": common.SHA256, "sha384": common.SHA384, "sha512": common.SHA512}
	out := map[string]*common.Checksum{}
	for k, v := range lock.Packages {
		algo, b64, ok := strings.Cut(v.Resolution.Integrity, "-")
		raw, err := base64.StdEncoding.DecodeString(b64)
		if a, known := algos[algo]; ok && known && err == nil {
			out[k] = &common.Checksum{Algorithm: a, Value: hex.EncodeToString(raw)}
		}
	}
	return out
}

// loadLicenseList returns the identifiers (licenses and exceptions) of the published SPDX license list.
func loadLicenseList() map[string]bool {
	ids := map[string]bool{}
	for _, u := range []string{"https://spdx.org/licenses/licenses.json", "https://spdx.org/licenses/exceptions.json"} {
		resp, err := getWithRetry(u)
		check(err)
		var l struct {
			Licenses []struct {
				LicenseID string `json:"licenseId"`
			} `json:"licenses"`
			Exceptions []struct {
				LicenseExceptionID string `json:"licenseExceptionId"`
			} `json:"exceptions"`
		}
		check(json.NewDecoder(resp.Body).Decode(&l))
		resp.Body.Close()
		for _, x := range l.Licenses {
			ids[x.LicenseID] = true
		}
		for _, x := range l.Exceptions {
			ids[x.LicenseExceptionID] = true
		}
	}
	if len(ids) < 400 {
		fatal("SPDX license list looks incomplete (%d ids)", len(ids))
	}
	return ids
}

func validExpression(expr string, ids map[string]bool) bool {
	if _, err := v3_0.ParseLicenseExpression(expr); err != nil {
		return false
	}
	for _, tok := range strings.Fields(strings.NewReplacer("(", " ", ")", " ").Replace(expr)) {
		switch tok {
		case "AND", "OR", "WITH":
			continue
		}
		if !ids[strings.TrimSuffix(tok, "+")] {
			return false
		}
	}
	return true
}

// dedupe makes equal suppliers and listed licenses share one element. The tools-golang conversion
// creates a new object for every use, which inflates the document (one Organization per package).
func dedupe(d *v3_0.Document) {
	sbom := d.RootElements[0].(*v3_0.SBOM)
	orgs := map[string]v3_0.AnyAgent{}
	for _, a := range d.CreationInfo.(*v3_0.CreationInfo).CreatedBy {
		if o, ok := a.(*v3_0.Organization); ok {
			orgs[o.Name] = o // the root package is supplied by its author
		}
	}
	licenses := map[string]v3_0.AnyElement{}
	for _, e := range sbom.Elements {
		switch x := e.(type) {
		case *v3_0.Package:
			if o, ok := x.SuppliedBy.(*v3_0.Organization); ok {
				if first, seen := orgs[o.Name]; seen {
					x.SuppliedBy = first
				} else {
					orgs[o.Name] = o
				}
			}
		case *v3_0.Relationship:
			if len(x.To) == 1 {
				if l, ok := x.To[0].(*v3_0.ListedLicense); ok {
					if first, seen := licenses[l.Name]; seen {
						x.To = v3_0.ElementList{first}
					} else {
						licenses[l.Name] = l
					}
				}
			}
		}
	}
}

func check(err error) {
	if err != nil {
		fatal("%v", err)
	}
}

func fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "sbomgen: "+format+"\n", a...)
	os.Exit(1)
}
