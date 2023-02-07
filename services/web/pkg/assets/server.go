package assets

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"

	"golang.org/x/net/html"
)

type fileServer struct {
	root http.FileSystem
}

func FileServer(root http.FileSystem) http.Handler {
	return &fileServer{root}
}

func (f *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := path.Clean(path.Join("/", r.URL.Path))
	r.URL.Path = upath

	fallbackIndex := func() {
		r.URL.Path = "/index.html"
		f.ServeHTTP(w, r)
	}

	asset, err := f.root.Open(upath)
	if err != nil {
		fallbackIndex()
		return
	}
	defer asset.Close()

	s, _ := asset.Stat()
	if s.IsDir() {
		fallbackIndex()
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(s.Name())))

	buf := new(bytes.Buffer)

	switch s.Name() {
	case "index.html", "oidc-callback.html", "oidc-silent-redirect.html":
		w.Header().Del("Expires")
		w.Header().Set("Cache-Control", "no-cache")
		_ = withBase(buf, asset, "/")
	default:
		_, _ = buf.ReadFrom(asset)
	}

	_, _ = w.Write(buf.Bytes())
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
