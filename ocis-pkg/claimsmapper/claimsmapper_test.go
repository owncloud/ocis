package claimsmapper

import (
	"testing"

	"github.com/test-go/testify/require"
)

func TestMapClaims(t *testing.T) {
	type innercase struct {
		input           string
		expectedNoMatch bool
		expectedSpaceID string
		expectedRole    string
	}
	var testCases = []struct {
		regexp  string
		mapping []string
		cases   []innercase
	}{
		{
			regexp: "some-string:moreinfo:([a-zA-Z0-9-]+):and-the-role-is-(.*)",
			cases: []innercase{
				{
					input:           "some-string:moreinfo:here-is-my-uuid:and-the-role-is-here",
					expectedSpaceID: "here-is-my-uuid",
					expectedRole:    "here",
				},
				{
					input:           "some-otherthing:moreinfo:here-is-my-uuid:and-the-role-is-not-here",
					expectedNoMatch: true,
				},
			},
		},
		{
			regexp:  "spaceid=([a-zA-Z0-9-]+),roleid=(.*)",
			mapping: []string{"overseer:manager", "worker:editor", "ghoul:viewer"},
			cases: []innercase{
				{
					input:           "spaceid=vault36,roleid=overseer",
					expectedSpaceID: "vault36",
					expectedRole:    "manager",
				},
				{
					input:           "spaceid=vault36,roleid=worker",
					expectedSpaceID: "vault36",
					expectedRole:    "editor",
				},
				{
					input:           "spaceid=vault36,roleid=ghoul",
					expectedSpaceID: "vault36",
					expectedRole:    "viewer",
				},
				{
					input:           "spaceid=vault36,roleid=radroach",
					expectedNoMatch: true,
				},
				{
					input:           "differentid=vault36,roleid=overseer",
					expectedNoMatch: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		cm := NewClaimsMapper(tc.regexp, tc.mapping)
		for _, c := range tc.cases {
			match, spaceID, role := cm.Exec(c.input)
			require.Equal(t, !c.expectedNoMatch, match, c.input)
			require.Equal(t, c.expectedSpaceID, spaceID, c.input)
			require.Equal(t, c.expectedRole, role, c.input)
		}

	}
}
