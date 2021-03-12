package cs3

//
//import (
//	"os"
//	"testing"
//
//	"github.com/owncloud/ocis/accounts/pkg/config"
//	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
//	. "github.com/owncloud/ocis/ocis-pkg/indexer/test"
//	"github.com/stretchr/testify/assert"
//)
//
//const cs3RootFolder = "/tmp/ocis/storage/users/data"
//
//func TestAutoincrementIndexAdd(t *testing.T) {
//	dataDir, err := WriteIndexTestData(Data, "ID", cs3RootFolder)
//	assert.NoError(t, err)
//	cfg := generateConfig()
//
//	sut := NewAutoincrementIndex(
//		option.WithTypeName(GetTypeFQN(User{})),
//		option.WithIndexBy("UID"),
//		option.WithDataURL(cfg.Repo.CS3.DataURL),
//		option.WithDataPrefix(cfg.Repo.CS3.DataPrefix),
//		option.WithJWTSecret(cfg.Repo.CS3.JWTSecret),
//		option.WithProviderAddr(cfg.Repo.CS3.ProviderAddr),
//	)
//
//	assert.NoError(t, sut.Init())
//
//	for i := 0; i < 5; i++ {
//		res, err := sut.Add("abcdefg-123", "")
//		assert.NoError(t, err)
//		t.Log(res)
//	}
//
//	_ = os.RemoveAll(dataDir)
//}
//
//func BenchmarkAutoincrementIndexAdd(b *testing.B) {
//	dataDir, err := WriteIndexTestData(Data, "ID", cs3RootFolder)
//	assert.NoError(b, err)
//	cfg := generateConfig()
//
//	sut := NewAutoincrementIndex(
//		option.WithTypeName(GetTypeFQN(User{})),
//		option.WithIndexBy("UID"),
//		option.WithDataURL(cfg.Repo.CS3.DataURL),
//		option.WithDataPrefix(cfg.Repo.CS3.DataPrefix),
//		option.WithJWTSecret(cfg.Repo.CS3.JWTSecret),
//		option.WithProviderAddr(cfg.Repo.CS3.ProviderAddr),
//	)
//
//	err = sut.Init()
//	assert.NoError(b, err)
//
//	for n := 0; n < b.N; n++ {
//		_, err := sut.Add("abcdefg-123", "")
//		if err != nil {
//			b.Error(err)
//		}
//		assert.NoError(b, err)
//	}
//
//	_ = os.RemoveAll(dataDir)
//}
//
//func generateConfig() config.Config {
//	return config.Config{
//		Repo: config.Repo{
//			Disk: config.Disk{
//				Path: "",
//			},
//			CS3: config.CS3{
//				ProviderAddr: "0.0.0.0:9215",
//				DataURL:      "http://localhost:9216",
//				DataPrefix:   "data",
//				JWTSecret:    "Pive-Fumkiu4",
//			},
//		},
//	}
//}
