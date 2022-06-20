package svc_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGraph(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Graph Suite")
}
