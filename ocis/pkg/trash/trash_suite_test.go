package trash_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTrash(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trash Suite")
}
