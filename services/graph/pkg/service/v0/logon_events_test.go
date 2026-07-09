package svc_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	mevents "go-micro.dev/v4/events"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
)

// fakeConsumer records whether the service subscribes for events.
type fakeConsumer struct{ calls int }

func (f *fakeConsumer) Consume(_ string, _ ...mevents.ConsumeOption) (<-chan mevents.Event, error) {
	f.calls++
	// never emit, so the dispatch goroutine (if started) stays idle until ctx is cancelled
	return make(chan mevents.Event), nil
}

var _ = Describe("logon event subscription", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		identityBackend *identitymocks.Backend
		consumer        *fakeConsumer
	)

	// newService constructs the graph service, which subscribes for logon events
	// during initialization depending on the UpdateUserLastSignInDate setting.
	newService := func(updateLastSignIn bool) {
		cfg := defaults.FullDefaultConfig()
		cfg.Commons = &shared.Commons{}
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}
		// avoid the LDAP backend construction path; inject a mock instead
		cfg.Identity.LDAP.CACert = ""
		cfg.Identity.LDAP.UpdateUserLastSignInDate = updateLastSignIn

		_, err := service.NewService(
			service.Config(cfg),
			service.Context(ctx),
			service.Logger(log.NewLogger()),
			service.WithIdentityBackend(identityBackend),
			service.EventsConsumer(consumer),
		)
		Expect(err).ToNot(HaveOccurred())
	}

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		identityBackend = &identitymocks.Backend{}
		consumer = &fakeConsumer{}
	})

	AfterEach(func() {
		// stop any dispatch goroutine started during service initialization
		cancel()
	})

	It("does not subscribe for logon events when UpdateUserLastSignInDate is disabled", func() {
		newService(false)
		Expect(consumer.calls).To(Equal(0), "Consume must not be called when UpdateUserLastSignInDate is false")
	})

	It("subscribes for logon events when UpdateUserLastSignInDate is enabled (default)", func() {
		newService(true)
		Expect(consumer.calls).To(Equal(1), "Consume must be called once when UpdateUserLastSignInDate is true")
	})
})
