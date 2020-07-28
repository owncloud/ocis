package data

// User holds the payload for a GetUser response
type User struct {
	// TODO needs better naming, clarify if we need a userid, a username or both
	UserID      string `json:"userid" xml:"userid"`
	Username    string `json:"username" xml:"username"`
	DisplayName string `json:"displayname" xml:"displayname"`
	Email       string `json:"email" xml:"email"`
	Enabled     bool   `json:"enabled" xml:"enabled"`
}

// SigningKey holds the Payload for a GetSigningKey response
type SigningKey struct {
	User       string `json:"user" xml:"user"`
	SigningKey string `json:"signing-key" xml:"signing-key"`
}
