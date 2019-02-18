# wazuh-notifier
wazuh alert notification command

## description

wazuh-notifier is alert send to slack channel.
There is a function to ignore the same notification for a certain time(default 1m).
## usage
```bash
$ cat alerts.json | wazuh-notifier -config path/to/config.toml
```

## config

```toml
endpoint = "https://example.com:55000/"
slack_token = "xxxxxxx"
cert = "/path/to/wazuh.crt"
key = "/path/to/wazuh.key"
[groups.example]
slack_channel = "xxxxx"
slack_mention = "xxxxx"
```
