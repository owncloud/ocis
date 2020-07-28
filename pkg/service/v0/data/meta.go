package data

// Meta holds response metadata
type Meta struct {
	Status       string `json:"status" xml:"status"`
	StatusCode   int    `json:"statuscode" xml:"statuscode"`
	Message      string `json:"message" xml:"message"`
	TotalItems   string `json:"totalitems,omitempty" xml:"totalitems,omitempty"`
	ItemsPerPage string `json:"itemsperpage,omitempty" xml:"itemsperpage,omitempty"`
}

// MetaOK is the default ok response
var MetaOK = Meta{Status: "ok", StatusCode: 100, Message: "OK"}

// MetaBadRequest is used for unknown errors
var MetaBadRequest = Meta{Status: "error", StatusCode: 400, Message: "Bad Request"}

// MetaServerError is returned on server errors
var MetaServerError = Meta{Status: "error", StatusCode: 996, Message: "Server Error"}

// MetaUnauthorized is returned on unauthorized requests
var MetaUnauthorized = Meta{Status: "error", StatusCode: 997, Message: "Unauthorised"}

// MetaNotFound is returned when trying to access not existing resources
var MetaNotFound = Meta{Status: "error", StatusCode: 998, Message: "Not Found"}

// MetaUnknownError is used for unknown errors
var MetaUnknownError = Meta{Status: "error", StatusCode: 999, Message: "Unknown Error"}
