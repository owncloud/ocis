package main

import (
	"ociswrapper/cmd"
	"ociswrapper/common"
)

func main() {
	cmd.Execute()

	common.Wg.Wait()
}
