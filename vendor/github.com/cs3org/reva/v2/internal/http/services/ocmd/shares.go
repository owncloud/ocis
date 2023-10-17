// Copyright 2018-2021 CERN
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

package ocmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"reflect"
	"strings"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ocmcore "github.com/cs3org/go-cs3apis/cs3/ocm/core/v1beta1"
	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
)

type sharesHandler struct {
	gatewayAddr string
}

func (h *sharesHandler) init(c *Config) {
	h.gatewayAddr = c.GatewaySvc
}

func (h *sharesHandler) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodPost:
			h.createShare(w, r)
		default:
			WriteError(w, r, APIErrorInvalidParameter, "Only POST method is allowed", nil)
		}
	})
}

func (h *sharesHandler) createShare(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := appctx.GetLogger(ctx)
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	var shareWith, meshProvider, resource, providerID, owner string
	var protocol map[string]interface{}
	if err == nil && contentType == "application/json" {
		defer r.Body.Close()
		reqBody, err := io.ReadAll(r.Body)
		if err == nil {
			reqMap := make(map[string]interface{})
			err = json.Unmarshal(reqBody, &reqMap)
			if err == nil {
				meshProvider = reqMap["meshProvider"].(string) // FIXME: get this from sharedBy string?
				shareWith, protocol = reqMap["shareWith"].(string), reqMap["protocol"].(map[string]interface{})
				resource, owner = reqMap["name"].(string), reqMap["owner"].(string)
				// Note that if an OCM request were to go directly from a Nextcloud server
				// to a Reva server, it will (incorrectly) sends an integer provider_id instead a string one.
				// This doesn't happen when using the sciencemesh-nextcloud app, but in order to make the OCM
				// test suite pass, this code works around that:
				if reflect.ValueOf(reqMap["providerId"]).Kind() == reflect.Float64 {
					providerID = fmt.Sprintf("%d", int(math.Round(reqMap["providerId"].(float64))))
				} else {
					providerID = reqMap["providerId"].(string)
				}
			} else {
				WriteError(w, r, APIErrorInvalidParameter, "could not parse json request body", nil)
			}
		}
	} else {
		var protocolJSON string
		shareWith, protocolJSON, meshProvider = r.FormValue("shareWith"), r.FormValue("protocol"), r.FormValue("meshProvider")
		resource, providerID, owner = r.FormValue("name"), r.FormValue("providerId"), r.FormValue("owner")
		err = json.Unmarshal([]byte(protocolJSON), &protocol)
		if err != nil {
			WriteError(w, r, APIErrorInvalidParameter, "invalid protocol parameters", nil)
		}
	}

	if resource == "" || providerID == "" || owner == "" {
		WriteError(w, r, APIErrorInvalidParameter, "missing details about resource to be shared", nil)
		return
	}
	if shareWith == "" || protocol["name"] == "" || meshProvider == "" {
		WriteError(w, r, APIErrorInvalidParameter, "missing request parameters", nil)
		return
	}

	gatewayClient, err := pool.GetGatewayServiceClient(h.gatewayAddr)
	if err != nil {
		WriteError(w, r, APIErrorServerError, "error getting storage grpc client", err)
		return
	}

	clientIP, err := utils.GetClientIP(r)
	if err != nil {
		WriteError(w, r, APIErrorServerError, fmt.Sprintf("error retrieving client IP from request: %s", r.RemoteAddr), err)
		return
	}
	providerInfo := ocmprovider.ProviderInfo{
		Domain: meshProvider,
		Services: []*ocmprovider.Service{
			{
				Host: clientIP,
			},
		},
	}

	providerAllowedResp, err := gatewayClient.IsProviderAllowed(ctx, &ocmprovider.IsProviderAllowedRequest{
		Provider: &providerInfo,
	})
	if err != nil {
		WriteError(w, r, APIErrorServerError, "error sending a grpc is provider allowed request", err)
		return
	}
	if providerAllowedResp.Status.Code != rpc.Code_CODE_OK {
		WriteError(w, r, APIErrorUnauthenticated, "provider not authorized", errors.New(providerAllowedResp.Status.Message))
		return
	}

	shareWithParts := strings.Split(shareWith, "@")
	userRes, err := gatewayClient.GetUser(ctx, &userpb.GetUserRequest{
		UserId: &userpb.UserId{OpaqueId: shareWithParts[0]}, SkipFetchingUserGroups: true,
	})
	if err != nil {
		WriteError(w, r, APIErrorServerError, "error searching recipient", err)
		return
	}
	if userRes.Status.Code != rpc.Code_CODE_OK {
		WriteError(w, r, APIErrorNotFound, "user not found", errors.New(userRes.Status.Message))
		return
	}

	var permissions conversions.Permissions
	var token string
	options, ok := protocol["options"].(map[string]interface{})
	if !ok {
		WriteError(w, r, APIErrorInvalidParameter, "protocol: webdav token not provided", nil)
		return
	}

	token, ok = options["sharedSecret"].(string)
	if !ok {
		token, ok = options["token"].(string)
		if !ok {
			WriteError(w, r, APIErrorInvalidParameter, "protocol: webdav token not provided", nil)
			return
		}
	}
	var role *conversions.Role
	pval, ok := options["permissions"].(int)
	if !ok {
		WriteError(w, r, APIErrorInvalidParameter, "permissions not provided", nil)
		return
	}

	permissions, err = conversions.NewPermissions(pval)
	if err != nil {
		WriteError(w, r, APIErrorInvalidParameter, err.Error(), nil)
		return
	}
	role = conversions.RoleFromOCSPermissions(permissions)

	val, err := json.Marshal(role.CS3ResourcePermissions())
	if err != nil {
		WriteError(w, r, APIErrorServerError, "could not encode role", nil)
		return
	}

	ownerID := &userpb.UserId{
		OpaqueId: owner,
		Idp:      meshProvider,
		Type:     userpb.UserType_USER_TYPE_PRIMARY,
	}
	createShareReq := &ocmcore.CreateOCMCoreShareRequest{
		Name:       resource,
		ProviderId: providerID,
		Owner:      ownerID,
		ShareWith:  userRes.User.GetId(),
		Protocol: &ocmcore.Protocol{
			Name: protocol["name"].(string),
			Opaque: &types.Opaque{
				Map: map[string]*types.OpaqueEntry{
					"permissions": {
						Decoder: "json",
						Value:   val,
					},
					"token": {
						Decoder: "plain",
						Value:   []byte(token),
					},
				},
			},
		},
	}
	createShareResponse, err := gatewayClient.CreateOCMCoreShare(ctx, createShareReq)
	if err != nil {
		WriteError(w, r, APIErrorServerError, "error sending a grpc create ocm core share request", err)
		return
	}
	if createShareResponse.Status.Code != rpc.Code_CODE_OK {
		if createShareResponse.Status.Code == rpc.Code_CODE_NOT_FOUND {
			WriteError(w, r, APIErrorNotFound, "not found", nil)
			return
		}
		WriteError(w, r, APIErrorServerError, "grpc create ocm core share request failed", errors.New(createShareResponse.Status.Message))
		return
	}

	timeCreated := createShareResponse.Created
	jsonOut, err := json.Marshal(
		map[string]string{
			"id":        createShareResponse.Id,
			"createdAt": time.Unix(int64(timeCreated.Seconds), int64(timeCreated.Nanos)).String(),
		},
	)
	if err != nil {
		WriteError(w, r, APIErrorServerError, "error marshalling share data", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonOut)
	if err != nil {
		WriteError(w, r, APIErrorServerError, "error writing shares data", err)
		return
	}

	log.Info().Msg("Share created.")
}
