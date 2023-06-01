package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"

	antivirus "github.com/owncloud/ocis/v2/services/antivirus/pkg/command"
	appprovider "github.com/owncloud/ocis/v2/services/app-provider/pkg/command"
	appregistry "github.com/owncloud/ocis/v2/services/app-registry/pkg/command"
	audit "github.com/owncloud/ocis/v2/services/audit/pkg/command"
	authbasic "github.com/owncloud/ocis/v2/services/auth-basic/pkg/command"
	authbearer "github.com/owncloud/ocis/v2/services/auth-bearer/pkg/command"
	authmachine "github.com/owncloud/ocis/v2/services/auth-machine/pkg/command"
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
	ocs "github.com/owncloud/ocis/v2/services/ocs/pkg/command"
	policies "github.com/owncloud/ocis/v2/services/policies/pkg/command"
	postprocessing "github.com/owncloud/ocis/v2/services/postprocessing/pkg/command"
	proxy "github.com/owncloud/ocis/v2/services/proxy/pkg/command"
	search "github.com/owncloud/ocis/v2/services/search/pkg/command"
	settings "github.com/owncloud/ocis/v2/services/settings/pkg/command"
	sharing "github.com/owncloud/ocis/v2/services/sharing/pkg/command"
	storagepubliclink "github.com/owncloud/ocis/v2/services/storage-publiclink/pkg/command"
	storageshares "github.com/owncloud/ocis/v2/services/storage-shares/pkg/command"
	storagesystem "github.com/owncloud/ocis/v2/services/storage-system/pkg/command"
	storageusers "github.com/owncloud/ocis/v2/services/storage-users/pkg/command"
	store "github.com/owncloud/ocis/v2/services/store/pkg/command"
	thumbnails "github.com/owncloud/ocis/v2/services/thumbnails/pkg/command"
	userlog "github.com/owncloud/ocis/v2/services/userlog/pkg/command"
	users "github.com/owncloud/ocis/v2/services/users/pkg/command"
	web "github.com/owncloud/ocis/v2/services/web/pkg/command"
	webdav "github.com/owncloud/ocis/v2/services/webdav/pkg/command"
	webfinger "github.com/owncloud/ocis/v2/services/webfinger/pkg/command"
)

var svccmds = []register.Command{
	func(cfg *config.Config) *cli.Command {
		// cfg.Antivirus.Commons = cfg.Commons // antivirus needs no commons atm
		return ServiceCommand(cfg, cfg.Antivirus.Service.Name, antivirus.GetCommands(cfg.Antivirus))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.AppProvider.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.AppProvider.Service.Name, appprovider.GetCommands(cfg.AppProvider))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.AppRegistry.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.AppRegistry.Service.Name, appregistry.GetCommands(cfg.AppRegistry))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Audit.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Audit.Service.Name, audit.GetCommands(cfg.Audit))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.AuthBasic.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.AuthBasic.Service.Name, authbasic.GetCommands(cfg.AuthBasic))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.AuthBearer.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.AuthBearer.Service.Name, authbearer.GetCommands(cfg.AuthBearer))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.AuthMachine.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.AuthMachine.Service.Name, authmachine.GetCommands(cfg.AuthMachine))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.EventHistory.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.EventHistory.Service.Name, eventhistory.GetCommands(cfg.EventHistory))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Frontend.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Frontend.Service.Name, frontend.GetCommands(cfg.Frontend))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Gateway.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Gateway.Service.Name, gateway.GetCommands(cfg.Gateway))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Graph.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Graph.Service.Name, graph.GetCommands(cfg.Graph))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Groups.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Groups.Service.Name, groups.GetCommands(cfg.Groups))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.IDM.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.IDM.Service.Name, idm.GetCommands(cfg.IDM))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.IDP.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.IDP.Service.Name, idp.GetCommands(cfg.IDP))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Invitations.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Invitations.Service.Name, invitations.GetCommands(cfg.Invitations))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Nats.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Nats.Service.Name, nats.GetCommands(cfg.Nats))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Notifications.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Notifications.Service.Name, notifications.GetCommands(cfg.Notifications))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.OCDav.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.OCDav.Service.Name, ocdav.GetCommands(cfg.OCDav))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.OCS.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.OCS.Service.Name, ocs.GetCommands(cfg.OCS))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Policies.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Policies.Service.Name, policies.GetCommands(cfg.Policies))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Postprocessing.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Postprocessing.Service.Name, postprocessing.GetCommands(cfg.Postprocessing))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Proxy.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Proxy.Service.Name, proxy.GetCommands(cfg.Proxy))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Search.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Search.Service.Name, search.GetCommands(cfg.Search))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Settings.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Settings.Service.Name, settings.GetCommands(cfg.Settings))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Sharing.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Sharing.Service.Name, sharing.GetCommands(cfg.Sharing))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.StoragePublicLink.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.StoragePublicLink.Service.Name, storagepubliclink.GetCommands(cfg.StoragePublicLink))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.StorageShares.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.StorageShares.Service.Name, storageshares.GetCommands(cfg.StorageShares))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.StorageSystem.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.StorageSystem.Service.Name, storagesystem.GetCommands(cfg.StorageSystem))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.StorageUsers.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.StorageUsers.Service.Name, storageusers.GetCommands(cfg.StorageUsers))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Store.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Store.Service.Name, store.GetCommands(cfg.Store))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Thumbnails.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Thumbnails.Service.Name, thumbnails.GetCommands(cfg.Thumbnails))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Userlog.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Userlog.Service.Name, userlog.GetCommands(cfg.Userlog))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Users.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Users.Service.Name, users.GetCommands(cfg.Users))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Web.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Web.Service.Name, web.GetCommands(cfg.Web))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.WebDAV.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.WebDAV.Service.Name, webdav.GetCommands(cfg.WebDAV))
	},
	func(cfg *config.Config) *cli.Command {
		cfg.Webfinger.Commons = cfg.Commons
		return ServiceCommand(cfg, cfg.Webfinger.Service.Name, webfinger.GetCommands(cfg.Webfinger))
	},
}

// ServiceCommand is the entry point for the all service commands.
func ServiceCommand(cfg *config.Config, servicename string, subcommands []*cli.Command) *cli.Command {
	return &cli.Command{
		Name:     servicename,
		Usage:    helper.SubcommandDescription(servicename),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			return nil
		},
		Subcommands: subcommands,
	}
}

func init() {
	for _, c := range svccmds {
		register.AddCommand(c)
	}
}
