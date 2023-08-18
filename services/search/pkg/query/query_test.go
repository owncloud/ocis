package query_test

import (
	"strings"
	"testing"

	"github.com/owncloud/ocis/v2/services/search/pkg/query"
)

func TestKqlToBleveQuery(t *testing.T) {
	want := "tag:ownCloud"
	r := strings.NewReader("tag:ownCloud")
	w := &strings.Builder{}

	if err := query.KqlToBleveQuery(r, w); err != nil {
		t.Fatal(err)
	}

	if w.String() != want {
		t.Fatalf("Compile mismatch \ngot: `%s` \nwant: `%s`", w.String(), want)
	}
}
