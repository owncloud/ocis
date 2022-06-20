package decorators

import (
	"context"

	thumbnailssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/thumbnails/v0"
)

// Interface acting as facade, holding all the interfaces that this
// thumbnails microservice is expecting to implement.
// For now, only the thumbnailssvc.ThumbnailServiceHandler is present,
// but a future configsvc.ConfigServiceHandler is expected to be added here
//
// This interface will also act as the base interface to implement
// a decorator pattern.
type DecoratedService interface {
	thumbnailssvc.ThumbnailServiceHandler
}

// Base type to implement the decorators. It will provide a basic implementation
// by delegating to the decoratedService
//
// Expected implementations will be like:
// ```
// type MyDecorator struct {
//   Decorator
//   myCustomOpts *opts
//   additionalSrv *srv
// }
//
// func NewMyDecorator(next DecoratedService, customOpts *customOpts) DecoratedService {
//   .....
//   return MyDecorator{
//     Decorator: Decorator{next: next},
//     myCustomOpts: opts,
//     additionalSrv: srv,
//   }
// }
// ```
type Decorator struct {
	next DecoratedService
}

// Base implementation for the GetThumbnail (for the thumbnailssvc).
// It will just delegate to the underlying decoratedService
//
// Your custom decorator is expected to overwrite this function,
// but it MUST call the underlying decoratedService at some point
// ```
// func (d MyDecorator) GetThumbnail(ctx context.Context, req *thumbnailssvc.GetThumbnailRequest, resp *thumbnailssvc.GetThumbnailResponse) error {
//   doSomething()
//   err := d.next.GetThumbnail(ctx, req, resp)
//   doAnotherThing()
//   return err
// }
// ```
func (deco Decorator) GetThumbnail(ctx context.Context, req *thumbnailssvc.GetThumbnailRequest, resp *thumbnailssvc.GetThumbnailResponse) error {
	return deco.next.GetThumbnail(ctx, req, resp)
}
