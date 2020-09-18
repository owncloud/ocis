package store

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"

	olog "github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis/settings/pkg/proto/v0"
	"github.com/stretchr/testify/assert"
)

var (
	einstein = "a4d07560-a670-4be9-8d60-9b547751a208"
	//marie    = "3c054db3-eec1-4ca4-b985-bc56dcf560cb"

	s = Store{
		dataPath: dataRoot,
		Logger:   logger,
	}

	logger = olog.NewLogger(
		olog.Color(true),
		olog.Pretty(true),
		olog.Level("info"),
	)

	bundles = []*proto.Bundle{
		{
			Id:          "f36db5e6-a03c-40df-8413-711c67e40b47",
			Type:        proto.Bundle_TYPE_ROLE,
			DisplayName: "test role - reads | update",
			Name:        "TEST_ROLE",
			Extension:   "ocis-settings",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_BUNDLE,
			},
			Settings: []*proto.Setting{
				{
					Name: "update",
					Value: &proto.Setting_PermissionValue{
						PermissionValue: &proto.Permission{
							Operation: proto.Permission_OPERATION_UPDATE,
						},
					},
				},
				{
					Name: "read",
					Value: &proto.Setting_PermissionValue{
						PermissionValue: &proto.Permission{
							Operation: proto.Permission_OPERATION_READ,
						},
					},
				},
			},
		},
		{
			Id:          "44f1a664-0a7f-461a-b0be-5b59e46bbc7a",
			Type:        proto.Bundle_TYPE_ROLE,
			DisplayName: "another",
			Name:        "ANOTHER_TEST_ROLE",
			Extension:   "ocis-settings",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_BUNDLE,
			},
			Settings: []*proto.Setting{
				{
					Name: "read",
					Value: &proto.Setting_PermissionValue{
						PermissionValue: &proto.Permission{
							Operation: proto.Permission_OPERATION_READ,
						},
					},
				},
			},
		},
	}
)

func init() {
	setupRoles()
}

func setupRoles() {
	for i := range bundles {
		if _, err := s.WriteBundle(bundles[i]); err != nil {
			log.Fatal(err)
		}
	}
}

func TestAssignmentUniqueness(t *testing.T) {
	var scenarios = []struct {
		name       string
		userID     string
		firstRole  string
		secondRole string
	}{
		{
			"roles assignments",
			einstein,
			"f36db5e6-a03c-40df-8413-711c67e40b47",
			"44f1a664-0a7f-461a-b0be-5b59e46bbc7a",
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			firstAssignment, err := s.WriteRoleAssignment(scenario.userID, scenario.firstRole)
			assert.NoError(t, err)
			assert.Equal(t, firstAssignment.RoleId, scenario.firstRole)
			assert.FileExists(t, filepath.Join(dataRoot, "assignments", firstAssignment.Id+".json"))

			list, err := s.ListRoleAssignments(scenario.userID)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(list))

			// creating another assignment shouldn't add another entry, as we support max one role per user.
			secondAssignment, err := s.WriteRoleAssignment(scenario.userID, scenario.secondRole)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(list))

			// assigning the second role should remove the old file and create a new one.
			list, err = s.ListRoleAssignments(scenario.userID)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(list))
			assert.Equal(t, secondAssignment.RoleId, scenario.secondRole)
			assert.NoFileExists(t, filepath.Join(dataRoot, "assignments", firstAssignment.Id+".json"))
			assert.FileExists(t, filepath.Join(dataRoot, "assignments", secondAssignment.Id+".json"))
		})
	}
	burnRoot()
}

func TestDeleteAssignment(t *testing.T) {
	var scenarios = []struct {
		name       string
		userID     string
		firstRole  string
		secondRole string
	}{
		{
			"roles assignments",
			einstein,
			"f36db5e6-a03c-40df-8413-711c67e40b47",
			"44f1a664-0a7f-461a-b0be-5b59e46bbc7a",
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			assignment, err := s.WriteRoleAssignment(scenario.userID, scenario.firstRole)
			assert.NoError(t, err)
			assert.Equal(t, assignment.RoleId, scenario.firstRole)
			assert.FileExists(t, filepath.Join(dataRoot, "assignments", assignment.Id+".json"))

			list, err := s.ListRoleAssignments(scenario.userID)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(list))

			err = s.RemoveRoleAssignment(assignment.Id)
			assert.NoError(t, err)

			list, err = s.ListRoleAssignments(scenario.userID)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(list))

			err = s.RemoveRoleAssignment(assignment.Id)
			merr := &os.PathError{}
			assert.Equal(t, true, errors.As(err, &merr))
		})
	}
	burnRoot()
}
