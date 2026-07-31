package command

import (
	"testing"

	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/reva/v2/pkg/publicshare/manager/json/persistence"
	"github.com/owncloud/reva/v2/pkg/utils"
)

// encodeEntry builds a persistence entry the same way the reva json public-share
// manager does: map{"share": <json-encoded PublicShare>, "password": <string>}.
func encodeEntry(t *testing.T, ps *link.PublicShare) map[string]interface{} {
	t.Helper()
	enc, err := utils.MarshalProtoV1ToJSON(ps)
	if err != nil {
		t.Fatalf("failed to marshal public share: %v", err)
	}
	return map[string]interface{}{
		"share":    string(enc),
		"password": "",
	}
}

func validShare(id, token string) *link.PublicShare {
	return &link.PublicShare{
		Id:    &link.PublicShareId{OpaqueId: id},
		Token: token,
		ResourceId: &provider.ResourceId{
			StorageId: "storage-1",
			SpaceId:   "space-1",
			OpaqueId:  "node-1",
		},
	}
}

func TestFindCorruptPublicShares(t *testing.T) {
	t.Run("nil resource_id is corrupt", func(t *testing.T) {
		ps := validShare("bad-nil", "tok-nil")
		ps.ResourceId = nil
		db := persistence.PublicShares{
			"good":   encodeEntry(t, validShare("good", "tok-good")),
			"bad-nil": encodeEntry(t, ps),
		}

		got := findCorruptPublicShares(db)
		if len(got) != 1 {
			t.Fatalf("expected 1 corrupt entry, got %d (%v)", len(got), got)
		}
		if got[0].ID != "bad-nil" {
			t.Errorf("expected corrupt id 'bad-nil', got %q", got[0].ID)
		}
		if got[0].Token != "tok-nil" {
			t.Errorf("expected token 'tok-nil', got %q", got[0].Token)
		}
	})

	t.Run("empty storage_id is corrupt", func(t *testing.T) {
		ps := validShare("bad-empty", "tok-empty")
		ps.ResourceId = &provider.ResourceId{StorageId: "", SpaceId: "s", OpaqueId: "o"}
		db := persistence.PublicShares{
			"bad-empty": encodeEntry(t, ps),
		}

		got := findCorruptPublicShares(db)
		if len(got) != 1 {
			t.Fatalf("expected 1 corrupt entry, got %d", len(got))
		}
		if got[0].ID != "bad-empty" {
			t.Errorf("expected corrupt id 'bad-empty', got %q", got[0].ID)
		}
	})

	t.Run("undecodable entry is corrupt", func(t *testing.T) {
		db := persistence.PublicShares{
			"not-a-record":  "just a string",
			"missing-share": map[string]interface{}{"password": ""},
			"bad-json":      map[string]interface{}{"share": "{not json", "password": ""},
		}

		got := findCorruptPublicShares(db)
		if len(got) != 3 {
			t.Fatalf("expected 3 corrupt entries, got %d (%v)", len(got), got)
		}
	})

	t.Run("all valid yields none", func(t *testing.T) {
		db := persistence.PublicShares{
			"a": encodeEntry(t, validShare("a", "tok-a")),
			"b": encodeEntry(t, validShare("b", "tok-b")),
		}

		got := findCorruptPublicShares(db)
		if len(got) != 0 {
			t.Fatalf("expected 0 corrupt entries, got %d (%v)", len(got), got)
		}
	})

	t.Run("empty db yields none", func(t *testing.T) {
		if got := findCorruptPublicShares(persistence.PublicShares{}); len(got) != 0 {
			t.Fatalf("expected 0 corrupt entries, got %d", len(got))
		}
	})

	t.Run("only corrupt entries are removed, valid ones preserved", func(t *testing.T) {
		bad := validShare("bad", "tok-bad")
		bad.ResourceId = nil
		db := persistence.PublicShares{
			"keep-1": encodeEntry(t, validShare("keep-1", "tok-1")),
			"bad":    encodeEntry(t, bad),
			"keep-2": encodeEntry(t, validShare("keep-2", "tok-2")),
		}

		for _, f := range findCorruptPublicShares(db) {
			delete(db, f.ID)
		}

		if len(db) != 2 {
			t.Fatalf("expected 2 entries remaining, got %d", len(db))
		}
		if _, ok := db["bad"]; ok {
			t.Error("corrupt entry 'bad' was not removed")
		}
		if _, ok := db["keep-1"]; !ok {
			t.Error("valid entry 'keep-1' was removed")
		}
		if _, ok := db["keep-2"]; !ok {
			t.Error("valid entry 'keep-2' was removed")
		}
	})
}

func TestDecodePublicShare(t *testing.T) {
	t.Run("valid entry decodes", func(t *testing.T) {
		entry := encodeEntry(t, validShare("x", "tok-x"))
		ps, ok := decodePublicShare(entry)
		if !ok {
			t.Fatal("expected decode to succeed")
		}
		if ps.GetToken() != "tok-x" {
			t.Errorf("expected token 'tok-x', got %q", ps.GetToken())
		}
		if ps.GetResourceId().GetStorageId() != "storage-1" {
			t.Errorf("expected storage_id 'storage-1', got %q", ps.GetResourceId().GetStorageId())
		}
	})

	t.Run("non-map entry fails", func(t *testing.T) {
		if _, ok := decodePublicShare("not a map"); ok {
			t.Error("expected decode to fail for non-map entry")
		}
	})

	t.Run("missing share key fails", func(t *testing.T) {
		if _, ok := decodePublicShare(map[string]interface{}{"password": ""}); ok {
			t.Error("expected decode to fail when 'share' key is missing")
		}
	})

	t.Run("invalid json fails", func(t *testing.T) {
		if _, ok := decodePublicShare(map[string]interface{}{"share": "{broken"}); ok {
			t.Error("expected decode to fail for invalid json")
		}
	})
}
