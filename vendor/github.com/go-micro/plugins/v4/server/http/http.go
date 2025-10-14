// Package http implements a go-micro.Server
package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sort"
	"sync"
	"time"

	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/codec"
	"go-micro.dev/v4/codec/jsonrpc"
	"go-micro.dev/v4/codec/protorpc"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go-micro.dev/v4/util/cmd"
)

var (
	defaultCodecs = map[string]codec.NewCodec{
		"application/json":         jsonrpc.NewCodec,
		"application/json-rpc":     jsonrpc.NewCodec,
		"application/protobuf":     protorpc.NewCodec,
		"application/proto-rpc":    protorpc.NewCodec,
		"application/octet-stream": protorpc.NewCodec,
	}
)

type httpServer struct {
	sync.Mutex
	opts         server.Options
	hd           server.Handler
	exit         chan chan error
	registerOnce sync.Once
	subscribers  map[*httpSubscriber][]broker.Subscriber
	// used for first registration
	registered bool
}

func init() {
	cmd.DefaultServers["http"] = NewServer
}

func (h *httpServer) newCodec(contentType string) (codec.NewCodec, error) {
	if cf, ok := h.opts.Codecs[contentType]; ok {
		return cf, nil
	}
	if cf, ok := defaultCodecs[contentType]; ok {
		return cf, nil
	}
	return nil, fmt.Errorf("Unsupported Content-Type: %s", contentType)
}

func (h *httpServer) getListener() net.Listener {
	if h.opts.Context == nil {
		return nil
	}

	if l, ok := h.opts.Context.Value(netListener{}).(net.Listener); ok && l != nil {
		return l
	}

	return nil
}

func (h *httpServer) Options() server.Options {
	h.Lock()
	opts := h.opts
	h.Unlock()
	return opts
}

func (h *httpServer) Init(opts ...server.Option) error {
	h.Lock()
	for _, o := range opts {
		o(&h.opts)
	}
	h.Unlock()
	return nil
}

func (h *httpServer) Handle(handler server.Handler) error {
	if _, ok := handler.Handler().(http.Handler); !ok {
		return errors.New("Handle requires http.Handler")
	}
	h.Lock()
	h.hd = handler
	h.Unlock()
	return nil
}

func (h *httpServer) NewHandler(handler interface{}, opts ...server.HandlerOption) server.Handler {
	options := server.HandlerOptions{
		Metadata: make(map[string]map[string]string),
	}

	for _, o := range opts {
		o(&options)
	}

	var eps []*registry.Endpoint

	if !options.Internal {
		for name, metadata := range options.Metadata {
			eps = append(eps, &registry.Endpoint{
				Name:     name,
				Metadata: metadata,
			})
		}
	}

	return &httpHandler{
		eps:  eps,
		hd:   handler,
		opts: options,
	}
}

func (h *httpServer) NewSubscriber(topic string, handler interface{}, opts ...server.SubscriberOption) server.Subscriber {
	return newSubscriber(topic, handler, opts...)
}

func (h *httpServer) Subscribe(sb server.Subscriber) error {
	sub, ok := sb.(*httpSubscriber)
	if !ok {
		return fmt.Errorf("invalid subscriber: expected *httpSubscriber")
	}
	if len(sub.handlers) == 0 {
		return fmt.Errorf("invalid subscriber: no handler functions")
	}

	if err := validateSubscriber(sb); err != nil {
		return err
	}

	h.Lock()
	defer h.Unlock()
	_, ok = h.subscribers[sub]
	if ok {
		return fmt.Errorf("subscriber %v already exists", h)
	}
	h.subscribers[sub] = nil
	return nil
}

func (h *httpServer) Register() error {
	h.Lock()
	opts := h.opts
	eps := h.hd.Endpoints()
	h.Unlock()

	service := serviceDef(opts)
	service.Endpoints = eps

	h.Lock()
	var subscriberList []*httpSubscriber
	for e := range h.subscribers {
		// Only advertise non internal subscribers
		if !e.Options().Internal {
			subscriberList = append(subscriberList, e)
		}
	}
	sort.Slice(subscriberList, func(i, j int) bool {
		return subscriberList[i].topic > subscriberList[j].topic
	})
	for _, e := range subscriberList {
		service.Endpoints = append(service.Endpoints, e.Endpoints()...)
	}
	h.Unlock()

	rOpts := []registry.RegisterOption{
		registry.RegisterTTL(opts.RegisterTTL),
	}

	h.registerOnce.Do(func() {
		opts.Logger.Logf(logger.InfoLevel, "Registering node: %s", opts.Name+"-"+opts.Id)
	})

	if err := opts.Registry.Register(service, rOpts...); err != nil {
		return err
	}

	h.Lock()
	defer h.Unlock()

	if h.registered {
		return nil
	}
	h.registered = true

	for sb := range h.subscribers {
		handler := h.createSubHandler(sb, opts)
		var subOpts []broker.SubscribeOption
		if queue := sb.Options().Queue; len(queue) > 0 {
			subOpts = append(subOpts, broker.Queue(queue))
		}

		if cx := sb.Options().Context; cx != nil {
			subOpts = append(subOpts, broker.SubscribeContext(cx))
		}

		if !sb.Options().AutoAck {
			subOpts = append(subOpts, broker.DisableAutoAck())
		}

		sub, err := opts.Broker.Subscribe(sb.Topic(), handler, subOpts...)
		if err != nil {
			return err
		}
		h.subscribers[sb] = []broker.Subscriber{sub}
	}
	return nil
}

func (h *httpServer) Deregister() error {
	h.Lock()
	opts := h.opts
	h.Unlock()

	opts.Logger.Logf(logger.InfoLevel, "Deregistering node: %s", opts.Name+"-"+opts.Id)

	service := serviceDef(opts)
	if err := opts.Registry.Deregister(service); err != nil {
		return err
	}

	h.Lock()

	if !h.registered {
		h.Unlock()
		return nil
	}
	h.registered = false

	wg := sync.WaitGroup{}
	for sb, subs := range h.subscribers {
		for _, sub := range subs {
			wg.Add(1)
			go func(s broker.Subscriber) {
				defer wg.Done()
				opts.Logger.Logf(logger.InfoLevel, "Unsubscribing from topic: %s", s.Topic())
				s.Unsubscribe()
			}(sub)
		}
		h.subscribers[sb] = nil
	}

	wg.Wait()
	h.Unlock()
	return nil
}

func (h *httpServer) Start() error {
	h.Lock()
	opts := h.opts
	hd := h.hd
	h.Unlock()

	var (
		ln  net.Listener
		err error
	)

	if l := h.getListener(); l != nil {
		ln = l
	} else if opts.TLSConfig != nil {
		ln, err = tls.Listen("tcp", opts.Address, opts.TLSConfig)
	} else {
		ln, err = net.Listen("tcp", opts.Address)
	}

	if err != nil {
		return err
	}

	opts.Logger.Logf(logger.InfoLevel, "Listening on %s", ln.Addr().String())

	h.Lock()
	h.opts.Address = ln.Addr().String()
	h.Unlock()

	handler, ok := hd.Handler().(http.Handler)
	if !ok {
		return errors.New("Server required http.Handler")
	}

	if err = opts.Broker.Connect(); err != nil {
		return err
	}

	// register
	if err = h.Register(); err != nil {
		return err
	}

	go http.Serve(ln, handler)

	go func() {
		t := new(time.Ticker)

		// only process if it exists
		if opts.RegisterInterval > time.Duration(0) {
			// new ticker
			t = time.NewTicker(opts.RegisterInterval)
		}

		// return error chan
		var ch chan error

	Loop:
		for {
			select {
			// register self on interval
			case <-t.C:
				if err := h.Register(); err != nil {
					opts.Logger.Log(logger.ErrorLevel, "Server register error: ", err)
				}
			// wait for exit
			case ch = <-h.exit:
				break Loop
			}
		}

		// deregister
		h.Deregister()

		opts.Broker.Disconnect()

		// Solve the problem of early exit
		ch <- ln.Close()
	}()

	return nil
}

func (h *httpServer) Stop() error {
	ch := make(chan error)
	h.exit <- ch
	return <-ch
}

func (h *httpServer) String() string {
	if h.opts.TLSConfig != nil {
		return "https"
	}
	return "http"
}

func newServer(opts ...server.Option) server.Server {
	return &httpServer{
		opts:        newOptions(opts...),
		exit:        make(chan chan error),
		subscribers: make(map[*httpSubscriber][]broker.Subscriber),
	}
}

func NewServer(opts ...server.Option) server.Server {
	return newServer(opts...)
}
