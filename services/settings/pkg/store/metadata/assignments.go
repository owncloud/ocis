// Package store implements the go-micro store interface
package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/gofrs/uuid"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/settings"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
)

// ListRoleAssignments loads and returns all role assignments matching the given assignment identifier.
func (s *Store) ListRoleAssignments(accountUUID string) ([]*settingsmsg.UserRoleAssignment, error) {
	// shortcut for service accounts
	for _, serviceAccountID := range s.cfg.ServiceAccountIDs {
		if accountUUID == serviceAccountID {
			return []*settingsmsg.UserRoleAssignment{
				{
					Id:          uuid.Must(uuid.NewV4()).String(), // should we hardcode this id too?
					AccountUuid: accountUUID,
					RoleId:      defaults.BundleUUIDServiceAccount,
				},
			}, nil
		}
	}
	s.Init()
	ctx := context.TODO()
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

// ListRoleAssignmentsByRole returns all role assignmentes matching the give roleID
func (s *Store) ListRoleAssignmentsByRole(roleID string) ([]*settingsmsg.UserRoleAssignment, error) {
	s.Init()
	ctx := context.TODO()
	cachedAssignments, err := s.assignmentsCache.List(ctx, roleID)
	switch err.(type) {
	case nil:
		// continue
	case errtypes.NotFound:
		return make([]*settingsmsg.UserRoleAssignment, 0), nil
	default:
		return nil, err
	}
	assignments := make([]*settingsmsg.UserRoleAssignment, 0, len(cachedAssignments))
	for id, v := range cachedAssignments {
		assignments = append(assignments,
			&settingsmsg.UserRoleAssignment{
				Id:          v.AssignmentID,
				AccountUuid: id,
				RoleId:      roleID,
			},
		)
	}

	return assignments, nil
}

// WriteRoleAssignment appends the given role assignment to the existing assignments of the respective account.
func (s *Store) WriteRoleAssignment(accountUUID, roleID string) (*settingsmsg.UserRoleAssignment, error) {
	s.Init()
	ctx := context.TODO()
	// as per https://github.com/owncloud/product/issues/103 "Each user can have exactly one role"
	err := s.mdc.Delete(ctx, accountPath(accountUUID))
	switch err.(type) {
	case nil:
		// continue
	case errtypes.NotFound:
		// already gone, continue
	default:
		return nil, err
	}

	// remove from cache
	r, err := s.ListBundles(settingsmsg.Bundle_TYPE_ROLE, []string{})
	if err != nil {
		return nil, err
	}
	for _, role := range r {
		err = s.assignmentsCache.Remove(ctx, role.GetId(), accountUUID)
		if err != nil {
			return nil, err
		}
	}

	err = s.mdc.MakeDirIfNotExist(ctx, accountPath(accountUUID))
	if err != nil {
		return nil, err
	}

	assignment := &settingsmsg.UserRoleAssignment{
		Id:          uuid.Must(uuid.NewV4()).String(),
		AccountUuid: accountUUID,
		RoleId:      roleID,
	}
	b, err := json.Marshal(assignment)
	if err != nil {
		return nil, err
	}
	err = s.mdc.SimpleUpload(ctx, assignmentPath(accountUUID, assignment.Id), b)
	if err != nil {
		return assignment, err
	}

	err = s.assignmentsCache.Add(ctx, roleID, assignment)
	if err != nil {
		return assignment, err
	}
	return assignment, err
}

// RemoveRoleAssignment deletes the given role assignment from the existing assignments of the respective account.
func (s *Store) RemoveRoleAssignment(assignmentID string) error {
	s.Init()
	ctx := context.TODO()
	accounts, err := s.mdc.ReadDir(ctx, accountsFolderLocation)
	switch err.(type) {
	case nil:
		// continue
	case errtypes.NotFound:
		return fmt.Errorf("assignmentID '%s' %w", assignmentID, settings.ErrNotFound)
	default:
		return err
	}

	// TODO: use indexer to avoid spamming Metadata service
	for _, accID := range accounts {
		assignmentIDs, err := s.mdc.ReadDir(ctx, accountPath(accID))
		if err != nil {
			// TODO: error?
			continue
		}

		for _, id := range assignmentIDs {
			if id == assignmentID {
				b, err := s.mdc.SimpleDownload(ctx, assignmentPath(accID, id))
				switch err.(type) {
				case nil:
					a := &settingsmsg.UserRoleAssignment{}
					if err = json.Unmarshal(b, a); err != nil {
						s.Logger.Error().Err(err).Str("assignmentid", id).Msg("failed to unmarshall assignment")
						// no return here, as we still want to delete the assignment
					} else if err = s.assignmentsCache.Remove(ctx, a.RoleId, accID); err != nil {
						s.Logger.Error().Err(err).Str("assignmentid", id).Msg("failed to remove assignment from cache")
						// no return here, as we still want to delete the assignment
					}
					// continue
				case errtypes.NotFound:
					continue
				default:
					s.Logger.Error().Err(err).Str("assignmentid", id).Msg("could not download assignment, for cache cleanup")
					// We're not returning here, as we still want to delete the assignment
				}

				// as per https://github.com/owncloud/product/issues/103 "Each user can have exactly one role"
				// we also have to delete the cached dir listing
				return s.mdc.Delete(ctx, accountPath(accID))
			}
		}
	}
	return fmt.Errorf("assignmentID '%s' %w", assignmentID, settings.ErrNotFound)
}

func accountPath(accountUUID string) string {
	return fmt.Sprintf("%s/%s", accountsFolderLocation, accountUUID)
}

func assignmentPath(accountUUID string, assignmentID string) string {
	return fmt.Sprintf("%s/%s/%s", accountsFolderLocation, accountUUID, assignmentID)
}
