package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

const yamlSource = "global_vars.yaml"

type ConfigVars struct {
	Variables []Variable `yaml:"variables"`
}

type Variable struct {
	Name              string    `yaml:"name"`
	Type              string    `yaml:"type"`
	DefaultValue      string    `yaml:"default_value"`
	Description       string    `yaml:"description"`
	DependendServices []Service `yaml:"dependend_services"`
	DoIgnore          bool      `yaml:"do_ignore"`
}

type Service struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

func GetRogueEnvs() {
	curdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fullYamlPath := filepath.Join(curdir, yamlSource)
	re := regexp.MustCompile(`"[A-z0-9_]{1,}"`)
	vars := &ConfigVars{}
	fmt.Printf("Reading existing variable definitions from %s\n", fullYamlPath)
	yfile, err := ioutil.ReadFile(fullYamlPath)
	if err == nil {
		err := yaml.Unmarshal(yfile, &vars)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := os.Chdir("../../"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Gathering variable definitions from source")
	out, err := exec.Command("bash", "-c", "grep -R os.Getenv | grep -v rogue-env.go |grep \\.go").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		r := strings.SplitN(l, ":\t", 2)
		if len(r) == 2 && r[0] != "" && r[1] != "" {
			fmt.Printf("Parsing %s\n", r[0])
			res := re.FindAll([]byte(r[1]), -1)
			for _, item := range res {
				AddUniqueToStruct(vars, Variable{Name: strings.Trim(string(item), "\""), Type: ""})
			}
		}
	}
	output, err := yaml.Marshal(vars)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Writing new variable definitions to %s\n", fullYamlPath)
	err = ioutil.WriteFile(fullYamlPath, output, 0666)
	if err != nil {
		log.Fatalf("could not write %s", fullYamlPath)
	}
	if err := os.Chdir(curdir); err != nil {
		log.Fatal(err)
	}
}

func AddUniqueToStruct(variables *ConfigVars, variable Variable) {
	for _, itm := range variables.Variables {
		if itm.Name == variable.Name {
			return
		}
	}
	variables.Variables = append(variables.Variables, variable)
}

func RenderGlobalVarsTemplate() {
	curdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fullYamlPath := filepath.Join(curdir, yamlSource)

	content, err := ioutil.ReadFile("../../docs/templates/ADOC_global.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	targetFolder := "../../docs/services/_includes/adoc/"

	vars := &ConfigVars{}
	fmt.Printf("Reading existing variable definitions from %s\n", fullYamlPath)
	yfile, err := ioutil.ReadFile(fullYamlPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yfile, &vars)
	if err != nil {
		log.Fatal(err)
	}

	targetFile, err := os.Create(filepath.Join(targetFolder, "global_configvars.adoc"))
	if err != nil {
		log.Fatalf("Failed to create target file: %s", err)
	}
	defer targetFile.Close()

	tpl := template.Must(template.New("").Parse(string(content)))
	if err = tpl.Execute(targetFile, *vars); err != nil {
		log.Fatalf("Failed to execute template: %s", err)
	}
}
