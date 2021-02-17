package svc

import (
	"net/http"
	"strings"

	"github.com/asim/go-micro/v3/client"
	"github.com/go-chi/chi"
	thumbnails "github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	"github.com/owncloud/ocis/webdav/pkg/config"
	thumbnail "github.com/owncloud/ocis/webdav/pkg/dav/thumbnails"
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
		mux:    m,
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Get("/remote.php/dav/files/{user}/*", svc.Thumbnail)
	})

	return svc
}

// Webdav defines implements the business logic for Service.
type Webdav struct {
	config *config.Config
	mux    *chi.Mux
}

// ServeHTTP implements the Service interface.
func (g Webdav) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Thumbnail implements the Service interface.
func (g Webdav) Thumbnail(w http.ResponseWriter, r *http.Request) {
	tr, err := thumbnail.NewRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	c := thumbnails.NewThumbnailService("com.owncloud.api.thumbnails", client.DefaultClient)
	rsp, err := c.GetThumbnail(r.Context(), &thumbnails.GetRequest{
		Filepath:      strings.TrimLeft(tr.Filepath, "/"),
		Filetype:      extensionToFiletype(tr.Filetype),
		Etag:          tr.Etag,
		Width:         int32(tr.Width),
		Height:        int32(tr.Height),
		Authorization: tr.Authorization,
		Username:      tr.Username,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if len(rsp.Thumbnail) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", rsp.GetMimetype())
	w.WriteHeader(http.StatusOK)
	w.Write(rsp.Thumbnail)
}

func extensionToFiletype(ext string) thumbnails.GetRequest_FileType {
	val, ok := thumbnails.GetRequest_FileType_value[strings.ToUpper(ext)]
	if !ok {
		return thumbnails.GetRequest_FileType(-1)
	}
	return thumbnails.GetRequest_FileType(val)
}
