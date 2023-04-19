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

package html

import (
	"net/http"
)

const (
	// ContinueExecution causes the execution of a panel to continue.
	ContinueExecution = true
	// AbortExecution causes the execution of a panel to be aborted.
	AbortExecution = false
)

// ExecutionResult is the type returned by the PreExecute function of PanelProvider.
type ExecutionResult = bool

// PanelProvider handles general panel tasks.
type PanelProvider interface {
	// GetActiveTemplate returns the name of the active template.
	GetActiveTemplate(*Session, string) string

	// PreExecute is called before the actual template is being executed.
	PreExecute(*Session, string, http.ResponseWriter, *http.Request) (ExecutionResult, error)
}

// PanelDataProvider is the function signature for panel data providers.
type PanelDataProvider = func(*Session) interface{}

// ContentProvider defines various methods for HTML content providers.
type ContentProvider interface {
	// GetTitle returns the title of the panel.
	GetTitle() string
	// GetCaption returns the caption which is displayed on the panel.
	GetCaption() string

	// GetContentJavaScript delivers additional JavaScript code.
	GetContentJavaScript() string
	// GetContentStyleSheet delivers additional stylesheet code.
	GetContentStyleSheet() string
	// GetContentBody delivers the actual body content.
	GetContentBody() string
}
