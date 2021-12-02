package cs3

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/token"
)

type metadataStorage struct {
	tokenManager      token.Manager
	storageProvider   provider.ProviderAPIClient
	dataGatewayClient *http.Client
}

func (r metadataStorage) uploadHelper(ctx context.Context, path string, content []byte) error {

	ref := provider.InitiateFileUploadRequest{
		Ref: &provider.Reference{
			Path: path,
		},
	}

	res, err := r.storageProvider.InitiateFileUpload(ctx, &ref)
	if err != nil {
		return err
	}

	var endpoint string

	for _, proto := range res.GetProtocols() {
		if proto.Protocol == "simple" {
			endpoint = proto.GetUploadEndpoint()
			break
		}
	}
	if endpoint == "" {
		return errors.New("metadata storage doesn't support the simple upload protocol")
	}

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(content))
	if err != nil {
		return err
	}
	resp, err := r.dataGatewayClient.Do(req)
	if err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		return err
	}
	return nil
}

func (r metadataStorage) downloadHelper(ctx context.Context, path string) (content []byte, err error) {

	ref := provider.InitiateFileDownloadRequest{
		Ref: &provider.Reference{
			Path: path,
		},
	}

	res, err := r.storageProvider.InitiateFileDownload(ctx, &ref)
	if err != nil {
		return []byte{}, err
	}

	var endpoint string

	for _, proto := range res.GetProtocols() {
		if proto.Protocol == "simple" {
			endpoint = proto.GetDownloadEndpoint()
			break
		}
	}
	if endpoint == "" {
		return []byte{}, errors.New("metadata storage doesn't support the simple download protocol")
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return []byte{}, err
	}
	resp, err := r.dataGatewayClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	//if resp.StatusCode != http.StatusOK {
	//	return []byte{}, &notFoundErr{}
	//}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if err = resp.Body.Close(); err != nil {
		return []byte{}, err
	}

	return b, nil
}
