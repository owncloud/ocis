package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	fmt.Println("Getting relevant packages")
	paths, err := filepath.Glob("../../*/pkg/config/defaultconfig.go")
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
	content, err := ioutil.ReadFile("extractor.go.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating intermediate go code")
	tpl := template.Must(template.New("").Parse(string(content)))
	os.Mkdir("output", 0700)
	runner, err := os.Create("output/runner.go")
	if err != nil {
		log.Fatal(err)
	}
	tpl.Execute(runner, paths)
	fmt.Println("Running intermediate go code")
	os.Chdir("output")
	os.Setenv("OCIS_BASE_DATA_PATH", "~/.ocis")
	out, err := exec.Command("go", "run", "runner.go").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
	fmt.Println("Cleaning up")
	os.Chdir("../")
	os.RemoveAll("output")
}
