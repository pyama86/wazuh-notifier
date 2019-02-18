package wazuh_notifier

import "github.com/mrtc0/wazuh"

type Wazuh struct {
	c   *Config
	api *wazuh.Client
}

func NewWazuh(c *Config) *Wazuh {
	options := []wazuh.ClientOption{}
	if c.Cert != "" && c.Key != "" {
		options = append(options, wazuh.SetClientCertificate(c.Cert, c.Key))
	}
	return &Wazuh{
		c:   c,
		api: wazuh.New(c.Endpoint, options...),
	}
}

func (w *Wazuh) getGroups(id string) ([]string, error) {
	agent, err := w.api.GetAnAgent(id)
	if err != nil {
		return nil, err
	}
	return agent.Group, nil
}
