package svc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/go-chi/chi/v5"
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

type schoolList struct {
	Value []*libregraph.EducationSchool
}

var _ = Describe("Schools", func() {
	var (
		svc                      service.Service
		ctx                      context.Context
		cfg                      *config.Config
		gatewayClient            *mocks.GatewayClient
		identityEducationBackend *identitymocks.EducationBackend

		rr *httptest.ResponseRecorder

		newSchool   *libregraph.EducationSchool
		currentUser = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				OpaqueId: "user",
			},
		}
	)

	BeforeEach(func() {

		identityEducationBackend = &identitymocks.EducationBackend{}
		gatewayClient = &mocks.GatewayClient{}
		newSchool = libregraph.NewEducationSchool()
		newSchool.SetId("school1")

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
			service.WithIdentityEducationBackend(identityEducationBackend),
		)
	})

	Describe("GetSchools", func() {
		It("handles invalid ODATA parameters", func() {
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/education/schools?§foo=bar", nil)
			svc.GetSchools(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles invalid sorting queries", func() {
			identityEducationBackend.On("GetSchools", ctx, mock.Anything).Return([]*libregraph.EducationSchool{newSchool}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/education/schools?$orderby=invalid", nil)
			svc.GetSchools(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			odataerr := libregraph.OdataError{}
			err = json.Unmarshal(data, &odataerr)
			Expect(err).ToNot(HaveOccurred())
			Expect(odataerr.Error.Code).To(Equal("invalidRequest"))
		})

		It("handles unknown backend errors", func() {
			identityEducationBackend.On("GetSchools", ctx, mock.Anything).Return(nil, errors.New("failed"))

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/education/schools", nil)
			svc.GetSchools(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			odataerr := libregraph.OdataError{}
			err = json.Unmarshal(data, &odataerr)
			Expect(err).ToNot(HaveOccurred())
			Expect(odataerr.Error.Code).To(Equal("generalException"))
		})

		It("handles backend errors", func() {
			identityEducationBackend.On("GetSchools", ctx, mock.Anything).Return(nil, errorcode.New(errorcode.AccessDenied, "access denied"))

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/education/schools", nil)
			svc.GetSchools(rr, r)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			odataerr := libregraph.OdataError{}
			err = json.Unmarshal(data, &odataerr)
			Expect(err).ToNot(HaveOccurred())
			Expect(odataerr.Error.Code).To(Equal("accessDenied"))
		})

		It("renders an empty list of schools", func() {
			identityEducationBackend.On("GetSchools", ctx, mock.Anything).Return([]*libregraph.EducationSchool{}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/education/schools", nil)
			svc.GetSchools(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := service.ListResponse{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Value).To(Equal([]interface{}{}))
		})

		It("renders a list of schools", func() {
			identityEducationBackend.On("GetSchools", ctx, mock.Anything).Return([]*libregraph.EducationSchool{newSchool}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/education/schools", nil)
			svc.GetSchools(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := schoolList{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))
			Expect(res.Value[0].GetId()).To(Equal("school1"))
		})
	})

	Describe("GetSchool", func() {
		It("handles missing or empty school id", func() {
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/education/schools", nil)
			svc.GetSchool(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))

			r = httptest.NewRequest(http.MethodGet, "/graph/v1.0/education/schools", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", "")
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, nil), chi.RouteCtxKey, rctx))
			svc.GetSchool(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		Context("with an existing school", func() {
			BeforeEach(func() {
				identityEducationBackend.On("GetSchool", mock.Anything, mock.Anything, mock.Anything).Return(newSchool, nil)
			})

			It("gets the school", func() {
				r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/education/schools/"+*newSchool.Id, nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("schoolID", *newSchool.Id)
				r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, nil), chi.RouteCtxKey, rctx))

				svc.GetSchool(rr, r)

				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("PostSchool", func() {
		It("handles invalid body", func() {
			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/education/schools/", bytes.NewBufferString("{invalid"))

			svc.PostSchool(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles missing display name", func() {
			newSchool = libregraph.NewEducationSchool()
			newSchool.SetSchoolNumber("0034")
			newSchoolJson, err := json.Marshal(newSchool)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/education/schools/", bytes.NewBuffer(newSchoolJson))

			svc.PostSchool(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles missing school number", func() {
			newSchool = libregraph.NewEducationSchool()
			newSchool.SetDisplayName("New School")
			newSchoolJson, err := json.Marshal(newSchool)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/education/schools/", bytes.NewBuffer(newSchoolJson))

			svc.PostSchool(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("disallows school create ids", func() {
			newSchool = libregraph.NewEducationSchool()
			newSchool.SetId("disallowed")
			newSchool.SetDisplayName("New School")
			newSchool.SetSchoolNumber("0034")
			newSchoolJson, err := json.Marshal(newSchool)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/education/schools/", bytes.NewBuffer(newSchoolJson))

			svc.PostSchool(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles backend errors", func() {
			newSchool = libregraph.NewEducationSchool()
			newSchool.SetDisplayName("New School")
			newSchool.SetSchoolNumber("0034")
			newSchoolJson, err := json.Marshal(newSchool)
			Expect(err).ToNot(HaveOccurred())

			identityEducationBackend.On("CreateSchool", mock.Anything, mock.Anything).Return(nil, errorcode.New(errorcode.AccessDenied, "access denied"))

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/education/schools/", bytes.NewBuffer(newSchoolJson))

			svc.PostSchool(rr, r)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("creates the school", func() {
			newSchool = libregraph.NewEducationSchool()
			newSchool.SetDisplayName("New School")
			newSchool.SetSchoolNumber("0034")
			newSchoolJson, err := json.Marshal(newSchool)
			Expect(err).ToNot(HaveOccurred())

			identityEducationBackend.On("CreateSchool", mock.Anything, mock.Anything).Return(newSchool, nil)

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/education/schools/", bytes.NewBuffer(newSchoolJson))

			svc.PostSchool(rr, r)

			Expect(rr.Code).To(Equal(http.StatusCreated))
		})
	})
	Describe("PatchSchool", func() {
		It("handles invalid body", func() {
			r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/education/schools/", bytes.NewBufferString("{invalid"))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", *newSchool.Id)
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PatchSchool(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles missing or empty school id", func() {
			r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/education/schools", nil)
			svc.PatchSchool(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))

			r = httptest.NewRequest(http.MethodPatch, "/graph/v1.0/education/schools", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", "")
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PatchSchool(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles malformed school id", func() {
			r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/education/schools", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", "school%id")
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PatchSchool(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("updates the school", func() {
			newSchool = libregraph.NewEducationSchool()
			newSchool.SetDisplayName("New School Name")
			newSchoolJson, err := json.Marshal(newSchool)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/education/schools/schoolid", bytes.NewBuffer(newSchoolJson))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", "school-id")
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))

			svc.PatchSchool(rr, r)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
		})
	})

	Describe("DeleteSchool", func() {
		Context("with an existing school", func() {
			BeforeEach(func() {
				identityEducationBackend.On("GetSchool", mock.Anything, mock.Anything, mock.Anything).Return(newSchool, nil)
			})
		})

		It("deletes the school", func() {
			identityEducationBackend.On("DeleteSchool", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/education/schools", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", *newSchool.Id)
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.DeleteSchool(rr, r)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
			identityEducationBackend.AssertNumberOfCalls(GinkgoT(), "DeleteSchool", 1)
		})
	})

	Describe("GetSchoolMembers", func() {
		It("gets the list of members", func() {
			user := libregraph.NewEducationUser()
			user.SetId("user")
			identityEducationBackend.On("GetSchoolMembers", mock.Anything, mock.Anything, mock.Anything).Return([]*libregraph.EducationUser{user}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/education/schools/{schoolID}/members", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", *newSchool.Id)
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetSchoolMembers(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			var members []*libregraph.User
			err = json.Unmarshal(data, &members)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(members)).To(Equal(1))
			Expect(members[0].GetId()).To(Equal("user"))
		})
	})

	Describe("PostSchoolMembers", func() {
		It("fails on invalid body", func() {
			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/education/schools/{schoolID}/members", bytes.NewBufferString("{invalid"))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", *newSchool.Id)
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PostSchoolMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("fails on missing member refs", func() {
			member := libregraph.NewMemberReference()
			data, err := json.Marshal(member)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/education/schools/{schoolID}/members", bytes.NewBuffer(data))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", *newSchool.Id)
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PostSchoolMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("fails on invalid member refs", func() {
			member := libregraph.NewMemberReference()
			member.SetOdataId("/invalidtype/user")
			data, err := json.Marshal(member)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/education/schools/{schoolID}/members", bytes.NewBuffer(data))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", *newSchool.Id)
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PostSchoolMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("adds a new member", func() {
			member := libregraph.NewMemberReference()
			member.SetOdataId("/users/user")
			data, err := json.Marshal(member)
			Expect(err).ToNot(HaveOccurred())
			identityEducationBackend.On("AddMembersToSchool", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/education/schools/{schoolID}/members", bytes.NewBuffer(data))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", *newSchool.Id)
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PostSchoolMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNoContent))

			identityEducationBackend.AssertNumberOfCalls(GinkgoT(), "AddMembersToSchool", 1)
		})
	})

	Describe("DeleteSchoolMembers", func() {
		It("handles missing or empty member id", func() {
			r := httptest.NewRequest(http.MethodDelete, "/graph/v1.0/education/schools/{schoolID}/members/{memberID}/$ref", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", *newSchool.Id)
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.DeleteSchoolMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
		It("handles missing or empty member id", func() {
			r := httptest.NewRequest(http.MethodDelete, "/graph/v1.0/education/schools/{schoolID}/members/{memberID}/$ref", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("memberID", "/users/user")
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.DeleteSchoolMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("deletes members", func() {
			identityEducationBackend.On("RemoveMemberFromSchool", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			r := httptest.NewRequest(http.MethodDelete, "/graph/v1.0/education/schools/{schoolID}/members/{memberID}/$ref", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("schoolID", *newSchool.Id)
			rctx.URLParams.Add("memberID", "/users/user1")
			r = r.WithContext(context.WithValue(ctxpkg.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.DeleteSchoolMember(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNoContent))

			identityEducationBackend.AssertNumberOfCalls(GinkgoT(), "RemoveMemberFromSchool", 1)
		})
	})
})
