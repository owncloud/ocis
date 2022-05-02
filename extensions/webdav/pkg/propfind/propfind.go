package propfind

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/owncloud/ocis/extensions/webdav/pkg/errors"
	"github.com/owncloud/ocis/extensions/webdav/pkg/prop"
)

const (
	_spaceTypeProject = "project"
)

type countingReader struct {
	n int
	r io.Reader
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += n
	return n, err
}

// Props represents properties related to a resource
// http://www.webdav.org/specs/rfc4918.html#ELEMENT_prop (for propfind)
type Props []xml.Name

// XML holds the xml representation of a propfind
// http://www.webdav.org/specs/rfc4918.html#ELEMENT_propfind
type XML struct {
	XMLName  xml.Name  `xml:"DAV: propfind"`
	Allprop  *struct{} `xml:"DAV: allprop"`
	Propname *struct{} `xml:"DAV: propname"`
	Prop     Props     `xml:"DAV: prop"`
	Include  Props     `xml:"DAV: include"`
}

// PropstatXML holds the xml representation of a propfind response
// http://www.webdav.org/specs/rfc4918.html#ELEMENT_propstat
type PropstatXML struct {
	// Prop requires DAV: to be the default namespace in the enclosing
	// XML. This is due to the standard encoding/xml package currently
	// not honoring namespace declarations inside a xmltag with a
	// parent element for anonymous slice elements.
	// Use of multistatusWriter takes care of this.
	Prop                []prop.PropertyXML `xml:"d:prop>_ignored_"`
	Status              string             `xml:"d:status"`
	Error               *errors.ErrorXML   `xml:"d:error"`
	ResponseDescription string             `xml:"d:responsedescription,omitempty"`
}

// ResponseXML holds the xml representation of a propfind response
type ResponseXML struct {
	XMLName             xml.Name         `xml:"d:response"`
	Href                string           `xml:"d:href"`
	Propstat            []PropstatXML    `xml:"d:propstat"`
	Status              string           `xml:"d:status,omitempty"`
	Error               *errors.ErrorXML `xml:"d:error"`
	ResponseDescription string           `xml:"d:responsedescription,omitempty"`
}

// MultiStatusResponseXML holds the xml representation of a multistatus propfind response
type MultiStatusResponseXML struct {
	XMLName xml.Name `xml:"d:multistatus"`
	XmlnsS  string   `xml:"xmlns:s,attr,omitempty"`
	XmlnsD  string   `xml:"xmlns:d,attr,omitempty"`
	XmlnsOC string   `xml:"xmlns:oc,attr,omitempty"`

	Responses []*ResponseXML `xml:"d:response"`
}

// ResponseUnmarshalXML is a workaround for https://github.com/golang/go/issues/13400
type ResponseUnmarshalXML struct {
	XMLName             xml.Name               `xml:"response"`
	Href                string                 `xml:"href"`
	Propstat            []PropstatUnmarshalXML `xml:"propstat"`
	Status              string                 `xml:"status,omitempty"`
	Error               *errors.ErrorXML       `xml:"d:error"`
	ResponseDescription string                 `xml:"responsedescription,omitempty"`
}

// MultiStatusResponseUnmarshalXML is a workaround for https://github.com/golang/go/issues/13400
type MultiStatusResponseUnmarshalXML struct {
	XMLName xml.Name `xml:"multistatus"`
	XmlnsS  string   `xml:"xmlns:s,attr,omitempty"`
	XmlnsD  string   `xml:"xmlns:d,attr,omitempty"`
	XmlnsOC string   `xml:"xmlns:oc,attr,omitempty"`

	Responses []*ResponseUnmarshalXML `xml:"response"`
}

// PropstatUnmarshalXML is a workaround for https://github.com/golang/go/issues/13400
type PropstatUnmarshalXML struct {
	// Prop requires DAV: to be the default namespace in the enclosing
	// XML. This is due to the standard encoding/xml package currently
	// not honoring namespace declarations inside a xmltag with a
	// parent element for anonymous slice elements.
	// Use of multistatusWriter takes care of this.
	Prop                []*prop.PropertyXML `xml:"prop"`
	Status              string              `xml:"status"`
	Error               *errors.ErrorXML    `xml:"d:error"`
	ResponseDescription string              `xml:"responsedescription,omitempty"`
}

// NewMultiStatusResponseXML returns a preconfigured instance of MultiStatusResponseXML
func NewMultiStatusResponseXML() *MultiStatusResponseXML {
	return &MultiStatusResponseXML{
		XmlnsD:  "DAV:",
		XmlnsS:  "http://sabredav.org/ns",
		XmlnsOC: "http://owncloud.org/ns",
	}
}

// ReadPropfind extracts and parses the propfind XML information from a Reader
// from https://github.com/golang/net/blob/e514e69ffb8bc3c76a71ae40de0118d794855992/webdav/xml.go#L178-L205
func ReadPropfind(r io.Reader) (pf XML, status int, err error) {
	c := countingReader{r: r}
	if err = xml.NewDecoder(&c).Decode(&pf); err != nil {
		if err == io.EOF {
			if c.n == 0 {
				// An empty body means to propfind allprop.
				// http://www.webdav.org/specs/rfc4918.html#METHOD_PROPFIND
				return XML{Allprop: new(struct{})}, 0, nil
			}
			err = errors.ErrInvalidPropfind
		}
		return XML{}, http.StatusBadRequest, err
	}

	if pf.Allprop == nil && pf.Include != nil {
		return XML{}, http.StatusBadRequest, errors.ErrInvalidPropfind
	}
	if pf.Allprop != nil && (pf.Prop != nil || pf.Propname != nil) {
		return XML{}, http.StatusBadRequest, errors.ErrInvalidPropfind
	}
	if pf.Prop != nil && pf.Propname != nil {
		return XML{}, http.StatusBadRequest, errors.ErrInvalidPropfind
	}
	if pf.Propname == nil && pf.Allprop == nil && pf.Prop == nil {
		// jfd: I think <d:prop></d:prop> is perfectly valid ... treat it as allprop
		return XML{Allprop: new(struct{})}, 0, nil
	}
	return pf, 0, nil
}

// UnmarshalXML appends the property names enclosed within start to pn.
//
// It returns an error if start does not contain any properties or if
// properties contain values. Character data between properties is ignored.
func (pn *Props) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := prop.Next(d)
		if err != nil {
			return err
		}
		switch e := t.(type) {
		case xml.EndElement:
			// jfd: I think <d:prop></d:prop> is perfectly valid ... treat it as allprop
			/*
				if len(*pn) == 0 {
					return fmt.Errorf("%s must not be empty", start.Name.Local)
				}
			*/
			return nil
		case xml.StartElement:
			t, err = prop.Next(d)
			if err != nil {
				return err
			}
			if _, ok := t.(xml.EndElement); !ok {
				return fmt.Errorf("unexpected token %T", t)
			}
			*pn = append(*pn, e.Name)
		}
	}
}
