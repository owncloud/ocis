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

package options

import (
	"path/filepath"
	"strings"

	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/sharedconf"
	"github.com/cs3org/reva/v2/pkg/storage/cache"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {

	// the gateway address
	GatewayAddr string `mapstructure:"gateway_addr"`

	// the metadata backend to use, currently supports `xattr` or `ini`
	MetadataBackend string `mapstructure:"metadata_backend"`

	// ocis fs works on top of a dir of uuid nodes
	Root string `mapstructure:"root"`

	// UserLayout describes the relative path from the storage's root node to the users home node.
	UserLayout string `mapstructure:"user_layout"`

	// propagate mtime changes as tmtime (tree modification time) to the parent directory when user.ocis.propagation=1 is set on a node
	TreeTimeAccounting bool `mapstructure:"treetime_accounting"`

	// propagate size changes as treesize
	TreeSizeAccounting bool `mapstructure:"treesize_accounting"`

	// permissions service to use when checking permissions
	PermissionsSVC           string `mapstructure:"permissionssvc"`
	PermissionsClientTLSMode string `mapstructure:"permissionssvc_tls_mode"`
	PermTLSMode              pool.TLSMode

	PersonalSpaceAliasTemplate string `mapstructure:"personalspacealias_template"`
	GeneralSpaceAliasTemplate  string `mapstructure:"generalspacealias_template"`

	AsyncFileUploads bool `mapstructure:"asyncfileuploads"`

	Events EventOptions `mapstructure:"events"`

	Tokens TokenOptions `mapstructure:"tokens"`

	StatCache         cache.Config `mapstructure:"statcache"`
	FileMetadataCache cache.Config `mapstructure:"filemetadatacache"`

	MaxAcquireLockCycles    int `mapstructure:"max_acquire_lock_cycles"`
	LockCycleDurationFactor int `mapstructure:"lock_cycle_duration_factor"`
	MaxConcurrency          int `mapstructure:"max_concurrency"`

	MaxQuota uint64 `mapstructure:"max_quota"`
}

// EventOptions are the configurable options for events
type EventOptions struct {
	NatsAddress          string `mapstructure:"natsaddress"`
	NatsClusterID        string `mapstructure:"natsclusterid"`
	TLSInsecure          bool   `mapstructure:"tlsinsecure"`
	TLSRootCACertificate string `mapstructure:"tlsrootcacertificate"`
	NumConsumers         int    `mapstructure:"numconsumers"`
}

// TokenOptions are the configurable option for tokens
type TokenOptions struct {
	DownloadEndpoint     string `mapstructure:"download_endpoint"`
	DataGatewayEndpoint  string `mapstructure:"datagateway_endpoint"`
	TransferSharedSecret string `mapstructure:"transfer_shared_secret"`
	TransferExpires      int64  `mapstructure:"transfer_expires"`
}

// New returns a new Options instance for the given configuration
func New(m map[string]interface{}) (*Options, error) {
	o := &Options{}
	if err := mapstructure.Decode(m, o); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}

	o.GatewayAddr = sharedconf.GetGatewaySVC(o.GatewayAddr)

	if o.MetadataBackend == "" {
		o.MetadataBackend = "xattrs"
	}

	if o.UserLayout == "" {
		o.UserLayout = "{{.Id.OpaqueId}}"
	}
	// ensure user layout has no starting or trailing /
	o.UserLayout = strings.Trim(o.UserLayout, "/")

	// c.DataDirectory should never end in / unless it is the root
	o.Root = filepath.Clean(o.Root)

	if o.PersonalSpaceAliasTemplate == "" {
		o.PersonalSpaceAliasTemplate = "{{.SpaceType}}/{{.User.Username}}"
	}

	if o.GeneralSpaceAliasTemplate == "" {
		o.GeneralSpaceAliasTemplate = "{{.SpaceType}}/{{.SpaceName | replace \" \" \"-\" | lower}}"
	}

	if o.PermissionsClientTLSMode != "" {
		var err error
		o.PermTLSMode, err = pool.StringToTLSMode(o.PermissionsClientTLSMode)
		if err != nil {
			return nil, err
		}
	} else {
		sharedOpt := sharedconf.GRPCClientOptions()
		var err error

		if o.PermTLSMode, err = pool.StringToTLSMode(sharedOpt.TLSMode); err != nil {
			return nil, err
		}
	}

	if o.MaxConcurrency <= 0 {
		o.MaxConcurrency = 100
	}

	return o, nil
}
