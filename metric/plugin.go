package metric

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-container-agent/cmdutil"
	cconfig "github.com/mackerelio/mackerel-container-agent/config"
)

const (
	pluginPrefix = "service."
)

type pluginGenerator struct {
	cconfig.MetricPlugin
}

// NewPluginGenerator creates a new plugin generator
func NewPluginGenerator(p *cconfig.MetricPlugin) Generator {
	return &pluginGenerator{*p}
}

// Generate generates metric values
func (g *pluginGenerator) Generate(ctx context.Context) (Values, error) {
	stdout, stderr, _, err := cmdutil.RunCommand(ctx, g.Command, g.User, g.Env, g.Timeout)

	if stderr != "" {
		logger.Infof("plugin %s (%s): %q", g.Name, g.Command, stderr)
	}
	if err != nil {
		return nil, fmt.Errorf("plugin %s (%s): %s", g.Name, g.Command, err)
	}

	values := make(Values)
	for _, line := range strings.Split(stdout, "\n") {
		// key, value, timestamp
		xs := strings.Fields(line)
		if len(xs) < 3 {
			continue
		}
		value, err := strconv.ParseFloat(xs[1], 64)
		if err != nil {
			continue
		}
		values[pluginPrefix+xs[0]] = value
	}

	return values, nil
}
