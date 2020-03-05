package storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const BasePath = "/home/corby/tmp/thumbnails/fs/"

type FileSystem struct {
}

func (s FileSystem) Get(key string) []byte {
	content, err := ioutil.ReadFile(BasePath + key)
	if err != nil {
		return nil
	}

	return content
}

func (s FileSystem) Set(key string, img []byte) error {
	folder := filepath.Dir(BasePath + key)
	if err := createFolderIfNotExists(folder); err != nil {
		return fmt.Errorf("error while creating folder %s", folder)
	}

	f, err := os.Create(BasePath + key)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer f.Close()
	_, err = f.Write(img)
	if err != nil {
		return err
	}
	return nil
}

func (s FileSystem) BuildKey(ctx StorageContext) string {
	etag := ctx.ETag
	filetype := ctx.Types[0]
	filename := strconv.Itoa(ctx.Width) + "x" + strconv.Itoa(ctx.Height) + "." + filetype

	key := new(bytes.Buffer)
	key.WriteString(etag[:2])
	key.WriteRune('/')
	key.WriteString(etag[2:4])
	key.WriteRune('/')
	key.WriteString(etag[4:])
	key.WriteRune('/')
	key.WriteString(filename)

	return key.String()
}

func createFolderIfNotExists(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err := os.MkdirAll(folder, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}
