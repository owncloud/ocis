package content

// Document wraps all resource meta fields,
// it is used as a content extraction result.
type Document struct {
	Title    string
	Name     string
	Content  string
	Size     uint64
	Mtime    string
	MimeType string
}
