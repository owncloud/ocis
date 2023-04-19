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

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	// Propagator is the default Reva propagator.
	Propagator      = propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	defaultProvider = revaDefaultTracerProvider{
		provider: trace.NewNoopTracerProvider(),
	}
)

type revaDefaultTracerProvider struct {
	mutex       sync.RWMutex
	initialized bool
	provider    trace.TracerProvider
}

type ctxKey struct{}

// ContextSetTracerProvider returns a copy of ctx with p associated.
func ContextSetTracerProvider(ctx context.Context, p trace.TracerProvider) context.Context {
	if tp, ok := ctx.Value(ctxKey{}).(trace.TracerProvider); ok {
		if tp == p {
			return ctx
		}
	}
	return context.WithValue(ctx, ctxKey{}, p)
}

// ContextGetTracerProvider returns the TracerProvider associated with the ctx.
// If no TracerProvider is associated is associated, the global default TracerProvider
// is returned
func ContextGetTracerProvider(ctx context.Context) trace.TracerProvider {
	if p, ok := ctx.Value(ctxKey{}).(trace.TracerProvider); ok {
		return p
	}
	return DefaultProvider()
}

// InitDefaultTracerProvider initializes a global default TracerProvider at a package level.
func InitDefaultTracerProvider(collectorEndpoint string, agentEndpoint string) {
	defaultProvider.mutex.Lock()
	defer defaultProvider.mutex.Unlock()
	if !defaultProvider.initialized {
		defaultProvider.provider = GetTracerProvider(true, collectorEndpoint, agentEndpoint, "reva default provider")
	}
	defaultProvider.initialized = true
}

// DefaultProvider returns the "global" default TracerProvider
func DefaultProvider() trace.TracerProvider {
	defaultProvider.mutex.RLock()
	defer defaultProvider.mutex.RUnlock()
	return defaultProvider.provider
}

// GetTracerProvider returns a new TracerProvider, configure for the specified service
func GetTracerProvider(enabled bool, collectorEndpoint string, agentEndpoint, serviceName string) trace.TracerProvider {
	if !enabled {
		return trace.NewNoopTracerProvider()
	}

	// default to 'reva' as service name if not set
	if serviceName == "" {
		serviceName = "reva"
	}

	var exp *jaeger.Exporter
	var err error

	if agentEndpoint != "" {
		var agentHost string
		var agentPort string

		agentHost, agentPort, err = parseAgentConfig(agentEndpoint)
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

	if collectorEndpoint != "" {
		exp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(collectorEndpoint)))
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
			semconv.ServiceNameKey.String(serviceName),
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
