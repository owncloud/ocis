package icapclient

import (
	"bufio"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// Response represents the icap server response data
type Response struct {
	StatusCode      int
	Status          string
	PreviewBytes    int
	Header          http.Header
	ContentRequest  *http.Request
	ContentResponse *http.Response
}

var (
	optionValues = map[string]bool{
		PreviewHeader:          true,
		MethodsHeader:          true,
		AllowHeader:            true,
		TransferPreviewHeader:  true,
		ServiceHeader:          true,
		ISTagHeader:            true,
		OptBodyTypeHeader:      true,
		MaxConnectionsHeader:   true,
		OptionsTTLHeader:       true,
		ServiceIDHeader:        true,
		TransferIgnoreHeader:   true,
		TransferCompleteHeader: true,
	}
)

// ReadResponse converts a Reader to a icapclient Response
func ReadResponse(b *bufio.Reader) (*Response, error) {

	resp := &Response{
		Header: make(map[string][]string),
	}

	scheme := ""
	httpMsg := ""
	for currentMsg, err := b.ReadString('\n'); err == nil || currentMsg != ""; currentMsg, err = b.ReadString('\n') { // keep reading the buffer message which is the http response message

		if isRequestLine(currentMsg) { // if the current message line if the first line of the message portion(request line)
			ss := strings.Split(currentMsg, " ")

			if len(ss) < 3 { // must contain 3 words, for example: "ICAP/1.0 200 OK" or "GET /something HTTP/1.1"
				return nil, errors.New(ErrInvalidTCPMsg + ":" + currentMsg)
			}

			// preparing the scheme below

			if ss[0] == ICAPVersion {
				scheme = SchemeICAP
				resp.StatusCode, resp.Status, err = getStatusWithCode(ss[1], strings.Join(ss[2:], " "))
				if err != nil {
					return nil, err
				}
				continue
			}

			if ss[0] == HTTPVersion {
				scheme = SchemeHTTPResp
				httpMsg = ""
			}

			if strings.TrimSpace(ss[2]) == HTTPVersion { // for a http request message if the scheme version is always at last, for example: GET /something HTTP/1.1
				scheme = SchemeHTTPReq
				httpMsg = ""
			}
		}

		// preparing the header for ICAP & contents for HTTP messages below

		if scheme == SchemeICAP {
			if currentMsg == LF || currentMsg == CRLF { // don't want to count the Line Feed as header
				continue
			}
			header, val := getHeaderVal(currentMsg)
			if header == PreviewHeader {
				pb, _ := strconv.Atoi(val)
				resp.PreviewBytes = pb
			}
			resp.Header.Add(header, val)
		}

		if scheme == SchemeHTTPReq {
			httpMsg += strings.TrimSpace(currentMsg) + CRLF
			bufferEmpty := b.Buffered() == 0
			if currentMsg == CRLF || bufferEmpty { // a CRLF indicates the end of a http message and the buffer check is just in case the buffer eneded with one last message instead of a CRLF
				var erR error
				resp.ContentRequest, erR = http.ReadRequest(bufio.NewReader(strings.NewReader(httpMsg)))
				if erR != nil {
					return nil, erR
				}
				continue
			}
		}

		if scheme == SchemeHTTPResp {
			httpMsg += strings.TrimSpace(currentMsg) + CRLF
			bufferEmpty := b.Buffered() == 0
			if currentMsg == CRLF || bufferEmpty {
				var erR error
				resp.ContentResponse, erR = http.ReadResponse(bufio.NewReader(strings.NewReader(httpMsg)), resp.ContentRequest)
				if erR != nil {
					return nil, erR
				}
				continue
			}

		}

	}

	return resp, nil

}
