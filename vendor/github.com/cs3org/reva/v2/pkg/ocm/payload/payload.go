package payload

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	// https://cs3org.github.io/OCM-API/docs.html?branch=develop&repo=OCM-API&user=cs3org#/paths/~1notifications/post
	// NotificationType one of "SHARE_ACCEPTED", "SHARE_DECLINED", "SHARE_CHANGE_PERMISSION", "SHARE_UNSHARED", "USER_REMOVED"
	SHARE_UNSHARED          = "SHARE_UNSHARED"
	SHARE_CHANGE_PERMISSION = "SHARE_CHANGE_PERMISSION"
)

// NotificationRequest is the request payload for the OCM API notifications endpoint.
// https://cs3org.github.io/OCM-API/docs.html?branch=develop&repo=OCM-API&user=cs3org#/paths/~1notifications/post
type NotificationRequest struct {
	NotificationType string `json:"notificationType" validate:"required"`
	ResourceType     string `json:"resourceType" validate:"required"`
	// Identifier to identify the shared resource at the provider side. This is unique per provider such that if the same resource is shared twice, this providerId will not be repeated.
	ProviderId string `json:"providerId" validate:"required"`
	// Optional additional parameters, depending on the notification and the resource type.
	Notification *Notification `json:"notification,omitempty"`
}

// Notification is the payload for the notification field in the NotificationRequest.
type Notification struct {
	// Owner        string `json:"owner,omitempty"`
	Grantee      string `json:"grantee,omitempty"`
	SharedSecret string `json:"sharedSecret,omitempty"`
}

// ToJSON returns the JSON io.Reader of the NotificationRequest.
func (r *NotificationRequest) ToJSON() (io.Reader, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(r); err != nil {
		return nil, err
	}
	return &b, nil
}

// ErrorMessageResponse is the response returned by the OCM API in case of an error.
// https://cs3org.github.io/OCM-API/docs.html?branch=develop&repo=OCM-API&user=cs3org#/paths/~1notifications/post
type ErrorMessageResponse struct {
	Message          string             `json:"message"`
	ValidationErrors []*ValidationError `json:"validationErrors,omitempty"`
}

// ValidationError is the payload for the validationErrors field in the ErrorMessageResponse.
type ValidationError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
