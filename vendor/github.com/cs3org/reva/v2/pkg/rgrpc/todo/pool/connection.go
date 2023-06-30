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

package pool

import (
	"crypto/tls"

	rtrace "github.com/cs3org/reva/v2/pkg/trace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	maxCallRecvMsgSize = 10240000
)

// NewConn creates a new connection to a grpc server
// with open census tracing support.
// TODO(labkode): make grpc tls configurable.
// TODO make maxCallRecvMsgSize configurable, raised from the default 4MB to be able to list 10k files
func NewConn(address string, opts ...Option) (*grpc.ClientConn, error) {

	options := ClientOptions{}
	if err := options.init(); err != nil {
		return nil, err
	}

	// then overwrite with supplied options
	for _, opt := range opts {
		opt(&options)
	}

	var cred credentials.TransportCredentials
	switch options.tlsMode {
	case TLSOff:
		cred = insecure.NewCredentials()
	case TLSInsecure:
		tlsConfig := tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		}
		cred = credentials.NewTLS(&tlsConfig)
	case TLSOn:
		if options.caCert != "" {
			var err error
			if cred, err = credentials.NewClientTLSFromFile(options.caCert, ""); err != nil {
				return nil, err
			}
		} else {
			// Use system's cert pool
			cred = credentials.NewTLS(&tls.Config{})
		}
	}

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(cred),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize),
		),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor(
			otelgrpc.WithTracerProvider(
				options.tracerProvider,
			),
			otelgrpc.WithPropagators(
				rtrace.Propagator,
			),
		)),
		grpc.WithUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(
				otelgrpc.WithTracerProvider(
					options.tracerProvider,
				),
				otelgrpc.WithPropagators(
					rtrace.Propagator,
				),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
