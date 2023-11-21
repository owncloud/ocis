package service

// FileReadyEvent is emitted when the postprocessing of a file is finished
type FileReadyEvent struct {
	ParentItemID string `json:"parentitemid"`
	ItemID       string `json:"itemid"`
}
