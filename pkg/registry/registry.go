// Package registry provides accessors to runtime services
package registry

import (
	"github.com/micro/go-micro/v2/store"
)

type Registry {
	Store store.Store
}