package proto

import (
	context "context"

	client "github.com/asim/go-micro/v3/client"
	empty "github.com/golang/protobuf/ptypes/empty"
)

// MockAccountsService can be used to write tests
/*
To creata a mock overwrite the functions of an instance like this:

```go
func mockAccSvc(retErr bool) proto.AccountsService {
	if retErr {
		return &proto.MockAccountsService{
			ListFunc: func(ctx context.Context, in *proto.ListAccountsRequest, opts ...client.CallOption) (out *proto.ListAccountsResponse, err error) {
				return nil, fmt.Errorf("error returned by mockAccountsService LIST")
			},
		}
	}

	return &proto.MockAccountsService{
		ListFunc: func(ctx context.Context, in *proto.ListAccountsRequest, opts ...client.CallOption) (out *proto.ListAccountsResponse, err error) {
			return &proto.ListAccountsResponse{
				Accounts: []*proto.Account{
					{
						Id: "yay",
					},
				},
			}, nil
		},
	}
}
```
*/
type MockAccountsService struct {
	ListFunc   func(ctx context.Context, in *ListAccountsRequest, opts ...client.CallOption) (*ListAccountsResponse, error)
	GetFunc    func(ctx context.Context, in *GetAccountRequest, opts ...client.CallOption) (*Account, error)
	CreateFunc func(ctx context.Context, in *CreateAccountRequest, opts ...client.CallOption) (*Account, error)
	UpdateFunc func(ctx context.Context, in *UpdateAccountRequest, opts ...client.CallOption) (*Account, error)
	DeleteFunc func(ctx context.Context, in *DeleteAccountRequest, opts ...client.CallOption) (*empty.Empty, error)
}

// ListAccounts will panic if the function has been called, but not mocked
func (m MockAccountsService) ListAccounts(ctx context.Context, in *ListAccountsRequest, opts ...client.CallOption) (*ListAccountsResponse, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, in, opts...)
	}

	panic("ListFunc was called in test but not mocked")
}

// GetAccount will panic if the function has been called, but not mocked
func (m MockAccountsService) GetAccount(ctx context.Context, in *GetAccountRequest, opts ...client.CallOption) (*Account, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, in, opts...)
	}

	panic("GetFunc was called in test but not mocked")
}

// CreateAccount will panic if the function has been called, but not mocked
func (m MockAccountsService) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...client.CallOption) (*Account, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, in, opts...)
	}

	panic("CreateFunc was called in test but not mocked")
}

// UpdateAccount will panic if the function has been called, but not mocked
func (m MockAccountsService) UpdateAccount(ctx context.Context, in *UpdateAccountRequest, opts ...client.CallOption) (*Account, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, in, opts...)
	}

	panic("UpdateFunc was called in test but not mocked")
}

// DeleteAccount will panic if the function has been called, but not mocked
func (m MockAccountsService) DeleteAccount(ctx context.Context, in *DeleteAccountRequest, opts ...client.CallOption) (*empty.Empty, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, in, opts...)
	}

	panic("DeleteFunc was called in test but not mocked")
}
