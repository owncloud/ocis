package content_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Content Suite")
}
