package main

import (
	"io/ioutil"
	"log"
	"os"

	"text/template"

	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/docs"
)

func main() {
	cfg := config.DefaultConfig()
	fields := docs.Display(*cfg)

	content, err := ioutil.ReadFile("docs/templates/CONFIGURATION.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	tpl := template.Must(template.New("").Parse(string(content)))
	tpl.Execute(os.Stdout, fields)
	// for _, f := range fields {
	// 	fmt.Printf("%s %s = %v\t%s\n", f.Name, f.Type, f.DefaultValue, f.Description)
	// }
}
