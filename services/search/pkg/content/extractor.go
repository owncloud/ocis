package content

import (
	"context"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

//go:generate mockery --name=Extractor
// Extractor is the interface that wraps the basic Extract method.
type Extractor interface {
	Extract(ctx context.Context, ref *provider.Reference, ri *provider.ResourceInfo) (Document, error)
}
