package storage

import (
	"context"

	accountsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/accounts/v0"
)

const (
	accountsFolder = "accounts"
	groupsFolder   = "groups"
)

// Repo defines the storage operations
type Repo interface {
	WriteAccount(ctx context.Context, a *accountsmsg.Account) (err error)
	LoadAccount(ctx context.Context, id string, a *accountsmsg.Account) (err error)
	LoadAccounts(ctx context.Context, a *[]*accountsmsg.Account) (err error)
	DeleteAccount(ctx context.Context, id string) (err error)
	WriteGroup(ctx context.Context, g *accountsmsg.Group) (err error)
	LoadGroup(ctx context.Context, id string, g *accountsmsg.Group) (err error)
	LoadGroups(ctx context.Context, g *[]*accountsmsg.Group) (err error)
	DeleteGroup(ctx context.Context, id string) (err error)
}
