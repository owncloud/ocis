package auth

import (
	"context"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	rstatus "github.com/owncloud/reva/v2/pkg/rgrpc/status"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func mfaResponse(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) (interface{}, error) {
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
}
