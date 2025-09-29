package shared

import (
	"context"
	"errors"
	"testing"

	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	libregraph "github.com/owncloud/libre-graph-api-go"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func okStatus() *cs3rpc.Status {
	return &cs3rpc.Status{Code: cs3rpc.Code_CODE_OK}
}

func status(code cs3rpc.Code, msg string) *cs3rpc.Status {
	return &cs3rpc.Status{Code: code, Message: msg}
}

func newSpace(id string, trashed bool) *storageprovider.StorageSpace {
	sp := &storageprovider.StorageSpace{
		Id:   &storageprovider.StorageSpaceId{OpaqueId: id},
		Root: nil,
	}
	if trashed {
		sp.Opaque = &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				_spaceStateTrashed: {Decoder: "plain", Value: []byte(_spaceStateTrashed)},
			},
		}
	}
	return sp
}

func TestEnsurePersonalSpace(t *testing.T) {
	t.Run("no-op when personal space exists and is active", func(t *testing.T) {
		gw := cs3mocks.NewGatewayAPIClient(t)
		gw.EXPECT().ListStorageSpaces(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.ListStorageSpacesResponse{
			Status:        okStatus(),
			StorageSpaces: []*storageprovider.StorageSpace{newSpace("ps1", false)},
		}, nil).Once()

		user := libregraph.NewUser("User One", "user1")
		user.SetId("user1")

		err := RestorePersonalSpace(context.Background(), gw, user.GetId())
		require.NoError(t, err)
	})

	t.Run("restores trashed personal space", func(t *testing.T) {
		gw := cs3mocks.NewGatewayAPIClient(t)
		sp := newSpace("ps2", true)

		gw.EXPECT().ListStorageSpaces(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.ListStorageSpacesResponse{
			Status:        okStatus(),
			StorageSpaces: []*storageprovider.StorageSpace{sp},
		}, nil).Once()

		gw.EXPECT().UpdateStorageSpace(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.UpdateStorageSpaceResponse{
			Status: okStatus(),
		}, nil).Once()

		user := libregraph.NewUser("User Two", "user2")
		user.SetId("user2")

		err := RestorePersonalSpace(context.Background(), gw, user.GetId())
		require.NoError(t, err)
	})

	t.Run("no-op when not found", func(t *testing.T) {
		gw := cs3mocks.NewGatewayAPIClient(t)
		gw.EXPECT().ListStorageSpaces(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.ListStorageSpacesResponse{
			Status:        okStatus(),
			StorageSpaces: []*storageprovider.StorageSpace{},
		}, nil).Once()

		user := libregraph.NewUser("User Three", "user3")
		user.SetId("user3")

		err := RestorePersonalSpace(context.Background(), gw, user.GetId())
		require.NoError(t, err)
	})

	t.Run("no-op when not found (empty spaces)", func(t *testing.T) {
		gw := cs3mocks.NewGatewayAPIClient(t)
		gw.EXPECT().ListStorageSpaces(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.ListStorageSpacesResponse{
			Status:        okStatus(),
			StorageSpaces: []*storageprovider.StorageSpace{},
		}, nil).Once()

		user := libregraph.NewUser("User Four", "user4")
		user.SetId("user4")

		err := RestorePersonalSpace(context.Background(), gw, user.GetId())
		require.NoError(t, err)
	})

	t.Run("propagates list error", func(t *testing.T) {
		gw := cs3mocks.NewGatewayAPIClient(t)
		gw.EXPECT().ListStorageSpaces(mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("boom")).Once()

		user := libregraph.NewUser("User Five", "user5")
		user.SetId("user5")

		err := RestorePersonalSpace(context.Background(), gw, user.GetId())
		require.Error(t, err)
	})

	t.Run("propagates non-ok update status", func(t *testing.T) {
		gw := cs3mocks.NewGatewayAPIClient(t)
		sp := newSpace("ps3", true)

		gw.EXPECT().ListStorageSpaces(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.ListStorageSpacesResponse{
			Status:        okStatus(),
			StorageSpaces: []*storageprovider.StorageSpace{sp},
		}, nil).Once()

		gw.EXPECT().UpdateStorageSpace(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.UpdateStorageSpaceResponse{
			Status: status(cs3rpc.Code_CODE_INVALID_ARGUMENT, "nope"),
		}, nil).Once()

		user := libregraph.NewUser("User Six", "user6")
		user.SetId("user6")

		err := RestorePersonalSpace(context.Background(), gw, user.GetId())
		require.Error(t, err)
	})
}

func TestDisablePersonalSpace(t *testing.T) {
	t.Run("no-op when personal space already trashed", func(t *testing.T) {
		gw := cs3mocks.NewGatewayAPIClient(t)
		sp := newSpace("ps4", true)
		gw.EXPECT().ListStorageSpaces(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.ListStorageSpacesResponse{
			Status:        okStatus(),
			StorageSpaces: []*storageprovider.StorageSpace{sp},
		}, nil).Once()

		err := DisablePersonalSpace(context.Background(), gw, "user1")
		require.NoError(t, err)
	})

	t.Run("deletes active personal space", func(t *testing.T) {
		gw := cs3mocks.NewGatewayAPIClient(t)
		sp := newSpace("ps5", false)
		gw.EXPECT().ListStorageSpaces(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.ListStorageSpacesResponse{
			Status:        okStatus(),
			StorageSpaces: []*storageprovider.StorageSpace{sp},
		}, nil).Once()

		gw.EXPECT().DeleteStorageSpace(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.DeleteStorageSpaceResponse{
			Status: okStatus(),
		}, nil).Once()

		err := DisablePersonalSpace(context.Background(), gw, "user1")
		require.NoError(t, err)
	})

	t.Run("propagates non-ok delete status", func(t *testing.T) {
		gw := cs3mocks.NewGatewayAPIClient(t)
		sp := newSpace("ps6", false)
		gw.EXPECT().ListStorageSpaces(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.ListStorageSpacesResponse{
			Status:        okStatus(),
			StorageSpaces: []*storageprovider.StorageSpace{sp},
		}, nil).Once()

		gw.EXPECT().DeleteStorageSpace(mock.Anything, mock.Anything, mock.Anything).Return(&storageprovider.DeleteStorageSpaceResponse{
			Status: status(cs3rpc.Code_CODE_INTERNAL, "fail"),
		}, nil).Once()

		err := DisablePersonalSpace(context.Background(), gw, "user1")
		require.Error(t, err)
	})

	t.Run("propagates list error", func(t *testing.T) {
		gw := cs3mocks.NewGatewayAPIClient(t)
		gw.EXPECT().ListStorageSpaces(mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("boom")).Once()

		err := DisablePersonalSpace(context.Background(), gw, "user1")
		require.Error(t, err)
	})
}
