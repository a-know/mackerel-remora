package agent

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

var logger = logging.GetLogger("agent")

// Agent interface
type Agent interface {
	Run([]string) error
}

// NewAgent creates a new Mackerel agent
func NewAgent(version, revision string) Agent {
	return &agent{version, revision}
}

type agent struct {
	version, revision string
}

func (a *agent) Run(_ []string) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP)
	for {
		errCh := make(chan error)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() { errCh <- a.start(ctx) }()
		select {
		case <-sigCh:
			cancel()
		case err := <-errCh:
			return err
		}
	}
}

func (a *agent) start(ctx context.Context) error {
	conf, err := config.Load(os.Getenv("MACKEREL_AGENT_CONFIG"))
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
	client.UserAgent = spec.BuildUserAgent(a.version, a.revision)

	var metricGenerators []metric.Generator
	for _, mp := range conf.ServiceMetricPlugins {
		metricGenerators = append(metricGenerators, metric.NewPluginGenerator(mp))
	}
	metricManager := metric.NewManager(metricGenerators, client)

	// specGenerators := pform.GetSpecGenerators()
	// specManager := spec.NewManager(specGenerators, client).
	// 	WithVersion(a.version, a.revision).
	// 	WithCustomIdentifier(customIdentifier)

	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	// defer signal.Stop(sigCh)

	// return run(ctx, client, metricManager, checkManager, specManager, pform, conf, sigCh)
	return nil
}
