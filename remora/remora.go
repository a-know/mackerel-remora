package remora

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	mackerel "github.com/mackerelio/mackerel-client-go"

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

	client := mackerel.NewClient(conf.Apikey)
	if conf.Apibase != "" {
		baseURL, err := url.Parse(conf.Apibase)
		if err != nil {
			return err
		}
		client.BaseURL = baseURL
	}
	client.UserAgent = spec.BuildUserAgent(r.version, r.revision)

	var metricGenerators []metric.Generator
	for _, mp := range conf.ServiceMetricPlugins {
		metricGenerators = append(metricGenerators, metric.NewPluginGenerator(mp))
	}
	metricManager := metric.NewManager(metricGenerators, client)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(sigCh)

	return run(ctx, client, metricManager, conf, sigCh)
}
