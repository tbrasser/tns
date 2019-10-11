package utils

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
	"io/ioutil"
	"os"
)

// SetupTracing installs jaeger tracing lib as opentracing tracer if the environment variable JAEGER_AGENT_HOST
// is set. Logging is controlled by JAEGER_REPORTER_LOG_SPANS env var.
func SetupTracing(serviceName string, logger log.Logger) func() error {
	config := newTracerFromEnv()
	trace := installJaeger(serviceName, config, jaegercfg.Logger(&GokitJaeger{Logger: logger}))
	return trace.Close
}

// Both NewTracerFromEnv and InstallJaeger github.com/weaveworks/common/tracing but modified so you can pass on a
// logger and enable logging of spans for debugging.

func newTracerFromEnv() *jaegercfg.Configuration {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		fmt.Printf("Could not load jaeger tracer configuration: %s\n", err.Error())
		os.Exit(1)
	}

	return cfg
}

func installJaeger(serviceName string, cfg *jaegercfg.Configuration, options ...jaegercfg.Option) io.Closer {
	if cfg.Sampler.SamplingServerURL == "" && cfg.Reporter.LocalAgentHostPort == "" {
		fmt.Printf("Jaeger tracer disabled: No trace report agent or config server specified\n")
		return ioutil.NopCloser(nil)
	}

	metricsFactory := jaegerprom.New()
	options = append(options, jaegercfg.Metrics(metricsFactory))
	closer, err := cfg.InitGlobalTracer(serviceName, options...)
	if err != nil {
		fmt.Printf("Could not initialize jaeger tracer: %s\n", err.Error())
		os.Exit(1)
	}
	return closer
}

type GokitJaeger struct {
	Logger log.Logger
}

func (logger *GokitJaeger) Error(msg string) {
	level.Error(logger.Logger).Log("msg", fmt.Sprintf(msg))
}

func (logger *GokitJaeger) Infof(msg string, args ...interface{}) {
	level.Info(logger.Logger).Log("msg", fmt.Sprintf(msg, args...))
}

