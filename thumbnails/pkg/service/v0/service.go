package svc

import (
	"context"
	merrors "github.com/asim/go-micro/v3/errors"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/token"
	"github.com/owncloud/ocis/ocis-pkg/log"
	v0proto "github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/imgsource"
	"google.golang.org/grpc/metadata"
	"image"
)

// NewService returns a service implementation for Service.
func NewService(opts ...Option) v0proto.ThumbnailServiceHandler {
	options := newOptions(opts...)
	logger := options.Logger
	resolutions, err := thumbnail.ParseResolutions(options.Config.Thumbnail.Resolutions)
	if err != nil {
		logger.Fatal().Err(err).Msg("resolutions not configured correctly")
	}
	svc := Thumbnail{
		serviceID: options.Config.Server.Namespace + "." + options.Config.Server.Name,
		manager: thumbnail.NewSimpleManager(
			resolutions,
			options.ThumbnailStorage,
			logger,
		),
		webdavSource: options.ImageSource,
		cs3Source:    options.CS3Source,
		logger:       logger,
		cs3Client:    options.CS3Client,
	}

	return svc
}

// Thumbnail implements the GRPC handler.
type Thumbnail struct {
	serviceID    string
	manager      thumbnail.Manager
	webdavSource imgsource.Source
	cs3Source    imgsource.Source
	logger       log.Logger
	cs3Client    gateway.GatewayAPIClient
}

// GetThumbnail retrieves a thumbnail for an image
func (g Thumbnail) GetThumbnail(ctx context.Context, req *v0proto.GetThumbnailRequest, rsp *v0proto.GetThumbnailResponse) error {
	_, ok := v0proto.GetThumbnailRequest_FileType_value[req.ThumbnailType.String()]
	if !ok {
		g.logger.Debug().Str("filetype", req.ThumbnailType.String()).Msg("unsupported filetype")
		return nil
	}
	encoder := thumbnail.EncoderForType(req.ThumbnailType.String())
	if encoder == nil {
		g.logger.Debug().Str("filetype", req.ThumbnailType.String()).Msg("unsupported filetype")
		return nil
	}
	sReq := &provider.StatRequest{
		Ref: &provider.Reference{
			Spec: &provider.Reference_Path{Path: "/home/" + req.Filepath},
		},
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return merrors.Unauthorized(g.serviceID, "authorization is missing")
	}
	auth, ok := md[token.TokenHeader]
	if !ok {
		return merrors.Unauthorized(g.serviceID, "authorization is missing")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, auth[0])
	sRes, err := g.cs3Client.Stat(ctx, sReq)
	if err != nil {
		g.logger.Error().Err(err).Msg("could stat file")
		return merrors.InternalServerError(g.serviceID, "could not stat file: %s", err.Error())
	}

	if sRes.Status.Code != rpc.Code_CODE_OK {
		g.logger.Error().Err(err).Msg("could not create Request")
		return merrors.InternalServerError(g.serviceID, "could not stat file: %s", err.Error())
	}

	tr := thumbnail.Request{
		Resolution: image.Rect(0, 0, int(req.Width), int(req.Height)),
		Encoder:    encoder,
		ETag:       sRes.GetInfo().GetChecksum().GetSum(),
	}

	thumb, ok := g.manager.Get(tr)
	if ok {
		rsp.Thumbnail = thumb
		rsp.Mimetype = tr.Encoder.MimeType()
		return nil
	}

	var img image.Image
	switch {
	case req.GetWebdavSource() != nil:
		src := req.GetWebdavSource()
		src.GetAuthorization()

		sCtx := imgsource.ContextSetAuthorization(ctx, src.GetAuthorization())
		img, err = g.webdavSource.Get(sCtx, src.GetUrl())
	case req.GetCs3Source() != nil:
		src := req.GetCs3Source()

		sCtx := imgsource.ContextSetAuthorization(ctx, auth[0])
		img, err = g.cs3Source.Get(sCtx, src.Path)
	default:
		g.logger.Error().Msg("no image source provided")
		return merrors.BadRequest(g.serviceID, "image source is missing")
	}
	if err != nil {
		return merrors.InternalServerError(g.serviceID, "could not get image from source: %v", err.Error())
	}
	if img == nil {
		return merrors.InternalServerError(g.serviceID, "could not get image from source")
	}
	if thumb, err = g.manager.Generate(tr, img); err != nil {
		return err
	}

	rsp.Thumbnail = thumb
	rsp.Mimetype = tr.Encoder.MimeType()
	return nil
}
