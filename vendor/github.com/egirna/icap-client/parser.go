package icapclient

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// getStatusWithCode prepares the status code and status text from two given strings
func getStatusWithCode(str1, str2 string) (int, string, error) {

	statusCode, err := strconv.Atoi(str1)

	if err != nil {
		return 0, "", err
	}

	status := strings.TrimSpace(str2)

	return statusCode, status, nil
}

// getHeaderVal parses the header and its value from a tcp message string
func getHeaderVal(str string) (string, string) {

	headerVals := strings.SplitN(str, ":", 2)
	header := headerVals[0]
	val := ""

	if len(headerVals) >= 2 {
		val = strings.TrimSpace(headerVals[1])
	}

	return header, val

}

// isRequestLine determines if the tcp message string is a request line, i.e the first line of the message or not
func isRequestLine(str string) bool {
	return strings.Contains(str, ICAPVersion) || strings.Contains(str, HTTPVersion)
}

// setEncapsulatedHeaderValue generates the Encapsulated  values and assigns to the ICAP request string
func setEncapsulatedHeaderValue(icapReqStr *string, httpReqStr, httpRespStr string) {
	encpVal := " "

	if strings.HasPrefix(*icapReqStr, MethodOPTIONS) { // if the request method is OPTIONS
		if httpReqStr == "" && httpRespStr == "" { // the most common case for OPTIONS method, no Encapsulated body
			encpVal += "null-body=0"
		} else {
			encpVal += "opt-body=0" // if there is an Encapsulated body
		}
	}

	if strings.HasPrefix(*icapReqStr, MethodREQMOD) || strings.HasPrefix(*icapReqStr, MethodRESPMOD) { // if the request method is RESPMOD or REQMOD
		re := regexp.MustCompile(DoubleCRLF)                // looking for the match of the string \r\n\r\n, as that is the expression that seperates each blocks, i.e headers and bodies
		reqIndices := re.FindAllStringIndex(httpReqStr, -1) // getting the offsets of the matches, tells us the starting/ending point of headers and bodies

		reqEndsAt := 0 // this is needed to calculate the response headers by adding the last offset of the request block
		if reqIndices != nil {
			encpVal += "req-hdr=0"
			reqEndsAt = reqIndices[0][1]
			if len(reqIndices) > 1 { // indicating there is a body present for the request block, as length would have been 1 for a single match of \r\n\r\n
				encpVal += fmt.Sprintf(", req-body=%d", reqIndices[0][1]) // assigning the starting point of the body
				reqEndsAt = reqIndices[1][1]
			} else if httpRespStr == "" {
				encpVal += fmt.Sprintf(", null-body=%d", reqIndices[0][1])
			}
			if httpRespStr != "" {
				encpVal += ", "
			}
		}

		respIndices := re.FindAllStringIndex(httpRespStr, -1)

		if respIndices != nil {
			encpVal += fmt.Sprintf("res-hdr=%d", reqEndsAt)
			if len(respIndices) > 1 {
				encpVal += fmt.Sprintf(", res-body=%d", reqEndsAt+respIndices[0][1])
			} else {
				encpVal += fmt.Sprintf(", null-body=%d", reqEndsAt+respIndices[0][1])
			}
		}

	}

	*icapReqStr = fmt.Sprintf(*icapReqStr, encpVal) // formatting the ICAP request Encapsulated header with the value
}

// replaceRequestURIWithActualURL replaces the just the escaped portion of the url with the entire URL in the dumped request message
func replaceRequestURIWithActualURL(str *string, uri, url string) {
	if uri == "" {
		uri = "/"
	}
	*str = strings.Replace(*str, uri, url, 1)
}

// addFullBodyInPreviewIndicator adds 0; ieof\r\n\r\n which indicates the entire body fitted in preview
func addFullBodyInPreviewIndicator(str *string) {
	*str = strings.TrimSuffix(*str, DoubleCRLF)
	*str += fullBodyEndIndicatorPreviewMode
}

// splitBodyAndHeader separates header and body from a http message
func splitBodyAndHeader(str string) (string, string, bool) {
	ss := strings.SplitN(str, DoubleCRLF, 2)

	if len(ss) < 2 || ss[1] == "" {
		return "", "", false
	}

	headerStr := ss[0]
	bodyStr := ss[1]

	return headerStr, bodyStr, true
}

// bodyAlreadyChunked determines if the http body is already chunked from the origin server or not
func bodyAlreadyChunked(str string) bool {
	_, bodyStr, ok := splitBodyAndHeader(str)

	if !ok {
		return false
	}

	r := regexp.MustCompile("\\r\\n0(\\r\\n)+$")
	return r.MatchString(bodyStr)

}

// parsePreviewBodyBytes parses the preview portion of the body and only keeps that in the message
func parsePreviewBodyBytes(str *string, pb int) {

	headerStr, bodyStr, ok := splitBodyAndHeader(*str)

	if !ok {
		return
	}

	bodyStr = bodyStr[:pb]

	*str = headerStr + DoubleCRLF + bodyStr
}

// addHexaBodyByteNotations adds the hexadecimal byte notaions in the messages
// for example: Hello World, becomes
// b
// Hello World
// 0
func addHexaBodyByteNotations(bodyStr *string) {

	bodyBytes := []byte(*bodyStr)

	*bodyStr = fmt.Sprintf("%x%s%s%s", len(bodyBytes), CRLF, *bodyStr, bodyEndIndicator)
}

// mergeHeaderAndBody merges the header and body of the http message togather
func mergeHeaderAndBody(src *string, headerStr, bodyStr string) {
	*src = headerStr + DoubleCRLF + bodyStr
}

func chunkBodyByBytes(bdyByte []byte, cl int) []byte {

	newBytes := []byte{}

	for i := 0; i < len(bdyByte); i += cl {
		end := i + cl
		if end > len(bdyByte) {
			end = len(bdyByte)
		}

		newBytes = append(newBytes, []byte(fmt.Sprintf("%x\r\n", len(bdyByte[i:end]))+string(bdyByte[i:end]))...)
	}

	newBytes = append(newBytes, []byte(bodyEndIndicator)...)

	return newBytes
}
