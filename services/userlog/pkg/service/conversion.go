package service

import (
	"bytes"
	"errors"
	"text/template"
	"time"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/leonelquinteros/gotext"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
)

var (
	_resourceTypeSpace = "storagespace"
	_resourceTypeShare = "share"

	// TODO: from config
	_pathToLocales = "/home/jkoberg/ocis/services/userlog/pkg/service/locales"
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

// ConvertEvent converts an eventhistory event to an OC10Notification
func (ul *UserlogService) ConvertEvent(event *ehmsg.Event, locale string) (OC10Notification, error) {
	etype, ok := ul.registeredEvents[event.Type]
	if !ok {
		// this should not happen
		return OC10Notification{}, errors.New("eventtype not registered")
	}

	einterface, err := etype.Unmarshal(event.Event)
	if err != nil {
		// this shouldn't happen either
		return OC10Notification{}, errors.New("cant unmarshal event")
	}

	switch ev := einterface.(type) {
	default:
		return OC10Notification{}, errors.New("unknown event type")
	// space related
	case events.SpaceDisabled:
		return ul.spaceMessage(event.Id, SpaceDisabled, ev.Executant, ev.ID.GetOpaqueId(), ev.Timestamp, locale)
	case events.SpaceDeleted:
		return ul.spaceDeletedMessage(event.Id, ev.Executant, ev.ID.GetOpaqueId(), ev.SpaceName, ev.Timestamp, locale)
	case events.SpaceShared:
		return ul.spaceMessage(event.Id, SpaceShared, ev.Executant, ev.ID.GetOpaqueId(), ev.Timestamp, locale)
	case events.SpaceUnshared:
		return ul.spaceMessage(event.Id, SpaceUnshared, ev.Executant, ev.ID.GetOpaqueId(), ev.Timestamp, locale)
	case events.SpaceMembershipExpired:
		return ul.spaceMessage(event.Id, SpaceMembershipExpired, ev.SpaceOwner, ev.SpaceID.GetOpaqueId(), ev.ExpiredAt, locale)

	// share related
	case events.ShareCreated:
		return ul.shareMessage(event.Id, ShareCreated, ev.Executant, ev.ItemID, ev.ShareID, utils.TSToTime(ev.CTime), locale)
	case events.ShareExpired:
		return ul.shareMessage(event.Id, ShareExpired, ev.ShareOwner, ev.ItemID, ev.ShareID, ev.ExpiredAt, locale)
	case events.ShareRemoved:
		return ul.shareMessage(event.Id, ShareRemoved, ev.Executant, ev.ItemID, ev.ShareID, ev.Timestamp, locale)
	}
}

func (ul *UserlogService) spaceDeletedMessage(eventid string, executant *user.UserId, spaceid string, spacename string, ts time.Time, locale string) (OC10Notification, error) {
	_, user, err := utils.Impersonate(executant, ul.gwClient, ul.cfg.MachineAuthAPIKey)
	if err != nil {
		return OC10Notification{}, err
	}

	subj, subjraw, msg, msgraw, err := ul.composeMessage(SpaceDeleted, locale, map[string]interface{}{
		"username":  user.GetDisplayName(),
		"spacename": spacename,
	})
	if err != nil {
		return OC10Notification{}, err
	}

	details := ul.getDetails(user, nil, nil, nil)
	details["space"] = map[string]string{
		"id":   spaceid,
		"name": spacename,
	}

	return OC10Notification{
		EventID:        eventid,
		Service:        ul.cfg.Service.Name,
		UserName:       user.GetUsername(),
		Timestamp:      ts.Format(time.RFC3339Nano),
		ResourceID:     spaceid,
		ResourceType:   _resourceTypeSpace,
		Subject:        subj,
		SubjectRaw:     subjraw,
		Message:        msg,
		MessageRaw:     msgraw,
		MessageDetails: details,
	}, nil
}

func (ul *UserlogService) spaceMessage(eventid string, nt NotificationTemplate, executant *user.UserId, spaceid string, ts time.Time, locale string) (OC10Notification, error) {
	ctx, user, err := utils.Impersonate(executant, ul.gwClient, ul.cfg.MachineAuthAPIKey)
	if err != nil {
		return OC10Notification{}, err
	}

	space, err := ul.getSpace(ctx, spaceid)
	if err != nil {
		return OC10Notification{}, err
	}

	subj, subjraw, msg, msgraw, err := ul.composeMessage(nt, locale, map[string]interface{}{
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
		Timestamp:      ts.Format(time.RFC3339Nano),
		ResourceID:     spaceid,
		ResourceType:   _resourceTypeSpace,
		Subject:        subj,
		SubjectRaw:     subjraw,
		Message:        msg,
		MessageRaw:     msgraw,
		MessageDetails: ul.getDetails(user, space, nil, nil),
	}, nil
}

func (ul *UserlogService) shareMessage(eventid string, nt NotificationTemplate, executant *user.UserId, resourceid *storageprovider.ResourceId, shareid *collaboration.ShareId, ts time.Time, locale string) (OC10Notification, error) {
	ctx, user, err := utils.Impersonate(executant, ul.gwClient, ul.cfg.MachineAuthAPIKey)
	if err != nil {
		return OC10Notification{}, err
	}

	info, err := ul.getResource(ctx, resourceid)
	if err != nil {
		return OC10Notification{}, err
	}

	subj, subjraw, msg, msgraw, err := ul.composeMessage(nt, locale, map[string]interface{}{
		"username":     user.GetDisplayName(),
		"resourcename": info.GetName(),
	})
	if err != nil {
		return OC10Notification{}, err
	}

	return OC10Notification{
		EventID:        eventid,
		Service:        ul.cfg.Service.Name,
		UserName:       user.GetUsername(),
		Timestamp:      ts.Format(time.RFC3339Nano),
		ResourceID:     storagespace.FormatResourceID(*info.GetId()),
		ResourceType:   _resourceTypeShare,
		Subject:        subj,
		SubjectRaw:     subjraw,
		Message:        msg,
		MessageRaw:     msgraw,
		MessageDetails: ul.getDetails(user, nil, info, shareid),
	}, nil
}

func (ul *UserlogService) composeMessage(nt NotificationTemplate, locale string, vars map[string]interface{}) (string, string, string, string, error) {
	subj, msg, err := ul.parseTemplate(nt, locale)
	if err != nil {
		return "", "", "", "", err
	}

	subject := ul.executeTemplate(subj, vars)

	subjectraw := ul.executeTemplate(subj, map[string]interface{}{
		"username":     "{user}",
		"spacename":    "{space}",
		"resourcename": "{resource}",
	})

	message := ul.executeTemplate(msg, vars)

	messageraw := ul.executeTemplate(msg, map[string]interface{}{
		"username":     "{user}",
		"spacename":    "{space}",
		"resourcename": "{resource}",
	})

	return subject, subjectraw, message, messageraw, nil

}

func (ul *UserlogService) getDetails(user *user.User, space *storageprovider.StorageSpace, item *storageprovider.ResourceInfo, shareid *collaboration.ShareId) map[string]interface{} {
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

	if shareid != nil {
		details["share"] = map[string]string{
			"id": shareid.GetOpaqueId(),
		}
	}

	return details
}

func (ul *UserlogService) parseTemplate(nt NotificationTemplate, locale string) (*template.Template, *template.Template, error) {
	// Create Locale with library path and language code and load domain '.../default.po'
	l := gotext.NewLocale(_pathToLocales, locale)
	l.AddDomain("default")

	subject, err := template.New("").Parse(l.Get(nt.Subject))
	if err != nil {
		return nil, nil, err
	}

	message, err := template.New("").Parse(l.Get(nt.Message))
	return subject, message, err

}

func (ul *UserlogService) executeTemplate(tpl *template.Template, vars map[string]interface{}) string {
	var writer bytes.Buffer
	if err := tpl.Execute(&writer, vars); err != nil {
		ul.log.Error().Err(err).Str("templateName", tpl.Name()).Msg("cannot execute template")
		return ""
	}

	return writer.String()
}
