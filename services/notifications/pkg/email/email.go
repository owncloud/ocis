// Package email implements utility for rendering the Email.
//
// The email package supports transifex translation for email templates.
package email

import (
	"bytes"
	"embed"
	"html"
	"html/template"
	"path/filepath"
)

var (
	//go:embed templates
	templatesFS embed.FS
)

// RenderEmailTemplate renders the email template for a new share
func RenderEmailTemplate(et MessageTemplate, locale string, emailTemplatePath string, translationPath string, vars map[string]interface{}) (string, string, error) {
	rawsub := ComposeMessage(et.Subject, locale, translationPath)
	// replace the body email placeholders with the values
	subject, err := executeRaw(rawsub, vars)
	if err != nil {
		return "", "", err
	}

	bodyPlaceholders := map[string]interface{}{}
	bodyPlaceholders["Greeting"] = ComposeMessage(et.Greeting, locale, translationPath)
	bodyPlaceholders["MessageBody"] = ComposeMessage(et.MessageBody, locale, translationPath)
	bodyPlaceholders["CallToAction"] = ComposeMessage(et.CallToAction, locale, translationPath)

	// replace the body email template placeholders with the translated template
	rawBody, err := executeEmailTemplate(emailTemplatePath, et.bodyTemplate, bodyPlaceholders)
	if err != nil {
		return "", "", err
	}
	// replace the body email placeholders with the values
	body, err := executeRaw(rawBody, vars)
	if err != nil {
		return "", "", err
	}
	return subject, body, nil
}

func executeEmailTemplate(emailTemplatePath, templateName string, vars map[string]interface{}) (string, error) {
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
	str, err := executeTemplate(tpl, vars)
	if err != nil {
		return "", err
	}
	return html.UnescapeString(str), err
}

func executeRaw(raw string, vars map[string]interface{}) (string, error) {
	tpl, err := template.New("").Parse(raw)
	if err != nil {
		return "", err
	}
	return executeTemplate(tpl, vars)
}

func executeTemplate(tpl *template.Template, vars map[string]interface{}) (string, error) {
	var writer bytes.Buffer
	if err := tpl.Execute(&writer, vars); err != nil {
		return "", err
	}
	return writer.String(), nil
}
