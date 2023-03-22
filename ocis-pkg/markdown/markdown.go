// Package markdown allows reading and editing Markdown files
package markdown

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Heading represents a markdown Heading
type Heading struct {
	Level   int    // the level of the heading. 1 means it's the H1
	Content string // the text of the heading
	Header  string // the heading itself
}

// MD represents a markdown file
type MD struct {
	Headings []Heading
}

// Bytes returns the markdown as []bytes, ignoring errors
func (md MD) Bytes() []byte {
	var b bytes.Buffer
	_, _ = md.WriteContent(&b)
	return b.Bytes()
}

// String returns the markdown as string, ignoring errors
func (md MD) String() string {
	var b strings.Builder
	_, _ = md.WriteContent(&b)
	return b.String()
}

// TocBytes returns the table of contents as []byte, ignoring errors
func (md MD) TocBytes() []byte {
	var b bytes.Buffer
	_, _ = md.WriteToc(&b)
	return b.Bytes()
}

// TocString returns the table of contents as string, ignoring errors
func (md MD) TocString() string {
	var b strings.Builder
	_, _ = md.WriteToc(&b)
	return b.String()
}

// WriteContent writes the MDs content to the given writer
func (md MD) WriteContent(w io.Writer) (int64, error) {
	max, written := len(md.Headings), int64(0)
	write := func(s string) error {
		n, err := w.Write([]byte(s))
		written += int64(n)
		return err
	}
	for i, h := range md.Headings {
		if err := write(strings.Repeat("#", h.Level) + " " + h.Header + "\n"); err != nil {
			return written, err
		}
		if err := write("\n"); err != nil {
			return written, err
		}
		if len(h.Content) > 0 {
			if err := write(h.Content); err != nil {
				return written, err
			}
			if i < max-1 {
				if err := write("\n"); err != nil {
					return written, err
				}
			}
		}
	}
	return written, nil
}

// WriteToc writes the table of contents to the given writer
func (md MD) WriteToc(w io.Writer) (int64, error) {
	var written int64
	for _, h := range md.Headings {
		if h.Level == 1 {
			// main title not in toc
			continue
		}
		link := fmt.Sprintf("#%s", strings.ToLower(strings.Replace(h.Header, " ", "-", -1)))
		s := fmt.Sprintf("%s* [%s](%s)\n", strings.Repeat("  ", h.Level-2), h.Header, link)
		n, err := w.Write([]byte(s))
		if err != nil {
			return written, err
		}
		written += int64(n)
	}
	return written, nil
}

// NewMD parses a new Markdown
func NewMD(b []byte) MD {
	var (
		md      MD
		heading Heading
		content strings.Builder
	)
	sendHeading := func() {
		if heading.Header != "" {
			heading.Content = content.String()
			md.Headings = append(md.Headings, heading)
			content = strings.Builder{}
		}
	}
	parts := strings.Split(string(b), "\n")
	for _, p := range parts {
		if p == "" {
			continue
		}
		if p[:1] == "#" { // this is a header
			sendHeading()
			heading = headingFromString(p)
		} else {
			// readd lost "\n"
			_, _ = content.WriteString(p + "\n")
		}
	}
	sendHeading()
	return md
}

func headingFromString(s string) Heading {
	i := strings.LastIndex(s, "#")
	levs, con := s[:i+1], s[i+1:]
	return Heading{
		Level:  len(levs),
		Header: strings.TrimPrefix(con, " "),
	}
}
