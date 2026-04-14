package mfa

import (
	"context"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/reva/v2/pkg/appctx"
	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/rgrpc"
	rstatus "github.com/owncloud/reva/v2/pkg/rgrpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

const (
	// defaultPriority places mfa before readonly (200) so the MFA check
	// runs first and returns a clear PermissionDenied rather than a readonly error.
	defaultPriority = 150
)

func init() {
	rgrpc.RegisterUnaryInterceptor("mfa", NewUnary)
}

// NewUnary returns a new unary interceptor that requires MFA to be satisfied
// for every gRPC call on the vault storage provider.
// Service accounts (UserType_USER_TYPE_SERVICE) are exempt because they are
// used for internal operations (postprocessing, event handling, etc.) that
// never carry an MFA claim.
func NewUnary(map[string]interface{}) (grpc.UnaryServerInterceptor, int, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := appctx.GetLogger(ctx)

		// Bypass for service accounts — they perform internal operations and
		// never carry an MFA claim.
		if u, ok := ctxpkg.ContextGetUser(ctx); ok {
			if u.GetId().GetType() == userpb.UserType_USER_TYPE_SERVICE {
				return handler(ctx, req)
			}
		}

		hasMFA, _ := ctxpkg.ContextGetMFA(ctx)
		if hasMFA {
			return handler(ctx, req)
		}

		log.Warn().Str("method", info.FullMethod).Msg("mfa: access denied, MFA required")

		const msg = "MFA required to access vault storage"
		switch req.(type) {
		case *provider.StatRequest:
			return &provider.StatResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.ListContainerRequest:
			return &provider.ListContainerResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.GetPathRequest:
			return &provider.GetPathResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.GetQuotaRequest:
			return &provider.GetQuotaResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.InitiateFileDownloadRequest:
			return &provider.InitiateFileDownloadResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.InitiateFileUploadRequest:
			return &provider.InitiateFileUploadResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.CreateContainerRequest:
			return &provider.CreateContainerResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.TouchFileRequest:
			return &provider.TouchFileResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.DeleteRequest:
			return &provider.DeleteResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.MoveRequest:
			return &provider.MoveResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.CreateHomeRequest:
			return &provider.CreateHomeResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.AddGrantRequest:
			return &provider.AddGrantResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.RemoveGrantRequest:
			return &provider.RemoveGrantResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.UpdateGrantRequest:
			return &provider.UpdateGrantResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.ListGrantsRequest:
			return &provider.ListGrantsResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.ListFileVersionsRequest:
			return &provider.ListFileVersionsResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.RestoreFileVersionRequest:
			return &provider.RestoreFileVersionResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.ListRecycleRequest:
			return &provider.ListRecycleResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.RestoreRecycleItemRequest:
			return &provider.RestoreRecycleItemResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.PurgeRecycleRequest:
			return &provider.PurgeRecycleResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.SetArbitraryMetadataRequest:
			return &provider.SetArbitraryMetadataResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		case *provider.UnsetArbitraryMetadataRequest:
			return &provider.UnsetArbitraryMetadataResponse{Status: rstatus.NewPermissionDenied(ctx, nil, msg)}, nil
		default:
			log.Debug().Str("method", info.FullMethod).Msg("mfa: blocking unknown request type")
			return nil, grpcstatus.Errorf(codes.PermissionDenied, "mfa: %s: %T", msg, req)
		}
	}, defaultPriority, nil
}
