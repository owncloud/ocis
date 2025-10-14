// Copyright 2018-2023 CERN
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
	"reflect"
	"strings"

	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ocmshare "github.com/owncloud/reva/v2/pkg/ocm/share"
	utils "github.com/owncloud/reva/v2/pkg/utils"
)

// Protocols is the list of protocols.
type Protocols []Protocol

// Protocol represents the way of access the resource
// in the OCM share.
type Protocol interface {
	// ToOCMProtocol convert the protocol to a ocm Protocol struct
	ToOCMProtocol(*typesv1beta1.Opaque) *ocm.Protocol
}

// protocols supported by the OCM API

// WebDAV contains the parameters for the WebDAV protocol.
type WebDAV struct {
	SharedSecret string   `json:"sharedSecret" validate:"required"`
	Permissions  []string `json:"permissions" validate:"required,dive,required,oneof=read write share"`
	URL          string   `json:"url" validate:"required"`
}

// ToOCMProtocol convert the protocol to a ocm Protocol struct.
func (w *WebDAV) ToOCMProtocol(o *typesv1beta1.Opaque) *ocm.Protocol {
	resourceType := utils.ReadPlainFromOpaque(o, "resourceType")
	perms := &ocm.SharePermissions{
		Permissions: &providerv1beta1.ResourcePermissions{},
	}
	for _, p := range w.Permissions {
		switch p {
		case "read":
			perms.Permissions.GetPath = true
			perms.Permissions.GetQuota = true
			perms.Permissions.InitiateFileDownload = true
			perms.Permissions.ListContainer = true
			perms.Permissions.ListRecycle = true
			perms.Permissions.Stat = true
		case "write":
			perms.Permissions.InitiateFileUpload = true
			perms.Permissions.RestoreRecycleItem = true
			if resourceType == "folder" {
				perms.Permissions.Move = true
				perms.Permissions.CreateContainer = true
				perms.Permissions.Delete = true
			}
		}
	}

	return ocmshare.NewWebDAVProtocol(w.URL, w.SharedSecret, perms)
}

// Webapp contains the parameters for the Webapp protocol.
type Webapp struct {
	URITemplate string `json:"uriTemplate" validate:"required"`
	ViewMode    string `json:"viewMode" validate:"required,dive,required,oneof=view read write"`
}

// ToOCMProtocol convert the protocol to a ocm Protocol struct.
func (w *Webapp) ToOCMProtocol(_ *typesv1beta1.Opaque) *ocm.Protocol {
	return ocmshare.NewWebappProtocol(w.URITemplate, utils.GetAppViewMode(w.ViewMode))
}

// Datatx contains the parameters for the Datatx protocol.
type Datatx struct {
	SharedSecret string `json:"sharedSecret" validate:"required"`
	SourceURI    string `json:"srcUri" validate:"required"`
	Size         uint64 `json:"size" validate:"required"`
}

// ToOCMProtocol convert the protocol to a ocm Protocol struct.
func (w *Datatx) ToOCMProtocol(_ *typesv1beta1.Opaque) *ocm.Protocol {
	return ocmshare.NewTransferProtocol(w.SourceURI, w.SharedSecret, w.Size)
}

var protocolImpl = map[string]reflect.Type{
	"webdav": reflect.TypeOf(WebDAV{}),
	"webapp": reflect.TypeOf(Webapp{}),
	"datatx": reflect.TypeOf(Datatx{}),
}

// UnmarshalJSON implements the Unmarshaler interface.
func (p *Protocols) UnmarshalJSON(data []byte) error {
	var prot map[string]json.RawMessage
	if err := json.Unmarshal(data, &prot); err != nil {
		return err
	}

	*p = []Protocol{}

	for name, d := range prot {
		var res Protocol

		// we do not support the OCM v1.0 properties for now, therefore just skip or bail out
		if name == "name" {
			continue
		}
		if name == "options" {
			var opt map[string]any
			if err := json.Unmarshal(d, &opt); err != nil || len(opt) > 0 {
				return fmt.Errorf("protocol options not supported: %s", string(d))
			}
			continue
		}
		ctype, ok := protocolImpl[name]
		if !ok {
			return fmt.Errorf("protocol %s not recognised", name)
		}
		res = reflect.New(ctype).Interface().(Protocol)
		if err := json.Unmarshal(d, &res); err != nil {
			return err
		}

		*p = append(*p, res)
	}
	return nil
}

// MarshalJSON implements the Marshaler interface.
func (p Protocols) MarshalJSON() ([]byte, error) {
	if len(p) == 0 {
		return nil, errors.New("no protocol defined")
	}
	d := make(map[string]any)
	for _, prot := range p {
		d[getProtocolName(prot)] = prot
	}
	// fill in the OCM v1.0 properties
	d["name"] = "multi"
	d["options"] = map[string]any{}
	return json.Marshal(d)
}

func getProtocolName(p Protocol) string {
	n := reflect.TypeOf(p).String()
	s := strings.Split(n, ".")
	return strings.ToLower(s[len(s)-1])
}
