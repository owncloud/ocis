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
func RenderEmailTemplate(mt MessageTemplate, locale string, emailTemplatePath string, translationPath string, vars map[string]interface{}) (string, string, error) {
	// translate a message
	mt.Subject = ComposeMessage(mt.Subject, locale, translationPath)
	mt.Greeting = ComposeMessage(mt.Greeting, locale, translationPath)
	mt.MessageBody = ComposeMessage(mt.MessageBody, locale, translationPath)
	mt.CallToAction = ComposeMessage(mt.CallToAction, locale, translationPath)

	// replace the body email placeholders with the values
	subject, err := executeRaw(mt.Subject, vars)
	if err != nil {
		return "", "", err
	}

	// replace the body email template placeholders with the translated template
	rawBody, err := executeEmailTemplate(emailTemplatePath, mt)
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

func executeEmailTemplate(emailTemplatePath string, mt MessageTemplate) (string, error) {
	var err error
	var tpl *template.Template
	// try to lookup the files in the filesystem
	tpl, err = template.ParseFiles(filepath.Join(emailTemplatePath, mt.bodyTemplate))
	if err != nil {
		// template has not been found in the fs, or path has not been specified => use embed templates
		tpl, err = template.ParseFS(templatesFS, filepath.Join("templates/", mt.bodyTemplate))
		if err != nil {
			return "", err
		}
	}
	str, err := executeTemplate(tpl, mt)
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

func executeTemplate(tpl *template.Template, vars any) (string, error) {
	var writer bytes.Buffer
	if err := tpl.Execute(&writer, vars); err != nil {
		return "", err
	}
	return writer.String(), nil
}
