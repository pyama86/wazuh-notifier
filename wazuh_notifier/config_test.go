package wazuh_notifier

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	c, err := NewConfig("test.toml")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if c.Endpoint != "http://example.com" {
		t.Errorf("config cannot parse endpoint")
	}
	if c.SlackToken != "example-token" {
		t.Errorf("config cannot parse slack token")
	}
	if c.Groups["example"].SlackChannel != "example-channel" {
		t.Errorf("config cannot parse slack chaneel")
	}
	if c.Groups["example"].SlackMention != "example-mention" {
		t.Errorf("config cannot parse slack mention")
	}
}
