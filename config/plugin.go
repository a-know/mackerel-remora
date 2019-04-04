package config

import (
	"time"

	"github.com/mackerelio/mackerel-container-agent/cmdutil"
	cconfig "github.com/mackerelio/mackerel-container-agent/config"
)

// ServiceMetricPlugin represents service-metric plugin
type ServiceMetricPlugin struct {
	Name    string
	Command cmdutil.Command
	User    string
	Env     cconfig.Env
	Timeout time.Duration
}
