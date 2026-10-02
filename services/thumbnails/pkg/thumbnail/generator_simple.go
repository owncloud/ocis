//go:build !enable_vips

package thumbnail

import (
	"context"
	"image"

	"github.com/kovidgoyal/imaging"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// SimpleGenerator is the default image generator and is used for all image types expect gif.
type SimpleGenerator struct {
	processor Processor
}

func NewSimpleGenerator(filetype, process string) (SimpleGenerator, error) {
	processor, err := ProcessorFor(process, filetype)
	if err != nil {
		return SimpleGenerator{}, err
	}
	return SimpleGenerator{processor: processor}, nil
}

// ProcessorID returns the processor identification.
func (g SimpleGenerator) ProcessorID() string {
	return g.processor.ID()
}

// Generate generates a alternative image version.
func (g SimpleGenerator) Generate(ctx context.Context, size image.Rectangle, img interface{}) (interface{}, error) {
	span := trace.SpanFromContext(ctx)
	newCtx, newSpan := span.TracerProvider().Tracer(tracerName).Start(
		ctx, spanNameGeneratorGenerate,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("ocis.thumbnails.generator.type", "SimpleGenerator"),
			attribute.Int("ocis.thumbnails.generator.generate.width", size.Dx()),
			attribute.Int("ocis.thumbnails.generator.generate.height", size.Dy()),
			attribute.String("ocis.thumbnails.generator.generate.processor", g.ProcessorID()),
		),
	)
	defer newSpan.End()

	m, ok := img.(image.Image)
	if !ok {
		return nil, errors.ErrInvalidType
	}

	return g.processor.Process(newCtx, m, size.Dx(), size.Dy(), imaging.Lanczos), nil
}

func (g SimpleGenerator) Dimensions(img interface{}) (image.Rectangle, error) {
	m, ok := img.(image.Image)
	if !ok {
		return image.Rectangle{}, errors.ErrInvalidType
	}
	return m.Bounds(), nil
}
