package email

import (
	"bytes"
	"embed"
	"html/template"
)

// go:embed templates/*.tmpl

// RenderEmailTemplate renders the email template for a new share
func RenderEmailTemplate(templateName string, templateVariables map[string]string) (string, error) {
	var fs embed.FS
	tpl, err := template.ParseFS(fs, templateName)
	if err != nil {
		return "", err
	}
	writer := bytes.NewBufferString("")
	err = tpl.Execute(writer, templateVariables)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}
