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

// fakeConsumer records whether StartListenForLogonEvents subscribes for events.
type fakeConsumer struct{ calls int }

func (f *fakeConsumer) Consume(_ string, _ ...mevents.ConsumeOption) (<-chan mevents.Event, error) {
	f.calls++
	// never emit, so the dispatch goroutine (if started) stays idle
	return make(chan mevents.Event), nil
}

var _ = Describe("StartListenForLogonEvents", func() {
	var (
		ctx             context.Context
		identityBackend *identitymocks.Backend
		consumer        *fakeConsumer
	)

	newService := func(updateLastSignIn bool) service.Graph {
		cfg := defaults.FullDefaultConfig()
		cfg.Commons = &shared.Commons{}
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}
		// avoid the LDAP backend construction path; inject a mock instead
		cfg.Identity.LDAP.CACert = ""
		cfg.Identity.LDAP.UpdateUserLastSignInDate = updateLastSignIn

		svc, err := service.NewService(
			service.Config(cfg),
			service.Context(ctx),
			service.Logger(log.NewLogger()),
			service.WithIdentityBackend(identityBackend),
			service.EventsConsumer(consumer),
		)
		Expect(err).ToNot(HaveOccurred())
		return svc
	}

	BeforeEach(func() {
		ctx = context.Background()
		identityBackend = &identitymocks.Backend{}
		consumer = &fakeConsumer{}
	})

	It("does not subscribe for logon events when UpdateUserLastSignInDate is disabled", func() {
		svc := newService(false)
		Expect(svc.StartListenForLogonEvents(ctx, log.NewLogger())).To(Succeed())
		Expect(consumer.calls).To(Equal(0), "Consume must not be called when UpdateUserLastSignInDate is false")
	})

	It("subscribes for logon events when UpdateUserLastSignInDate is enabled (default)", func() {
		svc := newService(true)
		Expect(svc.StartListenForLogonEvents(ctx, log.NewLogger())).To(Succeed())
		Expect(consumer.calls).To(BeNumerically(">=", 1), "Consume must be called when UpdateUserLastSignInDate is true")
	})
})
