package utf7_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtf7(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utf7 Suite")
}
