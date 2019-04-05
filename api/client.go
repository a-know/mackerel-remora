package api

import (
	"net/url"

	"github.com/a-know/mackerel-remora/config"
	mackerel "github.com/mackerelio/mackerel-client-go"
)

// Client represents a client of Mackerel API
type Client interface {
	PostServiceMetricValues(serviceName string, metricValues []*mackerel.MetricValue) error
}

// NewClient initialize and return Mackerel API Client
func NewClient(conf *config.Config) (*mackerel.Client, error) {
	client := mackerel.NewClient(conf.Apikey)
	if conf.Apibase != "" {
		baseURL, err := url.Parse(conf.Apibase)
		if err != nil {
			return nil, err
		}
		client.BaseURL = baseURL
	}
	return client, nil
}
