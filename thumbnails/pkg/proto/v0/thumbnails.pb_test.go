package proto_test

import (
	"testing"

	"github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	testDataName  string
	filepath      string
	filetype      proto.GetRequest_FileType
	etag          string
	width         int32
	height        int32
	authorization string
	expected      string
}

type TestResponse struct {
	testDataName string
	img          []byte
	mimetype     string
	expected     string
}

func TestRequestString(t *testing.T) {

	var tests = []TestRequest{
		{
			"ASCII",
			"Foo.jpg",
			proto.GetRequest_JPG,
			"33a64df551425fcc55e4d42a148795d9f25f89d4",
			24,
			24,
			"Basic SGVXaG9SZWFkc1RoaXM6SXNTdHVwaWQK",
			`filepath:"Foo.jpg" filetype:JPG etag:"33a64df551425fcc55e4d42a148795d9f25f89d4" width:24 height:24 authorization:"Basic SGVXaG9SZWFkc1RoaXM6SXNTdHVwaWQK" username:"user1"`,
		},
		{
			"UTF",
			"मिलन.jpg",
			proto.GetRequest_JPG,
			"33a64df551425fcc55e4d42a148795d9f25f89d4",
			24,
			24,
			"Basic SGVXaG9SZWFkc1RoaXM6SXNTdHVwaWQK",
			`filepath:"मिलन.jpg" filetype:JPG etag:"33a64df551425fcc55e4d42a148795d9f25f89d4" width:24 height:24 authorization:"Basic SGVXaG9SZWFkc1RoaXM6SXNTdHVwaWQK" username:"user1"`,
		},
		{
			"PNG",
			"Foo.png",
			proto.GetRequest_PNG,
			"33a64df551425fcc55e4d42a148795d9f25f89d4",
			24,
			24,
			"Basic SGVXaG9SZWFkc1RoaXM6SXNTdHVwaWQK",
			`filepath:"Foo.png" etag:"33a64df551425fcc55e4d42a148795d9f25f89d4" width:24 height:24 authorization:"Basic SGVXaG9SZWFkc1RoaXM6SXNTdHVwaWQK" username:"user1"`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {
			req := proto.GetRequest{
				Filepath:      testCase.filepath,
				Filetype:      testCase.filetype,
				Etag:          testCase.etag,
				Height:        testCase.height,
				Width:         testCase.width,
				Authorization: testCase.authorization,
				Username:      "user1",
			}
			assert.Equal(t, testCase.expected, req.String())
		})
	}
}

func TestResponseString(t *testing.T) {
	var tests = []TestResponse{
		{
			"ASCII",
			[]byte("image data"),
			"image/png",
			`thumbnail:"image data" mimetype:"image/png"`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testDataName, func(t *testing.T) {
			response := proto.GetResponse{
				Thumbnail: testCase.img,
				Mimetype:  testCase.mimetype,
			}

			assert.Equal(t, testCase.expected, response.String())
		})
	}
}
