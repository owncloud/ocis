package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type message struct {
	headers        textproto.MIMEHeader
	body           *bytes.Buffer
	writers        []*multipart.Writer
	parts          uint8
	cids           map[string]string
	charset        string
	encoding       encoding
	headerEncoding headerEncoding // Only None and Q are currently supported
}

func newMessage(email *Email) *message {
	return &message{
		headers:        email.headers,
		body:           new(bytes.Buffer),
		cids:           make(map[string]string),
		charset:        email.Charset,
		encoding:       email.Encoding,
		headerEncoding: email.HeaderEncoding}
}

func encodeHeader(text string, charset string, encoding headerEncoding, usedChars int) string {
	// create buffer
	buf := new(bytes.Buffer)

	// encode
	encoder := newEncoder(buf, charset, encoding, usedChars)
	encoder.encode([]byte(text))

	return buf.String()
}

// encodeQuotedParam encodes text like encodeHeader, for a parameter value
// written inside double quotes.
func encodeQuotedParam(text string, charset string, encoding headerEncoding, usedChars int) string {
	buf := new(bytes.Buffer)

	encoder := newEncoder(buf, charset, encoding, usedChars)
	encoder.quoted = true
	encoder.encode([]byte(text))

	return buf.String()
}

// getHeaders returns the message headers
func (msg *message) getHeaders() (headers string) {
	// if the date header isn't set, set it
	if date := msg.headers.Get("Date"); date == "" {
		msg.headers.Set("Date", time.Now().Format(time.RFC1123Z))
	}

	// encode and combine the headers
	for header, values := range msg.headers {
		encoded := encodeHeader(strings.Join(values, ", "), msg.charset, msg.headerEncoding, len(header)+2)
		headers += header + ": " + encoded + "\r\n"
	}

	headers = headers + "\r\n"

	return
}

// getCID gets the generated CID for the provided text
func (msg *message) getCID(text string) (cid string) {
	// set the date format to use
	const dateFormat = "20060102.150405"

	// get the cid if we have one
	cid, exists := msg.cids[text]
	if !exists {
		// generate a new cid
		cid = time.Now().Format(dateFormat) + "." + strconv.Itoa(len(msg.cids)+1) + "@mail.0"
		// save it
		msg.cids[text] = cid
	}

	return
}

// replaceCIDs replaces the CIDs found in a text string
// with generated ones
func (msg *message) replaceCIDs(input []byte) []byte {
	if len(msg.cids) == 0 {
		// One process replaceCIDs if we have anything to replace
		return input
	}

	text := string(input)

	// regular expression to find cids
	re := regexp.MustCompile(`(src|href)="cid:(.*?)"`)
	// replace all of the found cids with generated ones
	for _, matches := range re.FindAllStringSubmatch(text, -1) {
		cid := msg.getCID(matches[2])
		text = strings.Replace(text, "cid:"+matches[2], "cid:"+cid, -1)
	}

	return []byte(text)
}

// openMultipart creates a new part of a multipart message
func (msg *message) openMultipart(multipartType string) {
	// create a new multipart writer
	msg.writers = append(msg.writers, multipart.NewWriter(msg.body))
	// create the boundary
	contentType := "multipart/" + multipartType + ";\n \tboundary=" + msg.writers[msg.parts].Boundary()

	// if no existing parts, add header to main header group
	if msg.parts == 0 {
		msg.headers.Set("Content-Type", contentType)
	} else { // add header to multipart section
		header := make(textproto.MIMEHeader)
		header.Set("Content-Type", contentType)
		msg.writers[msg.parts-1].CreatePart(header)
	}

	msg.parts++
}

// closeMultipart closes a part of a multipart message
func (msg *message) closeMultipart() {
	if msg.parts > 0 {
		msg.writers[msg.parts-1].Close()
		msg.parts--
	}
}

// base64Encode base64 encodes the provided text with line wrapping
func base64Encode(text []byte) []byte {
	// create buffer
	buf := new(bytes.Buffer)

	// create base64 encoder that linewraps
	encoder := base64.NewEncoder(base64.StdEncoding, &base64LineWrap{writer: buf})

	// write the encoded text to buf
	encoder.Write(text)
	encoder.Close()

	return buf.Bytes()
}

// qpEncode uses the quoted-printable encoding to encode the provided text
func qpEncode(text []byte) []byte {
	// create buffer
	buf := new(bytes.Buffer)

	encoder := quotedprintable.NewWriter(buf)

	encoder.Write(text)
	encoder.Close()

	return buf.Bytes()
}

const maxLineChars = 76

type base64LineWrap struct {
	writer       io.Writer
	numLineChars int
}

func (e *base64LineWrap) Write(p []byte) (n int, err error) {
	n = 0
	// while we have more chars than are allowed
	for len(p)+e.numLineChars > maxLineChars {
		numCharsToWrite := maxLineChars - e.numLineChars
		// write the chars we can
		e.writer.Write(p[:numCharsToWrite])
		// write a line break
		e.writer.Write([]byte("\r\n"))
		// reset the line count
		e.numLineChars = 0
		// remove the chars that have been written
		p = p[numCharsToWrite:]
		// set the num of chars written
		n += numCharsToWrite
	}

	// write what is left
	e.writer.Write(p)
	e.numLineChars += len(p)
	n += len(p)

	return
}

func (msg *message) write(header textproto.MIMEHeader, body []byte, encoding encoding) {
	msg.writeHeader(header)
	msg.writeBody(body, encoding)
}

func (msg *message) writeHeader(headers textproto.MIMEHeader) {
	// if there are no parts add header to main headers
	if msg.parts == 0 {
		for header, value := range headers {
			msg.headers[header] = value
		}
	} else { // add header to multipart section
		msg.writers[msg.parts-1].CreatePart(headers)
	}
}

func (msg *message) writeBody(body []byte, encoding encoding) {
	// encode and write the body
	switch encoding {
	case EncodingQuotedPrintable:
		msg.body.Write(qpEncode(body))
	case EncodingBase64:
		msg.body.Write(base64Encode(body))
	default:
		msg.body.Write(body)
	}
}

func (msg *message) addBody(contentType string, body []byte) {
	body = msg.replaceCIDs(body)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", contentType+"; charset="+msg.charset)
	header.Set("Content-Transfer-Encoding", msg.encoding.string())
	msg.write(header, body, msg.encoding)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func (msg *message) addFiles(files []*File, inline bool) {
	encoding := EncodingBase64
	for _, file := range files {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Type",
			fmt.Sprintf("%s;\n \tname=\"%s\"",
				file.MimeType,
				encodeQuotedParam(escapeQuotes(file.Name), msg.charset, msg.headerEncoding, 6)))
		header.Set("Content-Transfer-Encoding", encoding.string())

		encodedFilename := encodeQuotedParam(escapeQuotes(file.Name), msg.charset, msg.headerEncoding, 10)

		if inline {
			header.Set("Content-Disposition", "inline;\n \tfilename=\""+encodedFilename+`"`)
			if len(file.ContentID) > 0 {
				header.Set("Content-ID", "<"+file.ContentID+">")
			} else {
				header.Set("Content-ID", "<"+msg.getCID(file.Name)+">")
			}
		} else {
			header.Set("Content-Disposition", "attachment;\n \tfilename=\""+encodedFilename+`"`)
		}

		msg.write(header, file.Data, encoding)
	}
}
