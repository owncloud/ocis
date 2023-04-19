package icapclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Request represents the icap client request data
type Request struct {
	Method                string
	URL                   *url.URL
	Header                http.Header
	HTTPRequest           *http.Request
	HTTPResponse          *http.Response
	ChunkLength           int
	PreviewBytes          int
	ctx                   *context.Context
	previewSet            bool
	bodyFittedInPreview   bool
	remainingPreviewBytes []byte
}

// NewRequest is the factory function for Request
func NewRequest(method, urlStr string, httpReq *http.Request, httpResp *http.Response) (*Request, error) {

	method = strings.ToUpper(method)

	u, err := url.Parse(urlStr)

	if err != nil {
		return nil, err
	}

	req := &Request{
		Method:       method,
		URL:          u,
		Header:       make(map[string][]string),
		HTTPRequest:  httpReq,
		HTTPResponse: httpResp,
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return req, nil
}

// DumpRequest returns the given request in its ICAP/1.x wire
// representation.
func DumpRequest(req *Request) ([]byte, error) {

	// Making the ICAP message block

	reqStr := fmt.Sprintf("%s %s %s%s", req.Method, req.URL.String(), ICAPVersion, CRLF)

	for headerName, vals := range req.Header {
		for _, val := range vals {
			reqStr += fmt.Sprintf("%s: %s%s", headerName, val, CRLF)
		}
	}

	reqStr += "Encapsulated: %s" + CRLF // will populate the Encapsulated header value after making the http Request & Response messages
	reqStr += CRLF

	// Making the HTTP Request message block

	httpReqStr := ""
	if req.HTTPRequest != nil {
		b, err := httputil.DumpRequestOut(req.HTTPRequest, true)

		if err != nil {
			return nil, err
		}

		httpReqStr += string(b)
		replaceRequestURIWithActualURL(&httpReqStr, req.HTTPRequest.URL.EscapedPath(), req.HTTPRequest.URL.String())

		if req.Method == MethodREQMOD {
			if req.previewSet {
				parsePreviewBodyBytes(&httpReqStr, req.PreviewBytes)
			}

			if !bodyAlreadyChunked(httpReqStr) {
				headerStr, bodyStr, ok := splitBodyAndHeader(httpReqStr)
				if ok {
					addHexaBodyByteNotations(&bodyStr)
					mergeHeaderAndBody(&httpReqStr, headerStr, bodyStr)
				}
			}

		}

		if httpReqStr != "" { // if the HTTP Request message block doesn't end with a \r\n\r\n, then going to add one by force for better calculation of byte offsets
			for !strings.HasSuffix(httpReqStr, DoubleCRLF) {
				httpReqStr += CRLF
			}
		}

	}

	// Making the HTTP Response message block

	httpRespStr := ""
	if req.HTTPResponse != nil {
		b, err := httputil.DumpResponse(req.HTTPResponse, true)

		if err != nil {
			return nil, err
		}

		httpRespStr += string(b)

		if req.previewSet {
			parsePreviewBodyBytes(&httpRespStr, req.PreviewBytes)
		}

		if !bodyAlreadyChunked(httpRespStr) {
			headerStr, bodyStr, ok := splitBodyAndHeader(httpRespStr)
			if ok {
				addHexaBodyByteNotations(&bodyStr)
				mergeHeaderAndBody(&httpRespStr, headerStr, bodyStr)
			}
		}

		if httpRespStr != "" && !strings.HasSuffix(httpRespStr, DoubleCRLF) { // if the HTTP Response message block doesn't end with a \r\n\r\n, then going to add one by force for better calculation of byte offsets
			httpRespStr += CRLF
		}

	}

	if encpVal := req.Header.Get(EncapsulatedHeader); encpVal != "" {
		reqStr = fmt.Sprintf(reqStr, encpVal)
	} else {
		//populating the Encapsulated header of the ICAP message portion
		setEncapsulatedHeaderValue(&reqStr, httpReqStr, httpRespStr)
	}

	// determining if the http message needs the full body fitted in the preview portion indicator or not
	if httpRespStr != "" && req.previewSet && req.bodyFittedInPreview {
		addFullBodyInPreviewIndicator(&httpRespStr)
	}

	if req.Method == MethodREQMOD && req.previewSet && req.bodyFittedInPreview {
		addFullBodyInPreviewIndicator(&httpReqStr)
	}

	data := []byte(reqStr + httpReqStr + httpRespStr)

	return data, nil
}

// SetContext sets a context for the ICAP request
func (r *Request) SetContext(ctx context.Context) { // TODO: make context take control over the whole operation
	r.ctx = &ctx
}
