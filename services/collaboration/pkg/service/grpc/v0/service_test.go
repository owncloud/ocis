package service_test

import (
	"context"
	"strconv"
	"time"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/reva/v2/pkg/rgrpc/status"
	"github.com/owncloud/reva/v2/pkg/utils"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/collaboration/mocks"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	service "github.com/owncloud/ocis/v2/services/collaboration/pkg/service/grpc/v0"
)

// Based on https://github.com/cs3org/reva/blob/b99ad4865401144a981d4cfd1ae28b5a018ea51d/pkg/token/manager/jwt/jwt.go#L82
func MintToken(u *userv1beta1.User, secret string, nowTime time.Time) string {
	scopes := make(map[string]*authpb.Scope)
	scopes["user"] = &authpb.Scope{
		Resource: &types.OpaqueEntry{
			Decoder: "json",
			Value:   []byte("{\"Path\":\"/\"}"),
		},
		Role: authpb.Role_ROLE_OWNER,
	}

	claims := jwt.MapClaims{
		"exp":   nowTime.Add(5 * time.Hour).Unix(),
		"iss":   "myself",
		"aud":   "reva",
		"iat":   nowTime.Unix(),
		"user":  u,
		"scope": scopes,
	}
	/*
		claims := claims{
			StandardClaims: jwt.RegisteredClaims{
				ExpiresAt: time.Now().Add(5 * time.Hour),
				Issuer:    "myself",
				Audience:  "reva",
				IssuedAt:  time.Now(),
			},
			User:  u,
			Scope: scopes,
		}
	*/

	t := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	tkn, _ := t.SignedString([]byte(secret))

	return tkn
}

var _ = Describe("Discovery", func() {
	var (
		cfg           *config.Config
		gatewayClient *cs3mocks.GatewayAPIClient
		srv           *service.Service
		srvTear       func()
	)

	BeforeEach(func() {
		cfg = &config.Config{
			TokenManager: &config.TokenManager{
				JWTSecret: "topsecret",
			},
		}
		gatewayClient = &cs3mocks.GatewayAPIClient{}

		gatewaySelector := mocks.NewSelectable[gatewayv1beta1.GatewayAPIClient](GinkgoT())
		gatewaySelector.On("Next").Return(gatewayClient, nil)

		srv, srvTear, _ = service.NewHandler(
			service.Logger(log.NopLogger()),
			service.Config(cfg),
			service.AppURLs(map[string]map[string]string{
				"view": {
					".pdf":  "https://test.server.prv/hosting/wopi/word/view",
					".djvu": "https://test.server.prv/hosting/wopi/word/view",
					".docx": "https://test.server.prv/hosting/wopi/word/view",
					".xls":  "https://test.server.prv/hosting/wopi/cell/view",
					".xlsb": "https://test.server.prv/hosting/wopi/cell/view",
				},
				"edit": {
					".docx":    "https://test.server.prv/hosting/wopi/word/edit",
					".invalid": "://test.server.prv/hosting/wopi/cell/edit",
				},
			}),
			service.GatewaySelector(gatewaySelector),
		)
	})

	AfterEach(func() {
		srvTear()
	})

	Describe("OpenInApp", func() {
		It("Invalid access token", func() {
			ctx := context.Background()

			cfg.Wopi.WopiSrc = "https://wopiserver.test.prv"

			req := &appproviderv1beta1.OpenInAppRequest{
				ResourceInfo: &providerv1beta1.ResourceInfo{
					Id: &providerv1beta1.ResourceId{
						StorageId: "myStorage",
						OpaqueId:  "storageOpaque001",
						SpaceId:   "SpaceA",
					},
					Path: "/path/to/file",
				},
				ViewMode:    appproviderv1beta1.ViewMode_VIEW_MODE_READ_WRITE,
				AccessToken: "goodAccessToken",
			}

			gatewayClient.On("WhoAmI", mock.Anything, mock.Anything).Times(1).Return(&gatewayv1beta1.WhoAmIResponse{
				Status: status.NewOK(ctx),
				User: &userv1beta1.User{
					Id: &userv1beta1.UserId{
						Idp:      "myIdp",
						OpaqueId: "opaque001",
						Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
					},
					Username: "username",
				},
			}, nil)

			resp, err := srv.OpenInApp(ctx, req)
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})

		DescribeTable(
			"Success",
			func(appName, path, lang string, disableChat bool, expectedAppUrl string) {
				ctx := context.Background()
				nowTime := time.Now()

				cfg.Wopi.WopiSrc = "https://wopiserver.test.prv"
				cfg.Wopi.Secret = "topsecret"
				cfg.Wopi.DisableChat = disableChat
				cfg.App.Name = appName
				cfg.App.Product = appName

				myself := &userv1beta1.User{
					Id: &userv1beta1.UserId{
						Idp:      "myIdp",
						OpaqueId: "opaque001",
						Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
					},
					Username: "username",
				}

				if path == "" {
					path = "/path/to/file.docx"
				}
				req := &appproviderv1beta1.OpenInAppRequest{
					ResourceInfo: &providerv1beta1.ResourceInfo{
						Id: &providerv1beta1.ResourceId{
							StorageId: "myStorage",
							OpaqueId:  "storageOpaque001",
							SpaceId:   "SpaceA",
						},
						Path: path,
					},
					ViewMode:    appproviderv1beta1.ViewMode_VIEW_MODE_READ_WRITE,
					AccessToken: MintToken(myself, cfg.Wopi.Secret, nowTime),
				}
				if lang != "" {
					req.Opaque = utils.AppendPlainToOpaque(req.Opaque, "lang", lang)
				}

				gatewayClient.On("WhoAmI", mock.Anything, mock.Anything).Times(1).Return(&gatewayv1beta1.WhoAmIResponse{
					Status: status.NewOK(ctx),
					User:   myself,
				}, nil)

				resp, err := srv.OpenInApp(ctx, req)
				Expect(err).To(Succeed())
				Expect(resp.GetStatus().GetCode()).To(Equal(rpcv1beta1.Code_CODE_OK))
				Expect(resp.GetAppUrl().GetMethod()).To(Equal("POST"))
				Expect(resp.GetAppUrl().GetAppUrl()).To(Equal(expectedAppUrl))
				Expect(resp.GetAppUrl().GetFormParameters()["access_token_ttl"]).To(Equal(strconv.FormatInt(nowTime.Add(5*time.Hour).Unix()*1000, 10)))
			},
			Entry("Microsoft chat no lang", "Microsoft", "", "", false, "https://test.server.prv/hosting/wopi/word/edit?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e"),
			Entry("Collabora chat no lang", "Collabora", "", "", false, "https://test.server.prv/hosting/wopi/word/view?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e"),
			Entry("OnlyOffice chat no lang", "OnlyOffice", "", "", false, "https://test.server.prv/hosting/wopi/word/edit?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e"),
			Entry("Microsoft chat lang", "Microsoft", "", "de", false, "https://test.server.prv/hosting/wopi/word/edit?UI_LLCC=de-DE&WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e"),
			Entry("Collabora chat lang", "Collabora", "", "de", false, "https://test.server.prv/hosting/wopi/word/view?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&lang=de-DE"),
			Entry("OnlyOffice chat lang", "OnlyOffice", "", "de", false, "https://test.server.prv/hosting/wopi/word/edit?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&ui=de-DE"),
			Entry("Microsoft no chat no lang", "Microsoft", "", "", true, "https://test.server.prv/hosting/wopi/word/edit?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&dchat=1"),
			Entry("Collabora no chat no lang", "Collabora", "", "", true, "https://test.server.prv/hosting/wopi/word/view?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&dchat=1"),
			Entry("OnlyOffice no chat no lang", "OnlyOffice", "", "", true, "https://test.server.prv/hosting/wopi/word/edit?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&dchat=1"),
			Entry("Microsoft no chat lang", "Microsoft", "", "de", true, "https://test.server.prv/hosting/wopi/word/edit?UI_LLCC=de-DE&WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&dchat=1"),
			Entry("Collabora no chat lang", "Collabora", "", "de", true, "https://test.server.prv/hosting/wopi/word/view?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&dchat=1&lang=de-DE"),
			Entry("OnlyOffice no chat lang", "OnlyOffice", "", "de", true, "https://test.server.prv/hosting/wopi/word/edit?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&dchat=1&ui=de-DE"),
			Entry("Microsoft no chat lang", "Microsoft", "/path/to/file.DOCX", "de", true, "https://test.server.prv/hosting/wopi/word/edit?UI_LLCC=de-DE&WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&dchat=1"),
			Entry("Collabora no chat lang", "Collabora", "/path/to/file.DOCX", "de", true, "https://test.server.prv/hosting/wopi/word/view?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&dchat=1&lang=de-DE"),
			Entry("OnlyOffice no chat lang", "OnlyOffice", "/path/to/file.DOCX", "de", true, "https://test.server.prv/hosting/wopi/word/edit?WOPISrc=https%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&dchat=1&ui=de-DE"),
		)
		It("Success with Wopi Proxy", func() {
			ctx := context.Background()
			nowTime := time.Now()

			cfg.Wopi.WopiSrc = "https://wopiserver.test.prv"
			cfg.Wopi.Secret = "topsecret"
			cfg.Wopi.ProxyURL = "https://office.proxy.test.prv"
			cfg.Wopi.ProxySecret = "your_supa_secret"
			cfg.App.Name = "Microsoft"

			myself := &userv1beta1.User{
				Id: &userv1beta1.UserId{
					Idp:      "myIdp",
					OpaqueId: "opaque001",
					Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
				},
				Username: "username",
			}

			req := &appproviderv1beta1.OpenInAppRequest{
				ResourceInfo: &providerv1beta1.ResourceInfo{
					Id: &providerv1beta1.ResourceId{
						StorageId: "myStorage",
						OpaqueId:  "storageOpaque001",
						SpaceId:   "SpaceA",
					},
					Path: "/path/to/file.docx",
				},
				ViewMode:    appproviderv1beta1.ViewMode_VIEW_MODE_READ_WRITE,
				AccessToken: MintToken(myself, cfg.Wopi.Secret, nowTime),
			}
			req.Opaque = utils.AppendPlainToOpaque(req.Opaque, "lang", "en")

			gatewayClient.On("WhoAmI", mock.Anything, mock.Anything).Times(1).Return(&gatewayv1beta1.WhoAmIResponse{
				Status: status.NewOK(ctx),
				User:   myself,
			}, nil)

			resp, err := srv.OpenInApp(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp.GetStatus().GetCode()).To(Equal(rpcv1beta1.Code_CODE_OK))
			Expect(resp.GetAppUrl().GetMethod()).To(Equal("POST"))
			Expect(resp.GetAppUrl().GetAppUrl()).To(Equal("https://test.server.prv/hosting/wopi/word/edit?UI_LLCC=en-US&WOPISrc=https%3A%2F%2Foffice.proxy.test.prv%2Fwopi%2Ffiles%2FeyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1IjoiaHR0cHM6Ly93b3Bpc2VydmVyLnRlc3QucHJ2L3dvcGkvZmlsZXMvIiwiZiI6IjJmNmVjMTg2OTZkZDEwMDgxMDY3NDliZDk0MTA2ZTVjZmFkNWMwOWUxNWRlN2I3NzA4OGQwMzg0M2U3MWI0M2UifQ.yfyLHZ18Z1MFOa6u7AP0LqfIiQ9X5AMkYauEZGhbCNs"))
			Expect(resp.GetAppUrl().GetFormParameters()["access_token_ttl"]).To(Equal(strconv.FormatInt(nowTime.Add(5*time.Hour).Unix()*1000, 10)))
		})
		It("Fail with invalid app url", func() {
			ctx := context.Background()
			nowTime := time.Now()

			cfg.Wopi.WopiSrc = "htttps://wopiserver.test.prv"
			cfg.Wopi.Secret = "topsecret"
			cfg.App.Name = "Microsoft"

			myself := &userv1beta1.User{
				Id: &userv1beta1.UserId{
					Idp:      "myIdp",
					OpaqueId: "opaque001",
					Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
				},
				Username: "username",
			}

			req := &appproviderv1beta1.OpenInAppRequest{
				ResourceInfo: &providerv1beta1.ResourceInfo{
					Id: &providerv1beta1.ResourceId{
						StorageId: "myStorage",
						OpaqueId:  "storageOpaque001",
						SpaceId:   "SpaceA",
					},
					Path: "/path/to/file.invalid",
				},
				ViewMode:    appproviderv1beta1.ViewMode_VIEW_MODE_READ_WRITE,
				AccessToken: MintToken(myself, cfg.Wopi.Secret, nowTime),
			}
			req.Opaque = utils.AppendPlainToOpaque(req.Opaque, "lang", "en")

			gatewayClient.On("WhoAmI", mock.Anything, mock.Anything).Times(1).Return(&gatewayv1beta1.WhoAmIResponse{
				Status: status.NewOK(ctx),
				User:   myself,
			}, nil)

			resp, err := srv.OpenInApp(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp.GetStatus().GetCode()).To(Equal(rpcv1beta1.Code_CODE_INVALID_ARGUMENT))
			Expect(resp.GetStatus().GetMessage()).To(Equal("OpenInApp: error parsing appUrl"))
		})
		It("Fail with invalid template id", func() {
			ctx := context.Background()
			nowTime := time.Now()

			cfg.Wopi.WopiSrc = "htttps://wopiserver.test.prv"
			cfg.Wopi.Secret = "topsecret"
			cfg.App.Name = "Microsoft"

			myself := &userv1beta1.User{
				Id: &userv1beta1.UserId{
					Idp:      "myIdp",
					OpaqueId: "opaque001",
					Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
				},
				Username: "username",
			}

			req := &appproviderv1beta1.OpenInAppRequest{
				ResourceInfo: &providerv1beta1.ResourceInfo{
					Id: &providerv1beta1.ResourceId{
						StorageId: "myStorage",
						OpaqueId:  "storageOpaque001",
						SpaceId:   "SpaceA",
					},
					Path: "/path/to/file.docx",
				},
				ViewMode:    appproviderv1beta1.ViewMode_VIEW_MODE_READ_WRITE,
				AccessToken: MintToken(myself, cfg.Wopi.Secret, nowTime),
			}
			req.Opaque = utils.AppendPlainToOpaque(req.Opaque, "lang", "en")
			req.Opaque = utils.AppendPlainToOpaque(req.Opaque, "template", "&file_id")

			gatewayClient.On("WhoAmI", mock.Anything, mock.Anything).Times(1).Return(&gatewayv1beta1.WhoAmIResponse{
				Status: status.NewOK(ctx),
				User:   myself,
			}, nil)

			resp, err := srv.OpenInApp(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp.GetStatus().GetCode()).To(Equal(rpcv1beta1.Code_CODE_INVALID_ARGUMENT))
			Expect(resp.GetStatus().GetMessage()).To(Equal("OpenInApp: error parsing templateID"))
		})
		It("Success with valid template id", func() {
			ctx := context.Background()
			nowTime := time.Now()

			cfg.Wopi.WopiSrc = "htttps://wopiserver.test.prv"
			cfg.Wopi.Secret = "topsecret"
			cfg.App.Name = "OnlyOffice"
			cfg.App.Product = "OnlyOffice"

			myself := &userv1beta1.User{
				Id: &userv1beta1.UserId{
					Idp:      "myIdp",
					OpaqueId: "opaque001",
					Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
				},
				Username: "username",
			}

			req := &appproviderv1beta1.OpenInAppRequest{
				ResourceInfo: &providerv1beta1.ResourceInfo{
					Id: &providerv1beta1.ResourceId{
						StorageId: "myStorage",
						OpaqueId:  "storageOpaque001",
						SpaceId:   "SpaceA",
					},
					Path: "/path/to/file.docx",
				},
				ViewMode:    appproviderv1beta1.ViewMode_VIEW_MODE_READ_WRITE,
				AccessToken: MintToken(myself, cfg.Wopi.Secret, nowTime),
			}
			req.Opaque = utils.AppendPlainToOpaque(req.Opaque, "lang", "en")
			req.Opaque = utils.AppendPlainToOpaque(req.Opaque, "template", "prodiderID$spaceID!opaqueID")

			gatewayClient.On("WhoAmI", mock.Anything, mock.Anything).Times(1).Return(&gatewayv1beta1.WhoAmIResponse{
				Status: status.NewOK(ctx),
				User:   myself,
			}, nil)

			resp, err := srv.OpenInApp(ctx, req)
			Expect(err).To(Succeed())
			Expect(resp.GetStatus().GetCode()).To(Equal(rpcv1beta1.Code_CODE_OK))
			Expect(resp.GetAppUrl().GetMethod()).To(Equal("POST"))
			Expect(resp.GetAppUrl().GetAppUrl()).To(Equal("https://test.server.prv/hosting/wopi/word/edit?WOPISrc=htttps%3A%2F%2Fwopiserver.test.prv%2Fwopi%2Ffiles%2F2f6ec18696dd1008106749bd94106e5cfad5c09e15de7b77088d03843e71b43e&ui=en-US"))
			Expect(resp.GetAppUrl().GetFormParameters()["access_token_ttl"]).To(Equal(strconv.FormatInt(nowTime.Add(5*time.Hour).Unix()*1000, 10)))
		})
	})
})
