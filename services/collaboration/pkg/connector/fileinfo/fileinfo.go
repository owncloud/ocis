package fileinfo

// FileInfo contains the properties of the file.
// Some properties refer to capabilities in the WOPI client, and capabilities
// that the WOPI server has.
//
// Specific implementations must allow json-encoding of their relevant
// properties because the object will be marshalled directly
type FileInfo interface {
	// SetProperties will set the properties of this FileInfo.
	// Keys should match any valid property that the FileInfo implementation
	// has. If a key doesn't match any property, it must be ignored.
	// The values must have its matching type for the target property,
	// otherwise panics might happen.
	//
	// This method should help to reduce the friction of using different
	// implementations with different properties. You can use the same map
	// for all the implementations knowing that the relevant properties for
	// each implementation will be set.
	SetProperties(props map[string]interface{})

	// GetTarget will return the target implementation (OnlyOffice, Collabora...).
	// This will help to identify the implementation we're using in an easy way.
	// Note that the returned value must be unique among all the implementations
	GetTarget() string
}
