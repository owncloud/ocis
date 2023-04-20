package icapclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// SetPreview sets the preview bytes in the icap header
func (r *Request) SetPreview(maxBytes int) error {

	bodyBytes := []byte{}

	previewBytes := 0

	// receiving the body bites to determine the preview bytes depending on the request ICAP method

	if r.Method == MethodREQMOD {
		if r.HTTPRequest == nil {
			return nil
		}
		if r.HTTPRequest.Body != nil {
			var err error
			bodyBytes, err = ioutil.ReadAll(r.HTTPRequest.Body)

			if err != nil {
				return err
			}

			defer r.HTTPRequest.Body.Close()
		}
	}

	if r.Method == MethodRESPMOD {
		if r.HTTPResponse == nil {
			return nil
		}

		if r.HTTPResponse.Body != nil {
			var err error
			bodyBytes, err = ioutil.ReadAll(r.HTTPResponse.Body)

			if err != nil {
				return err
			}

			defer r.HTTPResponse.Body.Close()
		}
	}

	previewBytes = len(bodyBytes)

	if previewBytes > 0 { // if the preview byte is 0 or less, there is no question of the body fitting insides
		r.bodyFittedInPreview = true
	}

	if previewBytes > maxBytes { // if the preview bytes is greater than what was mentioned by the ICAP Server(did not fit in the body)
		previewBytes = maxBytes
		r.bodyFittedInPreview = false
		r.remainingPreviewBytes = bodyBytes[maxBytes:] // storing the rest of the body byte which were not sent as preview for further operations
	}

	// returning the body back to the http message depending on the request method

	if r.Method == MethodREQMOD {
		r.HTTPRequest.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	}

	if r.Method == MethodRESPMOD {
		r.HTTPResponse.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	}

	// finally assinging the preview informations including setting the header

	r.Header.Set("Preview", strconv.Itoa(previewBytes))
	r.PreviewBytes = previewBytes
	r.previewSet = true

	return nil

}

// SetDefaultRequestHeaders assigns some of the headers with its default value if they are not set already
func (r *Request) SetDefaultRequestHeaders() {
	if _, exists := r.Header["Allow"]; !exists {
		r.Header.Add("Allow", "204") // assigning 204 by default if Allow not provided
	}
	if _, exists := r.Header["Host"]; !exists {
		hostName, _ := os.Hostname()
		r.Header.Add("Host", hostName)
	}
}

// ExtendHeader extends the current ICAP Request header with a new header
func (r *Request) ExtendHeader(hdr http.Header) error {
	for header, values := range hdr {

		if header == PreviewHeader && r.previewSet {
			continue
		}

		if header == EncapsulatedHeader {
			continue
		}

		for _, value := range values {
			if header == PreviewHeader {
				pb, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				if err := r.SetPreview(pb); err != nil {
					return err
				}
				continue
			}
			r.Header.Add(header, value)
		}
	}

	return nil
}
