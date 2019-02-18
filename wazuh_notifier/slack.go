package wazuh_notifier

import (
	"errors"
	"fmt"

	"github.com/nlopes/slack"
)

type Slack struct {
	c     *Config
	api   *slack.Client
	wazuh *Wazuh
}

func NewSlack(c *Config) *Slack {
	return &Slack{
		c:     c,
		api:   slack.New(c.SlackToken),
		wazuh: NewWazuh(c),
	}
}

func (s *Slack) Notify(a *Alert) error {
	color := "danger"
	if a.Rule.Level <= 4 {
		color = "good"
	} else if a.Rule.Level >= 5 && a.Rule.Level <= 7 {
		color = "warning"
	}

	agent := slack.AttachmentField{
		Title: "Agent",
		Value: fmt.Sprintf("(%s) - %s", a.Agent.ID, a.Agent.Name),
	}

	location := slack.AttachmentField{
		Title: "Location",
		Value: a.Location,
	}

	rule := slack.AttachmentField{
		Title: "Rule ID",
		Value: fmt.Sprintf("%s (Level %d)", a.Rule.ID, a.Rule.Level),
	}

	agent_url := slack.AttachmentField{
		Title: "Agent URL",
		Value: fmt.Sprintf("%sapp/wazuh#/agents?agent=%s", s.c.KibanaURL, a.Agent.ID),
	}

	groups, err := s.wazuh.getGroups(a.Agent.ID)
	if err != nil {
		return err
	}
	for _, g := range groups {
		gd, ok := s.c.Groups[g]
		if !ok || gd.SlackChannel == "" {
			continue
		}

		text := a.Message()
		if gd.SlackMention != "" {
			mid, err := s.mentionID(gd.SlackMention)
			if err != nil {
				return err
			}
			text = "<!subteam^" + mid + "|" + "@" + gd.SlackMention + ">" + " " + a.Message()
		}

		attachment := slack.Attachment{
			Color:   color,
			Title:   a.Rule.Description,
			Pretext: "Wazuh Alert",
			Text:    text,
			Fields: []slack.AttachmentField{
				agent,
				location,
				rule,
				agent_url,
			},
		}

		_, _, err := s.api.PostMessage(gd.SlackChannel,
			slack.MsgOptionAttachments(attachment),
			slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{LinkNames: 1}),
		)

		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Slack) mentionID(name string) (string, error) {
	userGroups, err := s.api.GetUserGroups()
	if err != nil {
		return "", err
	}
	for _, u := range userGroups {
		if u.Name == name {
			return u.ID, nil
		}
	}

	users, err := s.api.GetUsers()
	if err != nil {
		return "", err
	}
	for _, u := range users {
		if u.Name == name {
			return u.ID, nil
		}
	}
	return "", errors.New("User and Group Notfound")
}
