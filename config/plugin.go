package config

import (
	"time"

	"github.com/a-know/mackerel-remora/cmdutil"
)

// ServiceMetricPlugin represents service-metric plugin
type ServiceMetricPlugin struct {
	Name    string
	Command cmdutil.Command
	User    string
	Env     Env
	Timeout time.Duration
}
