package trash_test

import (
	"github.com/owncloud/ocis/v2/ocis/pkg/trash"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGatherData(t *testing.T) {
	testcases := []struct {
		Name     string
		Events   []interface{}
		Expected *trash.TrashDirs
	}{
		{
			Name: "no empty dirs",
			Events: []interface{}{
				trashDirs("linkpath", "nodepath"),
			},
			Expected: trashFunc(func(c *trash.TrashDirs) {
				c.LinkPath = "linkpath"
				c.NodePath = "nodepath"
			}),
		},
		{
			Name: "empty dirs",
			Events: []interface{}{
				trashDirs("linkpath", "nodepath"),
			},
			Expected: trashFunc(func(c *trash.TrashDirs) {
			}),
		},
	}

	for _, tc := range testcases {
		events := make(chan interface{})
		go func() {
			for _, e := range tc.Events {
				events <- e
			}
			close(events)
		}()

		td := trash.NewTrashDirs()
		td.GatherData(events)

		require.Equal(t, tc.Expected, td)
	}
}

func trashDirs(linkPath, nodePath string) interface{} {
	return trash.TrashDirs{
		LinkPath: linkPath,
		NodePath: nodePath,
	}
}

func trashFunc(f func(c *trash.TrashDirs)) *trash.TrashDirs {
	c := &trash.TrashDirs{}
	f(c)
	return c
}
