package svc_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/test-go/testify/mock"

	libregraph "github.com/owncloud/libre-graph-api-go"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
)

type groupList struct {
	Value []*libregraph.Group
}

var _ = Describe("Groups", func() {
	var (
		svc             service.Service
		ctx             context.Context
		cfg             *config.Config
		gatewayClient   *mocks.GatewayClient
		eventsPublisher mocks.Publisher
		identityBackend *identitymocks.Backend

		rr *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		identityBackend = &identitymocks.Backend{}
		rr = httptest.NewRecorder()
		ctx = context.Background()

		cfg = defaults.FullDefaultConfig()
		cfg.Identity.LDAP.CACert = "" // skip the startup checks, we don't use LDAP at all in this tests
		cfg.TokenManager.JWTSecret = "loremipsum"
		cfg.Commons = &shared.Commons{}
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}

		_ = ogrpc.Configure(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
		svc = service.NewService(
			service.Config(cfg),
			service.WithGatewayClient(gatewayClient),
			service.EventsPublisher(&eventsPublisher),
			service.WithIdentityBackend(identityBackend),
		)
	})

	Describe("GetGroups", func() {
		It("handles invalid ODATA parameters", func() {
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/groups?§foo=bar", nil)
			svc.GetGroups(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles unknown backend errors", func() {
			identityBackend.On("GetGroups", ctx, mock.Anything).Return(nil, errors.New("failed"))

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/groups", nil)
			svc.GetGroups(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			data, err := ioutil.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			odataerr := libregraph.OdataError{}
			err = json.Unmarshal(data, &odataerr)
			Expect(err).ToNot(HaveOccurred())
			Expect(odataerr.Error.Code).To(Equal("generalException"))
		})

		It("handles backend errors", func() {
			identityBackend.On("GetGroups", ctx, mock.Anything).Return(nil, errorcode.New(errorcode.AccessDenied, "access denied"))

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/groups", nil)
			svc.GetGroups(rr, r)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			data, err := ioutil.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			odataerr := libregraph.OdataError{}
			err = json.Unmarshal(data, &odataerr)
			Expect(err).ToNot(HaveOccurred())
			Expect(odataerr.Error.Code).To(Equal("accessDenied"))
		})

		It("renders an empty list of groups", func() {
			identityBackend.On("GetGroups", ctx, mock.Anything).Return([]*libregraph.Group{}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/groups", nil)
			svc.GetGroups(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := ioutil.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := service.ListResponse{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Value).To(Equal([]interface{}{}))
		})

		It("renders a list of groups", func() {
			group1 := libregraph.NewGroup()
			group1.SetId("group1")
			identityBackend.On("GetGroups", ctx, mock.Anything).Return([]*libregraph.Group{group1}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/groups", nil)
			svc.GetGroups(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := ioutil.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := groupList{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))
			Expect(res.Value[0].GetId()).To(Equal("group1"))
		})

		It("handles invalid sorting queries", func() {
			group1 := libregraph.NewGroup()
			group1.SetId("group1")
			identityBackend.On("GetGroups", ctx, mock.Anything).Return([]*libregraph.Group{group1}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/groups?$orderby=invalid", nil)
			svc.GetGroups(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			data, err := ioutil.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			odataerr := libregraph.OdataError{}
			err = json.Unmarshal(data, &odataerr)
			Expect(err).ToNot(HaveOccurred())
			Expect(odataerr.Error.Code).To(Equal("invalidRequest"))
		})
	})
})
