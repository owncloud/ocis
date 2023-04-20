package meshdirectory_web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var spaDist embed.FS

func ServeMeshDirectorySPA(w http.ResponseWriter, r *http.Request) {
    distFS, _ := fs.Sub(spaDist, "dist")
	server := http.FileServer(http.FS(distFS))
	server.ServeHTTP(w, r)
}
