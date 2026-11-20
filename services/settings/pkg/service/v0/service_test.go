package svc

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	v0 "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/settings/mocks"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	merrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
)

var (
	ctxWithUUID      = metadata.Set(context.Background(), middleware.AccountID, "61445573-4dbe-4d56-88dc-88ab47aceba7")
	ctxWithEmptyUUID = metadata.Set(context.Background(), middleware.AccountID, "")
	emptyCtx         = context.Background()

	scenarios = []struct {
		name        string
		accountUUID string
		ctx         context.Context
		expect      string
	}{
		{
			name:        "context with UUID; identifier = 'me'",
			ctx:         ctxWithUUID,
			accountUUID: "me",
			expect:      "61445573-4dbe-4d56-88dc-88ab47aceba7",
		},
		{
			name:        "context with empty UUID; identifier = 'me'",
			ctx:         ctxWithEmptyUUID,
			accountUUID: "me",
			expect:      "",
		},
		{
			name:        "context without UUID; identifier = 'me'",
			ctx:         emptyCtx,
			accountUUID: "me",
			expect:      "",
		},
		{
			name:        "context with UUID; identifier not 'me'",
			ctx:         ctxWithUUID,
			accountUUID: "",
			expect:      "",
		},
	}
)

func TestGetValidatedAccountUUID(t *testing.T) {
	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.name, func(t *testing.T) {
			got := getValidatedAccountUUID(scenario.ctx, scenario.accountUUID)
			assert.NotPanics(t, func() {
				getValidatedAccountUUID(emptyCtx, scenario.accountUUID)
			})
			assert.Equal(t, scenario.expect, got)
		})
	}
}

func TestEditOwnRoleAssignment(t *testing.T) {
	manager := &mocks.Manager{}
	svc := Service{
		manager: manager,
	}
	a := []*settingsmsg.UserRoleAssignment{}
	manager.On("ListRoleAssignments", mock.Anything).Return(a, nil)
	manager.On("WriteRoleAssignment", mock.Anything, mock.Anything).Return(nil, nil)
	// Creating an initial self assignment is expected to succeed for UserRole when no assignment exists yet
	req := v0.AssignRoleToUserRequest{
		AccountUuid: "61445573-4dbe-4d56-88dc-88ab47aceba7",
		RoleId:      defaults.BundleUUIDRoleUser,
	}
	res := v0.AssignRoleToUserResponse{}
	err := svc.AssignRoleToUser(ctxWithUUID, &req, &res)
	assert.Nil(t, err)

	// Creating an initial self assignment is expected to succeed for UserLightRole when no assignment exists yet
	req = v0.AssignRoleToUserRequest{
		AccountUuid: "61445573-4dbe-4d56-88dc-88ab47aceba7",
		RoleId:      defaults.BundleUUIDRoleUserLight,
	}
	res = v0.AssignRoleToUserResponse{}
	err = svc.AssignRoleToUser(ctxWithUUID, &req, &res)
	assert.Nil(t, err)

	// Creating an initial self assignment is expected to fail for non UserRole when no assignment exists yet
	req = v0.AssignRoleToUserRequest{
		AccountUuid: "61445573-4dbe-4d56-88dc-88ab47aceba7",
		RoleId:      defaults.BundleUUIDRoleAdmin,
	}
	res = v0.AssignRoleToUserResponse{}
	err = svc.AssignRoleToUser(ctxWithUUID, &req, &res)
	assert.NotNil(t, err)

	manager = &mocks.Manager{}
	svc = Service{
		manager: manager,
	}
	a = []*settingsmsg.UserRoleAssignment{
		{
			Id:          "00000000-0000-0000-0000-000000000001",
			AccountUuid: "61445573-4dbe-4d56-88dc-88ab47aceba7",
			RoleId:      "aceb15b8-7486-479f-ae32-c91118e07a39",
		},
	}
	editRolePermission := &settingsmsg.Permission{
		Operation:  settingsmsg.Permission_OPERATION_READWRITE,
		Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
	}
	manager.On("ListRoleAssignments", mock.Anything).Return(a, nil)
	manager.On("ReadPermissionByID", mock.Anything, mock.Anything).Return(editRolePermission, nil)

	// Creating an self assignment is expect to fail if there is already an assingment
	req = v0.AssignRoleToUserRequest{
		AccountUuid: "61445573-4dbe-4d56-88dc-88ab47aceba7",
		RoleId:      defaults.BundleUUIDRoleUser,
	}
	res = v0.AssignRoleToUserResponse{}
	err = svc.AssignRoleToUser(ctxWithUUID, &req, &res)
	assert.NotNil(t, err)

	manager.On("WriteRoleAssignment", mock.Anything, mock.Anything).Return(nil, nil)
	// Creating an assignment for somebody else is expected to succeed, give the right permissions
	req = v0.AssignRoleToUserRequest{
		AccountUuid: "00000000-0000-0000-0000-000000000000",
		RoleId:      "aceb15b8-7486-479f-ae32-c91118e07a39",
	}
	res = v0.AssignRoleToUserResponse{}
	err = svc.AssignRoleToUser(ctxWithUUID, &req, &res)
	assert.Nil(t, err)
}

func TestRemoveOwnRoleAssignment(t *testing.T) {
	manager := &mocks.Manager{}
	a := []*settingsmsg.UserRoleAssignment{
		{
			Id:          "00000000-0000-0000-0000-000000000001",
			AccountUuid: "61445573-4dbe-4d56-88dc-88ab47aceba7",
			RoleId:      "aceb15b8-7486-479f-ae32-c91118e07a39",
		},
	}
	editRolePermission := &settingsmsg.Permission{
		Operation:  settingsmsg.Permission_OPERATION_READWRITE,
		Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
	}
	manager.On("ReadPermissionByID", mock.Anything, mock.Anything).Return(editRolePermission, nil)
	manager.On("ListRoleAssignments", mock.Anything).Return(a, nil)
	svc := Service{
		manager: manager,
	}

	// Removing a role for oneself is expected to fail
	req := v0.RemoveRoleFromUserRequest{
		Id: "00000000-0000-0000-0000-000000000001",
	}
	err := svc.RemoveRoleFromUser(ctxWithUUID, &req, nil)
	assert.NotNil(t, err)

	manager = &mocks.Manager{}
	manager.On("ListRoleAssignments", mock.Anything).Return(nil, nil)
	manager.On("RemoveRoleAssignment", mock.Anything).Return(nil)
	manager.On("ReadPermissionByID", mock.Anything, mock.Anything).Return(editRolePermission, nil)
	svc = Service{
		manager: manager,
	}
	// Removing a role for someone else is expected to fail
	req = v0.RemoveRoleFromUserRequest{
		Id: "00000000-0000-0000-0000-000000000002",
	}
	err = svc.RemoveRoleFromUser(ctxWithUUID, &req, nil)
	assert.Nil(t, err)
}

func TestListPermissionsOfCurrentUser(t *testing.T) {
	manager := &mocks.Manager{}
	a := []*settingsmsg.UserRoleAssignment{
		{
			Id:          "00000000-0000-0000-0000-000000000001",
			AccountUuid: "61445573-4dbe-4d56-88dc-88ab47aceba7",
			RoleId:      "aceb15b8-7486-479f-ae32-c91118e07a39",
		},
	}
	manager.On("ListRoleAssignments", mock.Anything).Return(a, nil)
	b := &settingsmsg.Bundle{
		Id: "aceb15b8-7486-479f-ae32-c91118e07a39",
		Settings: []*settingsmsg.Setting{
			{
				Name: "some-permission",
			},
			{
				Name: "other-permission",
			},
			{
				Name: "duplicate-permission",
			},
			{
				Name: "duplicate-permission",
			},
		},
	}
	manager.On("ReadBundle", mock.Anything).Return(b, nil)
	svc := Service{
		manager: manager,
	}

	// Listing permissions for yourself
	req := v0.ListPermissionsRequest{
		AccountUuid: "61445573-4dbe-4d56-88dc-88ab47aceba7",
	}
	res := v0.ListPermissionsResponse{}
	err := svc.ListPermissions(ctxWithUUID, &req, &res)
	assert.NoError(t, err)
	assert.Len(t, res.Permissions, 3)
}

func TestListPermissionsOfOtherUser(t *testing.T) {
	manager := &mocks.Manager{}
	svc := Service{
		manager: manager,
	}

	// Listing permissions for another user produces a not found error
	req := v0.ListPermissionsRequest{
		AccountUuid: "66666666-4444-4444-8888-88ab47aceba7",
	}
	res := v0.ListPermissionsResponse{}
	err := svc.ListPermissions(ctxWithUUID, &req, &res)
	assert.Error(t, err)

	// assert the requested account uuid was not found
	merr, ok := merrors.As(err)
	assert.True(t, ok)
	assert.Equal(t, int32(http.StatusNotFound), merr.Code)
	assert.Contains(t, err.Error(), req.AccountUuid)
}

// tooManyAccountFilters returns one more account filter than validateUserRoleAssignmentFilters allows.
func tooManyAccountFilters() []*settingsmsg.UserRoleAssignmentFilter {
	filters := make([]*settingsmsg.UserRoleAssignmentFilter, 0, maxAccountRoleAssignmentFilters+1)
	for i := 0; i <= maxAccountRoleAssignmentFilters; i++ {
		filters = append(filters, &settingsmsg.UserRoleAssignmentFilter{
			Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ACCOUNT,
			Term: &settingsmsg.UserRoleAssignmentFilter_AccountUuid{
				AccountUuid: fmt.Sprintf("account-%d", i),
			},
		})
	}
	return filters
}

func TestListRoleAssignmentsFilteredValidation(t *testing.T) {
	manager := &mocks.Manager{}
	svc := Service{
		manager: manager,
	}

	tests := map[string]struct {
		req        *v0.ListRoleAssignmentsFilteredRequest
		statusCode int32
	}{
		"no filters": {
			req:        &v0.ListRoleAssignmentsFilteredRequest{},
			statusCode: http.StatusBadRequest,
		},
		"multiple filters": {
			req: &v0.ListRoleAssignmentsFilteredRequest{
				Filters: []*settingsmsg.UserRoleAssignmentFilter{
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ACCOUNT,
						Term: &settingsmsg.UserRoleAssignmentFilter_AccountUuid{
							AccountUuid: "uid",
						},
					},
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ROLE,
						Term: &settingsmsg.UserRoleAssignmentFilter_RoleId{
							RoleId: "rid",
						},
					},
				},
			},
			statusCode: http.StatusBadRequest,
		},
		"bad filtertype": {
			req: &v0.ListRoleAssignmentsFilteredRequest{
				Filters: []*settingsmsg.UserRoleAssignmentFilter{
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_UNKNOWN,
					},
				},
			},
			statusCode: http.StatusBadRequest,
		},
		"account filter without term": {
			req: &v0.ListRoleAssignmentsFilteredRequest{
				Filters: []*settingsmsg.UserRoleAssignmentFilter{
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ACCOUNT,
					},
				},
			},
			statusCode: http.StatusBadRequest,
		},
		"account filter with invalid term": {
			req: &v0.ListRoleAssignmentsFilteredRequest{
				Filters: []*settingsmsg.UserRoleAssignmentFilter{
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ACCOUNT,
						Term: &settingsmsg.UserRoleAssignmentFilter_AccountUuid{
							AccountUuid: "invalid-&*&^%$#",
						},
					},
				},
			},
			statusCode: http.StatusBadRequest,
		},
		"role filter without term": {
			req: &v0.ListRoleAssignmentsFilteredRequest{
				Filters: []*settingsmsg.UserRoleAssignmentFilter{
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ROLE,
					},
				},
			},
			statusCode: http.StatusBadRequest,
		},
		"role filter with invalid uuid": {
			req: &v0.ListRoleAssignmentsFilteredRequest{
				Filters: []*settingsmsg.UserRoleAssignmentFilter{
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ROLE,
						Term: &settingsmsg.UserRoleAssignmentFilter_RoleId{
							RoleId: "this is no uuid",
						},
					},
				},
			},
			statusCode: http.StatusBadRequest,
		},
		"mixed account and role filters, both otherwise valid": {
			req: &v0.ListRoleAssignmentsFilteredRequest{
				Filters: []*settingsmsg.UserRoleAssignmentFilter{
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ACCOUNT,
						Term: &settingsmsg.UserRoleAssignmentFilter_AccountUuid{
							AccountUuid: "61445573-4dbe-4d56-88dc-88ab47aceba7",
						},
					},
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ROLE,
						Term: &settingsmsg.UserRoleAssignmentFilter_RoleId{
							RoleId: "aceb15b8-7486-479f-ae32-c91118e07a39",
						},
					},
				},
			},
			statusCode: http.StatusBadRequest,
		},
		"two valid role filters": {
			req: &v0.ListRoleAssignmentsFilteredRequest{
				Filters: []*settingsmsg.UserRoleAssignmentFilter{
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ROLE,
						Term: &settingsmsg.UserRoleAssignmentFilter_RoleId{
							RoleId: "aceb15b8-7486-479f-ae32-c91118e07a39",
						},
					},
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ROLE,
						Term: &settingsmsg.UserRoleAssignmentFilter_RoleId{
							RoleId: "bceb15b8-7486-479f-ae32-c91118e07a39",
						},
					},
				},
			},
			statusCode: http.StatusBadRequest,
		},
		"too many account filters": {
			req: &v0.ListRoleAssignmentsFilteredRequest{
				Filters: tooManyAccountFilters(),
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := v0.ListRoleAssignmentsResponse{}
			err := svc.ListRoleAssignmentsFiltered(ctxWithUUID, test.req, &res)
			merr, ok := merrors.As(err)
			assert.True(t, ok)
			assert.Equal(t, int32(test.statusCode), merr.Code)
		})
	}
}

func TestListRoleAssignmentsFilteredByAccount(t *testing.T) {
	accountUUID := "61445573-4dbe-4d56-88dc-88ab47aceba7"

	tests := map[string]struct {
		result []*settingsmsg.UserRoleAssignment
		err    error
		status int32
	}{
		"handles manager error": {
			result: nil,
			err:    assert.AnError,
			status: http.StatusNotFound,
		},
		"succeeds with results": {
			result: []*settingsmsg.UserRoleAssignment{
				{
					Id:          "00000000-0000-0000-0000-000000000001",
					AccountUuid: accountUUID,
					RoleId:      "aceb15b8-7486-479f-ae32-c91118e07a39",
				},
			},
			err:    nil,
			status: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			manager := &mocks.Manager{}
			svc := Service{
				manager: manager,
			}
			manager.On("ListRoleAssignments", mock.Anything).Return(test.result, test.err)
			req := &v0.ListRoleAssignmentsFilteredRequest{
				Filters: []*settingsmsg.UserRoleAssignmentFilter{
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ACCOUNT,
						Term: &settingsmsg.UserRoleAssignmentFilter_AccountUuid{
							AccountUuid: accountUUID,
						},
					},
				},
			}
			res := v0.ListRoleAssignmentsResponse{}
			err := svc.ListRoleAssignmentsFiltered(ctxWithUUID, req, &res)
			switch test.err {
			case nil:
				assert.Nil(t, err)
			default:
				merr, ok := merrors.As(err)
				assert.True(t, ok)
				assert.Equal(t, int32(test.status), merr.Code)
			}
		})
	}
}

func TestListRoleAssignmentsFilteredByMultipleAccounts(t *testing.T) {
	accountA := "61445573-4dbe-4d56-88dc-88ab47aceba7"
	accountB := "72556684-5ecf-5e67-99ed-99bc58bdfc98"

	manager := &mocks.Manager{}
	manager.On("ListRoleAssignments", accountA).Return([]*settingsmsg.UserRoleAssignment{
		{Id: "00000000-0000-0000-0000-000000000001", AccountUuid: accountA, RoleId: "aceb15b8-7486-479f-ae32-c91118e07a39"},
	}, nil)
	manager.On("ListRoleAssignments", accountB).Return([]*settingsmsg.UserRoleAssignment{
		{Id: "00000000-0000-0000-0000-000000000002", AccountUuid: accountB, RoleId: "aceb15b8-7486-479f-ae32-c91118e07a39"},
	}, nil)
	svc := Service{manager: manager}

	req := &v0.ListRoleAssignmentsFilteredRequest{
		Filters: []*settingsmsg.UserRoleAssignmentFilter{
			{Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ACCOUNT, Term: &settingsmsg.UserRoleAssignmentFilter_AccountUuid{AccountUuid: accountA}},
			{Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ACCOUNT, Term: &settingsmsg.UserRoleAssignmentFilter_AccountUuid{AccountUuid: accountB}},
		},
	}
	res := v0.ListRoleAssignmentsResponse{}
	err := svc.ListRoleAssignmentsFiltered(ctxWithUUID, req, &res)
	assert.NoError(t, err)
	assert.Len(t, res.Assignments, 2)

	gotAccounts := []string{res.Assignments[0].GetAccountUuid(), res.Assignments[1].GetAccountUuid()}
	assert.ElementsMatch(t, []string{accountA, accountB}, gotAccounts)
}

func TestListRoleAssignmentsFilteredByAccountPartialFailure(t *testing.T) {
	accountGood := "61445573-4dbe-4d56-88dc-88ab47aceba7"
	accountBad := "72556684-5ecf-5e67-99ed-99bc58bdfc98"

	manager := &mocks.Manager{}
	manager.On("ListRoleAssignments", accountGood).Return([]*settingsmsg.UserRoleAssignment{
		{Id: "00000000-0000-0000-0000-000000000001", AccountUuid: accountGood, RoleId: "aceb15b8-7486-479f-ae32-c91118e07a39"},
	}, nil)
	manager.On("ListRoleAssignments", accountBad).Return(nil, assert.AnError)
	svc := Service{manager: manager}

	req := &v0.ListRoleAssignmentsFilteredRequest{
		Filters: []*settingsmsg.UserRoleAssignmentFilter{
			{Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ACCOUNT, Term: &settingsmsg.UserRoleAssignmentFilter_AccountUuid{AccountUuid: accountGood}},
			{Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ACCOUNT, Term: &settingsmsg.UserRoleAssignmentFilter_AccountUuid{AccountUuid: accountBad}},
		},
	}
	res := v0.ListRoleAssignmentsResponse{}
	err := svc.ListRoleAssignmentsFiltered(ctxWithUUID, req, &res)
	// one account erroring among several must not fail the whole batch or discard the others' results
	assert.NoError(t, err)
	assert.Len(t, res.Assignments, 1)
	assert.Equal(t, accountGood, res.Assignments[0].GetAccountUuid())
}

func TestListRoleAssignmentsFilteredByRole(t *testing.T) {
	roleID := "61445573-4dbe-4d56-88dc-88ab47aceba7"

	tests := map[string]struct {
		result []*settingsmsg.UserRoleAssignment
		err    error
		status int32
	}{
		"handles manager error": {
			result: nil,
			err:    assert.AnError,
			status: http.StatusNotFound,
		},
		"succeeds with results": {
			result: []*settingsmsg.UserRoleAssignment{
				{
					Id:          "00000000-0000-0000-0000-000000000001",
					AccountUuid: "aceb15b8-7486-479f-ae32-c91118e07a39",
					RoleId:      roleID,
				},
			},
			err:    nil,
			status: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			manager := &mocks.Manager{}
			svc := Service{
				manager: manager,
			}
			manager.On("ListRoleAssignmentsByRole", mock.Anything).Return(test.result, test.err)
			req := &v0.ListRoleAssignmentsFilteredRequest{
				Filters: []*settingsmsg.UserRoleAssignmentFilter{
					{
						Type: settingsmsg.UserRoleAssignmentFilter_TYPE_ROLE,
						Term: &settingsmsg.UserRoleAssignmentFilter_RoleId{
							RoleId: roleID,
						},
					},
				},
			}
			res := v0.ListRoleAssignmentsResponse{}
			err := svc.ListRoleAssignmentsFiltered(ctxWithUUID, req, &res)
			switch test.err {
			case nil:
				assert.Nil(t, err)
			default:
				merr, ok := merrors.As(err)
				assert.True(t, ok)
				assert.Equal(t, int32(test.status), merr.Code)
			}
		})
	}
}
