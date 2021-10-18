package svc

import (
	"context"
	"image"
	"net/url"
	"path"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/thumbnails/pkg/preprocessor"
	v0proto "github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/imgsource"
	"github.com/pkg/errors"
	merrors "go-micro.dev/v4/errors"
	"google.golang.org/grpc/metadata"
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
		serviceID:       options.Config.Server.Namespace + "." + options.Config.Server.Name,
		webdavNamespace: options.Config.Thumbnail.WebdavNamespace,
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
	serviceID       string
	webdavNamespace string
	manager         thumbnail.Manager
	webdavSource    imgsource.Source
	cs3Source       imgsource.Source
	logger          log.Logger
	cs3Client       gateway.GatewayAPIClient
}

// GetThumbnail retrieves a thumbnail for an image
func (g Thumbnail) GetThumbnail(ctx context.Context, req *v0proto.GetThumbnailRequest, rsp *v0proto.GetThumbnailResponse) error {
	_, ok := v0proto.GetThumbnailRequest_ThumbnailType_value[req.ThumbnailType.String()]
	if !ok {
		g.logger.Debug().Str("thumbnail_type", req.ThumbnailType.String()).Msg("unsupported thumbnail type")
		return nil
	}
	encoder := thumbnail.EncoderForType(req.ThumbnailType.String())
	if encoder == nil {
		g.logger.Debug().Str("thumbnail_type", req.ThumbnailType.String()).Msg("unsupported thumbnail type")
		return nil
	}

	var thumb []byte
	var err error
	switch {
	case req.GetWebdavSource() != nil:
		thumb, err = g.handleWebdavSource(ctx, req, encoder)
	case req.GetCs3Source() != nil:
		thumb, err = g.handleCS3Source(ctx, req, encoder)
	default:
		g.logger.Error().Msg("no image source provided")
		return merrors.BadRequest(g.serviceID, "image source is missing")
	}
	if err != nil {
		return err
	}

	rsp.Thumbnail = thumb
	rsp.Mimetype = encoder.MimeType()
	return nil
}

func (g Thumbnail) handleCS3Source(ctx context.Context, req *v0proto.GetThumbnailRequest, encoder thumbnail.Encoder) ([]byte, error) {
	src := req.GetCs3Source()
	sRes, err := g.stat(src.Path, src.Authorization)
	if err != nil {
		return nil, err
	}

	tr := thumbnail.Request{
		Resolution: image.Rect(0, 0, int(req.Width), int(req.Height)),
		Encoder:    encoder,
		Checksum:   sRes.GetInfo().GetChecksum().GetSum(),
	}

	thumb, ok := g.manager.Get(tr)
	if ok {
		return thumb, nil
	}

	ctx = imgsource.ContextSetAuthorization(ctx, src.Authorization)
	r, err := g.cs3Source.Get(ctx, src.Path)
	if err != nil {
		return nil, merrors.InternalServerError(g.serviceID, "could not get image from source: %s", err.Error())
	}
	defer r.Close() // nolint:errcheck
	pp := preprocessor.ForType(sRes.GetInfo().GetMimeType())
	img, err := pp.Convert(r)
	if img == nil || err != nil {
		return nil, merrors.InternalServerError(g.serviceID, "could not get image")
	}
	if thumb, err = g.manager.Generate(tr, img); err != nil {
		return nil, err
	}

	return thumb, nil
}

func (g Thumbnail) handleWebdavSource(ctx context.Context, req *v0proto.GetThumbnailRequest, encoder thumbnail.Encoder) ([]byte, error) {
	src := req.GetWebdavSource()
	imgURL, err := url.Parse(src.Url)
	if err != nil {
		return nil, errors.Wrap(err, "source url is invalid")
	}

	var auth, statPath string
	if src.IsPublicLink {
		q := imgURL.Query()
		var rsp *gateway.AuthenticateResponse
		if q.Get("signature") != "" && q.Get("expiration") != "" {
			// Handle pre-signed public links
			sig := q.Get("signature")
			exp := q.Get("expiration")
			rsp, err = g.cs3Client.Authenticate(ctx, &gateway.AuthenticateRequest{
				Type:         "publicshares",
				ClientId:     src.PublicLinkToken,
				ClientSecret: strings.Join([]string{"signature", sig, exp}, "|"),
			})
		} else {
			rsp, err = g.cs3Client.Authenticate(ctx, &gateway.AuthenticateRequest{
				Type:     "publicshares",
				ClientId: src.PublicLinkToken,
				// We pass an empty password because we expect non pre-signed public links
				// to not be password protected
				ClientSecret: "password|",
			})
		}

		if err != nil {
			return nil, merrors.InternalServerError(g.serviceID, "could not authenticate: %s", err.Error())
		}
		auth = rsp.Token
		statPath = path.Join("/public", src.PublicLinkToken, req.Filepath)
	} else {
		auth = src.RevaAuthorization
		statPath = path.Join(g.webdavNamespace, req.Filepath)
	}
	sRes, err := g.stat(statPath, auth)
	if err != nil {
		return nil, err
	}
	tr := thumbnail.Request{
		Resolution: image.Rect(0, 0, int(req.Width), int(req.Height)),
		Encoder:    encoder,
		Checksum:   sRes.GetInfo().GetChecksum().GetSum(),
	}
	thumb, ok := g.manager.Get(tr)
	if ok {
		return thumb, nil
	}

	if src.WebdavAuthorization != "" {
		ctx = imgsource.ContextSetAuthorization(ctx, src.WebdavAuthorization)
	}
	imgURL.RawQuery = ""
	r, err := g.webdavSource.Get(ctx, imgURL.String())
	if err != nil {
		return nil, merrors.InternalServerError(g.serviceID, "could not get image from source: %s", err.Error())
	}
	defer r.Close() // nolint:errcheck
	pp := preprocessor.ForType(sRes.GetInfo().GetMimeType())
	img, err := pp.Convert(r)
	if img == nil || err != nil {
		return nil, merrors.InternalServerError(g.serviceID, "could not get image")
	}
	if thumb, err = g.manager.Generate(tr, img); err != nil {
		return nil, err
	}

	return thumb, nil
}

func (g Thumbnail) stat(path, auth string) (*provider.StatResponse, error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(), revactx.TokenHeader, auth)

	req := &provider.StatRequest{
		Ref: &provider.Reference{
			Path: path,
		},
	}
	rsp, err := g.cs3Client.Stat(ctx, req)
	if err != nil {
		g.logger.Error().Err(err).Str("path", path).Msg("could not stat file")
		return nil, merrors.InternalServerError(g.serviceID, "could not stat file: %s", err.Error())
	}

	if rsp.Status.Code != rpc.Code_CODE_OK {
		switch rsp.Status.Code {
		case rpc.Code_CODE_NOT_FOUND:
			return nil, merrors.NotFound(g.serviceID, "could not stat file: %s", rsp.Status.Message)
		default:
			g.logger.Error().Str("status_message", rsp.Status.Message).Str("path", path).Msg("could not stat file")
			return nil, merrors.InternalServerError(g.serviceID, "could not stat file: %s", rsp.Status.Message)
		}
	}
	if rsp.Info.Type != provider.ResourceType_RESOURCE_TYPE_FILE {
		return nil, merrors.BadRequest(g.serviceID, "Unsupported file type")
	}
	if rsp.Info.GetChecksum().GetSum() == "" {
		g.logger.Error().Msg("resource info is missing checksum")
		return nil, merrors.NotFound(g.serviceID, "resource info is missing a checksum")
	}
	if !thumbnail.IsMimeTypeSupported(rsp.Info.MimeType) {
		return nil, merrors.NotFound(g.serviceID, "Unsupported file type")
	}
	return rsp, nil
}
