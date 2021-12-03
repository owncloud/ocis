package cs3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/cs3org/reva/pkg/token"
	"google.golang.org/grpc/metadata"
)

const (
	storageMountPath = "/meta"
)

type metadataStorage struct {
	tokenManager      token.Manager
	storageProvider   provider.ProviderAPIClient
	dataGatewayClient *http.Client
}

func (r metadataStorage) uploadHelper(ctx context.Context, uploadpath string, content []byte) error {

	ref := provider.InitiateFileUploadRequest{
		Ref: &provider.Reference{
			Path: path.Join(storageMountPath, uploadpath),
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

	md, _ := metadata.FromOutgoingContext(ctx)
	req.Header.Add(revactx.TokenHeader, md.Get(revactx.TokenHeader)[0])
	resp, err := r.dataGatewayClient.Do(req)
	if err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		return err
	}
	return nil
}

func (r metadataStorage) downloadHelper(ctx context.Context, downloadpath string) (content []byte, err error) {
	ref := provider.InitiateFileDownloadRequest{
		Ref: &provider.Reference{
			Path: path.Join(storageMountPath, downloadpath),
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

	md, _ := metadata.FromOutgoingContext(ctx)
	req.Header.Add(revactx.TokenHeader, md.Get(revactx.TokenHeader)[0])
	resp, err := r.dataGatewayClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, &notFoundErr{}
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if err = resp.Body.Close(); err != nil {
		return []byte{}, err
	}

	return b, nil
}

type notFoundErr struct {
	typ, id string
}

func (e notFoundErr) Error() string {
	return fmt.Sprintf("%s with id %s not found", e.typ, e.id)
}

// IsNotFoundErr can be returned by repo Load and Delete operations
func IsNotFoundErr(e error) bool {
	_, ok := e.(*notFoundErr)
	return ok
}
