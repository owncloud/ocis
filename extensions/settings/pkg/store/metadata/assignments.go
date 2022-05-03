// Package store implements the go-micro store interface
package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/owncloud/ocis/extensions/settings/pkg/store/defaults"
	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
)

// ListRoleAssignments loads and returns all role assignments matching the given assignment identifier.
func (s *Store) ListRoleAssignments(accountUUID string) ([]*settingsmsg.UserRoleAssignment, error) {
	if s.mdc == nil {
		return s.defaultRoleAssignments(accountUUID), nil
	}
	s.Init()
	ctx := context.TODO()
	assIDs, err := s.mdc.ReadDir(ctx, accountPath(accountUUID))
	if err != nil {
		return nil, err
	}

	ass := make([]*settingsmsg.UserRoleAssignment, 0, len(assIDs))
	for _, assID := range assIDs {
		b, err := s.mdc.SimpleDownload(ctx, assignmentPath(accountUUID, assID))
		if err != nil {
			return nil, err
		}

		a := &settingsmsg.UserRoleAssignment{}
		err = json.Unmarshal(b, a)
		if err != nil {
			return nil, err
		}

		ass = append(ass, a)
	}
	return ass, nil
}

// WriteRoleAssignment appends the given role assignment to the existing assignments of the respective account.
func (s *Store) WriteRoleAssignment(accountUUID, roleID string) (*settingsmsg.UserRoleAssignment, error) {
	s.Init()
	ctx := context.TODO()
	// as per https://github.com/owncloud/product/issues/103 "Each user can have exactly one role"
	_ = s.mdc.Delete(ctx, accountPath(accountUUID))
	// TODO: How to differentiate between 'not found' and other errors?

	err := s.mdc.MakeDirIfNotExist(ctx, accountPath(accountUUID))
	if err != nil {
		return nil, err
	}

	ass := &settingsmsg.UserRoleAssignment{
		Id:          uuid.Must(uuid.NewV4()).String(),
		AccountUuid: accountUUID,
		RoleId:      roleID,
	}
	b, err := json.Marshal(ass)
	if err != nil {
		return nil, err
	}
	return ass, s.mdc.SimpleUpload(ctx, assignmentPath(accountUUID, ass.Id), b)
}

// RemoveRoleAssignment deletes the given role assignment from the existing assignments of the respective account.
func (s *Store) RemoveRoleAssignment(assignmentID string) error {
	s.Init()
	ctx := context.TODO()
	accounts, err := s.mdc.ReadDir(ctx, accountsFolderLocation)
	if err != nil {
		return err
	}

	// TODO: use indexer to avoid spamming Metadata service
	for _, accID := range accounts {
		assIDs, err := s.mdc.ReadDir(ctx, accountPath(accID))
		if err != nil {
			// TODO: error?
			continue
		}

		for _, assID := range assIDs {
			if assID == assignmentID {
				return s.mdc.Delete(ctx, assignmentPath(accID, assID))
			}
		}
	}
	return fmt.Errorf("assignmentID '%s' not found", assignmentID)
}

func (s *Store) defaultRoleAssignments(accID string) []*settingsmsg.UserRoleAssignment {
	var assmnts []*settingsmsg.UserRoleAssignment
	for _, r := range defaults.DefaultRoleAssignments(s.cfg) {
		if r.AccountUuid == accID {
			assmnts = append(assmnts, r)
		}
	}
	return assmnts
}

func accountPath(accountUUID string) string {
	return fmt.Sprintf("%s/%s", accountsFolderLocation, accountUUID)
}

func assignmentPath(accountUUID string, assignmentID string) string {
	return fmt.Sprintf("%s/%s/%s", accountsFolderLocation, accountUUID, assignmentID)
}
