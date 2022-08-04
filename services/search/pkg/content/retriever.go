package content

import (
	"context"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"io"
)

//go:generate mockery --name=Retriever
// Retriever is the interface that wraps the basic Retrieve method. 🐕
// It requests and then returns a resource from the underlying storage.
type Retriever interface {
	Retrieve(ctx context.Context, ref *provider.Reference, owner *user.User) (io.ReadCloser, error)
}
