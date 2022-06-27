// Copyright 2018-2020 CERN
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

package registration

import "github.com/cs3org/reva/v2/pkg/siteacc/html"

// PanelTemplate is the content provider for the registration form.
type PanelTemplate struct {
	html.ContentProvider
}

// GetTitle returns the title of the panel.
func (template *PanelTemplate) GetTitle() string {
	return "ScienceMesh Site Administrator Account Registration"
}

// GetCaption returns the caption which is displayed on the panel.
func (template *PanelTemplate) GetCaption() string {
	return "Welcome to the ScienceMesh Site Administrator Account Registration!"
}

// GetContentJavaScript delivers additional JavaScript code.
func (template *PanelTemplate) GetContentJavaScript() string {
	return tplJavaScript
}

// GetContentStyleSheet delivers additional stylesheet code.
func (template *PanelTemplate) GetContentStyleSheet() string {
	return tplStyleSheet
}

// GetContentBody delivers the actual body content.
func (template *PanelTemplate) GetContentBody() string {
	return tplBody
}
