// Package email implements utility for rendering the Email.
//
// The email package supports transifex translation for email templates.
package email

import (
	"bytes"
	"embed"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/owncloud/ocis/v2/services/notifications/pkg/channels"
)

var (
	//go:embed templates
	templatesFS embed.FS
)

// RenderEmailTemplate renders the email template for a new share
func RenderEmailTemplate(mt MessageTemplate, locale string, emailTemplatePath string, translationPath string, vars map[string]interface{}) (*channels.Message, error) {
	textMt, err := NewTextTemplate(mt, locale, translationPath, vars)
	if err != nil {
		return nil, err
	}
	tpl, err := parseTemplate(emailTemplatePath, mt.textTemplate)
	if err != nil {
		return nil, err
	}
	textBody, err := emailTemplate(tpl, textMt)
	if err != nil {
		return nil, err
	}

	htmlMt, err := NewHTMLTemplate(mt, locale, translationPath, vars)
	if err != nil {
		return nil, err
	}
	htmlTpl, err := parseTemplate(emailTemplatePath, filepath.Join("html", "email.html.tmpl"))
	if err != nil {
		return nil, err
	}
	htmlBody, err := emailTemplate(htmlTpl, htmlMt)
	if err != nil {
		return nil, err
	}

	var data map[string][]byte
	data, err = readImages(emailTemplatePath)
	if err != nil {
		data, err = readFs()
		if err != nil {
			return nil, err
		}
	}
	return &channels.Message{
		Subject:      textMt.Subject,
		TextBody:     textBody,
		HTMLBody:     htmlBody,
		AttachInline: data,
	}, nil
}

func emailTemplate(tpl *template.Template, mt MessageTemplate) (string, error) {
	str, err := executeTemplate(tpl, map[string]interface{}{
		"Greeting":     template.HTML(strings.TrimSpace(mt.Greeting)),
		"MessageBody":  template.HTML(strings.TrimSpace(mt.MessageBody)),
		"CallToAction": template.HTML(strings.TrimSpace(mt.CallToAction)),
	})
	if err != nil {
		return "", err
	}
	return str, err
}

func parseTemplate(emailTemplatePath string, file string) (*template.Template, error) {
	var err error
	var tpl *template.Template
	// try to lookup the files in the filesystem
	tpl, err = template.ParseFiles(filepath.Join(emailTemplatePath, file))
	if err != nil {
		// template has not been found in the fs, or path has not been specified => use embed templates
		tpl, err = template.ParseFS(templatesFS, filepath.Join("templates", file))
		if err != nil {
			return nil, err
		}
	}
	return tpl, err
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

func readFs() (map[string][]byte, error) {
	dir := filepath.Join("templates", "html", "img")
	entries, err := templatesFS.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	list := make(map[string][]byte)
	for _, e := range entries {
		if !e.IsDir() {
			file, err := templatesFS.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				return nil, err
			}
			if !validateMime(file) {
				continue
			}
			list[e.Name()] = file
		}
	}
	return list, nil
}

func readImages(emailTemplatePath string) (map[string][]byte, error) {
	dir := filepath.Join(emailTemplatePath, "html", "img")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	list := make(map[string][]byte)
	for _, e := range entries {
		if !e.IsDir() {
			file, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				return nil, err
			}
			if !validateMime(file) {
				continue
			}
			list[e.Name()] = file
		}
	}
	return list, nil
}

// signature image formats signature https://go.dev/src/net/http/sniff.go #L:118
var signature = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

// validateMime validate the mime type of image file from its first few bytes
func validateMime(incipit []byte) bool {
	for s := range signature {
		if strings.HasPrefix(string(incipit), s) {
			return true
		}
	}
	return false
}
