package svc

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/owncloud/ocis-thumbnails/pkg/config"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails/imgsource"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails/storage"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Thumbnails(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Thumbnails{
		config: options.Config,
		mux:    m,
		manager: thumbnails.SimpleManager{
			Storage: storage.NewInMemoryStorage(),
		},
		source: imgsource.WebDav{
			Basepath: "http://localhost:9140/remote.php/webdav/",
		},
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Get("/thumbnails", svc.Thumbnails)
	})

	return svc
}

// Thumbnails defines implements the business logic for Service.
type Thumbnails struct {
	config  *config.Config
	mux     *chi.Mux
	manager thumbnails.Manager
	source  imgsource.Source
}

// ServeHTTP implements the Service interface.
func (g Thumbnails) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Thumbnails provides the endpoint to retrieve a thumbnail for an image
func (g Thumbnails) Thumbnails(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	width, _ := strconv.Atoi(query.Get("width"))
	height, _ := strconv.Atoi(query.Get("height"))
	fileType := query.Get("type")
	filePath := query.Get("file_path")

	encoder := thumbnails.EncoderForType(fileType)
	if encoder == nil {
		// TODO: better error responses
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can't encode that"))
		return
	}
	ctx := thumbnails.ThumbnailContext{
		Width:     width,
		Height:    height,
		ImagePath: filePath,
		Encoder:   encoder,
	}

	thumbnail := g.manager.GetStored(ctx)
	if thumbnail != nil {
		w.Write(thumbnail)
		return
	}

	auth := r.Header.Get("Authorization")

	sCtx := imgsource.NewContext()
	sCtx.Set(imgsource.WebDavAuth, auth)
	// TODO: clean up error handling
	img, err := g.source.Get(ctx.ImagePath, sCtx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if img == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("img is nil"))
		return
	}
	thumbnail, err = g.manager.Get(ctx, img)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(thumbnail)
}
