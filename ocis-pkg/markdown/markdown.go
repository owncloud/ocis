// Package markdown allows reading and editing Markdown files
package markdown

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"
)

// Heading represents a markdown Heading
type Heading struct {
	Content string
	Level   int
	Header  string
}

// MD represents a markdown file
type MD struct {
	Headings []Heading
}

// Bytes returns the markdown as []bytes to be written to a file
func (md MD) Bytes() []byte {
	b, num := bytes.NewBuffer(nil), len(md.Headings)
	for i, h := range md.Headings {
		b.Write([]byte(strings.Repeat("#", h.Level) + " " + h.Header + "\n"))
		b.Write([]byte("\n"))
		if len(h.Content) > 0 {
			b.Write([]byte(h.Content))
			if i < num-1 {
				b.Write([]byte("\n"))
			}
		}
	}
	return b.Bytes()
}

// Toc returns the table of contents as []byte
func (md MD) Toc() []byte {
	b := bytes.NewBuffer(nil)
	for _, h := range md.Headings {
		if h.Level == 1 {
			// main title not in toc
			continue
		}
		link := fmt.Sprintf("#%s", strings.ToLower(strings.Replace(h.Header, " ", "-", -1)))
		s := fmt.Sprintf("%s* [%s](%s)\n", strings.Repeat("  ", h.Level-2), h.Header, link)
		b.Write([]byte(s))
	}
	return b.Bytes()
}

// NewMD parses a new Markdown
func NewMD(b []byte) MD {
	md := MD{}
	var heading Heading
	parts := strings.Split(string(b), "\n")
	for _, p := range parts {
		if p == "" {
			continue
		}
		if p[:1] == "#" { // this is a header
			if heading.Header != "" {
				md.Headings = append(md.Headings, heading)
			}
			heading = Heading{}
			i := strings.LastIndex(p, "#")
			levs, con := p[:i+1], p[i+1:]
			heading.Header = strings.TrimPrefix(con, " ")
			heading.Level = len(levs)
		} else {
			heading.Content += p + "\n"
		}
	}
	if heading.Header != "" {
		md.Headings = append(md.Headings, heading)
	}

	return md
}

func main() {
	f, err := os.ReadFile("/home/jkoberg/ocis/services/antivirus/README.md")
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}

	md := NewMD(f)
	head := md.Headings[0]
	md.Headings = md.Headings[1:]

	tpl := template.Must(template.ParseFiles("index.tmpl"))
	b := bytes.NewBuffer(nil)
	err = tpl.Execute(b, map[string]interface{}{
		"ServiceName":  head.Header,
		"CreationTime": time.Now().Format(time.RFC3339Nano),
		"service":      "unknown",
		"Abstract":     head.Content,
		"TocTree":      string(md.Toc()),
		"Content":      string(md.Bytes()),
	})
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	err = os.WriteFile("test.md", b.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
}
