package storage

// Uncomment to test locally, requires started metadata-storage for now

//import (
//	"context"
//	"github.com/owncloud/ocis/accounts/pkg/config"
//	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//var cfg = &config.Config{
//	TokenManager: config.TokenManager{
//		JWTSecret: "Pive-Fumkiu4",
//	},
//	Repo: config.Repo{
//		CS3:  config.CS3{
//			ProviderAddr: "0.0.0.0:9215",
//		},
//	},
//}
//
//func TestCS3Repo_WriteAccount(t *testing.T) {
//	r, err := NewCS3Repo("hello", cfg)
//	assert.NoError(t, err)
//
//	err = r.WriteAccount(context.Background(), &proto.Account{
//		Id:             "fefef-egegweg-gegeg",
//		AccountEnabled: true,
//		DisplayName:    "Mike Jones",
//		Mail:           "mike@example.com",
//	})
//
//	assert.NoError(t, err)
//}
//
//func TestCS3Repo_LoadAccount(t *testing.T) {
//	r, err := NewCS3Repo("hello", cfg)
//	assert.NoError(t, err)
//
//	err = r.WriteAccount(context.Background(), &proto.Account{
//		Id:             "fefef-egegweg-gegeg",
//		AccountEnabled: true,
//		DisplayName:    "Mike Jones",
//		Mail:           "mike@example.com",
//	})
//
//	acc := &proto.Account{}
//	err = r.LoadAccount(context.Background(), "fefef-egegweg-gegeg", acc)
//
//	assert.NoError(t, err)
//	assert.Equal(t, "fefef-egegweg-gegeg", acc.Id)
//	assert.Equal(t, "Mike Jones", acc.DisplayName)
//	assert.Equal(t, "mike@example.com", acc.Mail)
//}
//
//func TestCS3Repo_DeleteAccount(t *testing.T) {
//	r, err := NewCS3Repo("hello", cfg)
//	assert.NoError(t, err)
//
//	err = r.WriteAccount(context.Background(), &proto.Account{
//		Id:             "delete-me-id",
//		AccountEnabled: true,
//		DisplayName:    "Mike Jones",
//		Mail:           "mike@example.com",
//	})
//
//	err = r.DeleteAccount(context.Background(), "delete-me-id")
//
//	assert.NoError(t, err)
//}
