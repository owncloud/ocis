package settings

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/*
var assets embed.FS

// Assets FS
var Assets http.FileSystem

func init() {
	embedFS, err := fs.Sub(assets, "assets")
	if err != nil {
		panic(err)
	}

	Assets = http.FS(embedFS)
}
