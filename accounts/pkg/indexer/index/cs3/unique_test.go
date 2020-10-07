package cs3

import (
	"context"
	"flag"
	"os"
	"path"
	"testing"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/indexer/option"
	. "github.com/owncloud/ocis/accounts/pkg/indexer/test"
	"github.com/owncloud/ocis/storage/pkg/command"
	mcfg "github.com/owncloud/ocis/storage/pkg/config"
	"github.com/stretchr/testify/assert"
)

var (
	ctx, cancelFunc = context.WithCancel(context.Background())
)

func setupMetadataStorage() {
	cfg := mcfg.New()
	app := cli.App{
		Name:     "storage-metadata-for-tests",
		Commands: []*cli.Command{command.StorageMetadata(cfg)},
	}

	_ = app.Command("storage-metadata").Run(cli.NewContext(&app, &flag.FlagSet{}, &cli.Context{Context: ctx}))
}

func TestCS3UniqueIndex_FakeSymlink(t *testing.T) {
	go setupMetadataStorage()
	defer cancelFunc()

	dataDir := WriteIndexTestDataCS3(t, TestData, "Id")
	cfg := config.Config{
		Repo: config.Repo{
			Disk: config.Disk{
				Path: "",
			},
			CS3: config.CS3{
				ProviderAddr: "0.0.0.0:9215",
				DataURL:      "http://localhost:9216",
				DataPrefix:   "data",
				JWTSecret:    "Pive-Fumkiu4",
			},
		},
	}

	sut := NewUniqueIndexWithOptions(
		option.WithTypeName(GetTypeFQN(TestUser{})),
		option.WithIndexBy("UserName"),
		option.WithFilesDir(path.Join(cfg.Repo.Disk.Path, "/meta")),
		option.WithDataDir(cfg.Repo.Disk.Path),
		option.WithDataURL(cfg.Repo.CS3.DataURL),
		option.WithDataPrefix(cfg.Repo.CS3.DataPrefix),
		option.WithJWTSecret(cfg.Repo.CS3.JWTSecret),
		option.WithProviderAddr(cfg.Repo.CS3.ProviderAddr),
	)

	err := sut.Init()
	assert.NoError(t, err)

	res, err := sut.Add("abcdefg-123", "mikey")
	assert.NoError(t, err)
	t.Log(res)

	resLookup, err := sut.Lookup("mikey")
	assert.NoError(t, err)
	t.Log(resLookup)

	err = sut.Update("abcdefg-123", "mikey", "mickeyX")
	assert.NoError(t, err)

	searchRes, err := sut.Search("m*")
	assert.NoError(t, err)
	assert.Len(t, searchRes, 1)
	assert.Equal(t, searchRes[0], "abcdefg-123")

	_ = os.RemoveAll(dataDir)
}

func TestCS3UniqueIndexSearch(t *testing.T) {
	go setupMetadataStorage()
	defer cancelFunc()

	dataDir := WriteIndexTestDataCS3(t, TestData, "Id")
	cfg := config.Config{
		Repo: config.Repo{
			Disk: config.Disk{
				Path: "",
			},
			CS3: config.CS3{
				ProviderAddr: "0.0.0.0:9215",
				DataURL:      "http://localhost:9216",
				DataPrefix:   "data",
				JWTSecret:    "Pive-Fumkiu4",
			},
		},
	}

	sut := NewUniqueIndexWithOptions(
		option.WithTypeName(GetTypeFQN(TestUser{})),
		option.WithIndexBy("UserName"),
		option.WithFilesDir(path.Join(cfg.Repo.Disk.Path, "/meta")),
		option.WithDataDir(cfg.Repo.Disk.Path),
		option.WithDataURL(cfg.Repo.CS3.DataURL),
		option.WithDataPrefix(cfg.Repo.CS3.DataPrefix),
		option.WithJWTSecret(cfg.Repo.CS3.JWTSecret),
		option.WithProviderAddr(cfg.Repo.CS3.ProviderAddr),
	)

	err := sut.Init()
	assert.NoError(t, err)

	_, err = sut.Add("hijklmn-456", "mikey")
	assert.NoError(t, err)

	_, err = sut.Add("ewf4ofk-555", "jacky")
	assert.NoError(t, err)

	res, err := sut.Search("*y")
	assert.NoError(t, err)
	t.Log(res)

	_ = os.RemoveAll(dataDir)
}
