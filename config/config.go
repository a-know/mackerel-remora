package config

import (
	"errors"
	"os"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/mackerelio/mackerel-container-agent/cmdutil"
	cconfig "github.com/mackerelio/mackerel-container-agent/config"
)

const (
	timeout     = 3 * time.Second
	defaultRoot = "/var/tmp/mackerel-remora"
)

// Config represents agent configuration
type Config struct {
	Apibase              string `yaml:"apibase"`
	Apikey               string `yaml:"apikey"`
	Root                 string `yaml:"root"`
	ServiceMetricPlugins map[string][]*cconfig.MetricPlugin
}

// yaml format example
//
// apibase: www.example.com
// apikey: qwertyuiop
// root: /var/tmp/mackerel-remora
// plugin:
//   servicemetrics:
//     <servicename>:
//       sample:
//         command: ruby /usr/local/bin/sample-plugin.rb
//         user: "sample-user"
//         env:
//           FOO: "FOO BAR"
//           QUX: 'QUX QUUX'

func parseConfig(data []byte) (*Config, error) {
	var conf struct {
		Config `yaml:",inline"`
		Plugin map[string]map[string]map[string]struct {
			Command        cmdutil.Command `yaml:"command"`
			User           string          `yaml:"user"`
			TimeoutSeconds int             `yaml:"timeoutSeconds"`
			Env            cconfig.Env     `yaml:"env"`
			Memo           string          `yaml:"memo"`
		} `yaml:"plugin"`
	}
	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	for serviceName, plugins := range conf.Plugin["servicemetrics"] {
		for settingName, plugin := range plugins {
			if plugin.Command.IsEmpty() {
				return nil, errors.New("specify command of service-metric plugin")
			}
			conf.Config.ServiceMetricPlugins[serviceName] = append(
				conf.Config.ServiceMetricPlugins[serviceName],
				&cconfig.MetricPlugin{
					Name:    settingName,
					Command: plugin.Command,
					User:    plugin.User,
					Env:     plugin.Env,
					Timeout: time.Duration(plugin.TimeoutSeconds) * time.Second,
				})
		}
	}
	return &conf.Config, nil
}

// Load loads agent configuration
func Load(location string) (*Config, error) {
	var conf *Config

	if location == "" {
		conf = &Config{}
	} else {
		data, err := fetch(location)
		if err != nil {
			return nil, err
		}

		conf, err = parseConfig(data)
		if err != nil {
			return nil, err
		}
	}

	if conf.Apibase == "" {
		conf.Apibase = os.Getenv("MACKEREL_APIBASE")
	}

	if conf.Apikey == "" {
		conf.Apikey = os.Getenv("MACKEREL_APIKEY")
	}

	if conf.Root == "" {
		conf.Root = defaultRoot
	}

	return conf, nil
}
