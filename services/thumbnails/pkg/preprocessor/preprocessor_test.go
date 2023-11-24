package preprocessor

// implement tests for the following structs and functions using ginkgo and gomega:
// ImageDecoder
// GifDecoder
// GgsDecoder
// Convert

import (
	"bytes"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"io"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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

		It("should return an error if the gif is invalid", func() {
			decoder := GifDecoder{}
			img, err := decoder.Convert(bytes.NewReader([]byte("not a gif")))
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
			decoder := GgsDecoder{}
			img, err := decoder.Convert(fileReader)
			Expect(err).ToNot(HaveOccurred())
			Expect(img).ToNot(BeNil())
		})

		It("should return an error if the ggs is invalid", func() {
			decoder := GgsDecoder{}
			img, err := decoder.Convert(bytes.NewReader([]byte("not a ggs")))
			Expect(err).To(HaveOccurred())
			Expect(img).To(BeNil())
		})
	})

	Describe("should decode audio", func() {
		//TODO: this needs test, but I have no idea what kind of audio
		//      I can leave here without causing legal trouble
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

		It("should return an ImageDecoder for unknown types", func() {
			decoder := ForType("unknown", nil)
			Expect(decoder).To(BeAssignableToTypeOf(ImageDecoder{}))
		})
	})
})
