package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

var targets = map[string]string{
	"example-config-generator.go.tmpl": "output/exampleconfig/example-config-generator.go",
	"extractor.go.tmpl":                "output/env/runner.go",
}

func main() {
	fmt.Println("Getting relevant packages")
	paths, err := filepath.Glob("../../extensions/*/pkg/config/defaults/defaultconfig.go")
	if err != nil {
		log.Fatal(err)
	}
	replacer := strings.NewReplacer(
		"../../", "github.com/owncloud/ocis/",
		"/defaultconfig.go", "",
	)
	for i := range paths {
		paths[i] = replacer.Replace(paths[i])
	}

	for template, output := range targets {
		GenerateIntermediateCode(template, output, paths)
		RunIntermediateCode(output)
	}
	fmt.Println("Cleaning up")
	os.RemoveAll("output")
}

func GenerateIntermediateCode(templatePath string, intermediateCodePath string, paths []string) {
	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating intermediate go code for " + intermediateCodePath + " using template " + templatePath)
	tpl := template.Must(template.New("").Parse(string(content)))
	os.MkdirAll(path.Dir(intermediateCodePath), 0700)
	runner, err := os.Create(intermediateCodePath)
	if err != nil {
		log.Fatal(err)
	}
	tpl.Execute(runner, paths)
}

func RunIntermediateCode(intermediateCodePath string) {
	fmt.Println("Running intermediate go code for " + intermediateCodePath)
	os.Setenv("OCIS_BASE_DATA_PATH", "~/.ocis")
	os.Setenv("OCIS_CONFIG_DIR", "~/.ocis/config")
	out, err := exec.Command("go", "run", intermediateCodePath).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
