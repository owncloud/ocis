package thumbnail

import (
	"bytes"
	"context"
	"image"
	"mime"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/errors"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/thumbnail/storage"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName       = "thumbnails"
	spanNameGenerate = "Manager.Generate"
	spanNameCheck    = "Manager.CheckThumbnail"
	spanNameGet      = "Manager.GetThumbnail"
)

var (
	// SupportedMimeTypes contains an all mimetypes which are supported by the thumbnailer.
	SupportedMimeTypes = map[string]struct{}{
		"image/png":                         {},
		"image/jpg":                         {},
		"image/jpeg":                        {},
		"image/gif":                         {},
		"image/bmp":                         {},
		"image/x-ms-bmp":                    {},
		"image/tiff":                        {},
		"text/plain":                        {},
		"audio/flac":                        {},
		"audio/mpeg":                        {},
		"audio/ogg":                         {},
		"application/vnd.geogebra.slides":   {},
		"application/vnd.geogebra.pinboard": {},
	}
)

// Request bundles information needed to generate a thumbnail for a file
type Request struct {
	Resolution image.Rectangle
	Encoder    Encoder
	Generator  Generator
	Checksum   string
}

// Manager is responsible for generating thumbnails
type Manager interface {
	// Generate creates a thumbnail and stores it.
	// The function returns a key with which the actual file can be retrieved.
	Generate(ctx context.Context, r Request, img interface{}) (string, error)
	// CheckThumbnail checks if a thumbnail with the requested attributes exists.
	// The function will return a status if the file exists and the key to the file.
	CheckThumbnail(ctx context.Context, r Request) (string, bool)
	// GetThumbnail will load the thumbnail from the storage and return its content.
	GetThumbnail(ctx context.Context, key string) ([]byte, error)
}

// NewSimpleManager creates a new instance of SimpleManager
func NewSimpleManager(resolutions Resolutions, storage storage.Storage, logger log.Logger, maxInputWidth, maxInputHeight int) SimpleManager {
	return SimpleManager{
		storage:      storage,
		logger:       logger,
		resolutions:  resolutions,
		maxDimension: image.Point{X: maxInputWidth, Y: maxInputHeight},
	}
}

// SimpleManager is a simple implementation of Manager
type SimpleManager struct {
	storage      storage.Storage
	logger       log.Logger
	resolutions  Resolutions
	maxDimension image.Point
}

// Generate creates a thumbnail and stores it
func (s SimpleManager) Generate(ctx context.Context, r Request, img interface{}) (string, error) {
	span := trace.SpanFromContext(ctx)
	newCtx, newSpan := span.TracerProvider().Tracer(tracerName).Start(
		ctx, spanNameGenerate,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("ocis.thumbnails.manager.type", "SimpleManager"),
			attribute.Int("ocis.thumbnails.manager.generate.req.width", r.Resolution.Dx()),
			attribute.Int("ocis.thumbnails.manager.generate.req.height", r.Resolution.Dy()),
			attribute.String("ocis.thumbnails.manager.generate.req.encoding", r.Encoder.MimeType()),
			attribute.String("ocis.thumbnails.manager.generate.req.generator", r.Generator.ProcessorID()),
		),
	)
	defer newSpan.End()

	var match image.Rectangle

	inputDimensions, err := r.Generator.Dimensions(img)
	if err != nil {
		return "", err
	}
	match = s.resolutions.ClosestMatch(r.Resolution, inputDimensions)

	// validate max input image dimensions - 6016x4000
	if inputDimensions.Size().X > s.maxDimension.X || inputDimensions.Size().Y > s.maxDimension.Y {
		return "", errors.ErrImageTooLarge
	}

	thumbnail, err := r.Generator.Generate(newCtx, match, img)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := r.Encoder.Encode(newCtx, buf, thumbnail); err != nil {
		return "", err
	}

	k := s.storage.BuildKey(newCtx, mapToStorageRequest(r))
	if err := s.storage.Put(newCtx, k, buf.Bytes()); err != nil {
		s.logger.Error().Err(err).Msg("could not store thumbnail")
		return "", err
	}
	return k, nil
}

// CheckThumbnail checks if a thumbnail with the requested attributes exists.
func (s SimpleManager) CheckThumbnail(ctx context.Context, r Request) (string, bool) {
	span := trace.SpanFromContext(ctx)
	newCtx, newSpan := span.TracerProvider().Tracer(tracerName).Start(
		ctx, spanNameCheck,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("ocis.thumbnails.manager.type", "SimpleManager"),
			attribute.Int("ocis.thumbnails.manager.generate.req.width", r.Resolution.Dx()),
			attribute.Int("ocis.thumbnails.manager.generate.req.height", r.Resolution.Dy()),
			attribute.String("ocis.thumbnails.manager.generate.req.encoding", r.Encoder.MimeType()),
			attribute.String("ocis.thumbnails.manager.generate.req.generator", r.Generator.ProcessorID()),
		),
	)
	defer newSpan.End()

	k := s.storage.BuildKey(newCtx, mapToStorageRequest(r))
	return k, s.storage.Stat(newCtx, k)
}

// GetThumbnail will load the thumbnail from the storage and return its content.
func (s SimpleManager) GetThumbnail(ctx context.Context, key string) ([]byte, error) {
	span := trace.SpanFromContext(ctx)
	newCtx, newSpan := span.TracerProvider().Tracer(tracerName).Start(
		ctx, spanNameGet,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("ocis.thumbnails.manager.type", "SimpleManager"),
		),
	)
	defer newSpan.End()

	return s.storage.Get(newCtx, key)
}

func mapToStorageRequest(r Request) storage.Request {
	return storage.Request{
		Checksum:       r.Checksum,
		Resolution:     r.Resolution,
		Types:          r.Encoder.Types(),
		Characteristic: r.Generator.ProcessorID(),
	}
}

// IsMimeTypeSupported validate if the mime type is supported
func IsMimeTypeSupported(m string) bool {
	mimeType, _, err := mime.ParseMediaType(m)
	if err != nil {
		return false
	}
	_, supported := SupportedMimeTypes[mimeType]
	return supported
}

// PrepareRequest prepare the request based on image parameters
func PrepareRequest(width, height int, tType, checksum, pID string) (Request, error) {
	generator, err := GeneratorFor(tType, pID)
	if err != nil {
		return Request{}, err
	}
	encoder, err := EncoderForType(tType)
	if err != nil {
		return Request{}, err
	}

	return Request{
		Resolution: image.Rect(0, 0, width, height),
		Generator:  generator,
		Encoder:    encoder,
		Checksum:   checksum,
	}, nil
}
