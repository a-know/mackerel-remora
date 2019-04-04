package config

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
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
	ServiceMetricPlugins []*ServiceMetricPlugin
}

// // Regexpwrapper wraps regexp.Regexp
// type Regexpwrapper struct {
// 	*regexp.Regexp
// }

// // UnmarshalText decodes regexp string
// func (r *Regexpwrapper) UnmarshalText(text []byte) error {
// 	var err error
// 	r.Regexp, err = regexp.Compile(string(text))
// 	return err
// }

// // HostStatus represents host status
// type HostStatus string

// // UnmarshalText decodes host status string
// func (s *HostStatus) UnmarshalText(text []byte) error {
// 	status := string(text)
// 	if status != mackerel.HostStatusWorking &&
// 		status != mackerel.HostStatusStandby &&
// 		status != mackerel.HostStatusMaintenance &&
// 		status != mackerel.HostStatusPoweroff {
// 		return fmt.Errorf("invalid host status: %q", status)
// 	}
// 	*s = HostStatus(status)
// 	return nil
// }

// yaml format
// plugin:
//   metrics:
//     sample:
//       command: ruby /usr/local/bin/sample-plugin.rb
//       user: "sample-user"
//       env:
//         FOO: "FOO BAR"
//         QUX: 'QUX QUUX'

func parseConfig(data []byte) (*Config, error) {
	var conf struct {
		Config `yaml:",inline"`
		Plugin map[string]map[string]struct {
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
	// for name, plugin := range conf.Plugin["metrics"] {
	for name, plugin := range conf.Plugin["service-metrics"] {
		if plugin.Command.IsEmpty() {
			return nil, errors.New("specify command of service-metric plugin")
		}
		conf.Config.ServiceMetricPlugins = append(conf.Config.ServiceMetricPlugins, &ServiceMetricPlugin{
			Name: name, Command: plugin.Command, User: plugin.User, Env: plugin.Env,
			Timeout: time.Duration(plugin.TimeoutSeconds) * time.Second,
		})
	}
	return &conf.Config, nil
}

// Load loads agent configuration
func Load(location string) (*Config, error) {
	var conf *Config

	if location == "" {
		conf = defaultConfig()
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

func fetch(location string) ([]byte, error) {
	u, err := url.Parse(location)
	if err != nil {
		return fetchFile(location)
	}

	switch u.Scheme {
	case "http", "https":
		return fetchHTTP(u)
	default:
		return fetchFile(u.Path)
	}
}

func fetchFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func fetchHTTP(u *url.URL) ([]byte, error) {
	cl := http.Client{
		Timeout: timeout,
	}
	resp, err := cl.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
