package registry

import (
	mRegistry "go-micro.dev/v4/registry"
)

func BuildGRPCService(serviceID, uuid, addr string, version string) *mRegistry.Service {
	node := &mRegistry.Node{
		Id:       serviceID + "-" + uuid,
		Address:  addr,
		Metadata: make(map[string]string),
	}

	node.Metadata["registry"] = GetRegistry().String()
	node.Metadata["server"] = "grpc"
	node.Metadata["transport"] = "grpc"
	node.Metadata["protocol"] = "grpc"

	return &mRegistry.Service{
		Name:      serviceID,
		Version:   version,
		Nodes:     []*mRegistry.Node{node},
		Endpoints: make([]*mRegistry.Endpoint, 0),
	}
}

func BuildHTTPService(serviceID, uuid, addr string, version string) *mRegistry.Service {
	node := &mRegistry.Node{
		Id:       serviceID + "-" + uuid,
		Address:  addr,
		Metadata: make(map[string]string),
	}

	node.Metadata["registry"] = GetRegistry().String()
	node.Metadata["server"] = "http"
	node.Metadata["transport"] = "http"
	node.Metadata["protocol"] = "http"

	return &mRegistry.Service{
		Name:      serviceID,
		Version:   version,
		Nodes:     []*mRegistry.Node{node},
		Endpoints: make([]*mRegistry.Endpoint, 0),
	}
}
