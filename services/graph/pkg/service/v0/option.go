package svc

import (
	"context"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/keycloak"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"go.opentelemetry.io/otel/trace"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Context                  context.Context
	Logger                   log.Logger
	Config                   *config.Config
	Middleware               []func(http.Handler) http.Handler
	RequireAdminMiddleware   func(http.Handler) http.Handler
	GatewaySelector          pool.Selectable[gateway.GatewayAPIClient]
	IdentityBackend          identity.Backend
	IdentityEducationBackend identity.EducationBackend
	RoleService              RoleService
	PermissionService        Permissions
	ValueService             settingssvc.ValueService
	RoleManager              *roles.Manager
	EventsPublisher          events.Publisher
	EventsConsumer           events.Consumer
	SearchService            searchsvc.SearchProviderService
	KeycloakClient           keycloak.Client
	EventHistoryClient       ehsvc.EventHistoryService
	TraceProvider            trace.TracerProvider
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Context provides a function to set the context option.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Logger provides a function to set the logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// Config provides a function to set the config option.
func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// Middleware provides a function to set the middleware option.
func Middleware(val ...func(http.Handler) http.Handler) Option {
	return func(o *Options) {
		o.Middleware = val
	}
}

// WithRequireAdminMiddleware provides a function to set the RequireAdminMiddleware option.
func WithRequireAdminMiddleware(val func(http.Handler) http.Handler) Option {
	return func(o *Options) {
		o.RequireAdminMiddleware = val
	}
}

// WithGatewaySelector provides a function to set the gateway client option.
func WithGatewaySelector(val pool.Selectable[gateway.GatewayAPIClient]) Option {
	return func(o *Options) {
		o.GatewaySelector = val
	}
}

// WithIdentityBackend provides a function to set the IdentityBackend option.
func WithIdentityBackend(val identity.Backend) Option {
	return func(o *Options) {
		o.IdentityBackend = val
	}
}

// WithIdentityEducationBackend provides a function to set the IdentityEducationBackend option.
func WithIdentityEducationBackend(val identity.EducationBackend) Option {
	return func(o *Options) {
		o.IdentityEducationBackend = val
	}
}

// WithRoleService provides a function to set the RoleService option.
func WithRoleService(val RoleService) Option {
	return func(o *Options) {
		o.RoleService = val
	}
}

// WithValueService provides a function to set the ValueService option.
func WithValueService(val settingssvc.ValueService) Option {
	return func(o *Options) {
		o.ValueService = val
	}
}

// WithSearchService provides a function to set the SearchService option.
func WithSearchService(val searchsvc.SearchProviderService) Option {
	return func(o *Options) {
		o.SearchService = val
	}
}

// PermissionService provides a function to set the PermissionService option.
func PermissionService(val settingssvc.PermissionService) Option {
	return func(o *Options) {
		o.PermissionService = val
	}
}

// RoleManager provides a function to set the RoleManager option.
func RoleManager(val *roles.Manager) Option {
	return func(o *Options) {
		o.RoleManager = val
	}
}

// EventsPublisher provides a function to set the EventsPublisher option.
func EventsPublisher(val events.Publisher) Option {
	return func(o *Options) {
		o.EventsPublisher = val
	}
}

// EventsConsumer provides a function to set the EventsConsumer option.
func EventsConsumer(val events.Consumer) Option {
	return func(o *Options) {
		o.EventsConsumer = val
	}
}

// KeycloakClient provides a function to set the KeycloakCient option.
func KeycloakClient(val keycloak.Client) Option {
	return func(o *Options) {
		o.KeycloakClient = val
	}
}

// EventHistoryClient provides a function to set the EventHistoryClient option.
func EventHistoryClient(val ehsvc.EventHistoryService) Option {
	return func(o *Options) {
		o.EventHistoryClient = val
	}
}

// TraceProvider provides a function to set the TraceProvider option.
func TraceProvider(val trace.TracerProvider) Option {
	return func(o *Options) {
		o.TraceProvider = val
	}
}
