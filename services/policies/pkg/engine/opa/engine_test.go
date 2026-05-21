package opa

import (
	"context"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine"
)

var _ = Describe("OPA Evaluate cache", func() {
	var (
		o       OPA
		regoDir string
		query   string
	)

	BeforeEach(func() {
		var err error
		regoDir = GinkgoT().TempDir()

		err = os.WriteFile(filepath.Join(regoDir, "test.rego"), []byte(`
package test

import future.keywords.if

default granted := true

granted = false if {
	input.stage == "block"
}
`), 0600)
		Expect(err).ToNot(HaveOccurred())

		query = "data.test.granted"
		logger := log.NewLogger()

		o, err = NewOPA(10*time.Second, logger, config.Engine{
			Policies: []string{filepath.Join(regoDir, "test.rego")},
		})
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns true for an allowed request", func() {
		result, err := o.Evaluate(context.Background(), query, engine.Environment{})
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeTrue())
	})

	It("returns false for a blocked request", func() {
		result, err := o.Evaluate(context.Background(), query, engine.Environment{Stage: "block"})
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeFalse())
	})

	It("populates the cache after the first call", func() {
		Expect(o.cache).To(BeEmpty())
		_, err := o.Evaluate(context.Background(), query, engine.Environment{})
		Expect(err).ToNot(HaveOccurred())
		Expect(o.cache).To(HaveLen(1))
	})

	It("does not grow the cache on repeated calls with the same query", func() {
		_, _ = o.Evaluate(context.Background(), query, engine.Environment{})
		_, _ = o.Evaluate(context.Background(), query, engine.Environment{})
		Expect(o.cache).To(HaveLen(1))
	})

	It("creates separate cache entries for different query strings", func() {
		err := os.WriteFile(filepath.Join(regoDir, "other.rego"), []byte(`
package other
default granted := true
`), 0600)
		Expect(err).ToNot(HaveOccurred())

		_, _ = o.Evaluate(context.Background(), query, engine.Environment{})
		_, _ = o.Evaluate(context.Background(), "data.other.granted", engine.Environment{})
		Expect(o.cache).To(HaveLen(2))
	})
})
