package email

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

const templatePath string = "../../email/templates"

// RenderEmailTemplate renders the email template for a new share
func RenderEmailTemplate(templateName string, templateVariables map[string]string) (string, error) {
	content, err := os.ReadFile(filepath.Join(templatePath, templateName))
	if err != nil {
		return "", err
	}
	tpl := template.Must(template.New("").Parse(string(content)))
	writer := bytes.NewBufferString("")
	err = tpl.Execute(writer, templateVariables)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}
