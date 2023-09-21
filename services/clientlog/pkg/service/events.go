package service

// FileReadyEvent is emitted when the postprocessing of a file is finished
type FileReadyEvent struct {
	ItemID string `json:"itemid"`
}
