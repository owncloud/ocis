package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

var targets = map[string]string{
	"templates/adoc-generator.go.tmpl":                      "output/adoc/adoc-generator.go",
	"templates/example-config-generator.go.tmpl":            "output/exampleconfig/example-config-generator.go",
	"templates/environment-variable-docs-generator.go.tmpl": "output/env/environment-variable-docs-generator.go",
	"templates/envar-delta-table.go.tmpl":                   "output/env/envvar-delta-table.go",
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
	err = os.RemoveAll("output")
	if err != nil {
		fmt.Println(err)
	}
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
	defaultConfigPath := "/etc/ocis"
	defaultDataPath := "/var/lib/ocis"
	os.Setenv("OCIS_BASE_DATA_PATH", defaultDataPath)
	os.Setenv("OCIS_CONFIG_DIR", defaultConfigPath)

	// Set AUTOMEMLIMIT_EXPERIMENT=system on non-Linux systems to avoid cgroups errors
	if runtime.GOOS != "linux" {
		os.Setenv("AUTOMEMLIMIT_EXPERIMENT", "system")
	}

	out, err := exec.Command("go", "run", intermediateCodePath).CombinedOutput()
	if err != nil {
		log.Fatal(string(out), err)
	}
	fmt.Println(string(out))
}
