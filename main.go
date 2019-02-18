package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/pyama86/wazuh-tailer/wazuh_notifier"
)

type Notifier interface {
	Notify(*wazuh_notifier.Alert) error
}

func newNotifier(c *wazuh_notifier.Config) Notifier {
	switch c.Notifier {
	case "slack":
		return wazuh_notifier.NewSlack(c)
	}
	return nil
}
func main() {
	path := flag.String("config", "/var/ossec/etc/wazuh_slack.toml", "config file path")
	flag.Parse()
	config, err := wazuh_notifier.NewConfig(*path)
	if err != nil {
		log.Fatal(err)
	}

	notifier := newNotifier(config)
	stdin := bufio.NewScanner(os.Stdin)
	if stdin.Scan() {
		a := wazuh_notifier.Alert{}
		err := json.Unmarshal(stdin.Bytes(), &a)
		if err != nil {
			log.Fatal(err)
		}
		err = notifier.Notify(&a)
		if err != nil {
			log.Fatal(err)
		}
	}
}
