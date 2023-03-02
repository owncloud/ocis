package service

import (
	"bytes"
	"context"
	"errors"
	"text/template"
	"time"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
)

var (
	_resourceTypeSpace = "storagespace"
)

// OC10Notification is the oc10 style representation of an event
// some fields are left out for simplicity
type OC10Notification struct {
	EventID        string                 `json:"notification_id"`
	Service        string                 `json:"app"`
	UserName       string                 `json:"user"`
	Timestamp      string                 `json:"datetime"`
	ResourceID     string                 `json:"object_id"`
	ResourceType   string                 `json:"object_type"`
	Subject        string                 `json:"subject"`
	SubjectRaw     string                 `json:"subjectRich"`
	Message        string                 `json:"message"`
	MessageRaw     string                 `json:"messageRich"`
	MessageDetails map[string]interface{} `json:"messageRichParameters"`
}

// SpaceDisabled converts a SpaceDisabled event to an OC10Notification
func (ul *UserlogService) SpaceDisabled(ctx context.Context, eventid string, ev events.SpaceDisabled) (OC10Notification, error) {
	user, err := ul.getUser(ctx, ev.Executant)
	if err != nil {
		return OC10Notification{}, err
	}

	space, err := ul.getSpace(ul.impersonate(user.GetId()), ev.ID.GetOpaqueId())
	if err != nil {
		return OC10Notification{}, err
	}

	subj, subjraw, msg, msgraw, err := ul.composeMessage(SpaceDisabled, map[string]string{
		"username":  user.GetDisplayName(),
		"spacename": space.GetName(),
	})
	if err != nil {
		return OC10Notification{}, err
	}

	return OC10Notification{
		EventID:        eventid,
		Service:        ul.cfg.Service.Name,
		UserName:       user.GetUsername(),
		Timestamp:      time.Now().Format(time.RFC3339Nano),
		ResourceID:     ev.ID.GetOpaqueId(),
		ResourceType:   _resourceTypeSpace,
		Subject:        subj,
		SubjectRaw:     subjraw,
		Message:        msg,
		MessageRaw:     msgraw,
		MessageDetails: ul.getDetails(user, space, nil),
	}, nil

}

func (ul *UserlogService) composeMessage(eventname string, vars map[string]string) (string, string, string, string, error) {
	tpl, ok := _templates[eventname]
	if !ok {
		return "", "", "", "", errors.New("unknown template name")
	}

	subject := ul.executeTemplate(tpl.Subject, vars)

	subjectraw := ul.executeTemplate(tpl.Subject, map[string]string{
		"username":  "{user}",
		"spacename": "{space}",
		"resource":  "{resource}",
	})

	message := ul.executeTemplate(tpl.Message, vars)

	messageraw := ul.executeTemplate(tpl.Message, map[string]string{
		"username":  "{user}",
		"spacename": "{space}",
		"resource":  "{resource}",
	})

	return subject, subjectraw, message, messageraw, nil

}

func (ul *UserlogService) getDetails(user *user.User, space *storageprovider.StorageSpace, item *storageprovider.ResourceInfo) map[string]interface{} {
	details := make(map[string]interface{})

	if user != nil {
		details["user"] = map[string]string{
			"id":          user.GetId().GetOpaqueId(),
			"name":        user.GetUsername(),
			"displayname": user.GetDisplayName(),
		}
	}

	if space != nil {
		details["space"] = map[string]string{
			"id":   space.GetId().GetOpaqueId(),
			"name": space.GetName(),
		}
	}

	if item != nil {
		details["resource"] = map[string]string{
			"id":   storagespace.FormatResourceID(*item.GetId()),
			"name": item.GetName(),
		}
	}

	return details
}

func (ul *UserlogService) executeTemplate(tpl *template.Template, vars map[string]string) string {
	var writer bytes.Buffer
	if err := tpl.Execute(&writer, vars); err != nil {
		ul.log.Error().Err(err).Str("templateName", tpl.Name()).Msg("cannot execute template")
		return ""
	}

	return writer.String()
}
