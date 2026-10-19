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

	Describe("GifDecoder", func() {
		var fileReader io.Reader
		BeforeEach(func() {
			fileContent, err := os.ReadFile("test_assets/noise.gif")
			Expect(err).ToNot(HaveOccurred())
			fileReader = bytes.NewReader(fileContent)
		})

		It("should decode a gif", func() {
			decoder := GifDecoder{}
			img, err := decoder.Convert(fileReader)
			Expect(err).ToNot(HaveOccurred())
			Expect(img).ToNot(BeNil())
		})

		It("should return an error if the gif is invalid", func() {
			decoder := GifDecoder{}
			img, err := decoder.Convert(bytes.NewReader([]byte("not a gif")))
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

		for name, disposal := range map[string]byte{
			"unspecified": 0,
			"none":        gif.DisposalNone,
			"background":  gif.DisposalBackground,
			"previous":    gif.DisposalPrevious,
		} {
			name, disposal := name, disposal
			It("should normalize "+name+" disposal for a static thumbnail", func() {
				palette := color.Palette{color.Black, color.White}
				frame := image.NewPaletted(image.Rect(0, 0, 2, 2), palette)
				frame.SetColorIndex(0, 0, 1)
				var source bytes.Buffer
				Expect(gif.EncodeAll(&source, &gif.GIF{
					Image:    []*image.Paletted{frame, frame},
					Delay:    []int{0, 0},
					Disposal: []byte{disposal, disposal},
				})).To(Succeed())

				img, err := (GifFirstFrameDecoder{}).Convert(&source)
				Expect(err).ToNot(HaveOccurred())
				decoded, ok := img.(*gif.GIF)
				Expect(ok).To(BeTrue())
				Expect(decoded.Image).To(HaveLen(1))
				Expect(decoded.Disposal).To(Equal([]byte{gif.DisposalNone}))

				generator, err := thumbnail.NewGifGenerator("gif", "resize")
				Expect(err).ToNot(HaveOccurred())
				generated, err := generator.Generate(image.Rect(0, 0, 2, 2), decoded)
				Expect(err).ToNot(HaveOccurred())

				var output bytes.Buffer
				Expect(gif.EncodeAll(&output, generated.(*gif.GIF))).To(Succeed())
				result, err := gif.DecodeAll(&output)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.Image).To(HaveLen(1))
				Expect(result.Disposal).To(Equal([]byte{gif.DisposalNone}))
				Expect(result.Image[0].ColorIndexAt(0, 0)).To(Equal(uint8(1)))
			})
		}

		It("should preserve the logical screen and frame bounds", func() {
			palette := color.Palette{color.Black, color.White}
			frame := image.NewPaletted(image.Rect(1, 1, 3, 3), palette)
			var source bytes.Buffer
			Expect(gif.EncodeAll(&source, &gif.GIF{
				Image:           []*image.Paletted{frame},
				Delay:           []int{0},
				Config:          image.Config{ColorModel: palette, Width: 4, Height: 5},
				BackgroundIndex: 1,
			})).To(Succeed())

			decoded, err := (GifFirstFrameDecoder{}).Convert(&source)
			Expect(err).ToNot(HaveOccurred())
			result, ok := decoded.(*gif.GIF)
			Expect(ok).To(BeTrue())
			Expect(result.Config.Width).To(Equal(4))
			Expect(result.Config.Height).To(Equal(5))
			Expect(result.BackgroundIndex).To(Equal(uint8(1)))
			Expect(result.Image[0].Bounds()).To(Equal(frame.Bounds()))

			generator, err := thumbnail.NewGifGenerator("gif", "resize")
			Expect(err).ToNot(HaveOccurred())
			_, err = generator.Generate(image.Rect(0, 0, 4, 5), result)
			Expect(err).ToNot(HaveOccurred())
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
		for _, test := range []struct {
			name   string
			opts   map[string]interface{}
			frames int
		}{
			{name: "by default", frames: 2},
			{name: "when unset", opts: map[string]interface{}{}, frames: 2},
			{name: "when disabled", opts: map[string]interface{}{"gifFirstFrameOnly": false}, frames: 2},
			{name: "when invalid", opts: map[string]interface{}{"gifFirstFrameOnly": "true"}, frames: 2},
			{name: "when enabled", opts: map[string]interface{}{"gifFirstFrameOnly": true}, frames: 1},
		} {
			test := test
			It("should select the expected gif decoder "+test.name, func() {
				palette := color.Palette{color.Black, color.White}
				frame := image.NewPaletted(image.Rect(0, 0, 2, 2), palette)
				var animatedGIF bytes.Buffer
				Expect(gif.EncodeAll(&animatedGIF, &gif.GIF{
					Image:    []*image.Paletted{frame, frame},
					Delay:    []int{5, 10},
					Disposal: []byte{gif.DisposalBackground, gif.DisposalPrevious},
				})).To(Succeed())

				decoded, err := ForType("image/gif", test.opts).Convert(bytes.NewReader(animatedGIF.Bytes()))
				Expect(err).ToNot(HaveOccurred())
				result := decoded.(*gif.GIF)
				Expect(result.Image).To(HaveLen(test.frames))
				if test.frames == 2 {
					Expect(result.Delay).To(Equal([]int{5, 10}))
					Expect(result.Disposal).To(Equal([]byte{gif.DisposalBackground, gif.DisposalPrevious}))
				}
			})
		}

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
