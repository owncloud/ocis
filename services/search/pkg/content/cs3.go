package content

import (
	"context"
	"crypto/tls"
	"fmt"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"io"
	"net/http"
)

type cs3 struct {
	client   gateway.GatewayAPIClient
	logger   log.Logger
	insecure bool
	secret   string
}

func newCS3Retriever(client gateway.GatewayAPIClient, logger log.Logger, insecure bool) cs3 {
	return cs3{
		client:   client,
		insecure: insecure,
		logger:   logger,
	}
}

// Retrieve downloads the file from a cs3 service
// The caller MUST make sure to close the returned ReadCloser
func (s cs3) Retrieve(ctx context.Context, rid *provider.ResourceId) (io.ReadCloser, error) {
	at, ok := contextGet(ctx, revactx.TokenHeader)
	if !ok {
		return nil, fmt.Errorf("context without %s", revactx.TokenHeader)
	}

	res, err := s.client.InitiateFileDownload(ctx, &provider.InitiateFileDownloadRequest{Ref: &provider.Reference{ResourceId: rid, Path: "."}})
	if err != nil {
		return nil, err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("could not load resoure: %s", res.Status.Message)
	}

	var ep, tt string
	for _, p := range res.Protocols {
		if p.Protocol == "spaces" {
			ep, tt = p.DownloadEndpoint, p.Token
			break
		}
	}
	if (ep == "" || tt == "") && len(res.Protocols) > 0 {
		ep, tt = res.Protocols[0].DownloadEndpoint, res.Protocols[0].Token
	}

	req, err := http.NewRequest("GET", ep, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(revactx.TokenHeader, at)
	req.Header.Set("X-Reva-Transfer", tt)

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
		return nil, fmt.Errorf("could not download resource. Request returned with statuscode %d ", cres.StatusCode)
	}

	return cres.Body, nil
}
