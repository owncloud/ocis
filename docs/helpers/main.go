package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "templates":
			RenderTemplates()
		case "rogue":
			GetRogueEnvs()
		case "globals":
			RenderGlobalVarsTemplate()
		case "service-index":
			GenerateServiceIndexMarkdowns()
		case "all":
			RenderTemplates()
			GetRogueEnvs()
			RenderGlobalVarsTemplate()
			GenerateServiceIndexMarkdowns()
		case "help":
			fallthrough
		default:
			fmt.Println("Usage: [templates|rogue|globals|service-index|all]")
		}
	} else {
		// Left here, even though present in the switch case, for backwards compatibility
		RenderTemplates()
		GetRogueEnvs()
		RenderGlobalVarsTemplate()
		GenerateServiceIndexMarkdowns()
	}
}
