package assets

import (
	"bytes"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/net/html"
)

type fileServer struct {
	fsys http.FileSystem
}

// FileServer defines the http handler for the embedded files
func FileServer(fsys fs.FS) http.Handler {
	return &fileServer{http.FS(fsys)}
}

func (f *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := path.Clean(path.Join("/", r.URL.Path))

	serveIndex := func() {
		// not every fs contains a file named index.html,
		// therefore, we need to check if the file exists
		indexPath := "/index.html"
		file, err := f.fsys.Open(indexPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		s, err := file.Stat()
		if err != nil || s.IsDir() {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(s.Name())))
		w.Header().Add("Vary", "Accept-Encoding")

		buf := new(bytes.Buffer)
		w.Header().Del("Expires")
		w.Header().Set("Cache-Control", "no-cache")
		if err := withBase(buf, file, "/"); err != nil {
			http.NotFound(w, r)
			return
		}
		if _, err := w.Write(buf.Bytes()); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if !isValid(f, path) {
		serveIndex()
		return
	}

	asset, err := f.fsys.Open(path)
	if err != nil {
		serveIndex()
		return
	}
	defer asset.Close()

	s, err := asset.Stat()
	if err != nil {
		serveIndex()
		return
	}
	if s.IsDir() {
		serveIndex()
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(s.Name())))
	w.Header().Add("Vary", "Accept-Encoding")

	buf := new(bytes.Buffer)

	switch s.Name() {
	case "index.html", "oidc-callback.html", "oidc-silent-redirect.html":
		w.Header().Del("Expires")
		w.Header().Set("Cache-Control", "no-cache")
		err = withBase(buf, asset, "/")
		if err != nil {
			http.NotFound(w, r)
			return
		}
	default:
		_, err := buf.ReadFrom(asset)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func isValid(f *fileServer, path string) bool {
	if dir, ok := f.fsys.(http.Dir); ok {
		rootAbs, err := filepath.Abs(string(dir))
		if err != nil {
			return false
		}
		rel := path
		if len(rel) > 0 && rel[0] == '/' {
			rel = rel[1:]
		}
		candidate := filepath.Join(rootAbs, filepath.FromSlash(rel))
		fi, err := os.Lstat(candidate)
		if err != nil {
			return false
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return false
		}
	}
	return true
}

func withBase(w io.Writer, r io.Reader, base string) error {
	doc, _ := html.Parse(r)
	var parse func(*html.Node)
	parse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "head" {
			n.InsertBefore(&html.Node{
				Type: html.ElementNode,
				Data: "base",
				Attr: []html.Attribute{
					{
						Key: "href",
						Val: base,
					},
				},
			}, n.FirstChild)

			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parse(c)
		}
	}
	parse(doc)

	return html.Render(w, doc)
}
