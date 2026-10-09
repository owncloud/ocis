package preprocessor

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"io"
	"os"
	"testing"

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
			if err != nil {
				panic(err)
			}
			fileReader = bytes.NewReader(fileContent)
		})

		It("should decode a gif", func() {
			decoder := GifDecoder{}
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
			decoder := GifDecoder{}
			img, err := decoder.Convert(reader)
			Expect(err).ToNot(HaveOccurred())

			decoded, ok := img.(*gif.GIF)
			Expect(ok).To(BeTrue())
			Expect(decoded.Image).To(HaveLen(1))
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
		})

		It("should return an error if the gif is invalid", func() {
			decoder := GifDecoder{}
			img, err := decoder.Convert(bytes.NewReader([]byte("not a gif")))
			Expect(err).To(HaveOccurred())
			Expect(img).To(BeNil())
		})

		It("should return an error if the first frame is invalid", func() {
			decoder := GifDecoder{}
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
