package proto_test

import (
	"testing"

	"github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	testDataName  string
	filepath      string
	filetype      proto.GetThumbnailRequest_ThumbnailType
	Checksum      string
	width         int32
	height        int32
	expected      proto.GetThumbnailRequest
}

type TestResponse struct {
	testDataName string
	img          []byte
	mimetype     string
	expected     proto.GetThumbnailResponse
}

func TestRequestString(t *testing.T) {

	var tests = []*TestRequest{
		{
			"ASCII",
			"Foo.jpg",
			proto.GetThumbnailRequest_JPG,
			"33a64df551425fcc55e4d42a148795d9f25f89d4",
			24,
			24,
			proto.GetThumbnailRequest{
				Filepath:      "Foo.jpg",
				ThumbnailType: proto.GetThumbnailRequest_JPG,
				Width:         24,
				Height:        24,
			},
		},
		{
			"UTF",
			"मिलन.jpg",
			proto.GetThumbnailRequest_JPG,
			"33a64df551425fcc55e4d42a148795d9f25f89d4",
			24,
			24,
			proto.GetThumbnailRequest{
				Filepath:      "\340\244\256\340\244\277\340\244\262\340\244\250.jpg",
				ThumbnailType: proto.GetThumbnailRequest_JPG,
				Width:         24,
				Height:        24,
			},
		},
		{
			"PNG",
			"Foo.png",
			proto.GetThumbnailRequest_PNG,
			"33a64df551425fcc55e4d42a148795d9f25f89d4",
			24,
			24,
			proto.GetThumbnailRequest{
				Filepath: "Foo.png",
				Width:    24,
				Height:   24,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {
			req := proto.GetThumbnailRequest{
				Filepath:      testCase.filepath,
				ThumbnailType: testCase.filetype,
				Height:        testCase.height,
				Width:         testCase.width,
			}
			assert.Equal(t, testCase.expected.String(), req.String())
		})
	}
}

func TestResponseString(t *testing.T) {
	var tests = []*TestResponse{
		{
			"ASCII",
			[]byte("image data"),
			"image/png",
			proto.GetThumbnailResponse{
				Thumbnail: []byte("image data"),
				Mimetype:  "image/png",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {
			response := proto.GetThumbnailResponse{
				Thumbnail: testCase.img,
				Mimetype:  testCase.mimetype,
			}

			assert.Equal(t, testCase.expected.String(), response.String())
		})
	}
}
