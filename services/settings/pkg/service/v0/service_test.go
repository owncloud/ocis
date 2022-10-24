package svc

import (
	"context"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	v0 "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/settings/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/mock"
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
	manager.On("ListRoleAssignments", mock.Anything).Return(a, nil)
	manager.On("ReadPermissionByID", mock.Anything, mock.Anything).Return(editRolePermission, nil)
	svc := Service{
		manager: manager,
	}

	// Creating an self assignment is expect to fail
	req := v0.AssignRoleToUserRequest{
		AccountUuid: "61445573-4dbe-4d56-88dc-88ab47aceba7",
		RoleId:      "aceb15b8-7486-479f-ae32-c91118e07a39",
	}
	res := v0.AssignRoleToUserResponse{}
	err := svc.AssignRoleToUser(ctxWithUUID, &req, &res)
	assert.NotNil(t, err)

	manager.On("WriteRoleAssignment", mock.Anything, mock.Anything).Return(nil, nil)
	// Creating an self assignment is expect to fail
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
