package preprocessor

import (
	"bufio"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/draw"
	"io"
	"mime"
	"strings"
)

const (
	fontSize         = 12
	spacing  float64 = 1.5
)

type FileConverter interface {
	Convert(r io.Reader) (image.Image, error)
}

type ImageDecoder struct{}

func (i ImageDecoder) Convert(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, errors.Wrap(err, `could not decode the image`)
	}
	return img, nil
}

type TxtToImageConverter struct{}

func (t TxtToImageConverter) Convert(r io.Reader) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	c := freetype.NewContext()
	// Ignoring the error since we are using the embedded Golang font.
	// This shouldn't return an error.
	f, _ := truetype.Parse(goregular.TTF)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)
	c.SetHinting(font.HintingFull)
	pt := freetype.Pt(10, 10+int(c.PointToFixed(fontSize)>>6))

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()
		cs := chunks(txt, 80)
		for _, s := range cs {
			_, err := c.DrawString(strings.TrimSpace(s), pt)
			if err != nil {
				return nil, err
			}
			pt.Y += c.PointToFixed(fontSize * spacing)
			if pt.Y.Round() >= img.Bounds().Dy() {
				return img, scanner.Err()
			}
		}

	}
	return img, scanner.Err()
}

// Code from https://stackoverflow.com/a/61469854
// Written By Igor Mikushkin
func chunks(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	length := 0
	for _, r := range s {
		chunk[length] = r
		length++
		if length == chunkSize {
			chunks = append(chunks, string(chunk))
			length = 0
		}
	}
	if length > 0 {
		chunks = append(chunks, string(chunk[:length]))
	}
	return chunks
}

func ForType(mimeType string) FileConverter {
	// We can ignore the error here because we parse it in IsMimeTypeSupported before and if it fails
	// return the service call. So we should only get here when the mimeType parses fine.
	mimeType, _, _ = mime.ParseMediaType(mimeType)
	switch mimeType {
	case "text/plain":
		return TxtToImageConverter{}
	default:
		return ImageDecoder{}
	}
}
