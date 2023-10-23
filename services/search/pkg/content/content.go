package content

import (
	"strings"

	"github.com/bbalet/stopwords"
)

func init() {
	stopwords.OverwriteWordSegmenter(`[^ ]+`)
}

// Document wraps all resource meta fields,
// it is used as a content extraction result.
type Document struct {
	Title    string
	Name     string
	Content  string
	Size     uint64
	Mtime    string
	MimeType string
	Tags     []string
}

func CleanString(content, langCode string) string {
	return strings.TrimSpace(stopwords.CleanString(content, langCode, true))
}
