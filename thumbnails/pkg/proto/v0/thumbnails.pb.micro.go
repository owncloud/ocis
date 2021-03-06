// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: thumbnails.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	math "math"
)

import (
	context "context"
	api "github.com/asim/go-micro/v3/api"
	client "github.com/asim/go-micro/v3/client"
	server "github.com/asim/go-micro/v3/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for ThumbnailService service

func NewThumbnailServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for ThumbnailService service

type ThumbnailService interface {
	// Generates the thumbnail and returns it.
	GetThumbnail(ctx context.Context, in *GetThumbnailRequest, opts ...client.CallOption) (*GetThumbnailResponse, error)
}

type thumbnailService struct {
	c    client.Client
	name string
}

func NewThumbnailService(name string, c client.Client) ThumbnailService {
	return &thumbnailService{
		c:    c,
		name: name,
	}
}

func (c *thumbnailService) GetThumbnail(ctx context.Context, in *GetThumbnailRequest, opts ...client.CallOption) (*GetThumbnailResponse, error) {
	req := c.c.NewRequest(c.name, "ThumbnailService.GetThumbnail", in)
	out := new(GetThumbnailResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ThumbnailService service

type ThumbnailServiceHandler interface {
	// Generates the thumbnail and returns it.
	GetThumbnail(context.Context, *GetThumbnailRequest, *GetThumbnailResponse) error
}

func RegisterThumbnailServiceHandler(s server.Server, hdlr ThumbnailServiceHandler, opts ...server.HandlerOption) error {
	type thumbnailService interface {
		GetThumbnail(ctx context.Context, in *GetThumbnailRequest, out *GetThumbnailResponse) error
	}
	type ThumbnailService struct {
		thumbnailService
	}
	h := &thumbnailServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&ThumbnailService{h}, opts...))
}

type thumbnailServiceHandler struct {
	ThumbnailServiceHandler
}

func (h *thumbnailServiceHandler) GetThumbnail(ctx context.Context, in *GetThumbnailRequest, out *GetThumbnailResponse) error {
	return h.ThumbnailServiceHandler.GetThumbnail(ctx, in, out)
}
