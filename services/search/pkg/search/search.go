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
	"github.com/cs3org/reva/v2/pkg/conversions"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storage/utils/grants"
	"github.com/cs3org/reva/v2/pkg/utils"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
)

var scopeRegex = regexp.MustCompile(`scope:\s*([^" "\n\r]*)`)

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
	if err != nil || gpRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, err
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

func getAuthContext(serviceAccountID string, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], secret string, logger log.Logger) (context.Context, error) {
	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		logger.Error().Err(err).Msg("could not get reva gatewayClient")
		return nil, err
	}

	return utils.GetServiceUserContext(serviceAccountID, gatewayClient, secret)
}

func statResource(ctx context.Context, ref *provider.Reference, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], logger log.Logger) (*provider.StatResponse, error) {
	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	res, err := gatewayClient.Stat(ctx, &provider.StatRequest{Ref: ref})
	if err != nil {
		logger.Error().Err(err).Msg("failed to stat the moved resource")
		return nil, err
	}
	switch res.GetStatus().GetCode() {
	case rpc.Code_CODE_OK:
		return res, nil
	case rpc.Code_CODE_NOT_FOUND:
		// Resource was moved or deleted in the meantime. ignore.
		return nil, err
	default:
		err := errors.New("failed to stat the moved resource")
		logger.Error().Interface("res", res).Msg(err.Error())
		return nil, err
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

// ParseScope extract a scope value from the query string and returns search, scope strings
func ParseScope(query string) (string, string) {
	match := scopeRegex.FindStringSubmatch(query)
	if len(match) >= 2 {
		cut := match[0]
		return strings.TrimSpace(strings.ReplaceAll(query, cut, "")), strings.TrimSpace(match[1])
	}
	return query, ""
}
