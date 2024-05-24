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

// assignStringTo will return a function whose parameter will be assigned
// to the provided key. The function will panic if the assignment isn't
// possible.
//
//	fn := AssignStringTo(&target)
//	fn(value)
//
// Is roughly equivalent to
//
//	target = value
//
// The reason for this method is to help the `SetProperties` method in order
// to provide a setter function for each property.
// Expected code for the `SetProperties` should be similar to
//
//	setters := map[string]func(value interface{}) {
//	  "Owner": AssignStringTo(&info.Owner),
//	  "DisplayName": AssignStringTo(&info.DisplayName),
//	  .....
//	}
//	for key, value := range props {
//	  fn := setters[key]
//	  fn(value)
//	}
//
// Further `assign*To` functions will be provided to be able to assign
// different data types
func assignStringTo(targetKey *string) func(value interface{}) {
	return func(value interface{}) {
		*targetKey = value.(string)
	}
}

// assignStringListTo will return a function whose parameter will be assigned
// to the provided key. The function will panic if the assignment isn't
// possible.
//
// See assignStringTo for more information
func assignStringListTo(targetKey *[]string) func(value interface{}) {
	return func(value interface{}) {
		*targetKey = value.([]string)
	}
}

// assignInt64To will return a function whose parameter will be assigned
// to the provided key. The function will panic if the assignment isn't
// possible.
//
// See assignStringTo for more information
func assignInt64To(targetKey *int64) func(value interface{}) {
	return func(value interface{}) {
		*targetKey = value.(int64)
	}
}

// assignIntTo will return a function whose parameter will be assigned
// to the provided key. The function will panic if the assignment isn't
// possible.
//
// See assignStringTo for more information
func assignIntTo(targetKey *int) func(value interface{}) {
	return func(value interface{}) {
		*targetKey = value.(int)
	}
}

// assignBoolTo will return a function whose parameter will be assigned
// to the provided key. The function will panic if the assignment isn't
// possible.
//
// See assignStringTo for more information
func assignBoolTo(targetKey *bool) func(value interface{}) {
	return func(value interface{}) {
		*targetKey = value.(bool)
	}
}
