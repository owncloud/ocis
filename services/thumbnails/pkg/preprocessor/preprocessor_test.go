package preprocessor

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"io"
	"os"
	"testing"

	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/thumbnail"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type singleByteReader struct {
	reader *bytes.Reader
}

func (r *singleByteReader) Read(p []byte) (int, error) {
	if len(p) > 1 {
		p = p[:1]
	}
	return r.reader.Read(p)
}

func TestImageDecoder(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "ImageDecoder Suite")
}

func TestGifFirstFrameDecoderDisposal(t *testing.T) {
	for name, disposal := range map[string]byte{
		"unspecified": 0,
		"none":        gif.DisposalNone,
		"background":  gif.DisposalBackground,
		"previous":    gif.DisposalPrevious,
	} {
		t.Run(name, func(t *testing.T) {
			palette := color.Palette{color.Black, color.White}
			frame := image.NewPaletted(image.Rect(0, 0, 2, 2), palette)
			frame.SetColorIndex(0, 0, 1)
			var source bytes.Buffer
			if err := gif.EncodeAll(&source, &gif.GIF{
				Image:    []*image.Paletted{frame, frame},
				Delay:    []int{0, 0},
				Disposal: []byte{disposal, disposal},
			}); err != nil {
				t.Fatal(err)
			}
			img, err := (GifFirstFrameDecoder{}).Convert(&source)
			if err != nil {
				t.Fatal(err)
			}
			decoded := img.(*gif.GIF)
			if len(decoded.Disposal) != len(decoded.Image) || len(decoded.Disposal) != 1 || decoded.Disposal[0] != gif.DisposalNone {
				t.Fatalf("unexpected single-frame disposal: %v", decoded.Disposal)
			}
			generator, err := thumbnail.NewGifGenerator("gif", "resize")
			if err != nil {
				t.Fatal(err)
			}
			// Exercise the generator's unguarded per-frame Disposal access.
			generated, err := generator.Generate(image.Rect(0, 0, 2, 2), decoded)
			if err != nil {
				t.Fatal(err)
			}
			var output bytes.Buffer
			if err := gif.EncodeAll(&output, generated.(*gif.GIF)); err != nil {
				t.Fatal(err)
			}
			result, err := gif.DecodeAll(&output)
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Image) != 1 || len(result.Disposal) != 1 || result.Disposal[0] != gif.DisposalNone {
				t.Fatalf("unexpected encoded thumbnail: %d frames, disposal %v", len(result.Image), result.Disposal)
			}
			if result.Image[0].ColorIndexAt(0, 0) != 1 {
				t.Fatal("first-frame pixel changed")
			}
		})
	}
}

func TestGifDecoderSelection(t *testing.T) {
	palette := color.Palette{color.Black, color.White}
	frame := image.NewPaletted(image.Rect(0, 0, 2, 2), palette)
	var source bytes.Buffer
	if err := gif.EncodeAll(&source, &gif.GIF{
		Image:    []*image.Paletted{frame, frame},
		Delay:    []int{5, 10},
		Disposal: []byte{gif.DisposalBackground, gif.DisposalPrevious},
	}); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name   string
		opts   map[string]interface{}
		frames int
	}{
		{"default", nil, 2},
		{"unset", map[string]interface{}{}, 2},
		{"disabled", map[string]interface{}{"gifFirstFrameOnly": false}, 2},
		{"invalid type", map[string]interface{}{"gifFirstFrameOnly": "true"}, 2},
		{"enabled", map[string]interface{}{"gifFirstFrameOnly": true}, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := ForType("image/gif", tt.opts).Convert(bytes.NewReader(source.Bytes()))
			if err != nil {
				t.Fatal(err)
			}
			result := decoded.(*gif.GIF)
			if len(result.Image) != tt.frames {
				t.Fatalf("got %d frames, want %d", len(result.Image), tt.frames)
			}
			if tt.frames == 2 && (result.Delay[0] != 5 || result.Delay[1] != 10 || result.Disposal[0] != gif.DisposalBackground || result.Disposal[1] != gif.DisposalPrevious) {
				t.Fatal("default decoder changed animation metadata")
			}
		})
	}
}

func TestGifFirstFrameDecoderLogicalScreen(t *testing.T) {
	palette := color.Palette{color.Black, color.White}
	frame := image.NewPaletted(image.Rect(1, 1, 3, 3), palette)
	var source bytes.Buffer
	if err := gif.EncodeAll(&source, &gif.GIF{
		Image:           []*image.Paletted{frame},
		Delay:           []int{0},
		Config:          image.Config{ColorModel: palette, Width: 4, Height: 5},
		BackgroundIndex: 1,
	}); err != nil {
		t.Fatal(err)
	}
	decoded, err := (GifFirstFrameDecoder{}).Convert(&source)
	if err != nil {
		t.Fatal(err)
	}
	result := decoded.(*gif.GIF)
	if result.Config.Width != 4 || result.Config.Height != 5 || result.BackgroundIndex != 1 || result.Image[0].Bounds() != frame.Bounds() {
		t.Fatalf("logical screen or frame bounds changed: %+v, background %d, bounds %v", result.Config, result.BackgroundIndex, result.Image[0].Bounds())
	}
	generator, err := thumbnail.NewGifGenerator("gif", "resize")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := generator.Generate(image.Rect(0, 0, 4, 5), result); err != nil {
		t.Fatal(err)
	}
}

var _ = Describe("ImageDecoder", func() {
	Describe("ImageDecoder", func() {
		var fileReader io.Reader
		BeforeEach(func() {
			fileContent, err := os.ReadFile("test_assets/noise.png")
			if err != nil {
				panic(err)
			}
			fileReader = bytes.NewReader(fileContent)
		})

		It("should decode an image", func() {
			decoder := ImageDecoder{}
			img, err := decoder.Convert(fileReader)
			Expect(err).ToNot(HaveOccurred())
			Expect(img).ToNot(BeNil())
		})

		It("should return an error if the image is invalid", func() {
			decoder := ImageDecoder{}
			img, err := decoder.Convert(bytes.NewReader([]byte("not an image")))
			Expect(err).To(HaveOccurred())
			Expect(img).To(BeNil())
		})
	})

	Describe("GifFirstFrameDecoder", func() {
		var fileReader io.Reader
		BeforeEach(func() {
			fileContent, err := os.ReadFile("test_assets/noise.gif")
			if err != nil {
				panic(err)
			}
			fileReader = bytes.NewReader(fileContent)
		})

		It("should decode a gif", func() {
			decoder := GifFirstFrameDecoder{}
			img, err := decoder.Convert(fileReader)
			Expect(err).ToNot(HaveOccurred())
			Expect(img).ToNot(BeNil())
		})

		It("should decode only the first frame of an animated gif", func() {
			palette := color.Palette{color.Black, color.White}
			first := image.NewPaletted(image.Rect(0, 0, 2, 2), palette)
			second := image.NewPaletted(image.Rect(0, 0, 2, 2), palette)
			second.SetColorIndex(0, 0, 1)

			var encoded bytes.Buffer
			err := gif.EncodeAll(&encoded, &gif.GIF{
				Image:           []*image.Paletted{first, second},
				Delay:           []int{0, 0},
				BackgroundIndex: 1,
				Config: image.Config{
					ColorModel: palette,
					Width:      2,
					Height:     2,
				},
			})
			Expect(err).ToNot(HaveOccurred())

			reader := &singleByteReader{reader: bytes.NewReader(encoded.Bytes())}
			decoder := GifFirstFrameDecoder{}
			img, err := decoder.Convert(reader)
			Expect(err).ToNot(HaveOccurred())

			decoded, ok := img.(*gif.GIF)
			Expect(ok).To(BeTrue())
			Expect(decoded.Image).To(HaveLen(1))
			Expect(decoded.Disposal).To(Equal([]byte{gif.DisposalNone}))
			Expect(decoded.Image[0].ColorIndexAt(0, 0)).To(Equal(uint8(0)))
			Expect(decoded.Config.Width).To(Equal(2))
			Expect(decoded.Config.Height).To(Equal(2))
			Expect(decoded.LoopCount).To(Equal(-1))
			Expect(decoded.BackgroundIndex).To(Equal(uint8(1)))
			Expect(reader.reader.Len()).To(BeNumerically(">", 0))

			var thumbnail bytes.Buffer
			Expect(gif.EncodeAll(&thumbnail, decoded)).To(Succeed())
			encodedThumbnail, err := gif.DecodeAll(&thumbnail)
			Expect(err).ToNot(HaveOccurred())
			Expect(encodedThumbnail.Image).To(HaveLen(1))
			Expect(encodedThumbnail.Disposal).To(Equal([]byte{gif.DisposalNone}))
		})

		It("should return an error if the gif is invalid", func() {
			decoder := GifFirstFrameDecoder{}
			img, err := decoder.Convert(bytes.NewReader([]byte("not a gif")))
			Expect(err).To(HaveOccurred())
			Expect(img).To(BeNil())
		})

		It("should return an error if the first frame is invalid", func() {
			decoder := GifFirstFrameDecoder{}
			validHeaderWithoutFrame := []byte("GIF89a\x01\x00\x01\x00\x00\x00\x00")
			img, err := decoder.Convert(bytes.NewReader(validHeaderWithoutFrame))
			Expect(err).To(HaveOccurred())
			Expect(img).To(BeNil())
		})
	})

	Describe("GgsDecoder", func() {
		var fileReader io.Reader
		BeforeEach(func() {
			fileContent, err := os.ReadFile("test_assets/ggs_test.ggs")
			if err != nil {
				panic(err)
			}
			fileReader = bytes.NewReader(fileContent)
		})

		It("should decode a ggs", func() {
			decoder := GgsDecoder{"_slide0/geogebra_thumbnail.png"}
			img, err := decoder.Convert(fileReader)
			Expect(err).ToNot(HaveOccurred())
			Expect(img).ToNot(BeNil())
		})

		It("should return an error if the ggs is invalid", func() {
			decoder := GgsDecoder{"_slide0/geogebra_thumbnail.png"}
			img, err := decoder.Convert(bytes.NewReader([]byte("not a ggs")))
			Expect(err).To(HaveOccurred())
			Expect(img).To(BeNil())
		})
	})

	Describe("should decode audio", func() {
		var fileReader io.Reader
		It("should decode an audio", func() {
			fileContent, err := os.ReadFile("test_assets/empty.mp3")
			if err != nil {
				panic(err)
			}
			fileReader = bytes.NewReader(fileContent)
			decoder := AudioDecoder{}
			img, err := decoder.Convert(fileReader)
			Expect(err).ToNot(HaveOccurred())
			Expect(img).ToNot(BeNil())
		})
		It("should decode an audio", func() {
			fileContent, err := os.ReadFile("test_assets/empty_no_image.mp3")
			if err != nil {
				panic(err)
			}
			fileReader = bytes.NewReader(fileContent)
			decoder := AudioDecoder{}
			img, err := decoder.Convert(fileReader)
			Expect(err).To(HaveOccurred())
			Expect(img).To(BeNil())
		})
		It("should return an error if the audio is invalid", func() {
			decoder := AudioDecoder{}
			img, err := decoder.Convert(bytes.NewReader([]byte("not an audio")))
			Expect(err).To(HaveOccurred())
			Expect(img).To(BeNil())
		})
	})

	Describe("should decode text", func() {
		var decoder TxtToImageConverter
		BeforeEach(func() {
			fontFaceOpts := &opentype.FaceOptions{
				Size:    12,
				DPI:     72,
				Hinting: font.HintingNone,
			}

			fontLoader, err := NewFontLoader("", fontFaceOpts)
			if err != nil {
				fontLoader, _ = NewFontLoader("", fontFaceOpts)
			}
			decoder = TxtToImageConverter{
				fontLoader: fontLoader,
			}
		})
		It("should decode a text", func() {
			img, err := decoder.Convert(bytes.NewReader([]byte("This is a test text")))
			Expect(err).ToNot(HaveOccurred())
			Expect(img).ToNot(BeNil())
		})
	})

	Describe("should decode bmp", func() {
		It("should decode a bmp", func() {
			fileContent, err := os.ReadFile("test_assets/bmp_test.bmp")
			if err != nil {
				panic(err)
			}
			decoder := BmpDecoder{}
			img, err := decoder.Convert(bytes.NewReader(fileContent))
			Expect(err).ToNot(HaveOccurred())
			Expect(img).ToNot(BeNil())
		})
		It("should return an error if the bmp is invalid", func() {
			decoder := BmpDecoder{}
			img, err := decoder.Convert(bytes.NewReader([]byte("not a bmp")))
			Expect(err).To(HaveOccurred())
			Expect(img).To(BeNil())
		})
	})

	Describe("test ForType", func() {
		It("should return an ImageDecoder for image types", func() {
			decoder := ForType("image/png", nil)
			Expect(decoder).To(BeAssignableToTypeOf(ImageDecoder{}))
		})

		It("should return an GifDecoder for gif types", func() {
			decoder := ForType("image/gif", nil)
			Expect(decoder).To(BeAssignableToTypeOf(GifDecoder{}))
		})

		It("should return an GgsDecoder for ggs types", func() {
			decoder := ForType("application/vnd.geogebra.ggs", nil)
			// This will not return the expected ggsDecoder, but an ImageDecoder since ggs contains an embedded png.
			Expect(decoder).To(BeAssignableToTypeOf(ImageDecoder{}))
		})

		It("should return an AudioDecoder for audio types", func() {
			decoder := ForType("audio/mpeg", nil)
			Expect(decoder).To(BeAssignableToTypeOf(AudioDecoder{}))
		})

		It("should return an TxtToImageConverter for text types", func() {
			decoder := ForType("text/plain", nil)
			Expect(decoder).To(BeAssignableToTypeOf(TxtToImageConverter{}))
		})

		It("should return a BmpDecoder for bmp types", func() {
			decoder := ForType("image/bmp", nil)
			Expect(decoder).To(BeAssignableToTypeOf(BmpDecoder{}))
		})

		It("should return a BmpDecoder for x-ms-bmp types", func() {
			decoder := ForType("image/x-ms-bmp", nil)
			Expect(decoder).To(BeAssignableToTypeOf(BmpDecoder{}))
		})

		It("should return an ImageDecoder for unknown types", func() {
			decoder := ForType("unknown", nil)
			Expect(decoder).To(BeAssignableToTypeOf(ImageDecoder{}))
		})
	})
})
