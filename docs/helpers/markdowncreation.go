package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/markdown"
)

var _configMarkdown = `{{< include file="services/_includes/%s-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/%s_configvars.md" >}}
`

// GenerateServiceIndexMarkdowns generates the _index.md files for the dev docu
func GenerateServiceIndexMarkdowns() {
	paths, err := filepath.Glob("../../services/*/README.md")
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range paths {
		service := filepath.Base(filepath.Dir(p))
		if err := generateMarkdown(p, service); err != nil {
			fmt.Printf("error generating markdown for %s: %s\n", service, err)
		}
	}
}

func generateMarkdown(filepath string, servicename string) error {
	f, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	md := markdown.NewMD(f)
	if len(md.Headings) == 0 || md.Headings[0].Level != 1 {
		return errors.New("readme has invalid format")
	}

	// we don't need the main title, we add in our template
	head := md.Headings[0]
	md.Headings = md.Headings[1:]
	md.Headings = append(md.Headings, markdown.Heading{
		Level:   2,
		Header:  "Example Yaml Config",
		Content: fmt.Sprintf(_configMarkdown, servicename, servicename),
	})

	tpl := template.Must(template.ParseFiles("templates/index.tmpl"))
	b := bytes.NewBuffer(nil)
	if err := tpl.Execute(b, map[string]interface{}{
		"ServiceName":  head.Header,
		"CreationTime": time.Now().Format(time.RFC3339Nano),
		"service":      servicename,
		"Abstract":     head.Content,
		"TocTree":      md.TocString(),
		"Content":      md.String(),
	}); err != nil {
		return err
	}

	path := fmt.Sprintf("../../docs/services/%s", servicename)

	if err := os.Mkdir(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	return os.WriteFile(path+"/_index.md", b.Bytes(), os.ModePerm)
}
