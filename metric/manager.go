package metric

import (
	"context"
	"time"

	"github.com/a-know/mackerel-remora/api"
	"github.com/mackerelio/golib/logging"
	mackerel "github.com/mackerelio/mackerel-client-go"
)

var logger = logging.GetLogger("metric")

// Manager in metric manager
type Manager struct {
	collector  *collector
	sender     *sender
	generators *map[string][]Generator
}

// NewManager creates metric manager instanace
func NewManager(generators map[string][]Generator, client api.Client) *Manager {
	return &Manager{
		collector:  newCollector(generators["ser"]),
		generators: &generators,
		sender:     newSender(client),
	}
}

// Run collect and send metrics
func (m *Manager) Run(ctx context.Context, serviceName string, interval time.Duration) (err error) {
	t := time.NewTicker(interval)
	errCh := make(chan error)
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-t.C:
			go func() {
				if err := m.collectAndPostValues(ctx, serviceName); err != nil {
					errCh <- err
				}
			}()
		case err = <-errCh:
			break loop
		}
	}
	return
}

func (m *Manager) collectAndPostValues(ctx context.Context, serviceName string) error {
	now := time.Now()
	generators := *m.generators
	m.collector = newCollector(generators[serviceName])
	values, err := m.collector.collect(ctx)
	if err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}
	var metricValues []*mackerel.MetricValue
	for name, value := range values {
		metricValues = append(metricValues, &mackerel.MetricValue{
			Name:  name,
			Time:  now.Unix(),
			Value: value,
		})
	}
	return m.sender.post(serviceName, metricValues)
}
