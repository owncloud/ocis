package engine_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pMessage "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/policies/v0"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine"
)

var _ = Describe("Engine", func() {
	DescribeTable("NewEnvironmentFromPB",
		func(incomingStage pMessage.Stage, outgoinStage engine.Stage) {
			pEnv := &pMessage.Environment{
				Stage: incomingStage,
			}

			env, err := engine.NewEnvironmentFromPB(pEnv)
			Expect(err).ToNot(HaveOccurred())

			Expect(env.Stage).To(Equal(outgoinStage))
		},
		Entry("http stage", pMessage.Stage_STAGE_HTTP, engine.StageHTTP),
		Entry("pp stage", pMessage.Stage_STAGE_PP, engine.StagePP),
	)
})
