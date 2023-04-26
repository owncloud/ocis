/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tika

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// ClientError is returned by Client's various parse methods and
// represents an error response from the Tika server. Example usage:
//
//    client := tika.NewClient(nil, tikaURL)
//    s, err := client.Parse(context.Background(), input)
//    var tikaErr tika.ClientError
//    if errors.As(err, &tikaErr) {
//        switch tikaErr.StatusCode {
//        case http.StatusUnsupportedMediaType, http.StatusUnprocessableEntity:
//            // Handle content related error
//        default:
//            // Handle possibly intermittent http error
//        }
//    } else if err != nil {
//        // Handle non-http error
//    }
type ClientError struct {
	// StatusCode is the HTTP status code returned by the Tika server.
	StatusCode int
}

func (e ClientError) Error() string {
	return fmt.Sprintf("response code %d", e.StatusCode)
}

// Client represents a connection to a Tika Server.
type Client struct {
	// url is the URL of the Tika Server, including the port (if necessary), but
	// not the trailing slash. For example, http://localhost:9998.
	url string
	// HTTPClient is the client that will be used to call the Tika Server. If no
	// client is specified, a default client will be used. Since http.Clients are
	// thread safe, the same client will be used for all requests by this Client.
	httpClient *http.Client
}

// NewClient creates a new Client. If httpClient is nil, the http.DefaultClient will be
// used.
func NewClient(httpClient *http.Client, urlString string) *Client {
	return &Client{httpClient: httpClient, url: urlString}
}

// A Parser represents a Tika Parser. To get a list of all Parsers, see Parsers().
type Parser struct {
	Name           string
	Decorated      bool
	Composite      bool
	Children       []Parser
	SupportedTypes []string
}

// MIMEType represents a Tika MIME Type. To get a list of all MIME Types, see
// MIMETypes.
type MIMEType struct {
	Alias     []string
	SuperType string
}

// A Detector represents a Tika Detector. Detectors are used to get the filetype
// of a file. To get a list of all Detectors, see Detectors().
type Detector struct {
	Name      string
	Composite bool
	Children  []Detector
}

// Translator represents the Java package of a Tika Translator.
type Translator string

// Translators available by default in Tika. You must configure all required
// authentication details in Tika Server (for example, an API key).
const (
	Lingo24Translator   Translator = "org.apache.tika.language.translate.Lingo24Translator"
	GoogleTranslator    Translator = "org.apache.tika.language.translate.GoogleTranslator"
	MosesTranslator     Translator = "org.apache.tika.language.translate.MosesTranslator"
	JoshuaTranslator    Translator = "org.apache.tika.language.translate.JoshuaTranslator"
	MicrosoftTranslator Translator = "org.apache.tika.language.translate.MicrosoftTranslator"
	YandexTranslator    Translator = "org.apache.tika.language.translate.YandexTranslator"
)

// XTIKAContent is the metadata field of the content of a file after recursive
// parsing. See ParseRecursive and MetaRecursive.
const XTIKAContent = "X-TIKA:content"

// call makes the given request to c and returns the response body.
// call returns an error and a nil reader if the response code is not 200 StatusOK.
func (c *Client) call(ctx context.Context, input io.Reader, method, path string, header http.Header) (io.ReadCloser, error) {
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, method, c.url+path, input)
	if err != nil {
		return nil, err
	}
	req.Header = header

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, ClientError{resp.StatusCode}
	}
	return resp.Body, nil
}

// callString makes the given request to c and returns the result as a string
// and error. callString returns an error if the response code is not 200 StatusOK.
func (c *Client) callString(ctx context.Context, input io.Reader, method, path string, header http.Header) (string, error) {
	body, err := c.call(ctx, input, method, path, header)
	if err != nil {
		return "", err
	}
	defer body.Close()

	b := &strings.Builder{}
	if _, err := io.Copy(b, body); err != nil {
		return "", err
	}

	return b.String(), nil
}

// Parse parses the given input, returning the body of the input as a string and an error.
// If the error is not nil, the body is undefined.
func (c *Client) Parse(ctx context.Context, input io.Reader) (string, error) {
	return c.ParseWithHeader(ctx, input, nil)
}

// ParseReader parses the given input, returning the body of the input as a reader and an error.
// If the error is nil, the returned reader must be closed, else, the reader is nil.
func (c *Client) ParseReader(ctx context.Context, input io.Reader) (io.ReadCloser, error) {
	return c.ParseReaderWithHeader(ctx, input, nil)
}

// ParseWithHeader parses the given input, returning the body of the input as a string and an error.
// If the error is not nil. the body is undefined.
// This function also accepts a header so the caller can specify things like `Accept`
func (c *Client) ParseWithHeader(ctx context.Context, input io.Reader, header http.Header) (string, error) {
	return c.callString(ctx, input, "PUT", "/tika", header)
}

// ParseReaderWithHeader parses the given input, returning the body of the input as a reader and an error.
// If the error is nil, the returned reader must be closed, else, the reader is nil.
// This function also accepts a header so the caller can specify things like `Accept`
func (c *Client) ParseReaderWithHeader(ctx context.Context, input io.Reader, header http.Header) (io.ReadCloser, error) {
	return c.call(ctx, input, "PUT", "/tika", header)
}

// ParseRecursive parses the given input and all embedded documents, returning a
// list of the contents of the input with one element per document. See
// MetaRecursive for access to all metadata fields. If the error is not nil, the
// result is undefined.
func (c *Client) ParseRecursive(ctx context.Context, input io.Reader) ([]string, error) {
	m, err := c.MetaRecursive(ctx, input)
	if err != nil {
		return nil, err
	}
	var r []string
	for _, d := range m {
		if content := d[XTIKAContent]; len(content) > 0 {
			r = append(r, content[0])
		}
	}
	return r, nil
}

// Meta parses the metadata from the given input, returning the metadata and an
// error. If the error is not nil, the metadata is undefined.
func (c *Client) Meta(ctx context.Context, input io.Reader) (string, error) {
	return c.MetaWithHeader(ctx, input, nil)
}

// MetaWithHeader parses the metadata from the given input, returning the metadata and an
// error. If the error is not nil, the metadata is undefined.
// This function also accepts a header so the caller can specify things like `Accept`
func (c *Client) MetaWithHeader(ctx context.Context, input io.Reader, header http.Header) (string, error) {
	return c.callString(ctx, input, "PUT", "/meta", header)
}

// MetaField parses the metadata from the given input and returns the given
// field. If the error is not nil, the result string is undefined.
func (c *Client) MetaField(ctx context.Context, input io.Reader, field string) (string, error) {
	return c.MetaFieldWithHeader(ctx, input, field, nil)
}

// MetaFieldWithHeader parses the metadata from the given input and returns the given
// field. If the error is not nil, the result string is undefined.
// This function also accepts a header so the caller can specify things like `Accept`
func (c *Client) MetaFieldWithHeader(ctx context.Context, input io.Reader, field string, header http.Header) (string, error) {
	return c.callString(ctx, input, "PUT", fmt.Sprintf("/meta/%v", field), header)
}

// Detect gets the mimetype of the given input, returning the mimetype and an
// error. If the error is not nil, the mimetype is undefined.
func (c *Client) Detect(ctx context.Context, input io.Reader) (string, error) {
	return c.callString(ctx, input, "PUT", "/detect/stream", nil)
}

// Language detects the language of the given input, returning the two letter
// language code and an error. If the error is not nil, the language is
// undefined.
func (c *Client) Language(ctx context.Context, input io.Reader) (string, error) {
	return c.callString(ctx, input, "PUT", "/language/stream", nil)
}

// LanguageString detects the language of the given string, returning the two letter
// language code and an error. If the error is not nil, the language is
// undefined.
func (c *Client) LanguageString(ctx context.Context, input string) (string, error) {
	r := strings.NewReader(input)
	return c.callString(ctx, r, "PUT", "/language/string", nil)
}

// MetaRecursive parses the given input and all embedded documents. The result
// is a list of maps from metadata key to value for each document. The content
// of each document is in the XTIKAContent field in text form. See
// ParseRecursive to just get the content of each document. If the error is not
// nil, the result list is undefined.
func (c *Client) MetaRecursive(ctx context.Context, input io.Reader) ([]map[string][]string, error) {
	return c.MetaRecursiveType(ctx, input, "text")
}

// MetaRecursiveType parses the given input and all embedded documents. The result
// is a list of maps from metadata key to value for each document. The content
// of each document is in the XTIKAContent field, and is of the type indicated
// by the contentType parameter An empty string can be passed in for a default
// type of XML. See ParseRecursive to just get the content of each document. If
// the error is not nil, the result list is undefined.
func (c *Client) MetaRecursiveType(ctx context.Context, input io.Reader, contentType string) ([]map[string][]string, error) {
	path := "/rmeta"
	if contentType != "" {
		path = fmt.Sprintf("/rmeta/%s", contentType)
	}
	body, err := c.call(ctx, input, "PUT", path, nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	var m []map[string]interface{}
	if err := json.NewDecoder(body).Decode(&m); err != nil {
		return nil, err
	}
	var r []map[string][]string
	for _, d := range m {
		doc := make(map[string][]string)
		r = append(r, doc)
		for k, v := range d {
			switch vt := v.(type) {
			case string:
				doc[k] = []string{vt}
			case []interface{}:
				for _, i := range vt {
					s, ok := i.(string)
					if !ok {
						return nil, fmt.Errorf("field %q has value %v and type %T, expected a string or []string", k, v, vt)
					}
					doc[k] = append(doc[k], s)
				}
			default:
				return nil, fmt.Errorf("field %q has value %v and type %v, expected a string or []string", k, v, reflect.TypeOf(v))
			}
		}
	}
	return r, nil
}

// Translate returns an error and the translated input from src language to
// dst language using t. If the error is not nil, the translation is undefined.
func (c *Client) Translate(ctx context.Context, input io.Reader, t Translator, src, dst string) (string, error) {
	return c.callString(ctx, input, "POST", fmt.Sprintf("/translate/all/%s/%s/%s", t, src, dst), nil)
}

// TranslateReader translates the given input from src language to dst language using t.
// It returns the translated document as a reader. If an error occurs, the reader is nil, else, the reader
// must be closed by the caller after usage.
func (c *Client) TranslateReader(ctx context.Context, input io.Reader, t Translator, src, dst string) (io.ReadCloser, error) {
	return c.call(ctx, input, "POST", fmt.Sprintf("/translate/all/%s/%s/%s", t, src, dst), nil)
}

// Version returns the default hello message from Tika server.
func (c *Client) Version(ctx context.Context) (string, error) {
	return c.callString(ctx, nil, "GET", "/version", nil)
}

var jsonHeader = http.Header{"Accept": []string{"application/json"}}

// callUnmarshal is like call, but unmarshals the JSON response into v.
func (c *Client) callUnmarshal(ctx context.Context, path string, v interface{}) error {
	body, err := c.call(ctx, nil, "GET", path, jsonHeader)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(v)
}

// Parsers returns the list of available parsers and an error. If the error is
// not nil, the list is undefined. To get all available parsers, iterate through
// the Children of every Parser.
func (c *Client) Parsers(ctx context.Context) (*Parser, error) {
	p := new(Parser)
	if err := c.callUnmarshal(ctx, "/parsers/details", p); err != nil {
		return nil, err
	}
	return p, nil
}

// MIMETypes returns a map from MIME Type name to MIMEType, or properties about
// that specific MIMEType.
func (c *Client) MIMETypes(ctx context.Context) (map[string]MIMEType, error) {
	mt := make(map[string]MIMEType)
	if err := c.callUnmarshal(ctx, "/mime-types", &mt); err != nil {
		return nil, err
	}
	return mt, nil
}

// Detectors returns the list of available Detectors for this server. To get all
// available detectors, iterate through the Children of every Detector.
func (c *Client) Detectors(ctx context.Context) (*Detector, error) {
	d := new(Detector)
	if err := c.callUnmarshal(ctx, "/detectors", d); err != nil {
		return nil, err
	}
	return d, nil
}
