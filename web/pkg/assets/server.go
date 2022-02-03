package assets

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

type fileServer struct {
	root http.FileSystem
}

func FileServer(root http.FileSystem) http.Handler {
	return &fileServer{root}
}

func (f *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	upath = path.Clean(upath)

	asset, err := f.root.Open(upath)
	if err != nil {
		r.URL.Path = "/index.html"
		f.ServeHTTP(w, r)
		return
	}
	defer asset.Close()

	s, _ := asset.Stat()
	if s.IsDir() {
		r.URL.Path = "/index.html"
		f.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(s.Name())))

	buf := new(bytes.Buffer)

	switch s.Name() {
	case "index.html", "oidc-callback.html", "oidc-silent-redirect.html":
		_ = withBase(buf, asset, "/")
	default:
		_, _ = buf.ReadFrom(asset)
	}

	w.Write(buf.Bytes())
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
					html.Attribute{
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
