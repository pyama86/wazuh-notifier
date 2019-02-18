package wazuh_notifier

import (
	"strings"

	"github.com/BurntSushi/toml"
)

func NewConfig(path string) (*Config, error) {
	var conf Config

	defaultConfig(&conf)

	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return nil, err
	}

	if conf.KibanaURL == "" {
		conf.KibanaURL = strings.Replace(conf.Endpoint, ":55000", "", 1)
	}
	return &conf, nil
}

type Config struct {
	Notifier   string
	Endpoint   string
	KibanaURL  string
	Cert       string
	Key        string
	SlackToken string                 `toml:"slack_token"`
	Groups     map[string]GroupConfig `toml:"groups"`
}

type GroupConfig struct {
	SlackChannel string `toml:"slack_channel"`
	SlackMention string `toml:"slack_mention"`
}

func defaultConfig(c *Config) {
	c.Notifier = "slack"
	c.Endpoint = "https://127.0.0.1:55000"
}
