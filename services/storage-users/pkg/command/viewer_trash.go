package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/reva/v2/pkg/conversions"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/storage/utils/grants"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

type viewerTrashResult string

const (
	_viewerTrashWouldFix viewerTrashResult = "would fix"
	_viewerTrashFixed    viewerTrashResult = "fixed"
	_viewerTrashFailed   viewerTrashResult = "failed"
	_viewerTrashCheck    viewerTrashResult = "check manually"
)

// viewerTrashClient is the part of the storage provider API that FixViewerTrash needs.
type viewerTrashClient interface {
	ListStorageSpaces(ctx context.Context, in *provider.ListStorageSpacesRequest, opts ...grpc.CallOption) (*provider.ListStorageSpacesResponse, error)
	ListGrants(ctx context.Context, in *provider.ListGrantsRequest, opts ...grpc.CallOption) (*provider.ListGrantsResponse, error)
	UpdateGrant(ctx context.Context, in *provider.UpdateGrantRequest, opts ...grpc.CallOption) (*provider.UpdateGrantResponse, error)
}

// viewerTrashEntry is one row of the report.
type viewerTrashEntry struct {
	spaceID   string
	spaceName string
	member    string
	result    viewerTrashResult
	reason    error // why it failed
}

// viewerTrashReport is what a run of FixViewerTrash found and did.
type viewerTrashReport struct {
	spacesChecked int
	entries       []viewerTrashEntry
}

func (r viewerTrashReport) countMatchingResults(result viewerTrashResult) int {
	n := 0
	for _, e := range r.entries {
		if string(e.result) == string(result) {
			n++
		}
	}
	return n
}

// FixViewerTrash removes the trash bin access from space members that got the
// "Can view" role before the role lost it (OCISDEV-1546).
//
// A member's permissions are copied into the storage when the member is added,
// so changing the role does not change existing members. This command rewrites
// those members through the running storage provider, logged in as the service
// account, so locking, caching and events behave like a role change in the web UI.
func FixViewerTrash(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "fix-viewer-trash",
		Usage: `remove the trash bin access from "Can view" space members that were added before the role lost it`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "Only show which members would be changed",
				Value:   true, // default must be true to avoid accidental changes
				Aliases: []string{"d"},
			},
			&cli.StringFlag{
				Name:    "space-id",
				Usage:   "Only check this space (omit to check all project spaces)",
				Aliases: []string{"s"},
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			dryRun := c.Bool("dry-run")

			if !dryRun && conversions.NewSpaceViewerRole().CS3ResourcePermissions().GetListRecycle() {
				return errors.New(`this oCIS version still gives the "Can view" role access to the trash bin; update oCIS first, otherwise the changed members lose the download button`)
			}

			// find the services through the registry, the same way the storage-users server does
			reg := registry.GetRegistry()
			gatewayClient, err := pool.GetGatewayServiceClient(cfg.Reva.Address, pool.WithRegistry(reg))
			if err != nil {
				return fmt.Errorf("error selecting gateway client: %w", err)
			}
			ctx, err := utils.GetServiceUserContext(cfg.ServiceAccount.ServiceAccountID, gatewayClient, cfg.ServiceAccount.ServiceAccountSecret)
			if err != nil {
				return fmt.Errorf("could not get service user context: %w", err)
			}

			// talk to this storage-users service only (storage-users or storage-users-vault)
			storageAddr := cfg.GRPC.Namespace + "." + cfg.Service.Name
			providerClient, err := pool.GetStorageProviderServiceClient(storageAddr, pool.WithRegistry(reg))
			if err != nil {
				return fmt.Errorf("error selecting storage provider client %s: %w", storageAddr, err)
			}
			spacesClient, err := pool.GetSpacesProviderServiceClient(storageAddr, pool.WithRegistry(reg))
			if err != nil {
				return fmt.Errorf("error selecting spaces provider client %s: %w", storageAddr, err)
			}
			client := struct {
				provider.ProviderAPIClient
				provider.SpacesAPIClient
			}{providerClient, spacesClient}

			fmt.Printf("checking the project spaces of %s (dry-run: %v)\n", storageAddr, dryRun)
			report, err := fixViewerTrash(ctx, client, c.String("space-id"), dryRun)
			if err != nil {
				return err
			}
			if err := printViewerTrashReport(os.Stdout, report, dryRun); err != nil {
				return err
			}

			if failed := report.countMatchingResults(_viewerTrashFailed); failed > 0 {
				return fmt.Errorf("%d member(s) could not be changed", failed)
			}
			return nil
		},
	}
}

// oldSpaceViewerPermissions returns the permissions that "Can view" members got
// before the role lost the trash bin access.
func oldSpaceViewerPermissions() *provider.ResourcePermissions {
	p := conversions.NewSpaceViewerRole().CS3ResourcePermissions()
	p.ListRecycle = true
	return p
}

// fixedSpaceViewerPermissions returns the permissions of the "Can view" role
// without the trash bin access.
func fixedSpaceViewerPermissions() *provider.ResourcePermissions {
	p := conversions.NewSpaceViewerRole().CS3ResourcePermissions()
	p.ListRecycle = false
	return p
}

// fixViewerTrash checks the members of the project spaces and removes the trash
// bin access from those with the old "Can view" permissions. If spaceIDFilter is
// empty, all project spaces are checked, otherwise only that space. With dryRun
// nothing is changed, the report only lists what would be changed.
func fixViewerTrash(ctx context.Context, client viewerTrashClient, spaceIDFilter string, dryRun bool) (viewerTrashReport, error) {
	report := viewerTrashReport{}

	res, err := client.ListStorageSpaces(ctx, &provider.ListStorageSpacesRequest{
		Filters: []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &provider.ListStorageSpacesRequest_Filter_SpaceType{SpaceType: _spaceTypeProject},
			},
		},
	})
	if err != nil {
		return report, fmt.Errorf("failed to list spaces: %w", err)
	}
	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return report, fmt.Errorf("failed to list spaces: %s", statusString(res.GetStatus()))
	}

	oldPerms := oldSpaceViewerPermissions()
	fixedPerms := fixedSpaceViewerPermissions()
	found := false

	for _, space := range res.GetStorageSpaces() {
		if spaceIDFilter != "" && !sameSpace(space.GetId().GetOpaqueId(), spaceIDFilter) {
			continue
		}
		found = true
		report.spacesChecked++

		// show the space ID in the "storage$space" form that the web UI and the
		// Graph API use, so it can be passed to --space-id as it is
		spaceEntry := viewerTrashEntry{
			spaceID:   storagespace.FormatStorageID(space.GetRoot().GetStorageId(), space.GetRoot().GetSpaceId()),
			spaceName: space.GetName(),
		}

		grantsRes, err := client.ListGrants(ctx, &provider.ListGrantsRequest{
			Ref: &provider.Reference{ResourceId: space.GetRoot()},
		})
		switch {
		case err != nil:
			spaceEntry.result, spaceEntry.reason = _viewerTrashFailed, fmt.Errorf("could not list members: %w", err)
			report.entries = append(report.entries, spaceEntry)
			continue
		case grantsRes.GetStatus().GetCode() != rpc.Code_CODE_OK:
			spaceEntry.result, spaceEntry.reason = _viewerTrashFailed, fmt.Errorf("could not list members: %s", statusString(grantsRes.GetStatus()))
			report.entries = append(report.entries, spaceEntry)
			continue
		}

		for _, g := range grantsRes.GetGrants() {
			perms := g.GetPermissions()
			entry := spaceEntry
			entry.member = granteeString(g.GetGrantee())

			switch {
			case grants.PermissionsEqual(perms, oldPerms):
				if dryRun {
					entry.result = _viewerTrashWouldFix
					break
				}
				if err := removeTrashAccess(ctx, client, space, g, fixedPerms); err != nil {
					entry.result, entry.reason = _viewerTrashFailed, err
					break
				}
				entry.result = _viewerTrashFixed
			case perms.GetListRecycle() && !perms.GetInitiateFileUpload():
				// A read-only member that can still list the trash bin, but whose
				// permissions are not exactly the old "Can view" ones. Leave it alone.
				entry.result = _viewerTrashCheck
			default:
				continue
			}
			report.entries = append(report.entries, entry)
		}
	}

	if spaceIDFilter != "" && !found {
		return report, fmt.Errorf("project space not found: %s", spaceIDFilter)
	}
	return report, nil
}

// removeTrashAccess writes the member back with the fixed permissions. It keeps
// who added the member and the expiration date.
func removeTrashAccess(ctx context.Context, client viewerTrashClient, space *provider.StorageSpace, g *provider.Grant, fixedPerms *provider.ResourcePermissions) error {
	// mark the request as a space member change, the storage needs the space
	// type to keep its member indexes right
	opaque := utils.AppendPlainToOpaque(nil, "spacegrant", "")
	opaque = utils.AppendPlainToOpaque(opaque, "spacetype", space.GetSpaceType())

	res, err := client.UpdateGrant(ctx, &provider.UpdateGrantRequest{
		Opaque: opaque,
		Ref:    &provider.Reference{ResourceId: space.GetRoot()},
		Grant: &provider.Grant{
			Grantee:     g.GetGrantee(),
			Permissions: fixedPerms,
			Expiration:  g.GetExpiration(),
			Creator:     g.GetCreator(),
		},
	})
	if err != nil {
		return err
	}
	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return errors.New(statusString(res.GetStatus()))
	}
	return nil
}

func statusString(s *rpc.Status) string {
	if s.GetMessage() == "" {
		return s.GetCode().String()
	}
	return s.GetCode().String() + " " + s.GetMessage()
}

// sameSpace compares two space IDs by their space part, so that
// "storage$space", "storage$space!node" and "space" all match.
func sameSpace(a, b string) bool {
	_, spaceA, _, errA := storagespace.SplitID(a)
	_, spaceB, _, errB := storagespace.SplitID(b)
	return errA == nil && errB == nil && spaceA == spaceB
}

func granteeString(g *provider.Grantee) string {
	switch g.GetType() {
	case provider.GranteeType_GRANTEE_TYPE_USER:
		return "user " + g.GetUserId().GetOpaqueId()
	case provider.GranteeType_GRANTEE_TYPE_GROUP:
		return "group " + g.GetGroupId().GetOpaqueId()
	default:
		return g.GetType().String()
	}
}

// printViewerTrashReport writes the report table and the summary line to w. The
// table is the only record of which members were changed, so a write error is
// returned and must not be ignored.
func printViewerTrashReport(w io.Writer, report viewerTrashReport, dryRun bool) error {
	if len(report.entries) > 0 {
		table := tablewriter.NewTable(w)
		table.Header("space id", "space name", "member", "result")
		for _, e := range report.entries {
			result := string(e.result)
			if e.reason != nil {
				result += ": " + e.reason.Error()
			}
			if err := table.Append([]string{e.spaceID, e.spaceName, e.member, result}); err != nil {
				return fmt.Errorf("could not add a row to the report: %w", err)
			}
		}
		if err := table.Render(); err != nil {
			return fmt.Errorf("could not print the report: %w", err)
		}
	}

	changed := report.countMatchingResults(_viewerTrashFixed)
	verb := "changed"
	if dryRun {
		changed = report.countMatchingResults(_viewerTrashWouldFix)
		verb = "would change"
	}
	if _, err := fmt.Fprintf(w, "spaces checked: %d, members %s: %d, failed: %d, check manually: %d\n",
		report.spacesChecked, verb, changed, report.countMatchingResults(_viewerTrashFailed), report.countMatchingResults(_viewerTrashCheck)); err != nil {
		return fmt.Errorf("could not print the summary: %w", err)
	}
	return nil
}
