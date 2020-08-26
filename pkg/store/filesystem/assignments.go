// Package store implements the go-micro store interface
package store

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
)

// ListRoleAssignments loads and returns all role assignments matching the given assignment identifier.
func (s Store) ListRoleAssignments(accountUUID string) ([]*proto.UserRoleAssignment, error) {
	var records []*proto.UserRoleAssignment
	assignmentsFolder := s.buildFolderPathForRoleAssignments(false)
	assignmentFiles, err := ioutil.ReadDir(assignmentsFolder)
	if err != nil {
		return records, nil
	}

	for _, assignmentFile := range assignmentFiles {
		record := proto.UserRoleAssignment{}
		err = s.parseRecordFromFile(&record, filepath.Join(assignmentsFolder, assignmentFile.Name()))
		if err == nil {
			if record.AccountUuid == accountUUID {
				records = append(records, &record)
			}
		}
	}

	return records, nil
}

// WriteRoleAssignment appends the given role assignment to the existing assignments of the respective account.
func (s Store) WriteRoleAssignment(accountUUID, roleID string) (*proto.UserRoleAssignment, error) {
	// as per https://github.com/owncloud/product/issues/103 "Each user can have exactly one role"
	list, err := s.ListRoleAssignments(accountUUID)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		filePath := s.buildFilePathForRoleAssignment(list[0].Id, true)
		if err := os.Remove(filePath); err != nil {
			return nil, err
		}
	}

	assignment := &proto.UserRoleAssignment{
		Id:          uuid.Must(uuid.NewV4()).String(),
		AccountUuid: accountUUID,
		RoleId:      roleID,
	}
	filePath := s.buildFilePathForRoleAssignment(assignment.Id, true)
	if err := s.writeRecordToFile(assignment, filePath); err != nil {
		return nil, err
	}

	s.Logger.Debug().Msgf("request contents written to file: %v", filePath)
	return assignment, nil
}

// RemoveRoleAssignment deletes the given role assignment from the existing assignments of the respective account.
func (s Store) RemoveRoleAssignment(assignmentID string) error {
	filePath := s.buildFilePathForRoleAssignment(assignmentID, false)
	return os.Remove(filePath)
}
