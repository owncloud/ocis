package search

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/reva/v2/pkg/conversions"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/storage/utils/grants"
	"github.com/owncloud/reva/v2/pkg/utils"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
)

var scopeRegex = regexp.MustCompile(`scope:\s*([^" "\n\r]*)`)
var vaultRegex = regexp.MustCompile(`vault:\s*(\S+)`)

// errResourceNotFound is returned by statResource when Stat answers
// CODE_NOT_FOUND: the resource was deleted or moved since it was listed. The
// gateway also answers NOT_FOUND when it cannot reach the storage provider, so
// a run of them can mean an outage (see IndexSpace).
var errResourceNotFound = errors.New("resource not found")

// errPathNotResolved is returned by ResolveReference when GetPath answers with
// a non-OK status. When the status is CODE_NOT_FOUND, the error also matches
// errResourceNotFound.
var errPathNotResolved = errors.New("could not resolve the resource's path")

// ResolveReference makes sure the path is relative to the space root
func ResolveReference(ctx context.Context, ref *provider.Reference, ri *provider.ResourceInfo, gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) (*provider.Reference, error) {
	if ref.GetResourceId().GetOpaqueId() == ref.GetResourceId().GetSpaceId() {
		return ref, nil
	}

	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	gpRes, err := gatewayClient.GetPath(ctx, &provider.GetPathRequest{
		ResourceId: ri.GetId(),
	})
	if err != nil {
		return nil, err
	}
	if code := gpRes.GetStatus().GetCode(); code != rpc.Code_CODE_OK {
		if code == rpc.Code_CODE_NOT_FOUND {
			// Deleted since the Stat, like a Stat that answers NOT_FOUND.
			return nil, fmt.Errorf("%w: %w: %s", errPathNotResolved, errResourceNotFound, gpRes.GetStatus().GetMessage())
		}
		return nil, fmt.Errorf("%w: %s %s", errPathNotResolved, code, gpRes.GetStatus().GetMessage())
	}
	return &provider.Reference{
		ResourceId: &provider.ResourceId{
			StorageId: ref.GetResourceId().GetStorageId(),
			SpaceId:   ref.GetResourceId().GetSpaceId(),
			OpaqueId:  ref.GetResourceId().GetSpaceId(),
		},
		Path: utils.MakeRelativePath(gpRes.GetPath()),
	}, nil
}

type matchArray []*searchmsg.Match

func (ma matchArray) Len() int {
	return len(ma)
}
func (ma matchArray) Swap(i, j int) {
	ma[i], ma[j] = ma[j], ma[i]
}
func (ma matchArray) Less(i, j int) bool {
	return ma[i].GetScore() > ma[j].GetScore()
}

func logDocCount(engine engine.Engine, logger log.Logger) {
	c, err := engine.DocCount()
	if err != nil {
		logger.Error().Err(err).Msg("error getting document count from the index")
	}
	logger.Debug().Interface("count", c).Msg("new document count")
}

func getAuthContext(ctx context.Context, serviceAccountID string, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], secret string, logger log.Logger) (context.Context, error) {
	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		logger.Error().Err(err).Msg("could not get reva gatewayClient")
		return nil, err
	}

	return utils.GetServiceUserContextWithContext(ctx, gatewayClient, serviceAccountID, secret)
}

func statResource(ctx context.Context, ref *provider.Reference, gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) (*provider.StatResponse, error) {
	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	res, err := gatewayClient.Stat(ctx, &provider.StatRequest{Ref: ref})
	if err != nil {
		return nil, err
	}
	switch code := res.GetStatus().GetCode(); code {
	case rpc.Code_CODE_OK:
		return res, nil
	case rpc.Code_CODE_NOT_FOUND:
		// Moved or deleted in the meantime, or a backend the gateway could
		// not reach: the gateway answers NOT_FOUND for both.
		return nil, errResourceNotFound
	default:
		return nil, fmt.Errorf("stat failed: %s %s", code, res.GetStatus().GetMessage())
	}
}

// NOTE: this converts CS3 to WebDAV permissions
// since conversions pkg is reva internal we have no other choice than to duplicate the logic
func convertToWebDAVPermissions(isShared, isMountpoint, isDir bool, p *provider.ResourcePermissions) string {
	if p == nil {
		return ""
	}
	var b strings.Builder
	if isShared {
		fmt.Fprintf(&b, "S")
	}
	if p.GetListContainer() &&
		p.GetListFileVersions() &&
		p.GetListRecycle() &&
		p.GetStat() &&
		p.GetGetPath() &&
		p.GetGetQuota() &&
		p.GetInitiateFileDownload() {
		fmt.Fprintf(&b, "R")
	}
	if isMountpoint {
		fmt.Fprintf(&b, "M")
	}
	if p.GetDelete() {
		fmt.Fprintf(&b, "D")
	}
	if p.GetInitiateFileUpload() &&
		p.GetRestoreFileVersion() &&
		p.GetRestoreRecycleItem() {
		fmt.Fprintf(&b, "NV")
		if !isDir {
			fmt.Fprintf(&b, "W")
		}
	}
	if isDir &&
		p.GetListContainer() &&
		p.GetStat() &&
		p.GetCreateContainer() &&
		p.GetInitiateFileUpload() {
		fmt.Fprintf(&b, "CK")
	}
	if grants.PermissionsEqual(p, conversions.NewSecureViewerRole().CS3ResourcePermissions()) {
		fmt.Fprintf(&b, "X")
	}
	return b.String()
}

// ParseVaultMode extracts "vault:<value>" from the query string and returns the
// cleaned query and whether vault mode was requested (value == "true").
func ParseVaultMode(query string) (string, bool) {
	match := vaultRegex.FindStringSubmatch(query)
	if len(match) < 2 {
		return query, false
	}
	cleaned := strings.TrimSpace(strings.ReplaceAll(query, match[0], ""))
	return cleaned, match[1] == "true"
}

// ParseScope extract a scope value from the query string and returns search, scope strings
func ParseScope(query string) (string, string) {
	match := scopeRegex.FindStringSubmatch(query)
	if len(match) >= 2 {
		cut := match[0]
		return strings.TrimSpace(strings.ReplaceAll(query, cut, "")), strings.TrimSpace(match[1])
	}
	return query, ""
}
