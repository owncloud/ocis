package event

import (
	"encoding/json"
	"time"
)

// PurgeTrashBin wraps all needed information to purge a trash-bin
type PurgeTrashBin struct {
	ExecutantID  string
	RemoveBefore time.Time
}

// Unmarshal to fulfill umarshaller interface
func (PurgeTrashBin) Unmarshal(v []byte) (interface{}, error) {
	e := PurgeTrashBin{}
	err := json.Unmarshal(v, &e)
	return e, err
}
