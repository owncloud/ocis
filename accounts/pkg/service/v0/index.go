package service

import (
	"context"
	"fmt"

	"github.com/owncloud/ocis/accounts/pkg/storage"

	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/indexer"
	"github.com/owncloud/ocis/ocis-pkg/indexer/config"
	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
)

// RebuildIndex deletes all indices (in memory and on storage) and rebuilds them from scratch.
func (s Service) RebuildIndex(ctx context.Context, request *proto.RebuildIndexRequest, response *proto.RebuildIndexResponse) error {
	if err := s.index.Reset(); err != nil {
		return fmt.Errorf("failed to delete index containers: %w", err)
	}

	c, err := configFromSvc(s.Config)
	if err != nil {
		return err
	}
	if err := recreateContainers(s.index, c); err != nil {
		return fmt.Errorf("failed to recreate index containers: %w", err)
	}

	if err := reindexDocuments(ctx, s.repo, s.index); err != nil {
		return fmt.Errorf("failed to reindex documents: %w", err)
	}

	return nil
}

// recreateContainers adds all indices to the indexer that we have for this service.
func recreateContainers(idx *indexer.Indexer, cfg *config.Config) error {
	// Accounts
	if err := idx.AddIndex(&proto.Account{}, "Id", "Id", "accounts", "non_unique", nil, true); err != nil {
		return err
	}

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
	if err := idx.AddIndex(&proto.Group{}, "OnPremisesSamAccountName", "Id", "groups", "unique", nil, false); err != nil {
		return err
	}

	if err := idx.AddIndex(&proto.Group{}, "DisplayName", "Id", "groups", "non_unique", nil, false); err != nil {
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
