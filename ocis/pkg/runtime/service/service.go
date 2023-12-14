package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/mohae/deepcopy"
	"github.com/olekukonko/tablewriter"
	"github.com/thejerf/suture/v4"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
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
	notifications "github.com/owncloud/ocis/v2/services/notifications/pkg/command"
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
	store "github.com/owncloud/ocis/v2/services/store/pkg/command"
	thumbnails "github.com/owncloud/ocis/v2/services/thumbnails/pkg/command"
	userlog "github.com/owncloud/ocis/v2/services/userlog/pkg/command"
	users "github.com/owncloud/ocis/v2/services/users/pkg/command"
	web "github.com/owncloud/ocis/v2/services/web/pkg/command"
	webdav "github.com/owncloud/ocis/v2/services/webdav/pkg/command"
	webfinger "github.com/owncloud/ocis/v2/services/webfinger/pkg/command"
)

var (
	// runset keeps track of which services to start supervised.
	runset map[string]struct{}
	// time to wait after starting the preliminary services
	_preliminaryDelay = 6 * time.Second
	// time to wait between starting service groups (preliminary, main, delayed)
	_startDelay = 2 * time.Second
)

type serviceFuncMap map[string]func(*ociscfg.Config) suture.Service

// Service represents a RPC service.
type Service struct {
	Supervisor       *suture.Supervisor
	Preliminary      serviceFuncMap
	ServicesRegistry serviceFuncMap
	Delayed          serviceFuncMap
	Additional       serviceFuncMap
	Log              log.Logger

	serviceToken map[string][]suture.ServiceToken
	context      context.Context
	cancel       context.CancelFunc
	cfg          *ociscfg.Config
}

// NewService returns a configured service with a controller and a default logger.
// When used as a library, flags are not parsed, and in order to avoid introducing a global state with init functions
// calls are done explicitly to loadFromEnv().
// Since this is the public constructor, options need to be added, at the moment only logging options
// are supported in order to match the running OwnCloud services structured log.
func NewService(options ...Option) (*Service, error) {
	opts := NewOptions()

	for _, f := range options {
		f(opts)
	}

	l := log.NewLogger(
		log.Color(opts.Config.Log.Color),
		log.Pretty(opts.Config.Log.Pretty),
		log.Level(opts.Config.Log.Level),
	)

	globalCtx, cancelGlobal := context.WithCancel(context.Background())

	s := &Service{
		Preliminary:      make(serviceFuncMap),
		ServicesRegistry: make(serviceFuncMap),
		Delayed:          make(serviceFuncMap),
		Additional:       make(serviceFuncMap),
		Log:              l,

		serviceToken: make(map[string][]suture.ServiceToken),
		context:      globalCtx,
		cancel:       cancelGlobal,
		cfg:          opts.Config,
	}

	// start nats first - it is used as service registry
	s.Preliminary[opts.Config.Nats.Service.Name] = NewSutureServiceBuilder(func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Nats.Context = ctx
		cfg.Nats.Commons = cfg.Commons
		return nats.Execute(cfg.Nats)
	})

	// populate services
	reg := func(name string, exec func(context.Context, *ociscfg.Config) error) {
		s.ServicesRegistry[name] = NewSutureServiceBuilder(exec)
	}
	reg(opts.Config.AppProvider.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AppProvider.Context = ctx
		cfg.AppProvider.Commons = cfg.Commons
		return appProvider.Execute(cfg.AppProvider)
	})
	reg(opts.Config.AppRegistry.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AppRegistry.Context = ctx
		cfg.AppRegistry.Commons = cfg.Commons
		return appRegistry.Execute(cfg.AppRegistry)
	})
	reg(opts.Config.AuthBasic.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AuthBasic.Context = ctx
		cfg.AuthBasic.Commons = cfg.Commons
		return authbasic.Execute(cfg.AuthBasic)
	})
	reg(opts.Config.AuthMachine.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AuthMachine.Context = ctx
		cfg.AuthMachine.Commons = cfg.Commons
		return authmachine.Execute(cfg.AuthMachine)
	})
	reg(opts.Config.AuthService.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.AuthService.Context = ctx
		cfg.AuthService.Commons = cfg.Commons
		return authservice.Execute(cfg.AuthService)
	})
	reg(opts.Config.Clientlog.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Clientlog.Context = ctx
		cfg.Clientlog.Commons = cfg.Commons
		return clientlog.Execute(cfg.Clientlog)
	})
	reg(opts.Config.EventHistory.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.EventHistory.Context = ctx
		cfg.EventHistory.Commons = cfg.Commons
		return eventhistory.Execute(cfg.EventHistory)
	})
	reg(opts.Config.Gateway.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Gateway.Context = ctx
		cfg.Gateway.Commons = cfg.Commons
		return gateway.Execute(cfg.Gateway)
	})
	reg(opts.Config.Graph.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Graph.Context = ctx
		cfg.Graph.Commons = cfg.Commons
		return graph.Execute(cfg.Graph)
	})
	reg(opts.Config.Groups.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Groups.Context = ctx
		cfg.Groups.Commons = cfg.Commons
		return groups.Execute(cfg.Groups)
	})
	reg(opts.Config.IDM.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.IDM.Context = ctx
		cfg.IDM.Commons = cfg.Commons
		return idm.Execute(cfg.IDM)
	})
	reg(opts.Config.Invitations.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Invitations.Context = ctx
		cfg.Invitations.Commons = cfg.Commons
		return invitations.Execute(cfg.Invitations)
	})
	reg(opts.Config.Notifications.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Notifications.Context = ctx
		cfg.Notifications.Commons = cfg.Commons
		return notifications.Execute(cfg.Notifications)
	})
	reg(opts.Config.OCDav.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.OCDav.Context = ctx
		cfg.OCDav.Commons = cfg.Commons
		return ocdav.Execute(cfg.OCDav)
	})
	reg(opts.Config.OCS.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.OCS.Context = ctx
		cfg.OCS.Commons = cfg.Commons
		return ocs.Execute(cfg.OCS)
	})
	reg(opts.Config.Postprocessing.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Postprocessing.Context = ctx
		cfg.Postprocessing.Commons = cfg.Commons
		return postprocessing.Execute(cfg.Postprocessing)
	})
	reg(opts.Config.Search.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Search.Context = ctx
		cfg.Search.Commons = cfg.Commons
		return search.Execute(cfg.Search)
	})
	reg(opts.Config.Settings.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Settings.Context = ctx
		cfg.Settings.Commons = cfg.Commons
		return settings.Execute(cfg.Settings)
	})
	reg(opts.Config.StoragePublicLink.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.StoragePublicLink.Context = ctx
		cfg.StoragePublicLink.Commons = cfg.Commons
		return storagepublic.Execute(cfg.StoragePublicLink)
	})
	reg(opts.Config.StorageShares.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.StorageShares.Context = ctx
		cfg.StorageShares.Commons = cfg.Commons
		return storageshares.Execute(cfg.StorageShares)
	})
	reg(opts.Config.StorageSystem.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.StorageSystem.Context = ctx
		cfg.StorageSystem.Commons = cfg.Commons
		return storageSystem.Execute(cfg.StorageSystem)
	})
	reg(opts.Config.StorageUsers.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.StorageUsers.Context = ctx
		cfg.StorageUsers.Commons = cfg.Commons
		return storageusers.Execute(cfg.StorageUsers)
	})
	reg(opts.Config.Store.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Store.Context = ctx
		cfg.Store.Commons = cfg.Commons
		return store.Execute(cfg.Store)
	})
	reg(opts.Config.Thumbnails.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Thumbnails.Context = ctx
		cfg.Thumbnails.Commons = cfg.Commons
		return thumbnails.Execute(cfg.Thumbnails)
	})
	reg(opts.Config.Userlog.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Userlog.Context = ctx
		cfg.Userlog.Commons = cfg.Commons
		return userlog.Execute(cfg.Userlog)
	})
	reg(opts.Config.Users.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Users.Context = ctx
		cfg.Users.Commons = cfg.Commons
		return users.Execute(cfg.Users)
	})
	reg(opts.Config.Web.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Web.Context = ctx
		cfg.Web.Commons = cfg.Commons
		return web.Execute(cfg.Web)
	})
	reg(opts.Config.WebDAV.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.WebDAV.Context = ctx
		cfg.WebDAV.Commons = cfg.Commons
		return webdav.Execute(cfg.WebDAV)
	})
	reg(opts.Config.Webfinger.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Webfinger.Context = ctx
		cfg.Webfinger.Commons = cfg.Commons
		return webfinger.Execute(cfg.Webfinger)
	})

	// populate optional services
	areg := func(name string, exec func(context.Context, *ociscfg.Config) error) {
		s.Additional[name] = NewSutureServiceBuilder(exec)
	}
	areg(opts.Config.Antivirus.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Antivirus.Context = ctx
		// cfg.Antivirus.Commons = cfg.Commons // antivirus holds no Commons atm
		return antivirus.Execute(cfg.Antivirus)
	})
	areg(opts.Config.Audit.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Audit.Context = ctx
		cfg.Audit.Commons = cfg.Commons
		return audit.Execute(cfg.Audit)
	})
	areg(opts.Config.Policies.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Policies.Context = ctx
		cfg.Policies.Commons = cfg.Commons
		return policies.Execute(cfg.Policies)
	})

	// populate delayed services
	dreg := func(name string, exec func(context.Context, *ociscfg.Config) error) {
		s.Delayed[name] = NewSutureServiceBuilder(exec)
	}
	dreg(opts.Config.Frontend.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Frontend.Context = ctx
		cfg.Frontend.Commons = cfg.Commons
		return frontend.Execute(cfg.Frontend)
	})
	dreg(opts.Config.IDP.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.IDP.Context = ctx
		cfg.IDP.Commons = cfg.Commons
		return idp.Execute(cfg.IDP)
	})
	dreg(opts.Config.Proxy.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Proxy.Context = ctx
		cfg.Proxy.Commons = cfg.Commons
		return proxy.Execute(cfg.Proxy)
	})
	dreg(opts.Config.Sharing.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.Sharing.Context = ctx
		cfg.Sharing.Commons = cfg.Commons
		return sharing.Execute(cfg.Sharing)
	})
	dreg(opts.Config.SSE.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.SSE.Context = ctx
		cfg.SSE.Commons = cfg.Commons
		return sse.Execute(cfg.SSE)
	})
	dreg(opts.Config.OCM.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
		cfg.OCM.Context = ctx
		cfg.OCM.Commons = cfg.Commons
		return ocm.Execute(cfg.OCM)
	})

	return s, nil
}

// Start a rpc service. By default, the package scope Start will run all default services to provide with a working
// oCIS instance.
func Start(o ...Option) error {
	// Start the runtime. Most likely this was called ONLY by the `ocis server` subcommand, but since we cannot protect
	// from the caller, the previous statement holds truth.

	// prepare a new rpc Service struct.
	s, err := NewService(o...)
	if err != nil {
		return err
	}

	// halt listens for interrupt signals and blocks.
	halt := make(chan os.Signal, 1)
	signal.Notify(halt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// tolerance controls backoff cycles from the supervisor.
	tolerance := 5
	totalBackoff := 0

	// Start creates its own supervisor. Running services under `ocis server` will create its own supervision tree.
	s.Supervisor = suture.New("ocis", suture.Spec{
		EventHook: func(e suture.Event) {
			if e.Type() == suture.EventTypeBackoff {
				totalBackoff++
				if totalBackoff == tolerance {
					halt <- os.Interrupt
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

	if err = rpc.Register(s); err != nil {
		if s != nil {
			s.Log.Fatal().Err(err)
		}
	}
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", net.JoinHostPort(s.cfg.Runtime.Host, s.cfg.Runtime.Port))
	if err != nil {
		s.Log.Fatal().Err(err)
	}

	defer func() {
		if r := recover(); r != nil {
			reason := strings.Builder{}
			if _, err = net.Dial("tcp", net.JoinHostPort(s.cfg.Runtime.Host, s.cfg.Runtime.Port)); err != nil {
				reason.WriteString("runtime address already in use")
			}

			fmt.Println(reason.String())
		}
	}()

	// prepare the set of services to run
	s.generateRunSet(s.cfg)

	// schedule preliminary services first
	scheduleServiceTokens(s, s.Preliminary)

	// there are reasons not to do this, but we have race conditions ourselves. Until we resolve them, mind the following disclaimer:
	// Calling ServeBackground will CORRECTLY start the supervisor running in a new goroutine. It is risky to directly run
	// go supervisor.Serve()
	// because that will briefly create a race condition as it starts up, if you try to .Add() services immediately afterward.
	// https://pkg.go.dev/github.com/thejerf/suture/v4@v4.0.0#Supervisor
	go s.Supervisor.ServeBackground(s.context)

	// trap will block on halt channel for interruptions.
	go trap(s, halt)

	// grace period for preliminary services to get up
	time.Sleep(_preliminaryDelay)

	// schedule services that we are sure don't have interdependencies.
	scheduleServiceTokens(s, s.ServicesRegistry)

	// schedule services that are optional
	scheduleServiceTokens(s, s.Additional)

	// add services with delayed execution.
	time.Sleep(_startDelay)
	scheduleServiceTokens(s, s.Delayed)

	return http.Serve(l, nil)
}

// scheduleServiceTokens adds service tokens to the service supervisor.
func scheduleServiceTokens(s *Service, funcSet serviceFuncMap) {
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
func (s *Service) generateRunSet(cfg *ociscfg.Config) {
	runset = make(map[string]struct{})
	if cfg.Runtime.Services != nil {
		for _, name := range cfg.Runtime.Services {
			runset[name] = struct{}{}
		}
		return
	}

	for name := range s.Preliminary {
		runset[name] = struct{}{}
	}

	for name := range s.ServicesRegistry {
		runset[name] = struct{}{}
	}

	for name := range s.Delayed {
		runset[name] = struct{}{}
	}

	// add additional services if explicitly added by config
	for _, name := range cfg.Runtime.Additional {
		runset[name] = struct{}{}
	}

	// remove services if explicitly excluded by config
	for _, name := range cfg.Runtime.Disabled {
		delete(runset, name)
	}
}

// List running processes for the Service Controller.
func (s *Service) List(_ struct{}, reply *string) error {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Service"})

	names := []string{}
	for t := range s.serviceToken {
		if len(s.serviceToken[t]) > 0 {
			names = append(names, t)
		}
	}

	sort.Strings(names)

	for n := range names {
		table.Append([]string{names[n]})
	}

	table.Render()
	*reply = tableString.String()
	return nil
}

// trap blocks on halt channel. When the runtime is interrupted it
// signals the controller to stop any supervised process.
func trap(s *Service, halt chan os.Signal) {
	<-halt
	s.cancel()
	for sName := range s.serviceToken {
		for i := range s.serviceToken[sName] {
			if err := s.Supervisor.Remove(s.serviceToken[sName][i]); err != nil {
				s.Log.Error().Err(err).Str("service", "runtime service").Msgf("terminating with signal: %v", s)
			}
		}
	}
	s.Log.Debug().Str("service", "runtime service").Msgf("terminating with signal: %v", s)
	os.Exit(0)
}
