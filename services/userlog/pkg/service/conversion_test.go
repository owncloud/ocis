package service_test

import (
	"context"
	"testing"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/reva/v2/pkg/events"

	"github.com/owncloud/ocis/v2/services/userlog/pkg/service"
)

// nil gateway selector proves the notification is built without a lookup.
func TestConvertEventUploadFailed(t *testing.T) {
	c := service.NewConverter(context.Background(), "en", nil, "userlog", "", "en")

	ev := events.UploadReady{
		Failed:   true,
		Filename: "big-upload.iso",
		ExecutingUser: &user.User{
			Id:       &user.UserId{OpaqueId: "user-1"},
			Username: "einstein",
		},
		ResourceID: &provider.ResourceId{
			StorageId: "storage-1",
			SpaceId:   "space-1",
			OpaqueId:  "item-1",
		},
	}

	n, err := c.ConvertEvent("event-1", ev)
	if err != nil {
		t.Fatalf("ConvertEvent returned error: %v", err)
	}

	if n.Subject != "Upload failed" {
		t.Errorf("unexpected subject: %q", n.Subject)
	}
	if want := "The upload of big-upload.iso could not be finalized. Please try again."; n.Message != want {
		t.Errorf("unexpected message: got %q want %q", n.Message, want)
	}
	if n.UserName != "einstein" {
		t.Errorf("unexpected username: %q", n.UserName)
	}
	if n.ResourceID != "storage-1$space-1!item-1" {
		t.Errorf("unexpected resource id: %q", n.ResourceID)
	}
}
