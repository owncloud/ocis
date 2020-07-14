package data

// User holds the payload for a GetUser response
type User struct {
	// TODO needs better naming, clarify if we need a userid, a username or both
	ID          string `json:"id" xml:"id"`
	DisplayName string `json:"display-name" xml:"display-name"`
	Email       string `json:"email" xml:"email"`
}

// SigningKey holds the Payload for a GetSigningKey response
type SigningKey struct {
	User       string `json:"user" xml:"user"`
	SigningKey string `json:"signing-key" xml:"signing-key"`
}
