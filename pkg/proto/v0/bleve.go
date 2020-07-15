package proto

// BleveRecord wraps the generated Record and adds a property that is used to distinguish documents in the index.
type BleveRecord struct {
	Record
	DatabaseTable string `json:"database_table"`
}
