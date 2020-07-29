package data

// Users holds user ids for the user listing
type Users struct {
	Users []string `json:"users" xml:"users>element"`
}

// User holds the payload for a GetUser response
type User struct {
	// TODO needs better naming, clarify if we need a userid, a username or both
	Enabled     bool   `json:"enabled" xml:"enabled"`
	UserID      string `json:"id" xml:"id"`
	Username    string `json:"username" xml:"username"`
	DisplayName string `json:"displayname" xml:"displayname"`
	Email       string `json:"email" xml:"email"`
	Quota       *Quota `json:"quota" xml:"quota"`
}

// Quota holds quota information
type Quota struct {
	Free       int64   `json:"free" xml:"free"`
	Used       int64   `json:"used" xml:"used"`
	Total      int64   `json:"total" xml:"total"`
	Relative   float32 `json:"relative" xml:"relative"`
	Definition string  `json:"definition" xml:"definition"`
}

// SigningKey holds the Payload for a GetSigningKey response
type SigningKey struct {
	User       string `json:"user" xml:"user"`
	SigningKey string `json:"signing-key" xml:"signing-key"`
}
