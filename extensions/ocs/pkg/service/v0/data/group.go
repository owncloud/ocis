package data

// Groups holds group ids for the groups listing
type Groups struct {
	Groups []string `json:"groups" xml:"groups>element"`
}
