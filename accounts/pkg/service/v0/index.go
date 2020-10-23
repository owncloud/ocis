package service

import (
	"context"

	"github.com/owncloud/ocis/accounts/pkg/storage"

	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/indexer"
	"github.com/owncloud/ocis/accounts/pkg/indexer/option"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
)

// RebuildIndex deletes all indices (in memory and on storage) and rebuilds them from scratch.
func (s Service) RebuildIndex(ctx context.Context, request *proto.RebuildIndexRequest, response *proto.RebuildIndexResponse) error {
	if err := s.index.Reset(); err != nil {
		return err
	}

	if err := recreateContainers(s.index, s.Config); err != nil {
		return err
	}

	return reindexDocuments(ctx, s.repo, s.index)
}

// recreateContainers adds all indices to the indexer that we have for this service.
func recreateContainers(idx *indexer.Indexer, cfg *config.Config) error {
	// Accounts
	if err := idx.AddIndex(&proto.Account{}, "DisplayName", "Id", "accounts", "non_unique", nil, true); err != nil {
		return err
	}
	if err := idx.AddIndex(&proto.Account{}, "Mail", "Id", "accounts", "unique", nil, true); err != nil {
		return err
	}

	if err := idx.AddIndex(&proto.Account{}, "OnPremisesSamAccountName", "Id", "accounts", "unique", nil, true); err != nil {
		return err
	}

	if err := idx.AddIndex(&proto.Account{}, "PreferredName", "Id", "accounts", "unique", nil, true); err != nil {
		return err
	}

	if err := idx.AddIndex(&proto.Account{}, "UidNumber", "Id", "accounts", "autoincrement", &option.Bound{
		Lower: cfg.Index.UID.Lower,
		Upper: cfg.Index.UID.Upper,
	}, false); err != nil {
		return err
	}

	// Groups
	if err := idx.AddIndex(&proto.Group{}, "OnPremisesSamAccountName", "Id", "groups", "unique", nil, true); err != nil {
		return err
	}

	if err := idx.AddIndex(&proto.Group{}, "DisplayName", "Id", "groups", "non_unique", nil, true); err != nil {
		return err
	}

	if err := idx.AddIndex(&proto.Group{}, "GidNumber", "Id", "groups", "autoincrement", &option.Bound{
		Lower: cfg.Index.GID.Lower,
		Upper: cfg.Index.GID.Upper,
	}, false); err != nil {
		return err
	}

	return nil
}

// reindexDocuments loads all existing documents and adds them to the index.
func reindexDocuments(ctx context.Context, repo storage.Repo, index *indexer.Indexer) error {
	accounts := make([]*proto.Account, 0)
	if err := repo.LoadAccounts(ctx, &accounts); err != nil {
		return err
	}
	for i := range accounts {
		_, err := index.Add(accounts[i])
		if err != nil {
			return err
		}
	}

	groups := make([]*proto.Group, 0)
	if err := repo.LoadGroups(ctx, &groups); err != nil {
		return err
	}
	for i := range groups {
		_, err := index.Add(groups[i])
		if err != nil {
			return err
		}
	}
	return nil
}
