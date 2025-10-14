package wopisrc_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWopisrc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wopisrc Suite")
}
