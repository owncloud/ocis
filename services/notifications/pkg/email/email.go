// Package email implements utility for rendering the Email.
//
// The email package supports transifex translation for email templates.
package email

import (
	"bytes"
	"embed"
	"errors"
	"html"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/owncloud/ocis/v2/services/notifications/pkg/channels"
)

var (
	//go:embed templates
	templatesFS embed.FS

	imgDir = filepath.Join("templates", "html", "img")
)

// RenderEmailTemplate is responsible to prepare a message which than can be used to notify the user via email.
func RenderEmailTemplate(mt MessageTemplate, locale, defaultLocale string, emailTemplatePath string, translationPath string, vars map[string]string) (*channels.Message, error) {
	textMt, err := NewTextTemplate(mt, locale, defaultLocale, translationPath, vars)
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
	htmlMt, err := NewHTMLTemplate(mt, locale, defaultLocale, translationPath, escapeStringMap(vars))
	if err != nil {
		return nil, err
	}
	htmlTpl, err := parseTemplate(emailTemplatePath, mt.htmlTemplate)
	if err != nil {
		return nil, err
	}
	htmlBody, err := emailTemplate(htmlTpl, htmlMt)
	if err != nil {
		return nil, err
	}
	var data map[string][]byte
	if emailTemplatePath != "" {
		data, err = readImages(emailTemplatePath)
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
	if emailTemplatePath != "" {
		return template.ParseFiles(filepath.Join(emailTemplatePath, file))
	}
	return template.ParseFS(templatesFS, filepath.Join(file))
}

func executeTemplate(tpl *template.Template, vars any) (string, error) {
	var writer bytes.Buffer
	if err := tpl.Execute(&writer, vars); err != nil {
		return "", err
	}
	return writer.String(), nil
}

func readImages(emailTemplatePath string) (map[string][]byte, error) {
	dir := filepath.Join(emailTemplatePath, imgDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return read(entries, os.DirFS(emailTemplatePath))
}

func read(entries []fs.DirEntry, fsys fs.FS) (map[string][]byte, error) {
	list := make(map[string][]byte)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		file, err := fs.ReadFile(fsys, filepath.Join(imgDir, e.Name()))
		if err != nil {
			// skip if the symbolic is a directory
			if errors.Is(err, syscall.EISDIR) {
				continue
			}
			return nil, err
		}
		if !validateMime(file) {
			continue
		}
		list[e.Name()] = file
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

func escapeStringMap(vars map[string]string) map[string]string {
	for k := range vars {
		vars[k] = html.EscapeString(vars[k])
	}
	return vars
}
