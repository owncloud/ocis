package data

// Users holds user ids for the user listing
type Users struct {
	Users []string `json:"users" xml:"users>element"`
}

// User holds the payload for a GetUser response
type User struct {
	Enabled           string `json:"enabled" xml:"enabled"`
	UserID            string `json:"id" xml:"id"` // UserID is mapped to the preferred_name attribute in accounts
	DisplayName       string `json:"display-name" xml:"display-name"`
	LegacyDisplayName string `json:"displayname" xml:"displayname"`
	Email             string `json:"email" xml:"email"`
	Quota             *Quota `json:"quota,omitempty" xml:"quota,omitempty"`
	UIDNumber         int64  `json:"uidnumber" xml:"uidnumber"`
	GIDNumber         int64  `json:"gidnumber" xml:"gidnumber"`
}

// Quota holds quota information
type Quota struct {
	Free       int64   `json:"free,omitempty" xml:"free,omitempty"`
	Used       int64   `json:"used,omitempty" xml:"used,omitempty"`
	Total      int64   `json:"total,omitempty" xml:"total,omitempty"`
	Relative   float32 `json:"relative,omitempty" xml:"relative,omitempty"`
	Definition string  `json:"definition,omitempty" xml:"definition,omitempty"`
}

// SigningKey holds the Payload for a GetSigningKey response
type SigningKey struct {
	User       string `json:"user" xml:"user"`
	SigningKey string `json:"signing-key" xml:"signing-key"`
}
