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
	"html/template"
	"net/http"
	"strings"

	"github.com/cs3org/reva/v2/pkg/siteacc/config"
	"github.com/cs3org/reva/v2/pkg/siteacc/data"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// TemplateID is the type for template identifiers.
type TemplateID = string

// Panel provides basic HTML panel functionality.
type Panel struct {
	conf *config.Configuration
	log  *zerolog.Logger

	name string

	provider PanelProvider

	templates map[TemplateID]*template.Template
}

const (
	pathParameterName = "path"
)

func (panel *Panel) initialize(name string, provider PanelProvider, conf *config.Configuration, log *zerolog.Logger) error {
	if name == "" {
		return errors.Errorf("no name provided")
	}
	panel.name = name

	if conf == nil {
		return errors.Errorf("no configuration provided")
	}
	panel.conf = conf

	if log == nil {
		return errors.Errorf("no logger provided")
	}
	panel.log = log

	if provider == nil {
		return errors.Errorf("no panel provider provided")
	}
	panel.provider = provider

	// Create space for the panel templates
	panel.templates = make(map[string]*template.Template, 5)

	return nil
}

func (panel *Panel) compile(provider ContentProvider) (string, error) {
	content := panelTemplate

	// Replace placeholders by the values provided by the content provider
	content = strings.ReplaceAll(content, "$(TITLE)", provider.GetTitle())
	content = strings.ReplaceAll(content, "$(CAPTION)", provider.GetCaption())

	content = strings.ReplaceAll(content, "$(CONTENT_JAVASCRIPT)", provider.GetContentJavaScript())
	content = strings.ReplaceAll(content, "$(CONTENT_STYLESHEET)", provider.GetContentStyleSheet())
	content = strings.ReplaceAll(content, "$(CONTENT_BODY)", provider.GetContentBody())

	return content, nil
}

// AddTemplate adds and compiles a new template.
func (panel *Panel) AddTemplate(name TemplateID, provider ContentProvider) error {
	name = panel.getFullTemplateName(name)

	if provider == nil {
		return errors.Errorf("no content provider provided")
	}

	content, err := panel.compile(provider)
	if err != nil {
		return errors.Wrapf(err, "error while compiling panel template %v", name)
	}

	tpl := template.New(name)
	panel.prepareTemplate(tpl)

	if _, err := tpl.Parse(content); err != nil {
		return errors.Wrapf(err, "error while parsing panel template %v", name)
	}
	panel.templates[name] = tpl

	return nil
}

// Execute generates the HTTP output of the panel and writes it to the response writer.
func (panel *Panel) Execute(w http.ResponseWriter, r *http.Request, session *Session, dataProvider PanelDataProvider) error {
	// Get the path query parameter; the panel provider may use this to determine the template to use
	path := r.URL.Query().Get(pathParameterName)

	actTpl := panel.provider.GetActiveTemplate(session, path)
	tplName := panel.getFullTemplateName(actTpl)
	tpl, ok := panel.templates[tplName]
	if !ok {
		return errors.Errorf("template %v not found", tplName)
	}

	// If a data provider is specified, use it to get additional template data
	var data interface{}
	if dataProvider != nil {
		data = dataProvider(session)
	}

	// Perform the pre-execution phase in which the panel provider can intercept the actual execution
	if state, err := panel.provider.PreExecute(session, actTpl, w, r); err == nil {
		if !state {
			return nil
		}
	} else {
		return errors.Wrapf(err, "pre-execution of template %v failed", tplName)
	}

	return tpl.Execute(w, data)
}

func (panel *Panel) prepareTemplate(tpl *template.Template) {
	// Add some custom helper functions to the template
	tpl.Funcs(template.FuncMap{
		"getServerAddress": func() string {
			return strings.TrimRight(panel.conf.Webserver.URL, "/")
		},
		"getSiteName": func(siteID string, fullName bool) string {
			siteName, _ := data.QuerySiteName(siteID, fullName, panel.conf.Mentix.URL, panel.conf.Mentix.DataEndpoint)
			return siteName
		},
	})
}

func (panel *Panel) getFullTemplateName(name string) string {
	return panel.name + "-" + name
}

// NewPanel creates a new panel.
func NewPanel(name string, provider PanelProvider, conf *config.Configuration, log *zerolog.Logger) (*Panel, error) {
	panel := &Panel{}
	if err := panel.initialize(name, provider, conf, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the panel")
	}
	return panel, nil
}
