package defaults

import "testing"

// TestDefaultConfigStoreNodes fails when the store has no default node.
// A nats-js-kv store without a node cannot connect.
func TestDefaultConfigStoreNodes(t *testing.T) {
	cfg := DefaultConfig()
	if len(cfg.Persistence.Nodes) == 0 {
		t.Fatal("the default persistence store has no node")
	}
}
