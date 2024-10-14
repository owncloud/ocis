//go:build enable_vips

package thumbnail

import (
	"image"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/errors"
)

// SimpleGenerator is the default image generator and is used for all image types expect gif.
type SimpleGenerator struct {
	crop    vips.Interesting
	size    vips.Size
	process string
}

func NewSimpleGenerator(filetype, process string) (SimpleGenerator, error) {
	switch strings.ToLower(process) {
	case "thumbnail":
		return SimpleGenerator{crop: vips.InterestingAttention, process: process, size: vips.SizeBoth}, nil
	case "fit":
		return SimpleGenerator{crop: vips.InterestingNone, process: process, size: vips.SizeBoth}, nil
	case "resize":
		return SimpleGenerator{crop: vips.InterestingNone, process: process, size: vips.SizeForce}, nil
	default:
		return SimpleGenerator{crop: vips.InterestingNone, process: process}, nil
	}
}

// ProcessorID returns the processor identification.
func (g SimpleGenerator) ProcessorID() string {
	return g.process
}

// Generate generates a alternative image version.
func (g SimpleGenerator) Generate(size image.Rectangle, img interface{}) (interface{}, error) {
	m, ok := img.(*vips.ImageRef)
	if !ok {
		return nil, errors.ErrInvalidType
	}

	if err := m.ThumbnailWithSize(size.Dx(), 0, g.crop, g.size); err != nil {
		return nil, err
	}

	if err := m.RemoveMetadata(); err != nil {
		return nil, err
	}

	return m, nil
}

func (g SimpleGenerator) Dimensions(img interface{}) (image.Rectangle, error) {
	m, ok := img.(*vips.ImageRef)
	if !ok {
		return image.Rectangle{}, errors.ErrInvalidType
	}
	return image.Rect(0, 0, m.Width(), m.Height()), nil
}
