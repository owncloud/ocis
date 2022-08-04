package content

import (
	"context"
	"crypto/tls"
	"fmt"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
)

type cs3 struct {
	client   gateway.GatewayAPIClient
	logger   log.Logger
	insecure bool
	secret   string
}

func newCS3Retriever(client gateway.GatewayAPIClient, logger log.Logger, secret string, insecure bool) cs3 {
	return cs3{
		client:   client,
		insecure: insecure,
		secret:   secret,
		logger:   logger,
	}
}

// Retrieve downloads the file from a cs3 service
// The caller MUST make sure to close the returned ReadCloser
func (s cs3) Retrieve(ctx context.Context, ref *provider.Reference, owner *user.User) (io.ReadCloser, error) {
	authRes, err := s.client.Authenticate(
		ctxpkg.ContextSetUser(ctx, owner),
		&gateway.AuthenticateRequest{
			Type:         "machine",
			ClientId:     "userid:" + owner.GetId().GetOpaqueId(),
			ClientSecret: s.secret,
		},
	)
	if err == nil && authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		err = errtypes.NewErrtypeFromStatus(authRes.Status)
	}
	if err != nil {
		s.logger.Error().Err(err).Interface("owner", owner).Interface("authRes", authRes).Msg("error using machine auth")
		return nil, err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, authRes.Token)
	res, err := s.client.InitiateFileDownload(ctx, &provider.InitiateFileDownloadRequest{Ref: ref})
	if err != nil {
		return nil, err
	}

	if res.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("could not load resoure: %s", res.Status.Message)
	}

	var ep, tk string
	for _, p := range res.Protocols {
		if p.Protocol == "spaces" {
			ep, tk = p.DownloadEndpoint, p.Token
			break
		}
	}
	if (ep == "" || tk == "") && len(res.Protocols) > 0 {
		ep, tk = res.Protocols[0].DownloadEndpoint, res.Protocols[0].Token
	}

	req, err := http.NewRequest("GET", ep, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(revactx.TokenHeader, authRes.Token)
	req.Header.Set("X-Reva-Transfer", tk)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: s.insecure},
		},
	}

	cres, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if cres.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get the resource \"%s\". Request returned with statuscode %d ", ref.Path, cres.StatusCode)
	}

	return cres.Body, nil
}
