// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package trace

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	// Propagator is the default Reva propagator.
	Propagator      = propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	defaultProvider = revaDefaultTracerProvider{}
)

type revaDefaultTracerProvider struct {
	mutex       sync.RWMutex
	initialized bool
}

// NewTracerProvider returns a new TracerProvider, configure for the specified service
func NewTracerProvider(opts ...Option) trace.TracerProvider {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	if options.TransportCredentials == nil {
		options.TransportCredentials = credentials.NewClientTLSFromCert(nil, "")
	}

	if !options.Enabled {
		return trace.NewNoopTracerProvider()
	}

	// default to 'reva' as service name if not set
	if options.ServiceName == "" {
		options.ServiceName = "reva"
	}

	switch options.Exporter {
	case "otlp":
		return getOtlpTracerProvider(options)
	default:
		return getJaegerTracerProvider(options)
	}
}

// SetDefaultTracerProvider sets the default trace provider
func SetDefaultTracerProvider(tp trace.TracerProvider) {
	otel.SetTracerProvider(tp)
	defaultProvider.mutex.Lock()
	defer defaultProvider.mutex.Unlock()
	defaultProvider.initialized = true
}

// InitDefaultTracerProvider initializes a global default jaeger TracerProvider at a package level.
//
// Deprecated: Use NewTracerProvider and SetDefaultTracerProvider to properly initialize a tracer provider with options
func InitDefaultTracerProvider(collector, endpoint string) {
	defaultProvider.mutex.Lock()
	defer defaultProvider.mutex.Unlock()
	if !defaultProvider.initialized {
		SetDefaultTracerProvider(getJaegerTracerProvider(Options{
			Enabled:     true,
			Collector:   collector,
			Endpoint:    endpoint,
			ServiceName: "reva default jaeger provider",
		}))
	}
}

// DefaultProvider returns the "global" default TracerProvider
// Currently used by the pool to get the global tracer
func DefaultProvider() trace.TracerProvider {
	defaultProvider.mutex.RLock()
	defer defaultProvider.mutex.RUnlock()
	return otel.GetTracerProvider()
}

// getJaegerTracerProvider returns a new TracerProvider, configure for the specified service
func getJaegerTracerProvider(options Options) trace.TracerProvider {
	var exp *jaeger.Exporter
	var err error

	if options.Endpoint != "" {
		var agentHost string
		var agentPort string

		agentHost, agentPort, err = parseAgentConfig(options.Endpoint)
		if err != nil {
			panic(err)
		}

		exp, err = jaeger.New(
			jaeger.WithAgentEndpoint(
				jaeger.WithAgentHost(agentHost),
				jaeger.WithAgentPort(agentPort),
			),
		)
		if err != nil {
			panic(err)
		}
	}

	if options.Collector != "" {
		exp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(options.Collector)))
		if err != nil {
			panic(err)
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(options.ServiceName),
			semconv.HostNameKey.String(hostname),
		)),
	)
}

func parseAgentConfig(ae string) (string, string, error) {
	u, err := url.Parse(ae)
	// as per url.go:
	// [...] Trying to parse a hostname and path
	// without a scheme is invalid but may not necessarily return an
	// error, due to parsing ambiguities.
	if err == nil && u.Hostname() != "" && u.Port() != "" {
		return u.Hostname(), u.Port(), nil
	}

	p := strings.Split(ae, ":")
	if len(p) != 2 {
		return "", "", fmt.Errorf(fmt.Sprintf("invalid agent endpoint `%s`. expected format: `hostname:port`", ae))
	}

	switch {
	case p[0] == "" && p[1] == "": // case ae = ":"
		return "", "", fmt.Errorf(fmt.Sprintf("invalid agent endpoint `%s`. expected format: `hostname:port`", ae))
	case p[0] == "":
		return "", "", fmt.Errorf(fmt.Sprintf("invalid agent endpoint `%s`. expected format: `hostname:port`", ae))
	}
	return p[0], p[1], nil
}

// getOtelTracerProvider returns a new TracerProvider, configure for the specified service
func getOtlpTracerProvider(options Options) trace.TracerProvider {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, options.Endpoint,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection to collector: %w", err))
	}
	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithGRPCConn(conn),
	)

	if err != nil {
		panic(err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", options.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)
}
