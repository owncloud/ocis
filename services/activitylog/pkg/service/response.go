package service

import (
	"time"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
)

// Translations
var (
	MessageResourceCreated   = l10n.Template("{user} created {resource}")
	MessageResourceTrashed   = l10n.Template("{user} trashed {resource}")
	MessageResourcePurged    = l10n.Template("{user} purged {resource}")
	MessageResourceMoved     = l10n.Template("{user} moved {resource}")
	MessageShareCreated      = l10n.Template("{user} shared {resource}")
	MessageShareUpdated      = l10n.Template("{user} updated share of {resource}")
	MessageShareDeleted      = l10n.Template("{user} deleted share of {resource}")
	MessageLinkCreated       = l10n.Template("{user} created link to {resource}")
	MessageLinkUpdated       = l10n.Template("{user} updated link to {resource}")
	MessageLinkDeleted       = l10n.Template("{user} deleted link to {resource}")
	MessageSpaceShared       = l10n.Template("{user} shared space {resource}")
	MessageSpaceShareUpdated = l10n.Template("{user} updated share of space {resource}")
	MessageSpaceUnshared     = l10n.Template("{user} unshared space {resource}")
)

// GetActivitiesResponse is the response on GET activities requests
type GetActivitiesResponse struct {
	Activities []Activity `json:"value"`
}

// Activity represents an activity as it is returned to the client
type Activity struct {
	ID       string   `json:"id"`
	Times    Times    `json:"times"`
	Template Template `json:"template"`
}

// Resource represents an item such as a file or folder
type Resource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Actor represents the user who performed the Action
type Actor struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// Times represents the timestamps of the Activity
type Times struct {
	RecordedTime time.Time `json:"recordedTime"`
}

// Template contains activity details
type Template struct {
	Message   string                 `json:"message"`
	Variables map[string]interface{} `json:"variables"`
}

// NewActivity creates a new activity
func NewActivity(message string, res Resource, user Actor, ts Times, eventID string) Activity {
	return Activity{
		ID:    eventID,
		Times: ts,
		Template: Template{
			Message: message,
			Variables: map[string]interface{}{
				"resource": res,
				"user":     user,
			},
		},
	}
}

// ResponseData returns the relevant response data for the activity
func (s *ActivitylogService) ResponseData(ref *provider.Reference, uid *user.UserId, username string, ts time.Time) (Resource, Actor, Times, error) {
	gwc, err := s.gws.Next()
	if err != nil {
		return Resource{}, Actor{}, Times{}, err
	}

	ctx, err := utils.GetServiceUserContext(s.cfg.ServiceAccount.ServiceAccountID, gwc, s.cfg.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		return Resource{}, Actor{}, Times{}, err
	}

	info, err := utils.GetResource(ctx, ref, gwc)
	if err != nil {
		return Resource{}, Actor{}, Times{}, err
	}

	if username == "" {
		u, err := utils.GetUser(uid, gwc)
		if err != nil {
			return Resource{}, Actor{}, Times{}, err
		}
		username = u.GetUsername()
	}

	return Resource{
			ID:   storagespace.FormatResourceID(*info.Id),
			Name: info.Path,
		}, Actor{
			ID:          uid.GetOpaqueId(),
			DisplayName: username,
		}, Times{
			RecordedTime: ts,
		}, nil

}
