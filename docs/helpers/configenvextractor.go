package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

var targets = map[string]string{
	"adoc-generator.go.tmpl":                      "output/adoc/adoc-generator.go",
	"example-config-generator.go.tmpl":            "output/exampleconfig/example-config-generator.go",
	"environment-variable-docs-generator.go.tmpl": "output/env/environment-variable-docs-generator.go",
}

// RenderTemplates does something with templates
func RenderTemplates() {
	fmt.Println("Getting relevant packages")
	paths, err := filepath.Glob("../../services/*/pkg/config/defaults/defaultconfig.go")
	if err != nil {
		log.Fatal(err)
	}
	replacer := strings.NewReplacer(
		"../../", "github.com/owncloud/ocis/v2/",
		"/defaultconfig.go", "",
	)
	for i := range paths {
		paths[i] = replacer.Replace(paths[i])
	}

	for template, output := range targets {
		generateIntermediateCode(template, output, paths)
		runIntermediateCode(output)
	}
	fmt.Println("Cleaning up")
	os.RemoveAll("output")
}

func generateIntermediateCode(templatePath string, intermediateCodePath string, paths []string) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating intermediate go code for " + intermediateCodePath + " using template " + templatePath)
	tpl := template.Must(template.New("").Parse(string(content)))
	err = os.MkdirAll(path.Dir(intermediateCodePath), 0700)
	if err != nil {
		log.Fatal(err)
	}
	runner, err := os.Create(intermediateCodePath)
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.Execute(runner, paths)
	if err != nil {
		log.Fatal(err)
	}
}

func runIntermediateCode(intermediateCodePath string) {
	fmt.Println("Running intermediate go code for " + intermediateCodePath)
	defaultPath := "~/.ocis"
	os.Setenv("OCIS_BASE_DATA_PATH", defaultPath)
	os.Setenv("OCIS_CONFIG_DIR", path.Join(defaultPath, "config"))
	out, err := exec.Command("go", "run", intermediateCodePath).CombinedOutput()
	if err != nil {
		log.Fatal(err, string(out))
	}
	fmt.Println(string(out))
}
