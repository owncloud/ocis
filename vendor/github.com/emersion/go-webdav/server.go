package webdav

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/emersion/go-webdav/internal"
)

// FileSystem is a WebDAV server backend.
type FileSystem interface {
	Open(name string) (io.ReadCloser, error)
	Stat(name string) (*FileInfo, error)
	Readdir(name string, recursive bool) ([]FileInfo, error)
	Create(name string) (io.WriteCloser, error)
	RemoveAll(name string) error
	Mkdir(name string) error
	Copy(name, dest string, recursive, overwrite bool) (created bool, err error)
	MoveAll(name, dest string, overwrite bool) (created bool, err error)
}

// Handler handles WebDAV HTTP requests. It can be used to create a WebDAV
// server.
type Handler struct {
	FileSystem FileSystem
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.FileSystem == nil {
		http.Error(w, "webdav: no filesystem available", http.StatusInternalServerError)
		return
	}

	b := backend{h.FileSystem}
	hh := internal.Handler{&b}
	hh.ServeHTTP(w, r)
}

// NewHTTPError creates a new error that is associated with an HTTP status code
// and optionally an error that lead to it. Backends can use this functions to
// return errors that convey some semantics (e.g. 404 not found, 403 access
// denied, etc) while also providing an (optional) arbitrary error context
// (intended for humans).
func NewHTTPError(statusCode int, cause error) error {
	return &internal.HTTPError{Code: statusCode, Err: cause}
}

type backend struct {
	FileSystem FileSystem
}

func (b *backend) Options(r *http.Request) (caps []string, allow []string, err error) {
	fi, err := b.FileSystem.Stat(r.URL.Path)
	if os.IsNotExist(err) {
		return nil, []string{http.MethodOptions, http.MethodPut, "MKCOL"}, nil
	} else if err != nil {
		return nil, nil, err
	}

	allow = []string{
		http.MethodOptions,
		http.MethodDelete,
		"PROPFIND",
		"COPY",
		"MOVE",
	}

	if !fi.IsDir {
		allow = append(allow, http.MethodHead, http.MethodGet, http.MethodPut)
	}

	return nil, allow, nil
}

func (b *backend) HeadGet(w http.ResponseWriter, r *http.Request) error {
	fi, err := b.FileSystem.Stat(r.URL.Path)
	if os.IsNotExist(err) {
		return &internal.HTTPError{Code: http.StatusNotFound, Err: err}
	} else if err != nil {
		return err
	}
	if fi.IsDir {
		return &internal.HTTPError{Code: http.StatusMethodNotAllowed}
	}

	f, err := b.FileSystem.Open(r.URL.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size, 10))
	if fi.MIMEType != "" {
		w.Header().Set("Content-Type", fi.MIMEType)
	}
	if !fi.ModTime.IsZero() {
		w.Header().Set("Last-Modified", fi.ModTime.UTC().Format(http.TimeFormat))
	}
	if fi.ETag != "" {
		w.Header().Set("ETag", internal.ETag(fi.ETag).String())
	}

	if rs, ok := f.(io.ReadSeeker); ok {
		// If it's an io.Seeker, use http.ServeContent which supports ranges
		http.ServeContent(w, r, r.URL.Path, fi.ModTime, rs)
	} else {
		if r.Method != http.MethodHead {
			io.Copy(w, f)
		}
	}
	return nil
}

func (b *backend) PropFind(r *http.Request, propfind *internal.PropFind, depth internal.Depth) (*internal.MultiStatus, error) {
	// TODO: use partial error Response on error

	fi, err := b.FileSystem.Stat(r.URL.Path)
	if os.IsNotExist(err) {
		return nil, &internal.HTTPError{Code: http.StatusNotFound, Err: err}
	} else if err != nil {
		return nil, err
	}

	var resps []internal.Response
	if depth != internal.DepthZero && fi.IsDir {
		children, err := b.FileSystem.Readdir(r.URL.Path, depth == internal.DepthInfinity)
		if err != nil {
			return nil, err
		}

		resps = make([]internal.Response, len(children))
		for i, child := range children {
			resp, err := b.propFindFile(propfind, &child)
			if err != nil {
				return nil, err
			}
			resps[i] = *resp
		}
	} else {
		resp, err := b.propFindFile(propfind, fi)
		if err != nil {
			return nil, err
		}

		resps = []internal.Response{*resp}
	}

	return internal.NewMultiStatus(resps...), nil
}

func (b *backend) propFindFile(propfind *internal.PropFind, fi *FileInfo) (*internal.Response, error) {
	props := make(map[xml.Name]internal.PropFindFunc)

	props[internal.ResourceTypeName] = func(*internal.RawXMLValue) (interface{}, error) {
		var types []xml.Name
		if fi.IsDir {
			types = append(types, internal.CollectionName)
		}
		return internal.NewResourceType(types...), nil
	}

	if !fi.IsDir {
		props[internal.GetContentLengthName] = func(*internal.RawXMLValue) (interface{}, error) {
			return &internal.GetContentLength{Length: fi.Size}, nil
		}

		if !fi.ModTime.IsZero() {
			props[internal.GetLastModifiedName] = func(*internal.RawXMLValue) (interface{}, error) {
				return &internal.GetLastModified{LastModified: internal.Time(fi.ModTime)}, nil
			}
		}

		if fi.MIMEType != "" {
			props[internal.GetContentTypeName] = func(*internal.RawXMLValue) (interface{}, error) {
				return &internal.GetContentType{Type: fi.MIMEType}, nil
			}
		}

		if fi.ETag != "" {
			props[internal.GetETagName] = func(*internal.RawXMLValue) (interface{}, error) {
				return &internal.GetETag{ETag: internal.ETag(fi.ETag)}, nil
			}
		}
	}

	return internal.NewPropFindResponse(fi.Path, propfind, props)
}

func (b *backend) PropPatch(r *http.Request, update *internal.PropertyUpdate) (*internal.Response, error) {
	// TODO: return a failed Response instead
	return nil, internal.HTTPErrorf(http.StatusForbidden, "webdav: PROPPATCH is unsupported")
}

func (b *backend) Put(r *http.Request) (*internal.Href, error) {
	wc, err := b.FileSystem.Create(r.URL.Path)
	if err != nil {
		return nil, err
	}
	defer wc.Close()

	if _, err := io.Copy(wc, r.Body); err != nil {
		return nil, err
	}

	return nil, wc.Close()
}

func (b *backend) Delete(r *http.Request) error {
	err := b.FileSystem.RemoveAll(r.URL.Path)
	if os.IsNotExist(err) {
		return &internal.HTTPError{Code: http.StatusNotFound, Err: err}
	}
	return err
}

func (b *backend) Mkcol(r *http.Request) error {
	if r.Header.Get("Content-Type") != "" {
		return internal.HTTPErrorf(http.StatusUnsupportedMediaType, "webdav: request body not supported in MKCOL request")
	}
	err := b.FileSystem.Mkdir(r.URL.Path)
	if os.IsNotExist(err) {
		return &internal.HTTPError{Code: http.StatusConflict, Err: err}
	}
	return err
}

func (b *backend) Copy(r *http.Request, dest *internal.Href, recursive, overwrite bool) (created bool, err error) {
	created, err = b.FileSystem.Copy(r.URL.Path, dest.Path, recursive, overwrite)
	if os.IsExist(err) {
		return false, &internal.HTTPError{http.StatusPreconditionFailed, err}
	}
	return created, err
}

func (b *backend) Move(r *http.Request, dest *internal.Href, overwrite bool) (created bool, err error) {
	created, err = b.FileSystem.MoveAll(r.URL.Path, dest.Path, overwrite)
	if os.IsExist(err) {
		return false, &internal.HTTPError{http.StatusPreconditionFailed, err}
	}
	return created, err
}

// BackendSuppliedHomeSet represents either a CalDAV calendar-home-set or a
// CardDAV addressbook-home-set. It should only be created via
// `caldav.NewCalendarHomeSet()` or `carddav.NewAddressbookHomeSet()`. Only to
// be used server-side, for listing a user's home sets as determined by the
// (external) backend.
type BackendSuppliedHomeSet interface {
	GetXMLName() xml.Name
}

// UserPrincipalBackend can determine the current user's principal URL for a
// given request context.
type UserPrincipalBackend interface {
	CurrentUserPrincipal(ctx context.Context) (string, error)
}

type Capability string

type ServePrincipalOptions struct {
	CurrentUserPrincipalPath string
	HomeSets                 []BackendSuppliedHomeSet
	Capabilities             []Capability
}

// ServePrincipal replies to requests for a principal URL.
func ServePrincipal(w http.ResponseWriter, r *http.Request, options *ServePrincipalOptions) {
	switch r.Method {
	case http.MethodOptions:
		caps := []string{"1", "3"}
		for _, c := range options.Capabilities {
			caps = append(caps, string(c))
		}
		allow := []string{http.MethodOptions, "PROPFIND", "REPORT", "DELETE", "MKCOL"}
		w.Header().Add("DAV", strings.Join(caps, ", "))
		w.Header().Add("Allow", strings.Join(allow, ", "))
		w.WriteHeader(http.StatusNoContent)
	case "PROPFIND":
		if err := servePrincipalPropfind(w, r, options); err != nil {
			internal.ServeError(w, err)
		}
	default:
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
	}
}

func servePrincipalPropfind(w http.ResponseWriter, r *http.Request, options *ServePrincipalOptions) error {
	var propfind internal.PropFind
	if err := internal.DecodeXMLRequest(r, &propfind); err != nil {
		return err
	}
	props := map[xml.Name]internal.PropFindFunc{
		internal.ResourceTypeName: func(*internal.RawXMLValue) (interface{}, error) {
			return internal.NewResourceType(principalName), nil
		},
		internal.CurrentUserPrincipalName: func(*internal.RawXMLValue) (interface{}, error) {
			return &internal.CurrentUserPrincipal{Href: internal.Href{Path: options.CurrentUserPrincipalPath}}, nil
		},
	}

	// TODO: handle Depth and more properties

	for _, homeSet := range options.HomeSets {
		hs := homeSet // capture variable for closure
		props[homeSet.GetXMLName()] = func(*internal.RawXMLValue) (interface{}, error) {
			return hs, nil
		}
	}

	resp, err := internal.NewPropFindResponse(r.URL.Path, &propfind, props)
	if err != nil {
		return err
	}

	ms := internal.NewMultiStatus(*resp)
	return internal.ServeMultiStatus(w, ms)
}
