package thumbnail

import "testing"

func TestEncoderForType(t *testing.T) {
	table := map[string]Encoder{
		"jpg":     JpegEncoder{},
		"JPG":     JpegEncoder{},
		"jpeg":    JpegEncoder{},
		"JPEG":    JpegEncoder{},
		"png":     PngEncoder{},
		"PNG":     PngEncoder{},
		"invalid": nil,
	}

	for k, v := range table {
		e, _ := EncoderForType(k)
		if e != v {
			t.Fail()
		}
	}
}
