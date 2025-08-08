package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var supportedFeatures uint64 = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	pgs.Init(
		pgs.DebugEnv("DEBUG"),
		pgs.SupportedFeatures(&supportedFeatures),
	).RegisterModule(
		MicroWeb(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}
