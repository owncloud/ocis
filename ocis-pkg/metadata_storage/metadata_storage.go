package metadataStorage

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/cs3org/reva/pkg/errtypes"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/utils"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"google.golang.org/grpc/metadata"
)

func NewMetadataStorage(providerAddr string) (s MetadataStorage, err error) {
	p, err := pool.GetStorageProviderServiceClient(providerAddr)
	if err != nil {
		return MetadataStorage{}, err
	}

	c := http.DefaultClient

	return MetadataStorage{
		storageProvider:   p,
		dataGatewayClient: c,
	}, nil
}

type MetadataStorage struct {
	storageProvider   provider.ProviderAPIClient
	dataGatewayClient *http.Client
	SpaceRoot         *provider.ResourceId
}

// init creates the metadata space
func (ms *MetadataStorage) Init(ctx context.Context, serviceUser config.ServiceUser) (err error) {
	// FIXME change CS3 api to allow sending a space id
	cssr, err := ms.storageProvider.CreateStorageSpace(ctx, &provider.CreateStorageSpaceRequest{
		Opaque: &typesv1beta1.Opaque{
			Map: map[string]*typesv1beta1.OpaqueEntry{
				"spaceid": {
					Decoder: "plain",
					Value:   []byte(serviceUser.UUID),
				},
			},
		},
		Owner: &user.User{
			Id: &user.UserId{
				OpaqueId: serviceUser.UUID,
			},
			Groups:    []string{},
			UidNumber: serviceUser.UID,
			GidNumber: serviceUser.GID,
		},
		Name: "Metadata",
		Type: "metadata",
	})
	switch {
	case err != nil:
		return err
	case cssr.Status.Code == v1beta11.Code_CODE_OK:
		ms.SpaceRoot = cssr.StorageSpace.Root
	case cssr.Status.Code == v1beta11.Code_CODE_ALREADY_EXISTS:
		// TODO make CreateStorageSpace return existing space?
		ms.SpaceRoot = &provider.ResourceId{StorageId: serviceUser.UUID, OpaqueId: serviceUser.UUID}
	default:
		return errtypes.NewErrtypeFromStatus(cssr.Status)
	}
	return nil
}

func (ms MetadataStorage) SimpleUpload(ctx context.Context, uploadpath string, content []byte) error {

	ref := provider.InitiateFileUploadRequest{
		Ref: &provider.Reference{
			ResourceId: ms.SpaceRoot,
			Path:       utils.MakeRelativePath(uploadpath),
		},
	}

	res, err := ms.storageProvider.InitiateFileUpload(ctx, &ref)
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
	resp, err := ms.dataGatewayClient.Do(req)
	if err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		return err
	}
	return nil
}

func (ms MetadataStorage) SimpleDownload(ctx context.Context, downloadpath string) (content []byte, err error) {
	ref := provider.InitiateFileDownloadRequest{
		Ref: &provider.Reference{
			ResourceId: ms.SpaceRoot,
			Path:       utils.MakeRelativePath(downloadpath),
		},
	}

	res, err := ms.storageProvider.InitiateFileDownload(ctx, &ref)
	if err != nil {
		return []byte{}, err
	}

	var endpoint string

	for _, proto := range res.GetProtocols() {
		if proto.Protocol == "spaces" {
			endpoint = proto.GetDownloadEndpoint()
			break
		}
	}
	if endpoint == "" {
		return []byte{}, errors.New("metadata storage doesn't support the spaces download protocol")
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return []byte{}, err
	}

	md, _ := metadata.FromOutgoingContext(ctx)
	req.Header.Add(revactx.TokenHeader, md.Get(revactx.TokenHeader)[0])
	resp, err := ms.dataGatewayClient.Do(req)
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
