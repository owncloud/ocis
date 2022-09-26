package email

import (
	"bytes"
	"embed"
	"html/template"
	"path/filepath"
)

var (
	//go:embed templates
	templatesFS embed.FS
)

// RenderEmailTemplate renders the email template for a new share
func RenderEmailTemplate(templateName string, templateVariables map[string]string, emailTemplatePath string) (string, error) {
	var err error
	var tpl *template.Template
	// try to lookup the files in the filesystem
	tpl, err = template.ParseFiles(filepath.Join(emailTemplatePath, templateName))
	if err != nil {
		// template has not been found in the fs, or path has not been specified => use embed templates
		tpl, err = template.ParseFS(templatesFS, filepath.Join("templates/", templateName))
		if err != nil {
			return "", err
		}
	}
	var writer bytes.Buffer
	err = tpl.Execute(&writer, templateVariables)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}
