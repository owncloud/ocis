package http

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"reflect"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	"github.com/owncloud/ocis/v2/ocis-pkg/cors"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	ohttp "github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/settings"
	"go-micro.dev/v4"
	"go-micro.dev/v4/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server initializes the http service and server.
func Server(opts ...Option) (ohttp.Service, error) {
	options := newOptions(opts...)

	service, err := ohttp.NewService(
		ohttp.TLSConfig(options.Config.HTTP.TLS),
		ohttp.Logger(options.Logger),
		ohttp.Name(options.Name),
		ohttp.Version(version.GetString()),
		ohttp.Address(options.Config.HTTP.Addr),
		ohttp.Namespace(options.Config.HTTP.Namespace),
		ohttp.Context(options.Context),
		ohttp.Flags(options.Flags...),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return ohttp.Service{}, fmt.Errorf("could not initialize http service: %w", err)
	}

	handle := options.ServiceHandler

	mux := chi.NewMux()

	mux.Use(middleware.GetOtelhttpMiddleware(options.Name, options.TraceProvider))
	mux.Use(chimiddleware.RealIP)
	mux.Use(chimiddleware.RequestID)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.Cors(
		cors.Logger(options.Logger),
		cors.AllowedOrigins(options.Config.HTTP.CORS.AllowedOrigins),
		cors.AllowedMethods(options.Config.HTTP.CORS.AllowedMethods),
		cors.AllowedHeaders(options.Config.HTTP.CORS.AllowedHeaders),
		cors.AllowCredentials(options.Config.HTTP.CORS.AllowCredentials),
	))
	mux.Use(middleware.ExtractAccountUUID(
		account.Logger(options.Logger),
		account.JWTSecret(options.Config.TokenManager.JWTSecret)),
	)

	mux.Use(middleware.Version(
		options.Name,
		version.GetString(),
	))

	mux.Use(middleware.Logger(
		options.Logger,
	))

	mux.Route(options.Config.HTTP.Root, func(r chi.Router) {
		registerHandlers(r, handle)
	})

	_ = chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		options.Logger.Debug().Str("method", method).Str("route", route).Int("middlewares", len(middlewares)).Msg("serving endpoint")
		return nil
	})

	micro.RegisterHandler(service.Server(), mux)

	return service, nil
}

type handlerFunc func(ctx context.Context, req interface{}, resp interface{}) error

func registerPostHandler(r chi.Router, path string, reqProto interface{}, respProto interface{}, handler handlerFunc) {
	r.Post(path, func(w http.ResponseWriter, r *http.Request) {
		req := reflect.New(reflect.TypeOf(reqProto).Elem()).Interface()
		resp := reflect.New(reflect.TypeOf(respProto).Elem()).Interface()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
			return
		}

		if err := protojson.Unmarshal(body, req.(proto.Message)); err != nil {
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
			return
		}

		if err := handler(r.Context(), req, resp); err != nil {
			handleError(w, err)
			return
		}

		// Check if response is empty (like emptypb.Empty)
		if reflect.TypeOf(respProto) == reflect.TypeOf(&emptypb.Empty{}) {
			render.Status(r, http.StatusNoContent)
			render.NoContent(w, r)
		} else {
			render.Status(r, http.StatusCreated)
			render.JSON(w, r, resp)
		}
	})
}

func wrapHandler[TReq, TResp any](handler func(context.Context, *TReq, *TResp) error) handlerFunc {
	return func(ctx context.Context, req interface{}, resp interface{}) error {
		return handler(ctx, req.(*TReq), resp.(*TResp))
	}
}

func registerHandlers(r chi.Router, h settings.ServiceHandler) {
	registerPostHandler(r, "/api/v0/settings/bundle-save", &settingssvc.SaveBundleRequest{}, &settingssvc.SaveBundleResponse{}, wrapHandler(h.SaveBundle))
	registerPostHandler(r, "/api/v0/settings/bundle-get", &settingssvc.GetBundleRequest{}, &settingssvc.GetBundleResponse{}, wrapHandler(h.GetBundle))
	registerPostHandler(r, "/api/v0/settings/bundles-list", &settingssvc.ListBundlesRequest{}, &settingssvc.ListBundlesResponse{}, wrapHandler(h.ListBundles))
	registerPostHandler(r, "/api/v0/settings/bundles-add-setting", &settingssvc.AddSettingToBundleRequest{}, &settingssvc.AddSettingToBundleResponse{}, wrapHandler(h.AddSettingToBundle))
	registerPostHandler(r, "/api/v0/settings/bundles-remove-setting", &settingssvc.RemoveSettingFromBundleRequest{}, &emptypb.Empty{}, wrapHandler(h.RemoveSettingFromBundle))

	registerPostHandler(r, "/api/v0/settings/values-save", &settingssvc.SaveValueRequest{}, &settingssvc.SaveValueResponse{}, wrapHandler(h.SaveValue))
	registerPostHandler(r, "/api/v0/settings/values-get", &settingssvc.GetValueRequest{}, &settingssvc.GetValueResponse{}, wrapHandler(h.GetValue))
	registerPostHandler(r, "/api/v0/settings/values-list", &settingssvc.ListValuesRequest{}, &settingssvc.ListValuesResponse{}, wrapHandler(h.ListValues))
	registerPostHandler(r, "/api/v0/settings/values-get-by-unique-identifiers", &settingssvc.GetValueByUniqueIdentifiersRequest{}, &settingssvc.GetValueResponse{}, wrapHandler(h.GetValueByUniqueIdentifiers))

	registerPostHandler(r, "/api/v0/settings/roles-list", &settingssvc.ListBundlesRequest{}, &settingssvc.ListBundlesResponse{}, wrapHandler(h.ListRoles))
	registerPostHandler(r, "/api/v0/settings/assignments-list", &settingssvc.ListRoleAssignmentsRequest{}, &settingssvc.ListRoleAssignmentsResponse{}, wrapHandler(h.ListRoleAssignments))
	registerPostHandler(r, "/api/v0/settings/assignments-list-filtered", &settingssvc.ListRoleAssignmentsFilteredRequest{}, &settingssvc.ListRoleAssignmentsResponse{}, wrapHandler(h.ListRoleAssignmentsFiltered))
	registerPostHandler(r, "/api/v0/settings/assignments-add", &settingssvc.AssignRoleToUserRequest{}, &settingssvc.AssignRoleToUserResponse{}, wrapHandler(h.AssignRoleToUser))
	registerPostHandler(r, "/api/v0/settings/assignments-remove", &settingssvc.RemoveRoleFromUserRequest{}, &emptypb.Empty{}, wrapHandler(h.RemoveRoleFromUser))

	registerPostHandler(r, "/api/v0/settings/permissions-list", &settingssvc.ListPermissionsRequest{}, &settingssvc.ListPermissionsResponse{}, wrapHandler(h.ListPermissions))
	registerPostHandler(r, "/api/v0/settings/permissions-list-by-resource", &settingssvc.ListPermissionsByResourceRequest{}, &settingssvc.ListPermissionsByResourceResponse{}, wrapHandler(h.ListPermissionsByResource))
	registerPostHandler(r, "/api/v0/settings/permissions-get-by-id", &settingssvc.GetPermissionByIDRequest{}, &settingssvc.GetPermissionByIDResponse{}, wrapHandler(h.GetPermissionByID))
}

func handleError(w http.ResponseWriter, err error) {
	if merr, ok := errors.As(err); ok {
		switch merr.Code {
		case http.StatusNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case http.StatusBadRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case http.StatusInternalServerError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
