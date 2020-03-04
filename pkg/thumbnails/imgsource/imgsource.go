package imgsource

import "image"

// Source defines the interface for image sources
type Source interface {
	Get(path string, ctx SourceContext) (image.Image, error)
}

// NewContext creates a new SourceContext instance
func NewContext() SourceContext {
	return SourceContext{
		m: make(map[string]interface{}),
	}
}

// SourceContext is used to pass source specific parameters
type SourceContext struct {
	m map[string]interface{}
}

// GetString tries to cast the value to a string
func (s SourceContext) GetString(key string) string {
	if s, ok := s.m[key].(string); ok {
		return s
	}
	return ""
}

// Set sets a value
func (s SourceContext) Set(key string, val interface{}) {
	s.m[key] = val
}
