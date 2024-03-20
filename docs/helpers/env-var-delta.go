package main

import (
	"fmt"
	"github.com/rogpeppe/go-internal/semver"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const envVarYamlSource = "env_vars.yaml"

var envVarOutPutTemplates = map[string]string{
	"added":      "templates/env-vars-added.md.tmpl",
	"removed":    "templates/env-vars-removed.md.tmpl",
	"deprecated": "templates/env-vars-deprecated.md.tmpl",
}

// ConfigField represents the env-var annotation in the code
type ConfigField struct {
	Name                string `yaml:"name"`
	DefaultValue        string `yaml:"defaultValue"`
	Type                string `yaml:"type"`
	Description         string `yaml:"description"`
	IntroductionVersion string `yaml:"introductionVersion"`
	DeprecationVersion  string `yaml:"deprecationVersion"`
	RemovalVersion      string `yaml:"removalVersion"`
	DeprecationInfo     string `yaml:"deprecationInfo"`
}

type TemplateData struct {
	StartVersion string
	EndVersion   string
	DeltaFields  []*ConfigField
}

// RenderEnvVarDeltaTable generates tables for env-var deltas
func RenderEnvVarDeltaTable(osArgs []string) {
	if !semver.IsValid(osArgs[2]) {
		log.Fatalf("Start version invalid semver: %s", osArgs[2])
	}
	if !semver.IsValid(osArgs[3]) {
		log.Fatalf("Target version invalid semver: %s", osArgs[3])
	}
	if semver.Compare(osArgs[2], osArgs[3]) >= 0 {
		log.Fatalf("Start version %s is not smaller than target version %s", osArgs[2], osArgs[3])
	}
	if semver.Compare(osArgs[2], "v5.0.0") < 0 {
		log.Fatalf("This tool does not support versions prior v5.0.0, (given %s)", osArgs[2])
	}
	startVersion := osArgs[2]
	endVersion := osArgs[3]
	fmt.Printf("Generating tables for env-var deltas between version %s and %s...\n", startVersion, endVersion)
	curdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fullYamlPath := filepath.Join(curdir, envVarYamlSource)
	configFields := make(map[string]*ConfigField)
	variableList := map[string][]*ConfigField{
		"added":      {},
		"removed":    {},
		"deprecated": {},
	}
	fmt.Printf("Reading existing variable definitions from %s\n", fullYamlPath)
	yfile, err := os.ReadFile(fullYamlPath)
	if err == nil {
		err := yaml.Unmarshal(yfile, configFields)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Success, found %d entries\n", len(configFields))
	for _, field := range configFields {
		if field.IntroductionVersion != "" &&
			field.IntroductionVersion != "pre5.0" &&
			!semver.IsValid(field.IntroductionVersion) &&
			field.IntroductionVersion[0] != 'v' {
			field.IntroductionVersion = "v" + field.IntroductionVersion
		}
		if field.IntroductionVersion != "pre5.0" && !semver.IsValid(field.IntroductionVersion) {
			fmt.Printf("Invalid semver for field %s: %s\n", field.Name, field.IntroductionVersion)
			os.Exit(1)
		}
		//fmt.Printf("Processing field %s dv: %s, iv: %s\n", field.Name, field.DeprecationVersion, field.IntroductionVersion)
		if semver.IsValid(field.RemovalVersion) && semver.Compare(startVersion, field.RemovalVersion) < 0 && semver.Compare(endVersion, field.RemovalVersion) >= 0 {
			variableList["removed"] = append(variableList["removed"], field)
		}
		if semver.IsValid(field.DeprecationVersion) && semver.Compare(startVersion, field.DeprecationVersion) <= 0 && semver.Compare(endVersion, field.DeprecationVersion) > 0 {
			variableList["deprecated"] = append(variableList["deprecated"], field)
		}
		if semver.IsValid(field.IntroductionVersion) && semver.Compare(startVersion, field.IntroductionVersion) <= 0 && semver.Compare(endVersion, field.IntroductionVersion) >= 0 {
			fmt.Printf("Adding field %s iv: %s\n", field.Name, field.IntroductionVersion)
			variableList["added"] = append(variableList["added"], field)
		}
	}
	for templateName, templatePath := range envVarOutPutTemplates {
		content, err := os.ReadFile(templatePath)
		if err != nil {
			log.Fatal(err)
		}
		tpl := template.Must(template.New(templateName).Parse(string(content)))
		err = os.MkdirAll("output/env-deltas", 0700)
		if err != nil {
			log.Fatal(err)
		}
		targetFile, err := os.Create(filepath.Join("output/env-deltas", fmt.Sprintf("%s-%s-%s.md", startVersion, endVersion, templateName)))
		if err != nil {
			log.Fatal(err)
		}
		err = tpl.Execute(targetFile, TemplateData{
			StartVersion: startVersion,
			EndVersion:   endVersion,
			DeltaFields:  variableList[templateName],
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
