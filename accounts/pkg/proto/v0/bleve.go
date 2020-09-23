package proto

// BleveAccount wraps the generated Account and adds a bleve type that is used to distinguish documents in the index
type BleveAccount struct {
	Account
	BleveType string `json:"bleve_type"`
}

// BleveGroup wraps the generated Group and adds a bleve type that is used to distinguish documents in the index
type BleveGroup struct {
	Group
	BleveType string `json:"bleve_type"`
}
