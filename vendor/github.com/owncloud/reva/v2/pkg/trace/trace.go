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
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
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
		return noop.NewTracerProvider()
	}

	// default to 'reva' as service name if not set
	if options.ServiceName == "" {
		options.ServiceName = "reva"
	}

	return getOtlpTracerProvider(options)
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
		SetDefaultTracerProvider(getOtlpTracerProvider(Options{
			Endpoint:    endpoint,
			ServiceName: "reva default otlp provider",
			Insecure:    true,
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

// getOtelTracerProvider returns a new TracerProvider, configure for the specified service
func getOtlpTracerProvider(options Options) trace.TracerProvider {
	transportCredentials := options.TransportCredentials
	if options.Insecure {
		transportCredentials = insecure.NewCredentials()
	}
	conn, err := grpc.NewClient(options.Endpoint,
		grpc.WithTransportCredentials(transportCredentials),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection to endpoint: %w", err))
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
