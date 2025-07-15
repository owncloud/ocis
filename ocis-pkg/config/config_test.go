package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Config", func() {
	It("Success generating the default config", func() {
		cfg := config.DefaultConfig()
		_, err := yaml.Marshal(cfg)
		Expect(err).To(BeNil())
	})
})
