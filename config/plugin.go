package config

import (
	"time"

	"github.com/mackerelio/mackerel-container-agent/cmdutil"
)

// ServiceMetricPlugin represents service-metric plugin
type ServiceMetricPlugin struct {
	Name    string
	Command cmdutil.Command
	User    string
	Env     Env
	Timeout time.Duration
}
