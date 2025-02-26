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

//go:build ceph
// +build ceph

package cephfs

import (
	"path/filepath"

	"github.com/cs3org/reva/v2/pkg/sharedconf"
)

// Options for the cephfs module
type Options struct {
	GatewaySvc   string `mapstructure:"gatewaysvc"`
	IndexPool    string `mapstructure:"index_pool"`
	Root         string `mapstructure:"root"`
	ShadowFolder string `mapstructure:"shadow_folder"`
	ShareFolder  string `mapstructure:"share_folder"`
	UploadFolder string `mapstructure:"uploads"`
	UserLayout   string `mapstructure:"user_layout"`

	DisableHome    bool   `mapstructure:"disable_home"`
	UserQuotaBytes uint64 `mapstructure:"user_quota_bytes"`
	HiddenDirs     map[string]bool
}

func (c *Options) fillDefaults() {
	c.GatewaySvc = sharedconf.GetGatewaySVC(c.GatewaySvc)

	if c.IndexPool == "" {
		c.IndexPool = "path_index"
	}

	if c.Root == "" {
		c.Root = "/home"
	} else {
		c.Root = addLeadingSlash(c.Root) //force absolute path in case leading "/" is omitted
	}

	if c.ShadowFolder == "" {
		c.ShadowFolder = "/.reva_hidden"
	} else {
		c.ShadowFolder = addLeadingSlash(c.ShadowFolder)
	}

	if c.ShareFolder == "" {
		c.ShareFolder = "/Shares"
	} else {
		c.ShareFolder = addLeadingSlash(c.ShareFolder)
	}

	if c.UploadFolder == "" {
		c.UploadFolder = ".uploads"
	}
	c.UploadFolder = filepath.Join(c.ShadowFolder, c.UploadFolder)

	if c.UserLayout == "" {
		c.UserLayout = "{{.Username}}"
	}

	c.HiddenDirs = map[string]bool{
		".":                                true,
		"..":                               true,
		removeLeadingSlash(c.ShadowFolder): true,
	}

	c.DisableHome = false // it is currently only home based

	if c.UserQuotaBytes == 0 {
		c.UserQuotaBytes = 50000000000
	}
}
