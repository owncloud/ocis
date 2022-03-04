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
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/ocis-pkg/log"
	thumbnailssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/thumbnails/v0"
	"github.com/owncloud/ocis/thumbnails/pkg/preprocessor"
	"github.com/owncloud/ocis/thumbnails/pkg/service/v0/decorators"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/imgsource"
	"github.com/pkg/errors"
	merrors "go-micro.dev/v4/errors"
	"google.golang.org/grpc/metadata"
)

// NewService returns a service implementation for Service.
func NewService(opts ...Option) decorators.DecoratedService {
	options := newOptions(opts...)
	logger := options.Logger
	resolutions, err := thumbnail.ParseResolutions(options.Config.Thumbnail.Resolutions)
	if err != nil {
		logger.Fatal().Err(err).Msg("resolutions not configured correctly")
	}
	svc := Thumbnail{
		serviceID: options.Config.GRPC.Namespace + "." + options.Config.Service.Name,
		manager: thumbnail.NewSimpleManager(
			resolutions,
			options.ThumbnailStorage,
			logger,
		),
		webdavSource: options.ImageSource,
		cs3Source:    options.CS3Source,
		logger:       logger,
		cs3Client:    options.CS3Client,
		preprocessorOpts: PreprocessorOpts{
			TxtFontFileMap: options.Config.Thumbnail.FontMapFile,
		},
	}

	return svc
}

// Thumbnail implements the GRPC handler.
type Thumbnail struct {
	serviceID        string
	manager          thumbnail.Manager
	webdavSource     imgsource.Source
	cs3Source        imgsource.Source
	logger           log.Logger
	cs3Client        gateway.GatewayAPIClient
	preprocessorOpts PreprocessorOpts
}

type PreprocessorOpts struct {
	TxtFontFileMap string
}

// GetThumbnail retrieves a thumbnail for an image
func (g Thumbnail) GetThumbnail(ctx context.Context, req *thumbnailssvc.GetThumbnailRequest, rsp *thumbnailssvc.GetThumbnailResponse) error {
	_, ok := thumbnailssvc.GetThumbnailRequest_ThumbnailType_value[req.ThumbnailType.String()]
	if !ok {
		g.logger.Debug().Str("thumbnail_type", req.ThumbnailType.String()).Msg("unsupported thumbnail type")
		return nil
	}
	generator, err := thumbnail.GeneratorForType(req.ThumbnailType.String())
	if err != nil {
		g.logger.Debug().Str("thumbnail_type", req.ThumbnailType.String()).Msg("unsupported thumbnail type")
		return nil
	}
	encoder, err := thumbnail.EncoderForType(req.ThumbnailType.String())
	if err != nil {
		g.logger.Debug().Str("thumbnail_type", req.ThumbnailType.String()).Msg("unsupported thumbnail type")
		return nil
	}

	var thumb []byte
	switch {
	case req.GetWebdavSource() != nil:
		thumb, err = g.handleWebdavSource(ctx, req, generator, encoder)
	case req.GetCs3Source() != nil:
		thumb, err = g.handleCS3Source(ctx, req, generator, encoder)
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

func (g Thumbnail) handleCS3Source(ctx context.Context, req *thumbnailssvc.GetThumbnailRequest, generator thumbnail.Generator, encoder thumbnail.Encoder) ([]byte, error) {
	src := req.GetCs3Source()
	sRes, err := g.stat(src.Path, src.Authorization)
	if err != nil {
		return nil, err
	}

	tr := thumbnail.Request{
		Resolution: image.Rect(0, 0, int(req.Width), int(req.Height)),
		Generator:  generator,
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
	ppOpts := map[string]interface{}{
		"fontFileMap": g.preprocessorOpts.TxtFontFileMap,
	}
	pp := preprocessor.ForType(sRes.GetInfo().GetMimeType(), ppOpts)
	img, err := pp.Convert(r)
	if img == nil || err != nil {
		return nil, merrors.InternalServerError(g.serviceID, "could not get image")
	}
	if thumb, err = g.manager.Generate(tr, img); err != nil {
		return nil, err
	}

	return thumb, nil
}

func (g Thumbnail) handleWebdavSource(ctx context.Context, req *thumbnailssvc.GetThumbnailRequest, generator thumbnail.Generator, encoder thumbnail.Encoder) ([]byte, error) {
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
		statPath = req.Filepath
	}
	sRes, err := g.stat(statPath, auth)
	if err != nil {
		return nil, err
	}
	tr := thumbnail.Request{
		Resolution: image.Rect(0, 0, int(req.Width), int(req.Height)),
		Generator:  generator,
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
	ppOpts := map[string]interface{}{
		"fontFileMap": g.preprocessorOpts.TxtFontFileMap,
	}
	pp := preprocessor.ForType(sRes.GetInfo().GetMimeType(), ppOpts)
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

	var ref *provider.Reference
	if strings.Contains(path, "!") {
		parts := strings.Split(path, "!")
		spaceID, path := parts[0], parts[1]
		ref = &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: spaceID,
				OpaqueId:  spaceID,
			},
			Path: path,
		}
	} else {
		ref = &provider.Reference{
			Path: path,
		}
	}

	req := &provider.StatRequest{Ref: ref}
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
