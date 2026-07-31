package autoprop

import (
	"context"

	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
)

func NewGoMicroClientCallWrapper() client.CallWrapper {
	return func(cf client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			newCtx := moveOcisMetaToGoMicroMetadata(ctx)
			return cf(newCtx, node, req, rsp, opts)
		}
	}
}

func NewGoMicroServerHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			newCtx := moveGoMicroMetadataToOcisMeta(ctx)
			return h(newCtx, req, rsp)
		}
	}
}

func NewGoMicroServerSubscriberWrapper() server.SubscriberWrapper {
	return func(next server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, msg server.Message) error {
			newCtx := moveGoMicroMetadataToOcisMeta(ctx)
			return next(newCtx, msg)
		}
	}
}

func NewGoMicroClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		w := &clientWrapper{
			Client: c,
		}
		return w
	}
}

type clientWrapper struct {
	client.Client
}

func (w *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	newCtx := moveOcisMetaToGoMicroMetadata(ctx)
	return w.Client.Call(newCtx, req, rsp, opts...)
}

func (w *clientWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	newCtx := moveOcisMetaToGoMicroMetadata(ctx)
	return w.Client.Stream(newCtx, req, opts...)
}

func (w *clientWrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	newCtx := moveOcisMetaToGoMicroMetadata(ctx)
	return w.Client.Publish(newCtx, p, opts...)
}
