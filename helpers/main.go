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
		case "env-var-delta-table":
			// This step is not covered by the all or default case, because it needs explicit arguments
			if len(os.Args) != 4 {
				fmt.Println("Needs two arguments: env-var-delta-table <first-version> <second-version>")
				fmt.Println("Example: env-var-delta-table v5.0.0 v6.0.0")
				fmt.Println("Will not generate usable results for versions Prior to v5.0.0")
			} else {
				RenderEnvVarDeltaTable(os.Args)
			}
		case "all":
			RenderTemplates()
			GetRogueEnvs()
			RenderGlobalVarsTemplate()
			GenerateServiceIndexMarkdowns()
		case "help":
			fallthrough
		default:
			fmt.Printf("Usage: %s [templates|rogue|globals|service-index|env-var-delta-table|all|help]\n", os.Args[0])
		}
	} else {
		// Left here, even though present in the switch case, for backwards compatibility
		RenderTemplates()
		GetRogueEnvs()
		RenderGlobalVarsTemplate()
		GenerateServiceIndexMarkdowns()
	}
}
