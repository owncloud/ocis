package svc

import (
	"io"
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
		w.WriteHeader(http.StatusBadRequest)
		mustWrite(g.log, w, []byte(err.Error()))
		return
	}

	c := thumbnails.NewThumbnailService("com.owncloud.api.thumbnails", grpc.DefaultClient)
	t := r.Header.Get("X-Access-Token")
	rsp, err := c.GetThumbnail(r.Context(), &thumbnails.GetThumbnailRequest{
		Filepath:      strings.TrimLeft(tr.Filepath, "/"),
		ThumbnailType: extensionToFiletype(strings.TrimLeft(tr.Extension, ".")),
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
		w.WriteHeader(http.StatusBadRequest)
		mustWrite(g.log, w, []byte(err.Error()))
		return
	}

	if len(rsp.Thumbnail) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", rsp.GetMimetype())
	w.WriteHeader(http.StatusOK)
	mustWrite(g.log, w, rsp.Thumbnail)
}

func (g Webdav) PublicThumbnail(w http.ResponseWriter, r *http.Request) {
	tr, err := requests.ParseThumbnailRequest(r)
	if err != nil {
		g.log.Error().Err(err).Msg("could not create Request")
		w.WriteHeader(http.StatusBadRequest)
		mustWrite(g.log, w, []byte(err.Error()))
		return
	}

	c := thumbnails.NewThumbnailService("com.owncloud.api.thumbnails", grpc.DefaultClient)
	rsp, err := c.GetThumbnail(r.Context(), &thumbnails.GetThumbnailRequest{
		Filepath:      strings.TrimLeft(tr.Filepath, "/"),
		ThumbnailType: extensionToFiletype(strings.TrimLeft(tr.Extension, ".")),
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
		w.WriteHeader(http.StatusBadRequest)
		mustWrite(g.log, w, []byte(err.Error()))
		return
	}

	if len(rsp.Thumbnail) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", rsp.GetMimetype())
	w.WriteHeader(http.StatusOK)
	mustWrite(g.log, w, rsp.Thumbnail)
}

func (g Webdav) PublicThumbnailHead(w http.ResponseWriter, r *http.Request) {
	tr, err := requests.ParseThumbnailRequest(r)
	if err != nil {
		g.log.Error().Err(err).Msg("could not create Request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := thumbnails.NewThumbnailService("com.owncloud.api.thumbnails", grpc.DefaultClient)
	rsp, err := c.GetThumbnail(r.Context(), &thumbnails.GetThumbnailRequest{
		Filepath:      strings.TrimLeft(tr.Filepath, "/"),
		ThumbnailType: extensionToFiletype(strings.TrimLeft(tr.Extension, ".")),
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(rsp.Thumbnail) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", rsp.GetMimetype())
	w.WriteHeader(http.StatusOK)
}

func extensionToFiletype(ext string) thumbnails.GetThumbnailRequest_FileType {
	switch strings.ToUpper(ext) {
	case "GIF", "PNG":
		return thumbnails.GetThumbnailRequest_PNG
	case "JPEG", "JPG":
		return thumbnails.GetThumbnailRequest_JPG
	default:
		return thumbnails.GetThumbnailRequest_FileType(-1)
	}
}

func mustWrite(logger log.Logger, w io.Writer, val []byte) {
	if _, err := w.Write(val); err != nil {
		logger.Error().Err(err).Msg("could not write response")
	}
}
