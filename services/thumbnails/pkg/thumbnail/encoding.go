package thumbnail

import (
	"context"
	"image/gif"
	"io"
	"strings"

	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	typePng  = "png"
	typeJpg  = "jpg"
	typeJpeg = "jpeg"
	typeGif  = "gif"
	typeGgs  = "ggs"
	typeGgp  = "ggp"

	//tracerName       = "thumbnails" // already defined in the thumbnails.go file
	spanNameEncode = "Encoder.Encode"
)

// Encoder encodes the thumbnail to a specific format.
type Encoder interface {
	// Encode encodes the image to a format.
	Encode(ctx context.Context, w io.Writer, img interface{}) error
	// Types returns the formats suffixes.
	Types() []string
	// MimeType returns the mimetype used by the encoder.
	MimeType() string
}

// GifEncoder encodes to gif
type GifEncoder struct{}

// Encode encodes the image to a gif format
func (e GifEncoder) Encode(ctx context.Context, w io.Writer, img interface{}) error {
	span := trace.SpanFromContext(ctx)
	_, newSpan := span.TracerProvider().Tracer(tracerName).Start(
		ctx, spanNameEncode,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("ocis.thumbnails.encoder.type", "GifEncoder"),
			attribute.String("ocis.thumbnails.encoder.mime", e.MimeType()),
		),
	)
	defer newSpan.End()

	g, ok := img.(*gif.GIF)
	if !ok {
		return errors.ErrInvalidType
	}
	return gif.EncodeAll(w, g)
}

// Types returns the supported types of the GifEncoder
func (e GifEncoder) Types() []string {
	return []string{typeGif}
}

// MimeType returns the mimetype used by the encoder.
func (e GifEncoder) MimeType() string {
	return "image/gif"
}

// EncoderForType returns the encoder for a given file type
// or nil if the type is not supported.
func EncoderForType(fileType string) (Encoder, error) {
	switch strings.ToLower(fileType) {
	case typePng, typeGgs, typeGgp:
		return PngEncoder{}, nil
	case typeJpg, typeJpeg:
		return JpegEncoder{}, nil
	case typeGif:
		return GifEncoder{}, nil
	default:
		return nil, errors.ErrNoEncoderForType
	}
}

// GetExtForMime return the supported extension by mime
func GetExtForMime(fileType string) string {
	ext := strings.TrimPrefix(strings.TrimSpace(strings.ToLower(fileType)), "image/")
	switch ext {
	case typeJpg, typeJpeg, typePng, typeGif:
		return ext
	case "application/vnd.geogebra.slides":
		return typeGgs
	case "application/vnd.geogebra.pinboard":
		return typeGgp
	default:
		return ""
	}
}
