package remora

import (
	"context"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/a-know/mackerel-remora/api"
	"github.com/a-know/mackerel-remora/config"
	"github.com/a-know/mackerel-remora/metric"
)

var (
	metricsInterval = time.Minute
)

func run(
	ctx context.Context,
	client api.Client,
	metricManager *metric.Manager,
	conf *config.Config,
	sigCh <-chan os.Signal,
) error {
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	defer cancel()

	var sig os.Signal
	eg.Go(func() error {
		select {
		case sig = <-sigCh:
			cancel()
			return nil
		case <-ctx.Done():
		}
		return nil
	})

	for serviceName := range conf.ServiceMetricPlugins {
		eg.Go(func() error {
			return metricManager.Run(ctx, serviceName, metricsInterval)
		})
	}

	err := eg.Wait()

	if sig != nil {
		logger.Infof("stop the remora: signal = %s", sig)
	}
	return err
}
