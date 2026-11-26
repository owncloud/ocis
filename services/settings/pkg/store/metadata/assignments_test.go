package store

import (
	"context"
	"log"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	olog "github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config/defaults"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/stretchr/testify/require"
)

var (
	einstein = "00000000-0000-0000-0000-000000000001"
	marie    = "00000000-0000-0000-0000-000000000002"
	moss     = "00000000-0000-0000-0000-000000000003"

	role1 = "11111111-1111-1111-1111-111111111111"
	role2 = "22222222-2222-2222-2222-222222222222"

	logger = olog.NewLogger(
		olog.Color(true),
		olog.Pretty(true),
		olog.Level("info"),
	)

	bundles = []*settingsmsg.Bundle{
		{
			Id:          "f36db5e6-a03c-40df-8413-711c67e40b47",
			Type:        settingsmsg.Bundle_TYPE_ROLE,
			DisplayName: "test role - reads | update",
			Name:        "TEST_ROLE",
			Extension:   "ocis-settings",
			Resource: &settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_BUNDLE,
			},
			Settings: []*settingsmsg.Setting{
				{
					Id:   "updateID",
					Name: "update",
					Value: &settingsmsg.Setting_PermissionValue{
						PermissionValue: &settingsmsg.Permission{
							Operation: settingsmsg.Permission_OPERATION_UPDATE,
						},
					},
					Resource: &settingsmsg.Resource{
						Type: settingsmsg.Resource_TYPE_SETTING,
					},
				},
				{
					Id:   "readID",
					Name: "read",
					Value: &settingsmsg.Setting_PermissionValue{
						PermissionValue: &settingsmsg.Permission{
							Operation: settingsmsg.Permission_OPERATION_READ,
						},
					},
					Resource: &settingsmsg.Resource{
						Type: settingsmsg.Resource_TYPE_BUNDLE,
					},
				},
			},
		},
		{
			Id:          "44f1a664-0a7f-461a-b0be-5b59e46bbc7a",
			Type:        settingsmsg.Bundle_TYPE_ROLE,
			DisplayName: "another",
			Name:        "ANOTHER_TEST_ROLE",
			Extension:   "ocis-settings",
			Resource: &settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_BUNDLE,
			},
			Settings: []*settingsmsg.Setting{
				{
					Id:   "readID",
					Name: "read",
					Value: &settingsmsg.Setting_PermissionValue{
						PermissionValue: &settingsmsg.Permission{
							Operation: settingsmsg.Permission_OPERATION_READ,
						},
					},
				},
			},
		},
	}
)

func initStore() *Store {
	s := &Store{
		Logger: logger,
		l:      &sync.Mutex{},
		cfg:    defaults.DefaultConfig(),
	}
	s.cfg.Commons = &shared.Commons{
		AdminUserID: uuid.Must(uuid.NewV4()).String(),
	}

	_ = NewMDC(s)
	return s
}

func setupRoles(s *Store) {
	for i := range bundles {
		if _, err := s.WriteBundle(bundles[i]); err != nil {
			log.Fatal("error initializing ", err)
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
		t.Run(scenario.name, func(t *testing.T) {
			s := initStore()
			setupRoles(s)
			firstAssignment, err := s.WriteRoleAssignment(scenario.userID, scenario.firstRole)
			require.NoError(t, err)
			require.Equal(t, firstAssignment.RoleId, scenario.firstRole)
			// TODO: check entry exists

			list, err := s.ListRoleAssignments(scenario.userID)
			require.NoError(t, err)
			require.Equal(t, 1, len(list))
			require.Equal(t, list[0].RoleId, scenario.firstRole)

			// creating another assignment shouldn't add another entry, as we support max one role per user.
			// assigning the second role should remove the old
			secondAssignment, err := s.WriteRoleAssignment(scenario.userID, scenario.secondRole)
			require.NoError(t, err)
			require.Equal(t, secondAssignment.RoleId, scenario.secondRole)

			list, err = s.ListRoleAssignments(scenario.userID)
			require.NoError(t, err)
			require.Equal(t, 1, len(list))
			require.Equal(t, list[0].RoleId, scenario.secondRole)
		})
	}
}

func TestListRoleAssignmentByRole(t *testing.T) {
	type assignment struct {
		userID string
		roleID string
	}

	var scenarios = []struct {
		name        string
		assignments []assignment
		queryRole   string
		numResults  int
	}{
		{
			name: "just 2 assignments",
			assignments: []assignment{
				{
					userID: einstein,
					roleID: role1,
				}, {
					userID: marie,
					roleID: role1,
				},
			},
			queryRole:  role1,
			numResults: 2,
		},
		{
			name: "no assignments match",
			assignments: []assignment{
				{
					userID: einstein,
					roleID: role1,
				}, {
					userID: marie,
					roleID: role1,
				},
			},
			queryRole:  role2,
			numResults: 0,
		},
		{
			name: "only one assignment matches",
			assignments: []assignment{
				{
					userID: einstein,
					roleID: role1,
				}, {
					userID: marie,
					roleID: role1,
				}, {
					userID: moss,
					roleID: role2,
				},
			},
			queryRole:  role2,
			numResults: 1,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			s := initStore()
			setupRoles(s)
			for _, a := range scenario.assignments {
				ass, err := s.WriteRoleAssignment(a.userID, a.roleID)
				require.NoError(t, err)
				require.Equal(t, ass.RoleId, a.roleID)
			}

			list, err := s.ListRoleAssignmentsByRole(scenario.queryRole)
			require.NoError(t, err)
			require.Equal(t, scenario.numResults, len(list))
			for _, ass := range list {
				require.Equal(t, ass.RoleId, scenario.queryRole)
			}
		})
	}
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
		t.Run(scenario.name, func(t *testing.T) {
			s := initStore()
			setupRoles(s)
			assignment, err := s.WriteRoleAssignment(scenario.userID, scenario.firstRole)
			require.NoError(t, err)
			require.Equal(t, assignment.RoleId, scenario.firstRole)
			// TODO: uncomment
			// require.True(t, mdc.IDExists(assignment.RoleId))

			list, err := s.ListRoleAssignments(scenario.userID)
			require.NoError(t, err)
			require.Equal(t, 1, len(list))
			require.Equal(t, assignment.Id, list[0].Id)

			err = s.RemoveRoleAssignment(assignment.Id)
			require.NoError(t, err)
			// TODO: uncomment
			// require.False(t, mdc.IDExists(assignment.RoleId))

			list, err = s.ListRoleAssignments(scenario.userID)
			require.NoError(t, err)
			require.Equal(t, 0, len(list))

			err = s.RemoveRoleAssignment(assignment.Id)
			require.Error(t, err)
			// TODO: do we want a custom error message?
		})
	}
}

// Each settings pod has its own CachedMDC. A role change on one pod does not
// clear the caches of the other pods. On the other pod, the cached listing of
// the account folder can then name an assignment file that is gone. The
// download of that file gives a real NotFound. ListRoleAssignments must not
// return an empty list in that case: the proxy then assigns the default role
// and so replaces the new role.
func TestListRoleAssignmentsWithStaleDirCache(t *testing.T) {
	backend := &MockedMetadataClient{data: make(map[string][]byte)}
	newPod := func() *Store {
		s := &Store{
			Logger: logger,
			l:      &sync.Mutex{},
			cfg:    defaults.DefaultConfig(),
		}
		s.cfg.Commons = &shared.Commons{
			AdminUserID: uuid.Must(uuid.NewV4()).String(),
		}
		require.NoError(t, s.initMetadataClient(&CachedMDC{next: backend, cfg: s.cfg, logger: logger}))
		return s
	}
	podA, podB := newPod(), newPod()

	// pod B writes the first role and caches the listing of the account folder
	_, err := podB.WriteRoleAssignment(einstein, role1)
	require.NoError(t, err)
	list, err := podB.ListRoleAssignments(einstein)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, role1, list[0].RoleId)

	// pod A changes the role, the old assignment file is gone
	_, err = podA.WriteRoleAssignment(einstein, role2)
	require.NoError(t, err)

	// the files cache entry on pod B is older than the dirs cache entry, so it
	// expires first
	cached, ok := podB.mdc.(*CachedMDC)
	require.True(t, ok)
	require.NoError(t, cached.removePrefix(cached.filesCache, accountPath(einstein)))

	list, err = podB.ListRoleAssignments(einstein)
	require.NoError(t, err)
	require.Len(t, list, 1, "a stale listing must not hide the current assignment")
	require.Equal(t, role2, list[0].RoleId)
}

// goneAccountMDC lists the account folder of gone, but that folder is gone,
// for example because the listing of the accounts folder is stale.
type goneAccountMDC struct {
	*MockedMetadataClient
	gone string
}

func (m *goneAccountMDC) ReadDir(ctx context.Context, id string) ([]string, error) {
	switch id {
	case accountPath(m.gone):
		return nil, errtypes.NotFound(id)
	case accountsFolderLocation:
		ids, err := m.MockedMetadataClient.ReadDir(ctx, id)
		return append(ids, m.gone), err
	default:
		return m.MockedMetadataClient.ReadDir(ctx, id)
	}
}

// One account folder that is gone must not hide the assignments of all
// other accounts.
func TestListRoleAssignmentsByRoleSkipsGoneAccount(t *testing.T) {
	s := &Store{
		Logger: logger,
		l:      &sync.Mutex{},
		cfg:    defaults.DefaultConfig(),
	}
	s.cfg.Commons = &shared.Commons{
		AdminUserID: uuid.Must(uuid.NewV4()).String(),
	}
	require.NoError(t, s.initMetadataClient(&goneAccountMDC{
		MockedMetadataClient: &MockedMetadataClient{data: make(map[string][]byte)},
		gone:                 moss,
	}))

	_, err := s.WriteRoleAssignment(einstein, role1)
	require.NoError(t, err)
	_, err = s.WriteRoleAssignment(marie, role1)
	require.NoError(t, err)

	list, err := s.ListRoleAssignmentsByRole(role1)
	require.NoError(t, err)
	require.Len(t, list, 2)
}
