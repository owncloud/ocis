//go:build enable_vips

package thumbnail

import (
	"bytes"
	"image"
	"image/draw"
	"os"
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/preprocessor"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/bmp"
)

// Start vips once for the package.
func TestMain(m *testing.M) {
	if err := vips.Startup(nil); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// Guards the txt-preview path: the preprocessor hands us *image.RGBA and the
// vips generator must import it without any encoded round-trip (BMP bytes
// need the ImageMagick loader, absent from minimal libvips runtimes).
func TestSimpleGenerator_TxtRGBA(t *testing.T) {
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

// Guards BMP sources: decoded in Go by BmpDecoder, imported without the
// magick loader that minimal libvips runtimes lack.
func TestSimpleGenerator_BmpSource(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	draw.Draw(src, src.Bounds(), image.White, image.Point{}, draw.Src)

	var buf bytes.Buffer
	require.NoError(t, bmp.Encode(&buf, src))

	img, err := preprocessor.ForType("image/bmp", nil).Convert(&buf)
	require.NoError(t, err)
	require.NotNil(t, img)

	gen, err := NewSimpleGenerator("jpg", "fit")
	require.NoError(t, err)

	dims, err := gen.Dimensions(img)
	require.NoError(t, err)
	require.Equal(t, image.Rect(0, 0, 32, 32), dims)

	got, err := gen.Generate(image.Rect(0, 0, 32, 32), img)
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
