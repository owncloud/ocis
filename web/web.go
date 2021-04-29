package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/*
//go:embed assets/js/*
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