package policy

import (
	"context"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// TODO use a cheap etcd kv store as store and load data from it...

// storage provide means to store Rego data. This is a temporary storage with minimal caching capabilities for showcasing
// the licensing PoC.
type storage struct {
	// db stores data that goes as input for rego
	client *clientv3.Client
}

type IStorage interface {
	UsersCount() int
}

// NewStorage returns a new IStorage
func NewStorage() IStorage {
	// swallow errors for the POC, as we don't expect everyone to have an etcd cluster running.
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return storage{}
	}

	return storage{
		client: cli,
	}
}

// Users returns a list of users in the system.
func (s storage) UsersCount() int {
	r, err := s.client.Get(context.TODO(), "users_count")
	if err != nil {
		return 0
	}
	v, err := strconv.Atoi(string(r.Kvs[0].Value))
	if err != nil {
		return 0
	}
	return v
}
