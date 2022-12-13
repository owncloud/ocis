package search

import (
	"context"
	"errors"
	"fmt"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
	"google.golang.org/grpc/metadata"
)

// ResolveReference makes sure the path is relative to the space root
func ResolveReference(ctx context.Context, ref *provider.Reference, ri *provider.ResourceInfo, gw gateway.GatewayAPIClient) (*provider.Reference, error) {
	if ref.GetResourceId().GetOpaqueId() == ref.GetResourceId().GetSpaceId() {
		return ref, nil
	}

	gpRes, err := gw.GetPath(ctx, &provider.GetPathRequest{
		ResourceId: ri.Id,
	})
	if err != nil || gpRes.Status.Code != rpc.Code_CODE_OK {
		return nil, err
	}
	return &provider.Reference{
		ResourceId: &provider.ResourceId{
			StorageId: ref.GetResourceId().GetStorageId(),
			SpaceId:   ref.GetResourceId().GetSpaceId(),
			OpaqueId:  ref.GetResourceId().GetSpaceId(),
		},
		Path: utils.MakeRelativePath(gpRes.Path),
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
	return ma[i].Score > ma[j].Score
}

func logDocCount(engine engine.Engine, logger log.Logger) {
	c, err := engine.DocCount()
	if err != nil {
		logger.Error().Err(err).Msg("error getting document count from the index")
	}
	logger.Debug().Interface("count", c).Msg("new document count")
}

func getAuthContext(owner *user.User, gw gateway.GatewayAPIClient, secret string, logger log.Logger) (context.Context, error) {
	ownerCtx := ctxpkg.ContextSetUser(context.Background(), owner)
	authRes, err := gw.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + owner.GetId().GetOpaqueId(),
		ClientSecret: secret,
	})

	if err == nil && authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		err = errtypes.NewErrtypeFromStatus(authRes.Status)
	}

	if err != nil {
		logger.Error().Err(err).Interface("owner", owner).Interface("authRes", authRes).Msg("error using machine auth")
		return nil, err
	}

	return metadata.AppendToOutgoingContext(ownerCtx, ctxpkg.TokenHeader, authRes.Token), nil
}

func statResource(ctx context.Context, ref *provider.Reference, gw gateway.GatewayAPIClient, logger log.Logger) (*provider.StatResponse, error) {
	res, err := gw.Stat(ctx, &provider.StatRequest{Ref: ref})
	if err != nil {
		logger.Error().Err(err).Msg("failed to stat the moved resource")
		return nil, err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		err := errors.New("failed to stat the moved resource")
		logger.Error().Interface("res", res).Msg(err.Error())
		return nil, err
	}

	return res, nil
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
	if p.ListContainer &&
		p.ListFileVersions &&
		p.ListRecycle &&
		p.Stat &&
		p.GetPath &&
		p.GetQuota &&
		p.InitiateFileDownload {
		fmt.Fprintf(&b, "R")
	}
	if isMountpoint {
		fmt.Fprintf(&b, "M")
	}
	if p.Delete &&
		p.PurgeRecycle {
		fmt.Fprintf(&b, "D")
	}
	if p.InitiateFileUpload &&
		p.RestoreFileVersion &&
		p.RestoreRecycleItem {
		fmt.Fprintf(&b, "NV")
		if !isDir {
			fmt.Fprintf(&b, "W")
		}
	}
	if isDir &&
		p.ListContainer &&
		p.Stat &&
		p.CreateContainer &&
		p.InitiateFileUpload {
		fmt.Fprintf(&b, "CK")
	}
	return b.String()
}
