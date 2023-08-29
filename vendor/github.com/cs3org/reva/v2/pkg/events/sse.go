package events

import (
	"encoding/json"
)

// SendSEE instructs the sse service to send a notification to a user
type SendSSE struct {
	UserID  string
	Type    string
	Message []byte
}

// Unmarshal to fulfill umarshaller interface
func (SendSSE) Unmarshal(v []byte) (interface{}, error) {
	e := SendSSE{}
	err := json.Unmarshal(v, &e)
	return e, err
}
