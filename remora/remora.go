package remora

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/a-know/mackerel-remora/api"
	"github.com/a-know/mackerel-remora/config"
	"github.com/a-know/mackerel-remora/metric"
	"github.com/a-know/mackerel-remora/spec"
	"github.com/mackerelio/golib/logging"
)

var logger = logging.GetLogger("remora")

// Remora interface
type Remora interface {
	Run([]string) error
}

// NewRemora creates a new Mackerel remora agent
func NewRemora(version, revision string) Remora {
	return &remora{version, revision}
}

type remora struct {
	version, revision string
}

func (r *remora) Run(_ []string) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP)
	for {
		errCh := make(chan error)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() { errCh <- r.start(ctx) }()
		select {
		case <-sigCh:
			cancel()
		case err := <-errCh:
			return err
		}
	}
}

func (r *remora) start(ctx context.Context) error {
	conf, err := config.Load(os.Getenv("MACKEREL_REMORA_CONFIG"))
	if err != nil {
		return err
	}

	client, err := api.NewClient(conf)
	if err != nil {
		return err
	}
	client.UserAgent = spec.BuildUserAgent(r.version, r.revision)

	metricGenerators := map[string][]metric.Generator{}
	for serviceName, plugins := range conf.ServiceMetricPlugins {
		for _, mp := range plugins {
			metricGenerators[serviceName] = append(metricGenerators[serviceName], metric.NewPluginGenerator(mp))
		}
	}
	metricManager := metric.NewManager(metricGenerators, client)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(sigCh)

	return run(ctx, client, metricManager, conf, sigCh)
}
