package config_test

import (
	gofig "github.com/gookit/config/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

var _ = Describe("Environment", func() {
	It("Succeed to parse a comma separated list in to a sting slice", func() {
		cfg := gofig.NewEmpty("test")
		err := cfg.Set("stringlist", "one,two,three")
		Expect(err).To(Not(HaveOccurred()))
		err = cfg.Set("stringlist2", "one ,two , t h r e e")
		Expect(err).To(Not(HaveOccurred()))
		var stringTest, stringTest2 []string
		eb := []shared.EnvBinding{
			{
				EnvVars:     []string{"stringlist"},
				Destination: &stringTest,
			},
			{
				EnvVars:     []string{"stringlist2"},
				Destination: &stringTest2,
			},
		}
		err = ociscfg.BindEnv(cfg, eb)
		Expect(err).To(Not(HaveOccurred()))
		Expect(stringTest).To(Equal([]string{"one", "two", "three"}))
		Expect(stringTest2).To(Equal([]string{"one", "two", "t h r e e"}))

	})
})
