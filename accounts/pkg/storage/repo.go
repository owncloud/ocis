package storage

import (
	"context"

	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
)

const (
	accountsFolder = "accounts"
	groupsFolder   = "groups"
)

type Accounts struct {
	Accounts []*proto.Account
}

type Groups struct {
	Groups []*proto.Group
}

// Repo defines the storage operations
type Repo interface {
	WriteAccount(ctx context.Context, a *proto.Account) (err error)
	LoadAccount(ctx context.Context, id string, a *proto.Account) (err error)
	LoadAccounts(ctx context.Context, a *Accounts) (err error)
	DeleteAccount(ctx context.Context, id string) (err error)
	WriteGroup(ctx context.Context, g *proto.Group) (err error)
	LoadGroup(ctx context.Context, id string, g *proto.Group) (err error)
	LoadGroups(ctx context.Context, g *Groups) (err error)
	DeleteGroup(ctx context.Context, id string) (err error)
}
