package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-chi/chi/v5"
	"github.com/owncloud/protoc-gen-microweb/examples/greeter/proto"
)

func main() {
	mux := chi.NewMux()

	proto.RegisterGreeterWeb(
		mux,
		&Greeter{},
	)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type Greeter struct{}

func (g *Greeter) Say(ctx context.Context, in *proto.SayRequest, out *proto.SayResponse) error {
	name := "World"

	if in.Name != "" {
		name = in.Name
	}

	out.Message = fmt.Sprintf("Hello %s!", name)
	return nil
}

func (g *Greeter) SayAnything(ctx context.Context, in *emptypb.Empty, out *proto.SayResponse) error {
	out.Message = "Saying Anything!"
	return nil
}
