package store

import (
	"testing"

	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/stretchr/testify/require"
)

func TestPermission(t *testing.T) {
	s := initStore()
	setupRoles(s)
	// bunldes are initialized within init func
	p, err := s.ReadPermissionByID("readID", []string{"f36db5e6-a03c-40df-8413-711c67e40b47"})
	require.NoError(t, err)
	require.Equal(t, settingsmsg.Permission_OPERATION_READ, p.Operation)

	p, err = s.ReadPermissionByName("read", []string{"f36db5e6-a03c-40df-8413-711c67e40b47"})
	require.NoError(t, err)
	require.Equal(t, settingsmsg.Permission_OPERATION_READ, p.Operation)

	pms, err := s.ListPermissionsByResource(&settingsmsg.Resource{
		Type: settingsmsg.Resource_TYPE_BUNDLE,
	}, []string{"f36db5e6-a03c-40df-8413-711c67e40b47"})
	require.NoError(t, err)
	require.Len(t, pms, 1)
	require.Equal(t, settingsmsg.Permission_OPERATION_READ, pms[0].Operation)

}
