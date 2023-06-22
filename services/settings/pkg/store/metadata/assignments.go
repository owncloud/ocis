// Package store implements the go-micro store interface
package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/gofrs/uuid"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
)

// ListRoleAssignments loads and returns all role assignments matching the given assignment identifier.
func (s *Store) ListRoleAssignments(accountUUID string) ([]*settingsmsg.UserRoleAssignment, error) {
	s.Init()
	ctx := context.TODO()

	b, err := s.mdc.SimpleDownload(ctx, accountAssignmentPath(accountUUID))
	switch err.(type) {
	case nil:
		a := &settingsmsg.UserRoleAssignment{}
		err = json.Unmarshal(b, a)
		if err != nil {
			return nil, err
		}
		return []*settingsmsg.UserRoleAssignment{a}, nil
	case errtypes.NotFound:
		// role assignments not found - migrate from old structure
		asslegacy, err := s.listStoreAssignmentsLegacy(ctx, accountUUID)
		if err != nil || len(asslegacy) == 0 {
			return asslegacy, err
		}

		a, err := s.WriteRoleAssignment(asslegacy[0].AccountUuid, asslegacy[0].RoleId)
		return []*settingsmsg.UserRoleAssignment{a}, err

	default:
		return nil, err
	}
}

// WriteRoleAssignment appends the given role assignment to the existing assignments of the respective account.
func (s *Store) WriteRoleAssignment(accountUUID, roleID string) (*settingsmsg.UserRoleAssignment, error) {
	s.Init()
	ctx := context.TODO()
	// as per https://github.com/owncloud/product/issues/103 "Each user can have exactly one role"
	ass := &settingsmsg.UserRoleAssignment{
		Id:          uuid.Must(uuid.NewV4()).String(),
		AccountUuid: accountUUID,
		RoleId:      roleID,
	}
	b, err := json.Marshal(ass)
	if err != nil {
		return nil, err
	}
	return ass, s.mdc.SimpleUpload(ctx, accountAssignmentPath(accountUUID), b)
}

// RemoveRoleAssignment deletes the given role assignment from the existing assignments of the respective account.
func (s *Store) RemoveRoleAssignment(accountUUID string) error { // <- BREAKING CHANGE - this needs the accountUUID now
	s.Init()
	ctx := context.TODO()
	return s.mdc.Delete(ctx, accountAssignmentPath(accountUUID))
}

func (s *Store) listStoreAssignmentsLegacy(ctx context.Context, accountUUID string) ([]*settingsmsg.UserRoleAssignment, error) {
	assIDs, err := s.mdc.ReadDir(ctx, accountPath(accountUUID))
	switch err.(type) {
	case nil:
		// continue
	case errtypes.NotFound:
		return make([]*settingsmsg.UserRoleAssignment, 0), nil
	default:
		return nil, err
	}

	ass := make([]*settingsmsg.UserRoleAssignment, 0, len(assIDs))
	for _, assID := range assIDs {
		b, err := s.mdc.SimpleDownload(ctx, assignmentPath(accountUUID, assID))
		switch err.(type) {
		case nil:
			// continue
		case errtypes.NotFound:
			continue
		default:
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

func accountPath(accountUUID string) string {
	return fmt.Sprintf("%s/%s", accountsFolderLocation, accountUUID)
}

func assignmentPath(accountUUID string, assignmentID string) string {
	return fmt.Sprintf("%s/%s/%s", accountsFolderLocation, accountUUID, assignmentID)
}

func accountAssignmentPath(accountUUID string) string {
	return fmt.Sprintf("%s/%s-assignments", accountsFolderLocation, accountUUID)
}
