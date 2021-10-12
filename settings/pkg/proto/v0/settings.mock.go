package proto

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"go-micro.dev/v4/client"
)

// MockBundleService can be used to write tests against the bundle service.
/*
To create a mock overwrite the functions of an instance like this:

```go
func mockBundleSvc(returnErr bool) proto.BundleService {
	if returnErr {
		return &proto.MockBundleService{
			ListBundlesFunc: func(ctx context.Context, in *proto.ListBundlesRequest, opts ...client.CallOption) (out *proto.ListBundlesResponse, err error) {
				return nil, fmt.Errorf("error returned by mockBundleSvc LIST")
			},
		}
	}

	return &proto.MockBundleService{
		ListBundlesFunc: func(ctx context.Context, in *proto.ListBundlesRequest, opts ...client.CallOption) (out *proto.ListBundlesResponse, err error) {
			return &proto.ListBundlesResponse{
				Bundles: []*proto.Bundle{
					{
						Id: "hello-there",
					},
				},
			}, nil
		},
	}
}
```
*/
type MockBundleService struct {
	ListBundlesFunc             func(ctx context.Context, req *ListBundlesRequest, opts ...client.CallOption) (*ListBundlesResponse, error)
	GetBundleFunc               func(ctx context.Context, req *GetBundleRequest, opts ...client.CallOption) (*GetBundleResponse, error)
	SaveBundleFunc              func(ctx context.Context, req *SaveBundleRequest, opts ...client.CallOption) (*SaveBundleResponse, error)
	AddSettingToBundleFunc      func(ctx context.Context, req *AddSettingToBundleRequest, opts ...client.CallOption) (*AddSettingToBundleResponse, error)
	RemoveSettingFromBundleFunc func(ctx context.Context, req *RemoveSettingFromBundleRequest, opts ...client.CallOption) (*empty.Empty, error)
}

// ListBundles will panic if the function has been called, but not mocked
func (m MockBundleService) ListBundles(ctx context.Context, req *ListBundlesRequest, opts ...client.CallOption) (*ListBundlesResponse, error) {
	if m.ListBundlesFunc != nil {
		return m.ListBundlesFunc(ctx, req, opts...)
	}
	panic("ListBundlesFunc was called in test but not mocked")
}

// GetBundle will panic if the function has been called, but not mocked
func (m MockBundleService) GetBundle(ctx context.Context, req *GetBundleRequest, opts ...client.CallOption) (*GetBundleResponse, error) {
	if m.GetBundleFunc != nil {
		return m.GetBundleFunc(ctx, req, opts...)
	}
	panic("GetBundleFunc was called in test but not mocked")
}

// SaveBundle will panic if the function has been called, but not mocked
func (m MockBundleService) SaveBundle(ctx context.Context, req *SaveBundleRequest, opts ...client.CallOption) (*SaveBundleResponse, error) {
	if m.SaveBundleFunc != nil {
		return m.SaveBundleFunc(ctx, req, opts...)
	}
	panic("SaveBundleFunc was called in test but not mocked")
}

// AddSettingToBundle will panic if the function has been called, but not mocked
func (m MockBundleService) AddSettingToBundle(ctx context.Context, req *AddSettingToBundleRequest, opts ...client.CallOption) (*AddSettingToBundleResponse, error) {
	if m.AddSettingToBundleFunc != nil {
		return m.AddSettingToBundleFunc(ctx, req, opts...)
	}
	panic("AddSettingToBundleFunc was called in test but not mocked")
}

// RemoveSettingFromBundle will panic if the function has been called, but not mocked
func (m MockBundleService) RemoveSettingFromBundle(ctx context.Context, req *RemoveSettingFromBundleRequest, opts ...client.CallOption) (*empty.Empty, error) {
	if m.RemoveSettingFromBundleFunc != nil {
		return m.RemoveSettingFromBundleFunc(ctx, req, opts...)
	}
	panic("RemoveSettingFromBundleFunc was called in test but not mocked")
}

// MockValueService can be used to write tests against the value service.
type MockValueService struct {
	ListValuesFunc                  func(ctx context.Context, req *ListValuesRequest, opts ...client.CallOption) (*ListValuesResponse, error)
	GetValueFunc                    func(ctx context.Context, req *GetValueRequest, opts ...client.CallOption) (*GetValueResponse, error)
	GetValueByUniqueIdentifiersFunc func(ctx context.Context, req *GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*GetValueResponse, error)
	SaveValueFunc                   func(ctx context.Context, req *SaveValueRequest, opts ...client.CallOption) (*SaveValueResponse, error)
}

// ListValues will panic if the function has been called, but not mocked
func (m MockValueService) ListValues(ctx context.Context, req *ListValuesRequest, opts ...client.CallOption) (*ListValuesResponse, error) {
	if m.ListValuesFunc != nil {
		return m.ListValuesFunc(ctx, req, opts...)
	}
	panic("ListValuesFunc was called in test but not mocked")
}

// GetValue will panic if the function has been called, but not mocked
func (m MockValueService) GetValue(ctx context.Context, req *GetValueRequest, opts ...client.CallOption) (*GetValueResponse, error) {
	if m.GetValueFunc != nil {
		return m.GetValueFunc(ctx, req, opts...)
	}
	panic("GetValueFunc was called in test but not mocked")
}

// GetValueByUniqueIdentifiers will panic if the function has been called, but not mocked
func (m MockValueService) GetValueByUniqueIdentifiers(ctx context.Context, req *GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*GetValueResponse, error) {
	if m.GetValueByUniqueIdentifiersFunc != nil {
		return m.GetValueByUniqueIdentifiersFunc(ctx, req, opts...)
	}
	panic("GetValueByUniqueIdentifiersFunc was called in test but not mocked")
}

// SaveValue will panic if the function has been called, but not mocked
func (m MockValueService) SaveValue(ctx context.Context, req *SaveValueRequest, opts ...client.CallOption) (*SaveValueResponse, error) {
	if m.SaveValueFunc != nil {
		return m.SaveValueFunc(ctx, req, opts...)
	}
	panic("SaveValueFunc was called in test but not mocked")
}

// MockRoleService will panic if the function has been called, but not mocked
type MockRoleService struct {
	ListRolesFunc           func(ctx context.Context, req *ListBundlesRequest, opts ...client.CallOption) (*ListBundlesResponse, error)
	ListRoleAssignmentsFunc func(ctx context.Context, req *ListRoleAssignmentsRequest, opts ...client.CallOption) (*ListRoleAssignmentsResponse, error)
	AssignRoleToUserFunc    func(ctx context.Context, req *AssignRoleToUserRequest, opts ...client.CallOption) (*AssignRoleToUserResponse, error)
	RemoveRoleFromUserFunc  func(ctx context.Context, req *RemoveRoleFromUserRequest, opts ...client.CallOption) (*empty.Empty, error)
}

// ListRoles will panic if the function has been called, but not mocked
func (m MockRoleService) ListRoles(ctx context.Context, req *ListBundlesRequest, opts ...client.CallOption) (*ListBundlesResponse, error) {
	if m.ListRolesFunc != nil {
		return m.ListRolesFunc(ctx, req, opts...)
	}
	panic("ListRolesFunc was called in test but not mocked")
}

// ListRoleAssignments will panic if the function has been called, but not mocked
func (m MockRoleService) ListRoleAssignments(ctx context.Context, req *ListRoleAssignmentsRequest, opts ...client.CallOption) (*ListRoleAssignmentsResponse, error) {
	if m.ListRoleAssignmentsFunc != nil {
		return m.ListRoleAssignmentsFunc(ctx, req, opts...)
	}
	panic("ListRoleAssignmentsFunc was called in test but not mocked")
}

// AssignRoleToUser will panic if the function has been called, but not mocked
func (m MockRoleService) AssignRoleToUser(ctx context.Context, req *AssignRoleToUserRequest, opts ...client.CallOption) (*AssignRoleToUserResponse, error) {
	if m.AssignRoleToUserFunc != nil {
		return m.AssignRoleToUserFunc(ctx, req, opts...)
	}
	panic("AssignRoleToUserFunc was called in test but not mocked")
}

// RemoveRoleFromUser will panic if the function has been called, but not mocked
func (m MockRoleService) RemoveRoleFromUser(ctx context.Context, req *RemoveRoleFromUserRequest, opts ...client.CallOption) (*empty.Empty, error) {
	if m.RemoveRoleFromUserFunc != nil {
		return m.RemoveRoleFromUserFunc(ctx, req, opts...)
	}
	panic("RemoveRoleFromUserFunc was called in test but not mocked")
}

// MockPermissionService will panic if the function has been called, but not mocked
type MockPermissionService struct {
	ListPermissionsByResourceFunc func(ctx context.Context, req *ListPermissionsByResourceRequest, opts ...client.CallOption) (*ListPermissionsByResourceResponse, error)
	GetPermissionByIDFunc         func(ctx context.Context, req *GetPermissionByIDRequest, opts ...client.CallOption) (*GetPermissionByIDResponse, error)
}

// ListPermissionsByResource will panic if the function has been called, but not mocked
func (m MockPermissionService) ListPermissionsByResource(ctx context.Context, req *ListPermissionsByResourceRequest, opts ...client.CallOption) (*ListPermissionsByResourceResponse, error) {
	if m.ListPermissionsByResourceFunc != nil {
		return m.ListPermissionsByResourceFunc(ctx, req, opts...)
	}
	panic("ListPermissionsByResourceFunc was called in test but not mocked")
}

// GetPermissionByID will panic if the function has been called, but not mocked
func (m MockPermissionService) GetPermissionByID(ctx context.Context, req *GetPermissionByIDRequest, opts ...client.CallOption) (*GetPermissionByIDResponse, error) {
	if m.GetPermissionByIDFunc != nil {
		return m.GetPermissionByIDFunc(ctx, req, opts...)
	}
	panic("GetPermissionByIDFunc was called in test but not mocked")
}
