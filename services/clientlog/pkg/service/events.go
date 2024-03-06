package service

// FileEvent is emitted when a file is uploaded/renamed/deleted/...
type FileEvent struct {
	ParentItemID string `json:"parentitemid"`
	ItemID       string `json:"itemid"`
}
