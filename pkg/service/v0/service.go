package svc

import (
	"context"
	"fmt"

	"github.com/owncloud/ocis-pkg/v2/log"
	v0proto "github.com/owncloud/ocis-thumbnails/pkg/proto/v0"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails/imgsource"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails/resolutions"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails/storage"
)

// NewService returns a service implementation for Service.
func NewService(opts ...Option) v0proto.ThumbnailServiceHandler {
	options := newOptions(opts...)
	logger := options.Logger
	resolutions, err := resolutions.New(options.Config.Thumbnail.Resolutions)
	if err != nil {
		logger.Fatal().Err(err).Msg("resolutions not configured correctly")
	}
	svc := Thumbnail{
		manager: thumbnails.NewSimpleManager(
			storage.NewFileSystemStorage(
				options.Config.Thumbnail.FileSystemStorage,
				logger,
			),
			logger,
		),
		resolutions: resolutions,
		source:      imgsource.NewWebDavSource(options.Config.Thumbnail.WebDavSource),
		logger:      logger,
	}

	return svc
}

// Thumbnail implements the GRPC handler.
type Thumbnail struct {
	manager     thumbnails.Manager
	resolutions resolutions.Resolutions
	source      imgsource.Source
	logger      log.Logger
}

// GetThumbnail retrieves a thumbnail for an image
func (g Thumbnail) GetThumbnail(ctx context.Context, req *v0proto.GetRequest, rsp *v0proto.GetResponse) error {
	encoder := thumbnails.EncoderForType(req.Filetype.String())
	if encoder == nil {
		// TODO: better error responses
		return fmt.Errorf("can't be encoded. filetype %s not supported", req.Filetype.String())
	}
	r := g.resolutions.ClosestMatch(int(req.Width), int(req.Height))
	tCtx := thumbnails.Context{
		Resolution: r,
		ImagePath:  req.Filepath,
		Encoder:    encoder,
		ETag:       req.Etag,
	}

	thumbnail := g.manager.GetStored(tCtx)
	if thumbnail != nil {
		rsp.Thumbnail = thumbnail
		rsp.Mimetype = tCtx.Encoder.MimeType()
		return nil
	}

	auth := req.Authorization
	sCtx := context.WithValue(ctx, imgsource.WebDavAuth, auth)
	// TODO: clean up error handling
	img, err := g.source.Get(sCtx, tCtx.ImagePath)
	if err != nil {
		return err
	}
	if img == nil {
		return fmt.Errorf("could not retrieve image")
	}
	thumbnail, err = g.manager.Get(tCtx, img)
	if err != nil {
		return err
	}

	rsp.Thumbnail = thumbnail
	rsp.Mimetype = tCtx.Encoder.MimeType()
	return nil
}
