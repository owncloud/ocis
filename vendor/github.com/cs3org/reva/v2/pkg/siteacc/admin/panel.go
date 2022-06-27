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

package admin

import (
	"net/http"

	"github.com/cs3org/reva/v2/pkg/siteacc/config"
	"github.com/cs3org/reva/v2/pkg/siteacc/data"
	"github.com/cs3org/reva/v2/pkg/siteacc/html"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Panel represents the web interface panel of the accounts service administration.
type Panel struct {
	html.PanelProvider
	html.ContentProvider

	htmlPanel *html.Panel
}

const (
	templateMain = "main"
)

func (panel *Panel) initialize(conf *config.Configuration, log *zerolog.Logger) error {
	// Create the internal HTML panel
	htmlPanel, err := html.NewPanel("admin-panel", panel, conf, log)
	if err != nil {
		return errors.Wrap(err, "unable to create the administration panel")
	}
	panel.htmlPanel = htmlPanel

	// Add all templates
	if err := panel.htmlPanel.AddTemplate(templateMain, panel); err != nil {
		return errors.Wrap(err, "unable to create the main template")
	}

	return nil
}

// GetActiveTemplate returns the name of the active template.
func (panel *Panel) GetActiveTemplate(*html.Session, string) string {
	return templateMain
}

// GetTitle returns the title of the htmlPanel.
func (panel *Panel) GetTitle() string {
	return "ScienceMesh Site Administrator Accounts Panel"
}

// GetCaption returns the caption which is displayed on the htmlPanel.
func (panel *Panel) GetCaption() string {
	return "ScienceMesh Site Administrator Accounts ({{.Accounts | len}})"
}

// GetContentJavaScript delivers additional JavaScript code.
func (panel *Panel) GetContentJavaScript() string {
	return tplJavaScript
}

// GetContentStyleSheet delivers additional stylesheet code.
func (panel *Panel) GetContentStyleSheet() string {
	return tplStyleSheet
}

// GetContentBody delivers the actual body content.
func (panel *Panel) GetContentBody() string {
	return tplBody
}

// PreExecute is called before the actual template is being executed.
func (panel *Panel) PreExecute(*html.Session, string, http.ResponseWriter, *http.Request) (html.ExecutionResult, error) {
	return html.ContinueExecution, nil
}

// Execute generates the HTTP output of the htmlPanel and writes it to the response writer.
func (panel *Panel) Execute(w http.ResponseWriter, r *http.Request, session *html.Session, accounts *data.Accounts) error {
	dataProvider := func(*html.Session) interface{} {
		type TemplateData struct {
			Accounts *data.Accounts
		}

		return TemplateData{
			Accounts: accounts,
		}
	}
	return panel.htmlPanel.Execute(w, r, session, dataProvider)
}

// NewPanel creates a new administration panel.
func NewPanel(conf *config.Configuration, log *zerolog.Logger) (*Panel, error) {
	panel := &Panel{}
	if err := panel.initialize(conf, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the administration panel")
	}
	return panel, nil
}
