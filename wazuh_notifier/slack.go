package wazuh_notifier

import (
	"errors"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/nlopes/slack"
	gocache "github.com/patrickmn/go-cache"
)

type Slack struct {
	c     *Config
	api   *slack.Client
	wazuh *Wazuh
	cache *gocache.Cache
}

func NewSlack(c *Config) *Slack {
	cache := gocache.New(time.Duration(c.IgnoreRepeatedMin)*time.Minute, 5*time.Minute)
	cache.LoadFile(c.IgnoreHistoryFile)
	return &Slack{
		c:     c,
		api:   slack.New(c.SlackToken),
		wazuh: NewWazuh(c),
		cache: cache,
	}
}

func (s *Slack) Notify(a *Alert) error {
	defer s.cache.DeleteExpired()
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
		log.Errorf("get group error agent_id %s : %s", a.Agent.ID, err.Error())
		return nil
	}
	for _, g := range groups {
		_, found := s.cache.Get(g + a.Rule.ID)
		if found {
			log.Infof("skip notify group %s, ruleid %s\n", g, a.Rule.ID)
			continue
		}

		gd, ok := s.c.Groups[g]
		if !ok || gd.SlackChannel == "" {
			continue
		}

		text := a.Message()
		if gd.SlackMention != "" {
			mid, err := s.mentionID(gd.SlackMention)
			if err != nil {
				return fmt.Errorf("get slack mention error mention %s : %s", gd.SlackMention, err.Error())
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
			return fmt.Errorf("send slack error %s channel %s", err, gd.SlackChannel)
		}
		log.Infof("notify slack to %s", gd.SlackChannel)
		s.cache.Set(g+a.Rule.ID, a.Rule.ID, cache.DefaultExpiration)
	}
	s.cache.SaveFile(s.c.IgnoreHistoryFile)
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
