package svc

import (
	"testing"
	"time"

	"github.com/CiscoM31/godata"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/stretchr/testify/assert"
)

type sortTest struct {
	Drives       []*libregraph.Drive
	Query        godata.GoDataRequest
	DrivesSorted []*libregraph.Drive
}

var Time1 = time.Date(2022, 02, 02, 15, 00, 00, 00, time.UTC)
var Time2 = time.Date(2022, 02, 03, 15, 00, 00, 00, time.UTC)
var Time3 *time.Time
var Time4 = time.Date(2022, 02, 05, 15, 00, 00, 00, time.UTC)
var Drives = []*libregraph.Drive{
	{Id: libregraph.PtrString("1"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Einstein"), LastModifiedDateTime: &Time1},
	{Id: libregraph.PtrString("2"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Marie"), LastModifiedDateTime: &Time2},
	{Id: libregraph.PtrString("3"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Admin"), LastModifiedDateTime: Time3},
	{Id: libregraph.PtrString("4"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Richard"), LastModifiedDateTime: &Time4},
}

var sortTests = []sortTest{
	{
		Drives: Drives,
		Query: godata.GoDataRequest{
			Query: &godata.GoDataQuery{
				OrderBy: &godata.GoDataOrderByQuery{
					OrderByItems: []*godata.OrderByItem{
						{Field: &godata.Token{Value: "name"}, Order: "asc"},
					},
				},
			},
		},
		DrivesSorted: []*libregraph.Drive{
			{Id: libregraph.PtrString("3"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Admin"), LastModifiedDateTime: Time3},
			{Id: libregraph.PtrString("1"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Einstein"), LastModifiedDateTime: &Time1},
			{Id: libregraph.PtrString("2"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Marie"), LastModifiedDateTime: &Time2},
			{Id: libregraph.PtrString("4"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Richard"), LastModifiedDateTime: &Time4},
		},
	},
	{
		Drives: Drives,
		Query: godata.GoDataRequest{
			Query: &godata.GoDataQuery{
				OrderBy: &godata.GoDataOrderByQuery{
					OrderByItems: []*godata.OrderByItem{
						{Field: &godata.Token{Value: "name"}, Order: "desc"},
					},
				},
			},
		},
		DrivesSorted: []*libregraph.Drive{
			{Id: libregraph.PtrString("4"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Richard"), LastModifiedDateTime: &Time4},
			{Id: libregraph.PtrString("2"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Marie"), LastModifiedDateTime: &Time2},
			{Id: libregraph.PtrString("1"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Einstein"), LastModifiedDateTime: &Time1},
			{Id: libregraph.PtrString("3"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Admin"), LastModifiedDateTime: Time3},
		},
	},
	{
		Drives: Drives,
		Query: godata.GoDataRequest{
			Query: &godata.GoDataQuery{
				OrderBy: &godata.GoDataOrderByQuery{
					OrderByItems: []*godata.OrderByItem{
						{Field: &godata.Token{Value: "lastModifiedDateTime"}, Order: "asc"},
					},
				},
			},
		},
		DrivesSorted: []*libregraph.Drive{
			{Id: libregraph.PtrString("3"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Admin"), LastModifiedDateTime: Time3},
			{Id: libregraph.PtrString("1"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Einstein"), LastModifiedDateTime: &Time1},
			{Id: libregraph.PtrString("2"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Marie"), LastModifiedDateTime: &Time2},
			{Id: libregraph.PtrString("4"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Richard"), LastModifiedDateTime: &Time4},
		},
	},
	{
		Drives: Drives,
		Query: godata.GoDataRequest{
			Query: &godata.GoDataQuery{
				OrderBy: &godata.GoDataOrderByQuery{
					OrderByItems: []*godata.OrderByItem{
						{Field: &godata.Token{Value: "lastModifiedDateTime"}, Order: "desc"},
					},
				},
			},
		},
		DrivesSorted: []*libregraph.Drive{
			{Id: libregraph.PtrString("4"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Richard"), LastModifiedDateTime: &Time4},
			{Id: libregraph.PtrString("2"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Marie"), LastModifiedDateTime: &Time2},
			{Id: libregraph.PtrString("1"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Einstein"), LastModifiedDateTime: &Time1},
			{Id: libregraph.PtrString("3"), DriveType: libregraph.PtrString("project"), Name: libregraph.PtrString("Admin"), LastModifiedDateTime: Time3},
		},
	},
}

func TestSort(t *testing.T) {
	graph := Graph{
		config:               nil,
		mux:                  nil,
		logger:               nil,
		identityBackend:      nil,
		gatewayClient:        nil,
		httpClient:           nil,
		spacePropertiesCache: nil,
	}
	for _, test := range sortTests {
		sorted := graph.sortSpaces(&test.Query, test.Drives)
		assert.Equal(t, test.DrivesSorted, sorted)
	}
}
