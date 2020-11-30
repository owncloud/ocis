package svc

import (
	"context"
	"image"

	"gopkg.in/square/go-jose.v2/jwt"

	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis/ocis-pkg/log"
	v0proto "github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/imgsource"
	"github.com/pkg/errors"
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
		source: options.ImageSource,
		logger: logger,
	}

	return svc
}

// Thumbnail implements the GRPC handler.
type Thumbnail struct {
	serviceID string
	manager   thumbnail.Manager
	source    imgsource.Source
	logger    log.Logger
}

// GetThumbnail retrieves a thumbnail for an image
func (g Thumbnail) GetThumbnail(ctx context.Context, req *v0proto.GetRequest, rsp *v0proto.GetResponse) error {
	encoder := thumbnail.EncoderForType(req.Filetype.String())
	if encoder == nil {
		g.logger.Debug().Str("filetype", req.Filetype.String()).Msg("unsupported filetype")
		return nil
	}

	auth := req.Authorization
	if auth == "" {
		return merrors.BadRequest(g.serviceID, "authorization is missing")
	}
	username, err := usernameFromAuthorization(auth)
	if err != nil {
		return merrors.InternalServerError(g.serviceID, "could not get username: %v", err.Error())
	}

	tr := thumbnail.Request{
		Resolution: image.Rect(0, 0, int(req.Width), int(req.Height)),
		Encoder:    encoder,
		ETag:       req.Etag,
		Username:   username,
	}

	thumbnail := g.manager.GetStored(tr)
	if thumbnail != nil {
		rsp.Thumbnail = thumbnail
		rsp.Mimetype = tr.Encoder.MimeType()
		return nil
	}

	sCtx := imgsource.ContextSetAuthorization(ctx, auth)
	img, err := g.source.Get(sCtx, req.Filepath)
	if err != nil {
		return merrors.InternalServerError(g.serviceID, "could not get image from source: %v", err.Error())
	}
	if img == nil {
		return merrors.InternalServerError(g.serviceID, "could not get image from source")
	}
	thumbnail, err = g.manager.Get(tr, img)
	if err != nil {
		return err
	}

	rsp.Thumbnail = thumbnail
	rsp.Mimetype = tr.Encoder.MimeType()
	return nil
}

func usernameFromAuthorization(auth string) (string, error) {
	tokenString := auth[len("Bearer "):] // strip the bearer prefix

	var claims map[string]interface{}
	token, err := jwt.ParseSigned(tokenString)
	if err != nil {
		return "", errors.Wrap(err, "could not parse auth token")
	}
	err = token.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		return "", errors.Wrap(err, "could not get claims from auth token")
	}

	identityMap := claims["kc.identity"].(map[string]interface{})
	username := identityMap["kc.i.un"].(string)

	return username, nil
}
