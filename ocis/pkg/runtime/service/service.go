package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	authapp "github.com/owncloud/ocis/v2/services/auth-app/pkg/command"
	"github.com/urfave/cli/v2"

	"github.com/cenkalti/backoff"
	"github.com/mohae/deepcopy"
	"github.com/olekukonko/tablewriter"
	notifications "github.com/owncloud/ocis/v2/services/notifications/pkg/command"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/owncloud/reva/v2/pkg/logger"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/thejerf/suture/v4"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	activitylog "github.com/owncloud/ocis/v2/services/activitylog/pkg/command"
	antivirus "github.com/owncloud/ocis/v2/services/antivirus/pkg/command"
	appProvider "github.com/owncloud/ocis/v2/services/app-provider/pkg/command"
	appRegistry "github.com/owncloud/ocis/v2/services/app-registry/pkg/command"
	audit "github.com/owncloud/ocis/v2/services/audit/pkg/command"
	authbasic "github.com/owncloud/ocis/v2/services/auth-basic/pkg/command"
	authmachine "github.com/owncloud/ocis/v2/services/auth-machine/pkg/command"
	authservice "github.com/owncloud/ocis/v2/services/auth-service/pkg/command"
	clientlog "github.com/owncloud/ocis/v2/services/clientlog/pkg/command"
	eventhistory "github.com/owncloud/ocis/v2/services/eventhistory/pkg/command"
	frontend "github.com/owncloud/ocis/v2/services/frontend/pkg/command"
	gateway "github.com/owncloud/ocis/v2/services/gateway/pkg/command"
	graph "github.com/owncloud/ocis/v2/services/graph/pkg/command"
	groups "github.com/owncloud/ocis/v2/services/groups/pkg/command"
	idm "github.com/owncloud/ocis/v2/services/idm/pkg/command"
	idp "github.com/owncloud/ocis/v2/services/idp/pkg/command"
	invitations "github.com/owncloud/ocis/v2/services/invitations/pkg/command"
	nats "github.com/owncloud/ocis/v2/services/nats/pkg/command"
	ocdav "github.com/owncloud/ocis/v2/services/ocdav/pkg/command"
	ocm "github.com/owncloud/ocis/v2/services/ocm/pkg/command"
	ocs "github.com/owncloud/ocis/v2/services/ocs/pkg/command"
	policies "github.com/owncloud/ocis/v2/services/policies/pkg/command"
	postprocessing "github.com/owncloud/ocis/v2/services/postprocessing/pkg/command"
	proxy "github.com/owncloud/ocis/v2/services/proxy/pkg/command"
	search "github.com/owncloud/ocis/v2/services/search/pkg/command"
	settings "github.com/owncloud/ocis/v2/services/settings/pkg/command"
	sharing "github.com/owncloud/ocis/v2/services/sharing/pkg/command"
	sse "github.com/owncloud/ocis/v2/services/sse/pkg/command"
	storagepublic "github.com/owncloud/ocis/v2/services/storage-publiclink/pkg/command"
	storageshares "github.com/owncloud/ocis/v2/services/storage-shares/pkg/command"
	storageSystem "github.com/owncloud/ocis/v2/services/storage-system/pkg/command"
	storageusers "github.com/owncloud/ocis/v2/services/storage-users/pkg/command"
	thumbnails "github.com/owncloud/ocis/v2/services/thumbnails/pkg/command"
	userlog "github.com/owncloud/ocis/v2/services/userlog/pkg/command"
	users "github.com/owncloud/ocis/v2/services/users/pkg/command"
	web "github.com/owncloud/ocis/v2/services/web/pkg/command"
	webdav "github.com/owncloud/ocis/v2/services/webdav/pkg/command"
	webfinger "github.com/owncloud/ocis/v2/services/webfinger/pkg/command"
)

var (
	// wait funcs run after the service group has been started.
	_waitFuncs = []func(*ociscfg.Config) error{pingNats, pingGateway, nil, wait(time.Second), nil}

	// Use the runner.DefaultInterruptDuration as defaults for the individual service shutdown timeouts.
	_defaultShutdownTimeoutDuration = runner.DefaultInterruptDuration
	// Use the runner.DefaultGroupInterruptDuration as defaults for the server interruption timeout.
	_defaultInterruptTimeoutDuration = runner.DefaultGroupInterruptDuration
)

type serviceFuncMap map[string]func(*ociscfg.Config) suture.Service

// Service represents a RPC service.
type Service struct {
	Supervisor *suture.Supervisor
	Services   []serviceFuncMap
	Additional serviceFuncMap
	Log        log.Logger

	mu           sync.Mutex
	serviceToken map[string][]suture.ServiceToken
	cfg          *ociscfg.Config
}

// NewService returns a configured service with a controller and a default logger.
// When used as a library, flags are not parsed, and in order to avoid introducing a global state with init functions
// calls are done explicitly to loadFromEnv().
// Since this is the public constructor, options need to be added, at the moment only logging options
// are supported in order to match the running OwnCloud services structured log.
func NewService(ctx context.Context, options ...Option) (*Service, error) {
	opts := NewOptions()

	for _, f := range options {
		f(opts)
	}

	l := log.NewLogger(
		log.Color(opts.Config.Log.Color),
		log.Pretty(opts.Config.Log.Pretty),
		log.Level(opts.Config.Log.Level),
	)

	s := &Service{
		Services:   make([]serviceFuncMap, len(_waitFuncs)),
		Additional: make(serviceFuncMap),
		Log:        l,

		serviceToken: make(map[string][]suture.ServiceToken),
		cfg:          opts.Config,
	}

	// run server command
	runServerCommand := func(ctx context.Context, server *cli.Command) error {
		cliCtx := &cli.Context{Context: ctx}
		if err := server.Before(cliCtx); err != nil {
			return err
		}
		return server.Action(cliCtx)
	}

	// populate services
	reg := func(priority int, name string, exec func(context.Context, *ociscfg.Config) error) {
		if s.Services[priority] == nil {
			s.Services[priority] = make(serviceFuncMap)
		}
		s.Services[priority][name] = NewSutureServiceBuilder(exec)
	}

	// nats is in priority group 0. It needs to start before all other services
	reg(0, opts.Config.Nats.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Nats.Context = ctx
		cfg.Nats.Commons = cfg.Commons
		return runServerCommand(ctx, nats.Server(cfg.Nats))
	})

	// gateway is in priority group 1. It needs to start before the reva services
	reg(1, opts.Config.Gateway.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Gateway.Context = ctx
		cfg.Gateway.Commons = cfg.Commons
		return runServerCommand(ctx, gateway.Server(cfg.Gateway))
	})

	// priority group 2 is empty for now

	// most services are in priority group 3
	reg(3, opts.Config.Activitylog.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Activitylog.Context = ctx
		cfg.Activitylog.Commons = cfg.Commons
		return runServerCommand(ctx, activitylog.Server(cfg.Activitylog))
	})
	reg(3, opts.Config.AppProvider.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AppProvider.Context = ctx
		cfg.AppProvider.Commons = cfg.Commons
		return runServerCommand(ctx, appProvider.Server(cfg.AppProvider))
	})
	reg(3, opts.Config.AppRegistry.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AppRegistry.Context = ctx
		cfg.AppRegistry.Commons = cfg.Commons
		return runServerCommand(ctx, appRegistry.Server(cfg.AppRegistry))
	})
	reg(3, opts.Config.AuthBasic.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AuthBasic.Context = ctx
		cfg.AuthBasic.Commons = cfg.Commons
		return runServerCommand(ctx, authbasic.Server(cfg.AuthBasic))
	})
	reg(3, opts.Config.AuthMachine.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AuthMachine.Context = ctx
		cfg.AuthMachine.Commons = cfg.Commons
		return runServerCommand(ctx, authmachine.Server(cfg.AuthMachine))
	})
	reg(3, opts.Config.AuthService.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AuthService.Context = ctx
		cfg.AuthService.Commons = cfg.Commons
		return runServerCommand(ctx, authservice.Server(cfg.AuthService))
	})
	reg(3, opts.Config.Clientlog.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Clientlog.Context = ctx
		cfg.Clientlog.Commons = cfg.Commons
		return runServerCommand(ctx, clientlog.Server(cfg.Clientlog))
	})
	reg(3, opts.Config.EventHistory.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.EventHistory.Context = ctx
		cfg.EventHistory.Commons = cfg.Commons
		return runServerCommand(ctx, eventhistory.Server(cfg.EventHistory))
	})
	reg(3, opts.Config.Graph.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Graph.Context = ctx
		cfg.Graph.Commons = cfg.Commons
		return runServerCommand(ctx, graph.Server(cfg.Graph))
	})
	reg(3, opts.Config.Groups.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Groups.Context = ctx
		cfg.Groups.Commons = cfg.Commons
		return runServerCommand(ctx, groups.Server(cfg.Groups))
	})
	reg(3, opts.Config.IDM.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.IDM.Context = ctx
		cfg.IDM.Commons = cfg.Commons
		return runServerCommand(ctx, idm.Server(cfg.IDM))
	})
	reg(3, opts.Config.OCDav.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.OCDav.Context = ctx
		cfg.OCDav.Commons = cfg.Commons
		return runServerCommand(ctx, ocdav.Server(cfg.OCDav))
	})
	reg(3, opts.Config.OCS.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.OCS.Context = ctx
		cfg.OCS.Commons = cfg.Commons
		return runServerCommand(ctx, ocs.Server(cfg.OCS))
	})
	reg(3, opts.Config.Postprocessing.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Postprocessing.Context = ctx
		cfg.Postprocessing.Commons = cfg.Commons
		return runServerCommand(ctx, postprocessing.Server(cfg.Postprocessing))
	})
	reg(3, opts.Config.Search.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Search.Context = ctx
		cfg.Search.Commons = cfg.Commons
		return runServerCommand(ctx, search.Server(cfg.Search))
	})
	reg(3, opts.Config.Settings.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Settings.Context = ctx
		cfg.Settings.Commons = cfg.Commons
		return runServerCommand(ctx, settings.Server(cfg.Settings))
	})
	reg(3, opts.Config.StoragePublicLink.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.StoragePublicLink.Context = ctx
		cfg.StoragePublicLink.Commons = cfg.Commons
		return runServerCommand(ctx, storagepublic.Server(cfg.StoragePublicLink))
	})
	reg(3, opts.Config.StorageShares.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.StorageShares.Context = ctx
		cfg.StorageShares.Commons = cfg.Commons
		return runServerCommand(ctx, storageshares.Server(cfg.StorageShares))
	})
	reg(3, opts.Config.StorageSystem.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.StorageSystem.Context = ctx
		cfg.StorageSystem.Commons = cfg.Commons
		return runServerCommand(ctx, storageSystem.Server(cfg.StorageSystem))
	})
	reg(3, opts.Config.StorageUsers.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.StorageUsers.Context = ctx
		cfg.StorageUsers.Commons = cfg.Commons
		return runServerCommand(ctx, storageusers.Server(cfg.StorageUsers))
	})
	reg(3, opts.Config.Thumbnails.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Thumbnails.Context = ctx
		cfg.Thumbnails.Commons = cfg.Commons
		return runServerCommand(ctx, thumbnails.Server(cfg.Thumbnails))
	})
	reg(3, opts.Config.Userlog.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Userlog.Context = ctx
		cfg.Userlog.Commons = cfg.Commons
		return runServerCommand(ctx, userlog.Server(cfg.Userlog))
	})
	reg(3, opts.Config.Users.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Users.Context = ctx
		cfg.Users.Commons = cfg.Commons
		return runServerCommand(ctx, users.Server(cfg.Users))
	})
	reg(3, opts.Config.Web.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Web.Context = ctx
		cfg.Web.Commons = cfg.Commons
		return runServerCommand(ctx, web.Server(cfg.Web))
	})
	reg(3, opts.Config.WebDAV.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.WebDAV.Context = ctx
		cfg.WebDAV.Commons = cfg.Commons
		return runServerCommand(ctx, webdav.Server(cfg.WebDAV))
	})
	reg(3, opts.Config.Webfinger.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Webfinger.Context = ctx
		cfg.Webfinger.Commons = cfg.Commons
		return runServerCommand(ctx, webfinger.Server(cfg.Webfinger))
	})
	reg(3, opts.Config.IDP.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.IDP.Context = ctx
		cfg.IDP.Commons = cfg.Commons
		return runServerCommand(ctx, idp.Server(cfg.IDP))
	})
	reg(3, opts.Config.Proxy.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Proxy.Context = ctx
		cfg.Proxy.Commons = cfg.Commons
		return runServerCommand(ctx, proxy.Server(cfg.Proxy))
	})
	reg(3, opts.Config.Sharing.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Sharing.Context = ctx
		cfg.Sharing.Commons = cfg.Commons
		return runServerCommand(ctx, sharing.Server(cfg.Sharing))
	})
	reg(3, opts.Config.SSE.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.SSE.Context = ctx
		cfg.SSE.Commons = cfg.Commons
		return runServerCommand(ctx, sse.Server(cfg.SSE))
	})
	reg(3, opts.Config.OCM.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.OCM.Context = ctx
		cfg.OCM.Commons = cfg.Commons
		return runServerCommand(ctx, ocm.Server(cfg.OCM))
	})

	// out of some unknown reason ci gets angry when frontend service starts in priority group 3
	// this is not reproducible locally, it can start when nats and gateway are already running
	// FIXME: find out why
	reg(4, opts.Config.Frontend.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Frontend.Context = ctx
		cfg.Frontend.Commons = cfg.Commons
		return runServerCommand(ctx, frontend.Server(cfg.Frontend))
	})

	// populate optional services
	areg := func(name string, exec func(context.Context, *ociscfg.Config) error) {
		s.Additional[name] = NewSutureServiceBuilder(exec)
	}
	areg(opts.Config.Antivirus.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Antivirus.Context = ctx
		// cfg.Antivirus.Commons = cfg.Commons // antivirus holds no Commons atm
		return runServerCommand(ctx, antivirus.Server(cfg.Antivirus))
	})
	areg(opts.Config.Audit.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Audit.Context = ctx
		cfg.Audit.Commons = cfg.Commons
		return runServerCommand(ctx, audit.Server(cfg.Audit))
	})
	areg(opts.Config.AuthApp.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AuthApp.Context = ctx
		cfg.AuthApp.Commons = cfg.Commons
		return runServerCommand(ctx, authapp.Server(cfg.AuthApp))
	})
	areg(opts.Config.Policies.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Policies.Context = ctx
		cfg.Policies.Commons = cfg.Commons
		return runServerCommand(ctx, policies.Server(cfg.Policies))
	})
	areg(opts.Config.Invitations.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Invitations.Context = ctx
		cfg.Invitations.Commons = cfg.Commons
		return runServerCommand(ctx, invitations.Server(cfg.Invitations))
	})
	areg(opts.Config.Notifications.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Notifications.Context = ctx
		cfg.Notifications.Commons = cfg.Commons
		return runServerCommand(ctx, notifications.Server(cfg.Notifications))
	})

	return s, nil
}

// Start a rpc service. By default, the package scope Start will run all default services to provide with a working
// oCIS instance.
func Start(ctx context.Context, o ...Option) error {
	// Start the runtime. Most likely this was called ONLY by the `ocis server` subcommand, but since we cannot protect
	// from the caller, the previous statement holds truth.

	// prepare a new rpc Service struct.
	s, err := NewService(ctx, o...)
	if err != nil {
		return err
	}

	// cancel the context when a signal is received.
	var cancel context.CancelFunc = func() {}
	if ctx == nil {
		ctx, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
		defer cancel()
	}

	// tolerance controls backoff cycles from the supervisor.
	tolerance := 5
	totalBackoff := 0

	// Start creates its own supervisor. Running services under `ocis server` will create its own supervision tree.
	s.Supervisor = suture.New("ocis", suture.Spec{
		EventHook: func(e suture.Event) {
			if e.Type() == suture.EventTypeBackoff {
				totalBackoff++
				if totalBackoff == tolerance {
					cancel()
				}
			}
			s.Log.Info().Str("event", e.String()).Msg(fmt.Sprintf("supervisor: %v", e.Map()["supervisor_name"]))
		},
		FailureThreshold: 5,
		FailureBackoff:   3 * time.Second,
	})

	if s.cfg.Commons == nil {
		s.cfg.Commons = &shared.Commons{
			Log: &shared.Log{},
		}
	}

	// prepare RPC server
	srv, err := newRPCServer(s)
	if err != nil {
		s.Log.Fatal().Err(err).Msg("could not create RPC server")
		return err
	}

	// prepare the set of services to run
	// runset keeps track of which services to start supervised.
	runset := s.generateRunSet(s.cfg)

	// There are reasons not to do this, but we have race conditions ourselves. Until we resolve them, mind the following disclaimer:
	// Calling ServeBackground will CORRECTLY start the supervisor running in a new goroutine. It is risky to directly run
	// go supervisor.Serve()
	// because that will briefly create a race condition as it starts up, if you try to .Add() services immediately afterward.
	// https://pkg.go.dev/github.com/thejerf/suture/v4@v4.0.0#Supervisor
	go s.Supervisor.ServeBackground(ctx)

	for i, service := range s.Services {
		scheduleServiceTokens(s, runset, service)
		if _waitFuncs[i] != nil {
			if err := _waitFuncs[i](s.cfg); err != nil {
				s.Log.Fatal().Err(err).Msg("wait func failed")
			}
		}
	}

	// schedule services that are optional
	scheduleServiceTokens(s, runset, s.Additional)

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Log.Fatal().Err(err).Msg("could not start rpc server")
		}
	}()

	// trapShutdownCtx will block on the context-done channel for interruptions.
	trapShutdownCtx(s, srv, ctx)
	return nil
}

// scheduleServiceTokens adds service tokens to the service supervisor.
func scheduleServiceTokens(s *Service, runset map[string]struct{}, funcSet serviceFuncMap) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for name := range runset {
		if _, ok := funcSet[name]; !ok {
			continue
		}

		swap := deepcopy.Copy(s.cfg)
		s.serviceToken[name] = append(s.serviceToken[name], s.Supervisor.Add(funcSet[name](swap.(*ociscfg.Config))))
	}
}

// generateRunSet interprets the cfg.Runtime.Services config option to cherry-pick which services to start using
// the runtime.
func (s *Service) generateRunSet(cfg *ociscfg.Config) map[string]struct{} {
	runset := make(map[string]struct{})
	if cfg.Runtime.Services != nil {
		for _, name := range cfg.Runtime.Services {
			runset[name] = struct{}{}
		}
		return runset
	}

	for _, service := range s.Services {
		for name := range service {
			runset[name] = struct{}{}
		}
	}

	// add additional services if explicitly added by config
	for _, name := range cfg.Runtime.Additional {
		runset[name] = struct{}{}
	}

	// remove services if explicitly excluded by config
	for _, name := range cfg.Runtime.Disabled {
		delete(runset, name)
	}
	return runset
}

// List running processes for the Service Controller.
func (s *Service) List(_ struct{}, reply *string) error {
	tableString := &strings.Builder{}
	table := tablewriter.NewTable(tableString)
	table.Header("Service")

	s.mu.Lock()
	names := []string{}
	for t := range s.serviceToken {
		if len(s.serviceToken[t]) > 0 {
			names = append(names, t)
		}
	}
	s.mu.Unlock()

	sort.Strings(names)

	for n := range names {
		table.Append([]string{names[n]})
	}

	table.Render()
	*reply = tableString.String()
	return nil
}

func trapShutdownCtx(s *Service, srv *http.Server, ctx context.Context) {
	<-ctx.Done()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), _defaultShutdownTimeoutDuration)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			s.Log.Error().Err(err).Msg("could not shutdown tcp listener")
			return
		}
		s.Log.Info().Msg("tcp listener shutdown")
	}()

	s.mu.Lock()
	for sName := range s.serviceToken {
		for i := range s.serviceToken[sName] {
			wg.Add(1)
			go func() {
				s.Log.Warn().Msgf("call supervisor RemoveAndWait for %s", sName)
				defer wg.Done()
				if err := s.Supervisor.RemoveAndWait(s.serviceToken[sName][i], _defaultShutdownTimeoutDuration); err != nil && !errors.Is(err, suture.ErrSupervisorNotRunning) {
					s.Log.Error().Err(err).Str("service", sName).Msgf("terminating with signal: %+v", s)
				}
				s.Log.Warn().Msgf("done supervisor RemoveAndWait for %s", sName)
			}()
		}
	}
	s.mu.Unlock()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-time.After(_defaultInterruptTimeoutDuration):
		s.Log.Fatal().Msg("ocis graceful shutdown timeout reached, terminating")
	case <-done:
		s.Log.Info().Msg("all ocis services gracefully stopped")
		return
	}
}

// pingNats will attempt to connect to nats, blocking until a connection is established
func pingNats(cfg *ociscfg.Config) error {
	// We need to get a natsconfig from somewhere. We can use any one.
	evcfg := cfg.Postprocessing.Postprocessing.Events
	_, err := stream.NatsFromConfig("initial", true, stream.NatsConfig(evcfg))
	return err
}

func pingGateway(cfg *ociscfg.Config) error {
	// init grpc connection
	_, err := ogrpc.NewClient()
	if err != nil {
		return err
	}

	b := backoff.NewExponentialBackOff()
	o := func() error {
		n := b.NextBackOff()
		_, err := pool.GetGatewayServiceClient(cfg.Reva.Address)
		if err != nil && n > time.Second {
			logger.New().Error().Err(err).Msgf("can't connect to gateway service, retrying in %s", n)
		}
		return err
	}

	err = backoff.Retry(o, b)
	return err
}

func wait(d time.Duration) func(cfg *ociscfg.Config) error {
	return func(cfg *ociscfg.Config) error {
		time.Sleep(d)
		return nil
	}
}

// newRPCServer creates an HTTP server to expose the "s" service's methods using RPC.
// The host and port for the server are taken from the service configuration (s.cfg.Runtime)
func newRPCServer(s *Service) (*http.Server, error) {
	rpcSrv := rpc.NewServer()
	if err := rpcSrv.Register(s); err != nil {
		return nil, err
	}
	rpcSrv.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	// srv.HandleHTTP will register the handlers in the http.DefaultServeMux

	srv := &http.Server{
		Addr: net.JoinHostPort(s.cfg.Runtime.Host, s.cfg.Runtime.Port), // this is always tcp
		// Handler must be nil to use the http.DefaultServeMux
		ReadHeaderTimeout: 5 * time.Second,
	}
	return srv, nil
}
