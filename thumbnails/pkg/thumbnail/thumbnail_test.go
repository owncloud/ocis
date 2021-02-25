package thumbnail

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/storage"
)

type NoOpManager struct {
	storage.Storage
}

func (m NoOpManager) BuildKey(r storage.Request) string {
	return ""
}

func (m NoOpManager) Set(username, key string, thumbnail []byte) error {
	return nil
}

func BenchmarkGet(b *testing.B) {

	sut := NewSimpleManager(
		Resolutions{},
		NoOpManager{},
		log.NewLogger(),
	)

	res, _ := ParseResolution("32x32")
	req := Request{
		Resolution: res,
		ETag:       "1872ade88f3013edeb33decd74a4f947",
	}
	cwd, _ := os.Getwd()
	p := filepath.Join(cwd, "../../testdata/oc.png")
	f, _ := os.Open(p)
	defer f.Close()
	img, ext, _ := image.Decode(f)
	req.Encoder = EncoderForType(ext)
	for i := 0; i < b.N; i++ {
		_, _ = sut.Get(req, img)
	}
}
