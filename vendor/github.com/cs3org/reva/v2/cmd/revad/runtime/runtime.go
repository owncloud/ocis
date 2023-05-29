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

package runtime

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/cs3org/reva/v2/cmd/revad/internal/grace"
	"github.com/cs3org/reva/v2/pkg/logger"
	"github.com/cs3org/reva/v2/pkg/registry"
	"github.com/cs3org/reva/v2/pkg/rgrpc"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/cs3org/reva/v2/pkg/sharedconf"
	rtrace "github.com/cs3org/reva/v2/pkg/trace"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Run runs a reva server with the given config file and pid file.
func Run(mainConf map[string]interface{}, pidFile, logLevel string) {
	logConf := parseLogConfOrDie(mainConf["log"], logLevel)
	logger := initLogger(logConf)
	RunWithOptions(mainConf, pidFile, WithLogger(logger))
}

// RunWithOptions runs a reva server with the given config file, pid file and options.
func RunWithOptions(mainConf map[string]interface{}, pidFile string, opts ...Option) {
	options := newOptions(opts...)
	parseSharedConfOrDie(mainConf["shared"])
	coreConf := parseCoreConfOrDie(mainConf["core"])

	if err := registry.Init(options.Registry); err != nil {
		panic(err)
	}

	run(mainConf, coreConf, options.Logger, pidFile)
}

type coreConf struct {
	MaxCPUs            string `mapstructure:"max_cpus"`
	TracingEnabled     bool   `mapstructure:"tracing_enabled"`
	TracingInsecure    bool   `mapstructure:"tracing_insecure"`
	TracingExporter    string `mapstructure:"tracing_exporter"`
	TracingEndpoint    string `mapstructure:"tracing_endpoint"`
	TracingCollector   string `mapstructure:"tracing_collector"`
	TracingServiceName string `mapstructure:"tracing_service_name"`

	// TracingService specifies the service. i.e OpenCensus, OpenTelemetry, OpenTracing...
	TracingService string `mapstructure:"tracing_service"`
}

func run(mainConf map[string]interface{}, coreConf *coreConf, logger *zerolog.Logger, filename string) {
	host, _ := os.Hostname()
	logger.Info().Msgf("host info: %s", host)

	tp := initTracing(coreConf)
	initCPUCount(coreConf, logger)

	servers := initServers(mainConf, logger, tp)
	watcher, err := initWatcher(logger, filename)
	if err != nil {
		log.Panic(err)
	}
	listeners := initListeners(watcher, servers, logger)

	start(mainConf, servers, listeners, logger, watcher)
}

func initListeners(watcher *grace.Watcher, servers map[string]grace.Server, log *zerolog.Logger) map[string]net.Listener {
	listeners, err := watcher.GetListeners(servers)
	if err != nil {
		log.Error().Err(err).Msg("error getting sockets")
		watcher.Exit(1)
	}
	return listeners
}

func initWatcher(log *zerolog.Logger, filename string) (*grace.Watcher, error) {
	watcher, err := handlePIDFlag(log, filename)
	// TODO(labkode): maybe pidfile can be created later on? like once a server is going to be created?
	if err != nil {
		log.Error().Err(err).Msg("error creating grace watcher")
		os.Exit(1)
	}
	return watcher, err
}

func initServers(mainConf map[string]interface{}, log *zerolog.Logger, tp trace.TracerProvider) map[string]grace.Server {
	servers := map[string]grace.Server{}
	if isEnabledHTTP(mainConf) {
		s, err := getHTTPServer(mainConf["http"], log, tp)
		if err != nil {
			log.Error().Err(err).Msg("error creating http server")
			os.Exit(1)
		}
		servers["http"] = s
	}

	if isEnabledGRPC(mainConf) {
		s, err := getGRPCServer(mainConf["grpc"], log, tp)
		if err != nil {
			log.Error().Err(err).Msg("error creating grpc server")
			os.Exit(1)
		}
		servers["grpc"] = s
	}

	if len(servers) == 0 {
		log.Info().Msg("nothing to do, no grpc/http enabled_services declared in config")
		os.Exit(1)
	}
	return servers
}

func initTracing(conf *coreConf) trace.TracerProvider {
	if conf.TracingEnabled {
		opts := []rtrace.Option{
			rtrace.WithExporter(conf.TracingExporter),
			rtrace.WithEndpoint(conf.TracingEndpoint),
			rtrace.WithCollector(conf.TracingCollector),
			rtrace.WithServiceName(conf.TracingServiceName),
		}
		if conf.TracingEnabled {
			opts = append(opts, rtrace.WithEnabled())
		}
		if conf.TracingInsecure {
			opts = append(opts, rtrace.WithInsecure())
		}
		tp := rtrace.NewTracerProvider(opts...)
		rtrace.SetDefaultTracerProvider(tp)
		return tp
	}
	return rtrace.DefaultProvider()
}

func initCPUCount(conf *coreConf, log *zerolog.Logger) {
	ncpus, err := adjustCPU(conf.MaxCPUs)
	if err != nil {
		log.Error().Err(err).Msg("error adjusting number of cpus")
		os.Exit(1)
	}
	// log.Info().Msgf("%s", getVersionString())
	log.Info().Msgf("running on %d cpus", ncpus)
}

func initLogger(conf *logConf) *zerolog.Logger {
	log, err := newLogger(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating logger, exiting ...")
		os.Exit(1)
	}
	return log
}

func handlePIDFlag(l *zerolog.Logger, pidFile string) (*grace.Watcher, error) {
	var opts []grace.Option
	opts = append(opts, grace.WithPIDFile(pidFile))
	opts = append(opts, grace.WithLogger(l.With().Str("pkg", "grace").Logger()))
	w := grace.NewWatcher(opts...)
	err := w.WritePID()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func start(mainConf map[string]interface{}, servers map[string]grace.Server, listeners map[string]net.Listener, log *zerolog.Logger, watcher *grace.Watcher) {
	if isEnabledHTTP(mainConf) {
		go func() {
			if err := servers["http"].(*rhttp.Server).Start(listeners["http"]); err != nil {
				log.Error().Err(err).Msg("error starting the http server")
				watcher.Exit(1)
			}
		}()
	}
	if isEnabledGRPC(mainConf) {
		go func() {
			if err := servers["grpc"].(*rgrpc.Server).Start(listeners["grpc"]); err != nil {
				log.Error().Err(err).Msg("error starting the grpc server")
				watcher.Exit(1)
			}
		}()
	}
	watcher.TrapSignals()
}

func newLogger(conf *logConf) (*zerolog.Logger, error) {
	// TODO(labkode): use debug level rather than info as default until reaching a stable version.
	// Helps having smaller development files.
	if conf.Level == "" {
		conf.Level = zerolog.DebugLevel.String()
	}

	var opts []logger.Option
	opts = append(opts, logger.WithLevel(conf.Level))

	w, err := getWriter(conf.Output)
	if err != nil {
		return nil, err
	}

	opts = append(opts, logger.WithWriter(w, logger.Mode(conf.Mode)))

	l := logger.New(opts...)
	sub := l.With().Int("pid", os.Getpid()).Logger()
	return &sub, nil
}

func getWriter(out string) (io.Writer, error) {
	if out == "stderr" || out == "" {
		return os.Stderr, nil
	}

	if out == "stdout" {
		return os.Stdout, nil
	}

	fd, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = errors.Wrap(err, "error creating log file: "+out)
		return nil, err
	}

	return fd, nil
}

func getGRPCServer(conf interface{}, l *zerolog.Logger, tp trace.TracerProvider) (*rgrpc.Server, error) {
	sub := l.With().Str("pkg", "rgrpc").Logger()
	s, err := rgrpc.NewServer(conf, sub, tp)
	if err != nil {
		err = errors.Wrap(err, "main: error creating grpc server")
		return nil, err
	}
	return s, nil
}

func getHTTPServer(conf interface{}, l *zerolog.Logger, tp trace.TracerProvider) (*rhttp.Server, error) {
	sub := l.With().Str("pkg", "rhttp").Logger()
	s, err := rhttp.New(conf, sub, tp)
	if err != nil {
		err = errors.Wrap(err, "main: error creating http server")
		return nil, err
	}
	return s, nil
}

// adjustCPU parses string cpu and sets GOMAXPROCS
// according to its value. It accepts either
// a number (e.g. 3) or a percent (e.g. 50%).
// Default is to use all available cores.
func adjustCPU(cpu string) (int, error) {
	var numCPU int

	availCPU := runtime.NumCPU()

	if cpu != "" {
		if strings.HasSuffix(cpu, "%") {
			// Percent
			var percent float32
			pctStr := cpu[:len(cpu)-1]
			pctInt, err := strconv.Atoi(pctStr)
			if err != nil || pctInt < 1 || pctInt > 100 {
				return 0, fmt.Errorf("invalid CPU value: percentage must be between 1-100")
			}
			percent = float32(pctInt) / 100
			numCPU = int(float32(availCPU) * percent)
		} else {
			// Number
			num, err := strconv.Atoi(cpu)
			if err != nil || num < 1 {
				return 0, fmt.Errorf("invalid CPU value: provide a number or percent greater than 0")
			}
			numCPU = num
		}
	} else {
		numCPU = availCPU
	}

	if numCPU > availCPU || numCPU == 0 {
		numCPU = availCPU
	}

	runtime.GOMAXPROCS(numCPU)
	return numCPU, nil
}

func parseCoreConfOrDie(v interface{}) *coreConf {
	c := &coreConf{}
	if err := mapstructure.Decode(v, c); err != nil {
		fmt.Fprintf(os.Stderr, "error decoding core config: %s\n", err.Error())
		os.Exit(1)
	}

	// tracing defaults to enabled if not explicitly configured
	if v == nil {
		c.TracingEnabled = true
		c.TracingEndpoint = "localhost:6831"
	} else if _, ok := v.(map[string]interface{})["tracing_enabled"]; !ok {
		c.TracingEnabled = true
		c.TracingEndpoint = "localhost:6831"
	}

	return c
}

func parseSharedConfOrDie(v interface{}) {
	if err := sharedconf.Decode(v); err != nil {
		fmt.Fprintf(os.Stderr, "error decoding shared config: %s\n", err.Error())
		os.Exit(1)
	}
}

func parseLogConfOrDie(v interface{}, logLevel string) *logConf {
	c := &logConf{}
	if err := mapstructure.Decode(v, c); err != nil {
		fmt.Fprintf(os.Stderr, "error decoding log config: %s\n", err.Error())
		os.Exit(1)
	}

	// if mode is not set, we use console mode, easier for devs
	if c.Mode == "" {
		c.Mode = "console"
	}

	// Give priority to the log level passed through the command line.
	if logLevel != "" {
		c.Level = logLevel
	}

	return c
}

type logConf struct {
	Output string `mapstructure:"output"`
	Mode   string `mapstructure:"mode"`
	Level  string `mapstructure:"level"`
}

func isEnabledHTTP(conf map[string]interface{}) bool {
	return isEnabled("http", conf)
}

func isEnabledGRPC(conf map[string]interface{}) bool {
	return isEnabled("grpc", conf)
}

func isEnabled(key string, conf map[string]interface{}) bool {
	if a, ok := conf[key]; ok {
		if b, ok := a.(map[string]interface{}); ok {
			if c, ok := b["services"]; ok {
				if d, ok := c.(map[string]interface{}); ok {
					if len(d) > 0 {
						return true
					}
				}
			}
		}
	}
	return false
}
