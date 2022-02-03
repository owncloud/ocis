package registry

import (
	"time"

	revareg "github.com/cs3org/reva/pkg/registry"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/registry"
	microreg "go-micro.dev/v4/registry"
)

func GetRevaRegistry() RevaRegistry {
	return RevaRegistry{
		r: GetRegistry(),
	}
}

type RevaRegistry struct {
	r microreg.Registry
}

func (r RevaRegistry) Add(s revareg.Service) error {
	nodes := []*microreg.Node{}

	for _, n := range s.Nodes() {
		meta := n.Metadata()

		meta["broker"] = broker.String()
		meta["registry"] = microreg.String()

		node := microreg.Node{
			Address:  n.Address(),
			Id:       n.ID(),
			Metadata: n.Metadata(),
		}
		nodes = append(nodes, &node)
	}

	svc := &microreg.Service{
		Name:    s.Name(),
		Version: version.String,
		Nodes:   nodes,
	}
	rOpts := []registry.RegisterOption{registry.RegisterTTL(time.Minute)}

	return r.r.Register(svc, rOpts...)
}

func (r RevaRegistry) GetService(serviceName string) (revareg.Service, error) {

	services, err := r.r.GetService(serviceName)
	if err != nil {
		return Service{}, err
	}

	nodes := []Node{}

	for _, svc := range services {
		for _, nd := range svc.Nodes {
			n := Node{
				address:  nd.Address,
				id:       nd.Id,
				metadata: nd.Metadata,
			}
			nodes = append(nodes, n)
		}
	}

	svc := Service{
		name:  serviceName,
		nodes: nodes,
	}

	return svc, nil

}

type Service struct {
	name  string
	nodes []Node
}

func (s Service) Name() string {
	return s.name
}

func (s Service) Nodes() []revareg.Node {
	nodes := []revareg.Node{}
	for _, n := range s.nodes {
		nodes = append(nodes, n)
	}
	return nodes
}

type Node struct {
	address  string
	metadata map[string]string
	id       string
}

func (n Node) Address() string {
	return n.address
}

func (n Node) Metadata() map[string]string {
	return n.metadata
}

func (n Node) ID() string {
	return n.id
}
