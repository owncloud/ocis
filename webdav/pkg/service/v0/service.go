package svc

import (
	"encoding/xml"
	merrors "github.com/asim/go-micro/v3/errors"
	"github.com/go-chi/render"
	"net/http"
	"path"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"

	"github.com/go-chi/chi"
	thumbnails "github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	"github.com/owncloud/ocis/webdav/pkg/config"
	"github.com/owncloud/ocis/webdav/pkg/dav/requests"
)

const (
	TokenHeader = "X-Access-Token"
)

var (
	codesEnum = map[int]string{
		http.StatusBadRequest:       "Sabre\\DAV\\Exception\\BadRequest",
		http.StatusUnauthorized:     "Sabre\\DAV\\Exception\\NotAuthenticated",
		http.StatusNotFound:         "Sabre\\DAV\\Exception\\NotFound",
		http.StatusMethodNotAllowed: "Sabre\\DAV\\Exception\\MethodNotAllowed",
	}
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Thumbnail(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Webdav{
		config: options.Config,
		log:    options.Logger,
		mux:    m,
		thumbnailsClient: thumbnails.NewThumbnailService("com.owncloud.api.thumbnails", grpc.DefaultClient),
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Get("/remote.php/dav/files/{user}/*", svc.Thumbnail)
		r.Get("/remote.php/dav/public-files/{token}/*", svc.PublicThumbnail)
		r.Head("/remote.php/dav/public-files/{token}/*", svc.PublicThumbnailHead)
	})

	return svc
}

// Webdav defines implements the business logic for Service.
type Webdav struct {
	config *config.Config
	log    log.Logger
	mux    *chi.Mux
	thumbnailsClient thumbnails.ThumbnailService
}

// ServeHTTP implements the Service interface.
func (g Webdav) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Thumbnail implements the Service interface.
func (g Webdav) Thumbnail(w http.ResponseWriter, r *http.Request) {
	tr, err := requests.ParseThumbnailRequest(r)
	if err != nil {
		g.log.Error().Err(err).Msg("could not create Request")
		renderError(w, r, errBadRequest(err.Error()))
		return
	}

	t := r.Header.Get(TokenHeader)
	rsp, err := g.thumbnailsClient.GetThumbnail(r.Context(), &thumbnails.GetThumbnailRequest{
		Filepath:      strings.TrimLeft(tr.Filepath, "/"),
		ThumbnailType: extensionToThumbnailType(strings.TrimLeft(tr.Extension, ".")),
		Width:         tr.Width,
		Height:        tr.Height,
		Source: &thumbnails.GetThumbnailRequest_Cs3Source{
			Cs3Source: &thumbnails.CS3Source{
				Path:          path.Join("/home", tr.Filepath),
				Authorization: t,
			},
		},
	})
	if err != nil {
		g.log.Error().Err(err).Msg("could not get thumbnail")
		e := merrors.Parse(err.Error())
		switch e.Code {
		case http.StatusNotFound:
			renderError(w, r, errNotFound(notFoundMsg(tr.Filename)))
		case http.StatusBadRequest:
			renderError(w, r, errBadRequest(err.Error()))
		default:
			renderError(w, r, errInternalError(err.Error()))
		}
		return
	}

	if len(rsp.Thumbnail) == 0 {
		renderError(w, r, errNotFound(""))
		return
	}

	g.mustRender(w, r, newThumbnailResponse(rsp))
}

func (g Webdav) PublicThumbnail(w http.ResponseWriter, r *http.Request) {
	tr, err := requests.ParseThumbnailRequest(r)
	if err != nil {
		g.log.Error().Err(err).Msg("could not create Request")
		renderError(w, r, errBadRequest(err.Error()))
		return
	}

	rsp, err := g.thumbnailsClient.GetThumbnail(r.Context(), &thumbnails.GetThumbnailRequest{
		Filepath:      strings.TrimLeft(tr.Filepath, "/"),
		ThumbnailType: extensionToThumbnailType(strings.TrimLeft(tr.Extension, ".")),
		Width:         tr.Width,
		Height:        tr.Height,
		Source: &thumbnails.GetThumbnailRequest_WebdavSource{
			WebdavSource: &thumbnails.WebdavSource{
				Url:             g.config.OcisPublicURL + r.URL.RequestURI(),
				IsPublicLink:    true,
				PublicLinkToken: tr.PublicLinkToken,
			},
		},
	})
	if err != nil {
		g.log.Error().Err(err).Msg("could not get thumbnail")
		e := merrors.Parse(err.Error())
		switch e.Code {
		case http.StatusNotFound:
			renderError(w, r, errNotFound(notFoundMsg(tr.Filename)))
		case http.StatusBadRequest:
			renderError(w, r, errBadRequest(err.Error()))
		default:
			renderError(w, r, errInternalError(err.Error()))
		}
		return
	}

	if len(rsp.Thumbnail) == 0 {
		renderError(w, r, errNotFound(""))
		return
	}

	g.mustRender(w, r, newThumbnailResponse(rsp))
}

func (g Webdav) PublicThumbnailHead(w http.ResponseWriter, r *http.Request) {
	tr, err := requests.ParseThumbnailRequest(r)
	if err != nil {
		g.log.Error().Err(err).Msg("could not create Request")
		renderError(w, r, errBadRequest(err.Error()))
		return
	}

	rsp, err := g.thumbnailsClient.GetThumbnail(r.Context(), &thumbnails.GetThumbnailRequest{
		Filepath:      strings.TrimLeft(tr.Filepath, "/"),
		ThumbnailType: extensionToThumbnailType(strings.TrimLeft(tr.Extension, ".")),
		Width:         tr.Width,
		Height:        tr.Height,
		Source: &thumbnails.GetThumbnailRequest_WebdavSource{
			WebdavSource: &thumbnails.WebdavSource{
				Url:             g.config.OcisPublicURL + r.URL.RequestURI(),
				IsPublicLink:    true,
				PublicLinkToken: tr.PublicLinkToken,
			},
		},
	})
	if err != nil {
		e := merrors.Parse(err.Error())
		switch e.Code {
		case http.StatusNotFound:
			renderError(w, r, errNotFound(notFoundMsg(tr.Filename)))
		case http.StatusBadRequest:
			g.log.Error().Err(err).Msg("could not get thumbnail")
			renderError(w, r, errBadRequest(err.Error()))
		default:
			g.log.Error().Err(err).Msg("could not get thumbnail")
			renderError(w, r, errInternalError(err.Error()))
		}
		return
	}

	if len(rsp.Thumbnail) == 0 {
		renderError(w, r, errNotFound(""))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func extensionToThumbnailType(ext string) thumbnails.GetThumbnailRequest_ThumbnailType {
	switch strings.ToUpper(ext) {
	case "GIF", "PNG":
		return thumbnails.GetThumbnailRequest_PNG
	default:
		return thumbnails.GetThumbnailRequest_JPG
	}
}

func (g Webdav) mustRender(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
	if err := render.Render(w, r, renderer); err != nil {
		g.log.Err(err).Msg("failed to write response")
	}
}

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_error
type errResponse struct {
	HTTPStatusCode int      `json:"-" xml:"-"`
	XMLName        xml.Name `xml:"d:error"`
	Xmlnsd         string   `xml:"xmlns:d,attr"`
	Xmlnss         string   `xml:"xmlns:s,attr"`
	Exception      string   `xml:"s:exception"`
	Message        string   `xml:"s:message"`
	InnerXML       []byte   `xml:",innerxml"`
}

func newErrResponse(statusCode int, msg string) *errResponse {
	rsp := &errResponse{
		HTTPStatusCode: statusCode,
		Xmlnsd:         "DAV",
		Xmlnss:         "http://sabredav.org/ns",
		Exception:      codesEnum[statusCode],
	}
	if msg != "" {
		rsp.Message = msg
	}
	return rsp
}

func errInternalError(msg string) *errResponse {
	return newErrResponse(http.StatusInternalServerError, msg)
}

func errBadRequest(msg string) *errResponse {
	return newErrResponse(http.StatusBadRequest, msg)
}

func errNotFound(msg string) *errResponse {
	return newErrResponse(http.StatusNotFound, msg)
}

type thumbnailResponse struct {
	contentType string
	thumbnail   []byte
}

func (t *thumbnailResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", t.contentType)
	_, err := w.Write(t.thumbnail)
	return err
}

func newThumbnailResponse(rsp *thumbnails.GetThumbnailResponse) *thumbnailResponse {
	return &thumbnailResponse{
		contentType: rsp.Mimetype,
		thumbnail:   rsp.Thumbnail,
	}
}

func renderError(w http.ResponseWriter, r *http.Request, err *errResponse) {
	render.Status(r, err.HTTPStatusCode)
	render.XML(w, r, err)
}

func notFoundMsg(name string) string {
	return "File with name " + name + " could not be located"
}
