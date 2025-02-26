package gowebdav

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	pathpkg "path"
	"strconv"
	"strings"
	"time"
)

const XInhibitRedirect = "X-Gowebdav-Inhibit-Redirect"

var defaultProps = []string{"displayname", "resourcetype", "getcontentlength", "getcontenttype", "getetag", "getlastmodified"}

// Client defines our structure
type Client struct {
	root        string
	headers     http.Header
	interceptor func(method string, rq *http.Request)
	c           *http.Client
	auth        Authorizer
}

// NewClient creates a new instance of client
func NewClient(uri, user, pw string) *Client {
	return NewAuthClient(uri, NewAutoAuth(user, pw))
}

// NewAuthClient creates a new client instance with a custom Authorizer
func NewAuthClient(uri string, auth Authorizer) *Client {
	c := &http.Client{
		CheckRedirect: func(rq *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return ErrTooManyRedirects
			}
			if via[0].Header.Get(XInhibitRedirect) != "" {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	return &Client{root: FixSlash(uri), headers: make(http.Header), interceptor: nil, c: c, auth: auth}
}

// SetHeader lets us set arbitrary headers for a given client
func (c *Client) SetHeader(key, value string) {
	c.headers.Add(key, value)
}

// SetInterceptor lets us set an arbitrary interceptor for a given client
func (c *Client) SetInterceptor(interceptor func(method string, rq *http.Request)) {
	c.interceptor = interceptor
}

// SetTimeout exposes the ability to set a time limit for requests
func (c *Client) SetTimeout(timeout time.Duration) {
	c.c.Timeout = timeout
}

// SetTransport exposes the ability to define custom transports
func (c *Client) SetTransport(transport http.RoundTripper) {
	c.c.Transport = transport
}

// SetJar exposes the ability to set a cookie jar to the client.
func (c *Client) SetJar(jar http.CookieJar) {
	c.c.Jar = jar
}

// Connect connects to our dav server
func (c *Client) Connect() error {
	rs, err := c.options("/")
	if err != nil {
		return err
	}

	err = rs.Body.Close()
	if err != nil {
		return err
	}

	if rs.StatusCode != 200 {
		return NewPathError("Connect", c.root, rs.StatusCode)
	}

	return nil
}

type Props struct {
	m map[xml.Name]interface{}
}

func (c *Props) GetString(key xml.Name) string {
	return fmt.Sprintf("%v", c.m[key])
}

func (c *Props) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	c.m = make(map[xml.Name]interface{})

	type flatProp struct {
		XMLName xml.Name
		Content string `xml:",innerxml"`
	}
	type flatProps struct {
		Props []flatProp `xml:",any"`
	}
	parsedProps := flatProps{}
	if err := d.DecodeElement(&parsedProps, &start); err != nil {
		return err
	}
	for _, v := range parsedProps.Props {
		c.m[v.XMLName] = v.Content
	}
	return nil
}

type propstat struct {
	Status string `xml:"DAV: status"`

	Props Props `xml:"DAV: prop"`
}

func (p *propstat) Type() string {
	if p.Props.GetString(xml.Name{Space: "DAV:", Local: "resourcetype"}) == "" {
		return "file"
	}
	return "collection"
}

func (p *propstat) Name() string {
	return p.Props.GetString(xml.Name{Space: "DAV:", Local: "displayname"})
}

func (p *propstat) ContentType() string {
	return p.Props.GetString(xml.Name{Space: "DAV:", Local: "getcontenttype"})
}

func (p *propstat) Size() int64 {
	if n, e := strconv.ParseInt(p.Props.GetString(xml.Name{Space: "DAV:", Local: "getcontentlength"}), 10, 64); e == nil {
		return n
	}
	return 0
}

func (p *propstat) ETag() string {
	return p.Props.GetString(xml.Name{Space: "DAV:", Local: "getetag"})
}

func (p *propstat) Modified() time.Time {
	if t, e := time.Parse(time.RFC1123, p.Props.GetString(xml.Name{Space: "DAV:", Local: "getlastmodified"})); e == nil {
		return t
	}
	return time.Unix(0, 0)
}

func (p *propstat) Lock() string {
	return p.Props.GetString(xml.Name{Space: "DAV:", Local: "lockdiscovery"})
}

func (p *propstat) StatusCode() int {
	parts := strings.Split(p.Status, " ")
	if len(parts) < 2 {
		return -1
	}

	code, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1
	}

	return code
}

type response struct {
	Href      string     `xml:"DAV: href"`
	Propstats []propstat `xml:"DAV: propstat"`
}

func getPropstat(r *response, statuses []string) *propstat {
	for _, prop := range r.Propstats {
		for _, status := range statuses {
			if strings.Contains(prop.Status, status) {
				return &prop
			}
		}
	}
	return nil
}

// ReadDir reads the contents of a remote directory
func (c *Client) ReadDir(path string) ([]FileInfo, error) {
	return c.ReadDirWithProps(path, defaultProps)
}

// ReadDirWithProps reads the contents of the directory at the given path, along with the specified properties.
func (c *Client) ReadDirWithProps(path string, props []string) ([]FileInfo, error) {
	files := make([]FileInfo, 0)
	skipSelf := true
	parse := func(resp interface{}) error {
		r := resp.(*response)

		if skipSelf {
			skipSelf = false
			if p := getPropstat(r, []string{"200", "425"}); p != nil && p.Type() == "collection" {
				r.Propstats = nil
				return nil
			}
			return NewPathError("ReadDir", path, 405)
		}

		if p := getPropstat(r, []string{"200", "425"}); p != nil {
			var name string
			if ps, err := url.PathUnescape(r.Href); err == nil {
				name = pathpkg.Base(ps)
			} else {
				name = p.Name()
			}
			files = append(files, *newFile(path, name, p))
		}

		r.Propstats = nil
		return nil
	}

	propXML := "<d:propfind xmlns:d='DAV:'>"
	switch {
	case len(props) > 0:
		propXML += "<d:prop>"
		for _, prop := range props {
			propXML += "<d:" + prop + "/>"
		}
		propXML += "</d:prop>"
	default:
		propXML += "<allprop/>"
	}
	propXML += "</d:propfind>"

	err := c.propfind(path, false,
		propXML,
		&response{},
		parse)

	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			err = NewPathErrorErr("ReadDir", path, err)
		}
	}
	return files, err
}

// Stat returns the file stats for a specified path with the default properties
func (c *Client) Stat(path string) (FileInfo, error) {
	return c.StatWithProps(path, defaultProps)
}

// StatWithProps returns the FileInfo for the specified path along with the specified properties.
func (c *Client) StatWithProps(path string, props []string) (FileInfo, error) {
	var f *File
	parse := func(resp interface{}) error {
		r := resp.(*response)
		if p := getPropstat(r, []string{"200", "425"}); p != nil && f == nil {
			f = newFile(".", path, p)
		} else {
			return NewPathError("StatWithProps", path, 404)
		}

		r.Propstats = nil
		return nil
	}

	propXML := "<d:propfind xmlns:d='DAV:'>"
	switch {
	case len(props) > 0:
		propXML += "<d:prop>"
		for _, prop := range props {
			propXML += "<d:" + prop + "/>"
		}
		propXML += "</d:prop>"
	default:
		propXML += "<allprop/>"
	}
	propXML += "</d:propfind>"

	err := c.propfind(path, true,
		propXML,
		&response{},
		parse)

	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return nil, NewPathErrorErr("StatWithProps", path, err)
		}
		return nil, err
	}

	if f == nil {
		return nil, NewPathError("StatWithProps", path, 404)
	}
	return *f, err
}

// Remove removes a remote file
func (c *Client) Remove(path string) error {
	return c.RemoveAll(path)
}

// RemoveAll removes remote files
func (c *Client) RemoveAll(path string) error {
	rs, err := c.req("DELETE", path, nil, nil)
	if err != nil {
		return NewPathError("Remove", path, 400)
	}
	err = rs.Body.Close()
	if err != nil {
		return err
	}

	if rs.StatusCode == 200 || rs.StatusCode == 204 || rs.StatusCode == 404 {
		return nil
	}

	return NewPathError("Remove", path, rs.StatusCode)
}

// Mkdir makes a directory
func (c *Client) Mkdir(path string, _ os.FileMode) (err error) {
	path = FixSlashes(path)
	status, err := c.mkcol(path)
	if err != nil {
		return
	}
	if status == 201 {
		return nil
	}

	return NewPathError("Mkdir", path, status)
}

// MkdirAll like mkdir -p, but for webdav
func (c *Client) MkdirAll(path string, _ os.FileMode) (err error) {
	path = FixSlashes(path)
	status, err := c.mkcol(path)
	if err != nil {
		return
	}
	if status == 201 {
		return nil
	}
	if status == 409 {
		paths := strings.Split(path, "/")
		sub := "/"
		for _, e := range paths {
			if e == "" {
				continue
			}
			sub += e + "/"
			status, err = c.mkcol(sub)
			if err != nil {
				return
			}
			if status != 201 {
				return NewPathError("MkdirAll", sub, status)
			}
		}
		return nil
	}

	return NewPathError("MkdirAll", path, status)
}

// Rename moves a file from A to B
func (c *Client) Rename(oldpath, newpath string, overwrite bool) error {
	return c.copymove("MOVE", oldpath, newpath, overwrite)
}

// Copy copies a file from A to B
func (c *Client) Copy(oldpath, newpath string, overwrite bool) error {
	return c.copymove("COPY", oldpath, newpath, overwrite)
}

// GetLock gets a lock
func (c *Client) GetLock(path string) (string, error) {
	fi, err := c.Stat(path)
	if err != nil {
		return "", err
	}

	f, ok := fi.(File)
	if !ok {
		// This won't happen unless implementation is changed
		return "", errors.New("Stat did not return a File")
	}

	return f.Lock(), nil
}

// Lock locks a file
func (c *Client) Lock(path string, token string) error {
	return c.lock(path, token, false)
}

// RefreshLock refreshes a lock
func (c *Client) RefreshLock(path string, token string) error {
	return c.lock(path, token, true)
}

// Unlock unlocks a file
func (c *Client) Unlock(path string, token string) error {
	return c.unlock(path, token)
}

// Read reads the contents of a remote file
func (c *Client) Read(path string) ([]byte, error) {
	var stream io.ReadCloser
	var err error

	if stream, err = c.ReadStream(path); err != nil {
		return nil, err
	}
	defer stream.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(stream)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ReadStream reads the stream for a given path
func (c *Client) ReadStream(path string) (io.ReadCloser, error) {
	rs, err := c.req("GET", path, nil, nil)
	if err != nil {
		return nil, NewPathErrorErr("ReadStream", path, err)
	}

	if rs.StatusCode == 200 {
		return rs.Body, nil
	}

	rs.Body.Close()
	return nil, NewPathError("ReadStream", path, rs.StatusCode)
}

// ReadStreamRange reads the stream representing a subset of bytes for a given path,
// utilizing HTTP Range Requests if the server supports it.
// The range is expressed as offset from the start of the file and length, for example
// offset=10, length=10 will return bytes 10 through 19.
//
// If the server does not support partial content requests and returns full content instead,
// this function will emulate the behavior by skipping `offset` bytes and limiting the result
// to `length`.
func (c *Client) ReadStreamRange(path string, offset, length int64) (io.ReadCloser, error) {
	rs, err := c.req("GET", path, nil, func(r *http.Request) {
		if length > 0 {
			r.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", offset, offset+length-1))
		} else {
			r.Header.Add("Range", fmt.Sprintf("bytes=%d-", offset))
		}
	})
	if err != nil {
		return nil, NewPathErrorErr("ReadStreamRange", path, err)
	}

	if rs.StatusCode == http.StatusPartialContent {
		// server supported partial content, return as-is.
		return rs.Body, nil
	}

	// server returned success, but did not support partial content, so we have the whole
	// stream in rs.Body
	if rs.StatusCode == 200 {
		// discard first 'offset' bytes.
		if _, err := io.Copy(io.Discard, io.LimitReader(rs.Body, offset)); err != nil {
			return nil, NewPathErrorErr("ReadStreamRange", path, err)
		}

		// return a io.ReadCloser that is limited to `length` bytes.
		return &limitedReadCloser{rc: rs.Body, remaining: int(length)}, nil
	}

	rs.Body.Close()
	return nil, NewPathError("ReadStream", path, rs.StatusCode)
}

// Write writes data to a given path
func (c *Client) Write(path string, data []byte, _ os.FileMode) (err error) {
	s, err := c.put(path, bytes.NewReader(data), "")
	if err != nil {
		return
	}

	switch s {

	case 200, 201, 204:
		return nil

	case 404, 409:
		err = c.createParentCollection(path)
		if err != nil {
			return
		}

		s, err = c.put(path, bytes.NewReader(data), "")
		if err != nil {
			return
		}
		if s == 200 || s == 201 || s == 204 {
			return
		}
	}

	return NewPathError("Write", path, s)
}

// WriteStream writes a stream
func (c *Client) WriteStream(path string, stream io.Reader, _ os.FileMode, locktoken string) (err error) {

	err = c.createParentCollection(path)
	if err != nil {
		return err
	}

	s, err := c.put(path, stream, locktoken)
	if err != nil {
		return err
	}

	switch s {
	case 200, 201, 204:
		return nil

	default:
		return NewPathError("WriteStream", path, s)
	}
}
