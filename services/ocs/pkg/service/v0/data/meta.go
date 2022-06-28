package data

// Meta holds response metadata
type Meta struct {
	Status       string `json:"status" xml:"status"`
	StatusCode   int    `json:"statuscode" xml:"statuscode"`
	Message      string `json:"message" xml:"message"`
	TotalItems   string `json:"totalitems,omitempty" xml:"totalitems,omitempty"`
	ItemsPerPage string `json:"itemsperpage,omitempty" xml:"itemsperpage,omitempty"`
}

// MetaOK is the default ok response with code 100
var MetaOK = Meta{Status: "ok", StatusCode: 100, Message: "OK"}

// MetaFailure is a failure response with code 101
var MetaFailure = Meta{Status: "", StatusCode: 101, Message: "Failure"}

// MetaInvalidInput is an error response with code 102
var MetaInvalidInput = Meta{Status: "", StatusCode: 102, Message: "Invalid Input"}

// MetaForbidden is an error response with code 104
var MetaForbidden = Meta{Status: "", StatusCode: 104, Message: "Forbidden"}

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

// MessageUserNotFound is  used when a user can not be found
var MessageUserNotFound = "The requested user could not be found"

// MessageGroupNotFound is used when a group can not be found
var MessageGroupNotFound = "The requested group could not be found"
