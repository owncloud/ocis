package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

const yamlSource = "extended_vars.yaml"

// ConfigVars is the main yaml source
type ConfigVars struct {
	Variables []Variable `yaml:"variables"`
}

// Variable contains all information about one rogue envvar
type Variable struct {
	// These field structs are automatically filled:
	// RawName can be the name of the envvar or the name of its var
	RawName string `yaml:"rawname"`
	// Path to the envvar with linenumber
	Path string `yaml:"path"`
	// FoundInCode indicates if the variable is still found in the codebase. TODO: delete immediately?
	FoundInCode bool `yaml:"foundincode"`
	// Name is equal to RawName but will not be overwritten in consecutive runs
	Name string `yaml:"name"`

	// These field structs need manual filling:
	// Type of the envvar
	Type string `yaml:"type"`
	// DefaultValue of the envvar
	DefaultValue string `yaml:"default_value"`
	// Description of what this envvar does
	Description string `yaml:"description"`
	// Ignore this envvar when creating docs?
	Ignore bool `yaml:"do_ignore"`

	// For simplicity ignored for now:
	// DependendServices []Service `yaml:"dependend_services"`
}

// GetRogueEnvs extracts the rogue envs from the code
func GetRogueEnvs() {
	curdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fullYamlPath := filepath.Join(curdir, yamlSource)
	re := regexp.MustCompile(`os.Getenv\(([^\)]+)\)`)
	vars := &ConfigVars{}
	fmt.Printf("Reading existing variable definitions from %s\n", fullYamlPath)
	yfile, err := os.ReadFile(fullYamlPath)
	if err == nil {
		err := yaml.Unmarshal(yfile, &vars)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := os.Chdir("../../"); err != nil {
		log.Fatal(err)
	}

	// Option to process a specific directory, for faster development iteration
	processDir := os.Getenv("DIR")
	if processDir == "" {
		processDir = "."
	}

	fmt.Printf("Processing directory: %s\n", processDir)

	// Old command `grep -RHn os.Getenv $DIR --exclude-dir=vendor | grep -v extendedEnv.go | grep \\.go`
	// New command, 10x faster than direct grep:
	// 	C_ALL=C - Sets locale to C for consistent output
	// 	find %s -type f -name '*.go' - Finds all .go files in the specified directory
	// 	xargs -0 -P $(getconf _NPROCESSORS_ONLN) - Processes files in parallel using all available CPU cores
	// 	grep -F -Hn 'os.Getenv' - Searches for literal "os.Getenv" with filename and line number output
	grepCmd := fmt.Sprintf("C_ALL=C find %s -type f -name '*.go' -not -path '*/vendor/*' -print0 | xargs -0 -P $(getconf _NPROCESSORS_ONLN) grep -F -Hn 'os.Getenv' | grep -v extendedEnv.go", processDir)
	out, err := exec.Command("sh", "-c", grepCmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")

	// find current vars
	currentVars := make(map[string]Variable)
	totalLines := len(lines)
	fmt.Printf("Processing %d lines...\n", totalLines)

	for _, l := range lines {

		if l == "" {
			continue
		}

		fmt.Printf("Parsing %s\n", l)
		r := strings.SplitN(l, ":", 3)
		if len(r) != 3 || r[0] == "" || r[2] == "" {
			continue
		}

		// Remove ./ prefix from path if it exists
		path := strings.TrimPrefix(r[0], "./")
		path = path + ":" + r[1] // Reconstruct path:line
		content := r[2]

		res := re.FindAllSubmatch([]byte(content), -1)
		if len(res) < 1 {
			fmt.Printf("  No envvar found in content: %s\n", content)
			continue
		}

		for _, m := range res {
			name := strings.Trim(string(m[1]), "\"")
			fmt.Printf("  Found envvar: %s at %s\n", name, path)
			currentVars[path+name] = Variable{
				RawName:     name,
				Path:        path,
				FoundInCode: true,
				Name:        name,
			}
		}
	}

	fmt.Printf("Found %d current variables\n", len(currentVars))
	// adjust existing vars
	for i, v := range vars.Variables {
		_, ok := currentVars[v.Path+v.RawName]
		if !ok {
			vars.Variables[i].FoundInCode = false
			continue
		}

		vars.Variables[i].FoundInCode = true
		delete(currentVars, v.Path+v.RawName)
	}

	// add new envvars
	for _, v := range currentVars {
		vars.Variables = append(vars.Variables, v)
	}

	less := func(i, j int) bool {
		return vars.Variables[i].Name < vars.Variables[j].Name
	}

	sort.Slice(vars.Variables, less)

	output, err := yaml.Marshal(vars)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Writing new variable definitions to %s\n", fullYamlPath)
	err = os.WriteFile(fullYamlPath, output, 0666)
	if err != nil {
		log.Fatalf("could not write %s", fullYamlPath)
	}
	if err := os.Chdir(curdir); err != nil {
		log.Fatal(err)
	}
}

// RenderGlobalVarsTemplate renders the global vars template
func RenderGlobalVarsTemplate() {
	curdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fullYamlPath := filepath.Join(curdir, yamlSource)

	content, err := os.ReadFile("../../docs/templates/ADOC_extended.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	targetFolder := "../../docs/services/_includes/adoc/"

	vars := &ConfigVars{}
	fmt.Printf("Reading existing variable definitions from %s\n", fullYamlPath)
	yfile, err := os.ReadFile(fullYamlPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yfile, &vars)
	if err != nil {
		log.Fatal(err)
	}

	targetFile, err := os.Create(filepath.Join(targetFolder, "extended_configvars.adoc"))
	if err != nil {
		log.Fatalf("Failed to create target file: %s", err)
	}
	defer targetFile.Close()

	tpl := template.Must(template.New("").Parse(string(content)))
	if err = tpl.Execute(targetFile, *vars); err != nil {
		log.Fatalf("Failed to execute template: %s", err)
	}
}
