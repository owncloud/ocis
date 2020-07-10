package proto

// Bleve uses a private bleveClassifier interface to determine the type of a struct
// see https://github.com/blevesearch/bleve/blob/master/mapping/mapping.go#L32-L38

type BleveAccount struct {
	Account
	BleveType string `json:"bleve_type"`
}

type BleveGroup struct {
	Group
	BleveType string `json:"bleve_type"`
}
