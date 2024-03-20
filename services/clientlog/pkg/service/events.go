package service

// FileEvent is emitted when a file is uploaded/renamed/deleted/...
type FileEvent struct {
	ParentItemID string `json:"parentitemid"`
	ItemID       string `json:"itemid"`
	SpaceID      string `json:"spaceid"`
	InitiatorID  string `json:"initiatorid"`
	Etag         string `json:"etag"`
}
