package proto_test

import (
	"bytes"
	"context"
	"image"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/owncloud/ocis/thumbnails/pkg/proto/v0"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/imgsource"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/storage"
	"github.com/stretchr/testify/assert"

	svc "github.com/owncloud/ocis/thumbnails/pkg/service/v0"
)

var service = grpc.Service{}

func init() {
	service = grpc.NewService(
		grpc.Namespace("com.owncloud.api"),
		grpc.Name("thumbnails"),
		grpc.Address("localhost:9992"),
	)

	cfg := config.New()
	cfg.Thumbnail.Resolutions = []string{"16x16", "32x32", "64x64", "128x128"}

	wd, _ := os.Getwd()
	fsCfg := config.FileSystemSource{
		BasePath: filepath.Join(wd, "../../../testdata/"),
	}
	err := proto.RegisterThumbnailServiceHandler(
		service.Server(),
		svc.NewService(
			svc.Config(cfg),
			svc.ThumbnailStorage(storage.NewInMemoryStorage()),
			svc.ThumbnailSource(imgsource.NewFileSystemSource(fsCfg)),
		),
	)
	if err != nil {
		log.Fatalf("could not register ThumbnailHandler: %v", err)
	}
	 if err := service.Server().Start(); err != nil {
	 	log.Fatalf("could not start server: %v", err)
	 }
}

func TestGetThumbnailInvalidImage(t *testing.T) {
	req := proto.GetRequest{
		Filepath: "invalid.png",
		Filetype: proto.GetRequest_PNG,
		Etag:     "33a64df551425fcc55e4d42a148795d9f25f89d4",
		Height:   32,
		Width:    32,
		Username: "user1",
	}
	client := service.Client()
	cl := proto.NewThumbnailService("com.owncloud.api.thumbnails", client)
	_, err := cl.GetThumbnail(context.Background(), &req)

	assert.NotNil(t, err)
}

func TestGetThumbnail(t *testing.T) {
	req := proto.GetRequest{
		Filepath:      "oc.png",
		Filetype:      proto.GetRequest_PNG,
		Etag:          "33a64df551425fcc55e4d42a148795d9f25f89d4",
		Height:        32,
		Width:         32,
		Authorization: "Bearer eyJhbGciOiJQUzI1NiIsImtpZCI6IiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJwaG9lbml4IiwiZXhwIjoxNTkwNTc1Mzk4LCJqdGkiOiJqUEw5c1A3UUEzY0diYi1yRnhkSjJCWnFPc1BDTDg1ZyIsImlhdCI6MTU5MDU3NDc5OCwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6OTIwMCIsInN1YiI6Ilh0U2lfbWl5V1NCLXBrdkdueFBvQzVBNGZsaWgwVUNMZ3ZVN2NMd2ptakNLWDdGWW4ySFdrNnJSQ0V1eTJHNXFBeV95TVFjX0ZLOWFORmhVTXJYMnBRQGtvbm5lY3QiLCJrYy5pc0FjY2Vzc1Rva2VuIjp0cnVlLCJrYy5hdXRob3JpemVkU2NvcGVzIjpbIm9wZW5pZCIsInByb2ZpbGUiLCJlbWFpbCJdLCJrYy5pZGVudGl0eSI6eyJrYy5pLmRuIjoiRWluc3RlaW4iLCJrYy5pLmlkIjoiY249ZWluc3RlaW4sb3U9dXNlcnMsZGM9ZXhhbXBsZSxkYz1vcmciLCJrYy5pLnVuIjoiZWluc3RlaW4ifSwia2MucHJvdmlkZXIiOiJpZGVudGlmaWVyLWxkYXAifQ.FSDe4vzwYpHbNfckBON5EI-01MS_dYFxenddqfJPzjlAEMEH2FFn2xQHCsxhC7wSxivhjV7Z5eRoNUR606keA64Tjs8pJBNECSptBMmE_xfAlc6X5IFILgDnR5bBu6Z2hhu-dVj72Hcyvo_X__OeWekYu7oyoXW41Mw3ayiUAwjCAzV3WPOAJ_r0zbW68_m29BgH3BoSxaF6lmjStIIAIyw7IBZ2QXb_FvGouknmfeWlGL9lkFPGL_dYKwjWieG947nY4Kg8IvHByEbw-xlY3L2EdA7Q8ZMbqdX7GzjtEIVYvCT4-TxWRcmB3SmO-Z8CVq27NHlKm3aZ0k2PS8Ga1w",
		Username:      "user1",
	}
	client := service.Client()
	cl := proto.NewThumbnailService("com.owncloud.api.thumbnails", client)
	rsp, err := cl.GetThumbnail(context.Background(), &req)
	if err != nil {
		log.Fatalf("error %s", err.Error())
	}
	assert.NotEmpty(t, rsp.GetThumbnail())

	img, _, _ := image.Decode(bytes.NewReader(rsp.GetThumbnail()))
	assert.Equal(t, 32, img.Bounds().Size().X)

	assert.Equal(t, "image/png", rsp.GetMimetype())
}
