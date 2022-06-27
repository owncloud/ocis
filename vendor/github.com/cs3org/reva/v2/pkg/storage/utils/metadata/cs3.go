// Copyright 2018-2022 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package metadata

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	"google.golang.org/grpc/metadata"
)

// CS3 represents a metadata storage with a cs3 storage backend
type CS3 struct {
	providerAddr      string
	gatewayAddr       string
	serviceUser       *user.User
	machineAuthAPIKey string
	dataGatewayClient *http.Client
	SpaceRoot         *provider.ResourceId
}

// NewCS3Storage returns a new cs3 storage instance
func NewCS3Storage(gwAddr, providerAddr, serviceUserID, serviceUserIDP, machineAuthAPIKey string) (s Storage, err error) {
	c := http.DefaultClient

	return &CS3{
		providerAddr:      providerAddr,
		gatewayAddr:       gwAddr,
		dataGatewayClient: c,
		machineAuthAPIKey: machineAuthAPIKey,
		serviceUser: &user.User{
			Id: &user.UserId{
				OpaqueId: serviceUserID,
				Idp:      serviceUserIDP,
			},
		},
	}, nil
}

// Backend returns the backend name of the storage
func (cs3 *CS3) Backend() string {
	return "cs3"
}

// Init creates the metadata space
func (cs3 *CS3) Init(ctx context.Context, spaceid string) (err error) {
	client, err := cs3.providerClient()
	if err != nil {
		return err
	}

	ctx, err = cs3.getAuthContext(ctx)
	if err != nil {
		return err
	}
	// FIXME change CS3 api to allow sending a space id
	cssr, err := client.CreateStorageSpace(ctx, &provider.CreateStorageSpaceRequest{
		Opaque: &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"spaceid": {
					Decoder: "plain",
					Value:   []byte(spaceid),
				},
			},
		},
		Owner: cs3.serviceUser,
		Name:  "Metadata",
		Type:  "metadata",
	})
	switch {
	case err != nil:
		return err
	case cssr.Status.Code == rpc.Code_CODE_OK:
		cs3.SpaceRoot = cssr.StorageSpace.Root
	case cssr.Status.Code == rpc.Code_CODE_ALREADY_EXISTS:
		// TODO make CreateStorageSpace return existing space?
		cs3.SpaceRoot = &provider.ResourceId{StorageId: spaceid, OpaqueId: spaceid}
	default:
		return errtypes.NewErrtypeFromStatus(cssr.Status)
	}
	return nil
}

// SimpleUpload uploads a file to the metadata storage
func (cs3 *CS3) SimpleUpload(ctx context.Context, uploadpath string, content []byte) error {
	client, err := cs3.providerClient()
	if err != nil {
		return err
	}
	ctx, err = cs3.getAuthContext(ctx)
	if err != nil {
		return err
	}

	ref := provider.InitiateFileUploadRequest{
		Ref: &provider.Reference{
			ResourceId: cs3.SpaceRoot,
			Path:       utils.MakeRelativePath(uploadpath),
		},
	}

	res, err := client.InitiateFileUpload(ctx, &ref)
	if err != nil {
		return err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return status.NewErrorFromCode(res.Status.Code, "cs3 metadata SimpleUpload")
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
	req.Header.Add(ctxpkg.TokenHeader, md.Get(ctxpkg.TokenHeader)[0])
	resp, err := cs3.dataGatewayClient.Do(req)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// SimpleDownload reads a file from the metadata storage
func (cs3 *CS3) SimpleDownload(ctx context.Context, downloadpath string) (content []byte, err error) {
	client, err := cs3.providerClient()
	if err != nil {
		return nil, err
	}
	ctx, err = cs3.getAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	dreq := provider.InitiateFileDownloadRequest{
		Ref: &provider.Reference{
			ResourceId: cs3.SpaceRoot,
			Path:       utils.MakeRelativePath(downloadpath),
		},
	}

	res, err := client.InitiateFileDownload(ctx, &dreq)
	if err != nil {
		return []byte{}, errtypes.NotFound(dreq.Ref.Path)
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
	req.Header.Add(ctxpkg.TokenHeader, md.Get(ctxpkg.TokenHeader)[0])
	resp, err := cs3.dataGatewayClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, errtypes.NotFound(dreq.Ref.Path)
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

// Delete deletes a path
func (cs3 *CS3) Delete(ctx context.Context, path string) error {
	client, err := cs3.providerClient()
	if err != nil {
		return err
	}
	ctx, err = cs3.getAuthContext(ctx)
	if err != nil {
		return err
	}

	res, err := client.Delete(ctx, &provider.DeleteRequest{
		Ref: &provider.Reference{
			ResourceId: cs3.SpaceRoot,
			Path:       utils.MakeRelativePath(path),
		},
	})
	if err != nil {
		return err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return errtypes.NewErrtypeFromStatus(res.Status)
	}

	return nil
}

// ReadDir returns the entries in a given directory
func (cs3 *CS3) ReadDir(ctx context.Context, path string) ([]string, error) {
	client, err := cs3.providerClient()
	if err != nil {
		return nil, err
	}
	ctx, err = cs3.getAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	relPath := utils.MakeRelativePath(path)
	res, err := client.ListContainer(ctx, &provider.ListContainerRequest{
		Ref: &provider.Reference{
			ResourceId: cs3.SpaceRoot,
			Path:       relPath,
		},
	})

	if err != nil {
		return nil, err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return nil, errtypes.NewErrtypeFromStatus(res.Status)
	}

	entries := []string{}
	for _, ri := range res.Infos {
		entries = append(entries, ri.Path)
	}
	return entries, nil
}

// MakeDirIfNotExist will create a root node in the metadata storage. Requires an authenticated context.
func (cs3 *CS3) MakeDirIfNotExist(ctx context.Context, folder string) error {
	client, err := cs3.providerClient()
	if err != nil {
		return err
	}
	ctx, err = cs3.getAuthContext(ctx)
	if err != nil {
		return err
	}

	var rootPathRef = &provider.Reference{
		ResourceId: cs3.SpaceRoot,
		Path:       utils.MakeRelativePath(folder),
	}

	resp, err := client.Stat(ctx, &provider.StatRequest{
		Ref: rootPathRef,
	})

	if err != nil {
		return err
	}

	switch {
	case err != nil:
		return err
	case resp.Status.Code == rpc.Code_CODE_OK:
		// nothing to do in this case
	case resp.Status.Code == rpc.Code_CODE_NOT_FOUND:
		r, err := client.CreateContainer(ctx, &provider.CreateContainerRequest{
			Ref: rootPathRef,
		})

		if err != nil {
			return err
		}

		if r.Status.Code != rpc.Code_CODE_OK {
			return errtypes.NewErrtypeFromStatus(r.Status)
		}
	default:
		return errtypes.NewErrtypeFromStatus(resp.Status)
	}

	return nil
}

// CreateSymlink creates a symlink
func (cs3 *CS3) CreateSymlink(ctx context.Context, oldname, newname string) error {
	if _, err := cs3.ResolveSymlink(ctx, newname); err == nil {
		return os.ErrExist
	}

	return cs3.SimpleUpload(ctx, newname, []byte(oldname))
}

// ResolveSymlink resolves a symlink
func (cs3 *CS3) ResolveSymlink(ctx context.Context, name string) (string, error) {
	b, err := cs3.SimpleDownload(ctx, name)
	if err != nil {
		if errors.Is(err, errtypes.NotFound("")) {
			return "", os.ErrNotExist
		}
		return "", err
	}

	return string(b), err
}

func (cs3 *CS3) providerClient() (provider.ProviderAPIClient, error) {
	return pool.GetStorageProviderServiceClient(cs3.providerAddr)
}

func (cs3 *CS3) getAuthContext(ctx context.Context) (context.Context, error) {
	client, err := pool.GetGatewayServiceClient(cs3.gatewayAddr)
	if err != nil {
		return nil, err
	}

	authCtx := ctxpkg.ContextSetUser(context.Background(), cs3.serviceUser)
	authRes, err := client.Authenticate(authCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + cs3.serviceUser.Id.OpaqueId,
		ClientSecret: cs3.machineAuthAPIKey,
	})
	if err != nil {
		return nil, err
	}
	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, errtypes.NewErrtypeFromStatus(authRes.GetStatus())
	}
	authCtx = metadata.AppendToOutgoingContext(authCtx, ctxpkg.TokenHeader, authRes.Token)
	return authCtx, nil
}
