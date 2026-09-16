//go:build enable_vips

package thumbnail

import (
	"bytes"
	"image"
	"image/draw"
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/stretchr/testify/require"
)

// Guards the txt-preview path: the preprocessor hands us *image.RGBA and the
// vips generator must import it without any encoded round-trip (BMP bytes
// need the ImageMagick loader, absent from minimal libvips runtimes).
func TestSimpleGenerator_TxtRGBA(t *testing.T) {
	require.NoError(t, vips.Startup(nil))
	defer vips.Shutdown()

	src := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(src, src.Bounds(), image.White, image.Point{}, draw.Src)

	gen, err := NewSimpleGenerator("png", "fit")
	require.NoError(t, err)

	dims, err := gen.Dimensions(src)
	require.NoError(t, err)
	require.Equal(t, image.Rect(0, 0, 640, 480), dims)

	got, err := gen.Generate(image.Rect(0, 0, 32, 32), src)
	require.NoError(t, err)
	m, ok := got.(*vips.ImageRef)
	require.True(t, ok, "expected *vips.ImageRef, got %T", got)
	require.Equal(t, 32, m.Width())

	var pngBuf bytes.Buffer
	require.NoError(t, PngEncoder{}.Encode(&pngBuf, m))
	require.NotEmpty(t, pngBuf.Bytes())

	var jpgBuf bytes.Buffer
	require.NoError(t, JpegEncoder{}.Encode(&jpgBuf, m))
	require.NotEmpty(t, jpgBuf.Bytes())
}
