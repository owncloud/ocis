package external

import (
	"context"
	"testing"

	"github.com/micro/go-micro/v2/registry"
	"github.com/owncloud/ocis-pkg/v2/log"
)

func TestRegisterGRPCEndpoint(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	err := RegisterGRPCEndpoint(ctx, "test", "1234", "192.168.0.1:777", log.Logger{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	s, err := registry.GetService("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(s) != 1 {
		t.Errorf("Expected exactly one service to be returned got %v", len(s))
	}

	if len(s[0].Nodes) != 1 {
		t.Errorf("Expected exactly one node to be returned got %v", len(s[0].Nodes))
	}

	testSvc := s[0]
	if testSvc.Name != "test" {
		t.Errorf("Expected service name to be 'test' got %v", s[0].Name)
	}

	testNode := testSvc.Nodes[0]

	if testNode.Address != "192.168.0.1:777" {
		t.Errorf("Expected node address to be '192.168.0.1:777' got %v", testNode.Address)
	}

	if testNode.Id != "test-1234" {
		t.Errorf("Expected node id to be 'test-1234' got %v", testNode.Id)
	}

	cancel()

	s, err = registry.GetService("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(s) != 0 {
		t.Errorf("Deregister on cancelation failed. Result-length should be zero, got %v", len(s))
	}
}
