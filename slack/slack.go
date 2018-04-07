package slack

import (
	"time"

	"github.com/nlopes/slack"
)

type slackClient struct {
	channel string
	client  *slack.Client
}

const (
	retryCount = 5
)

var client *slackClient

func Configure(token string, channel string) {
	client = &slackClient{
		channel: channel,
		client:  slack.New(token),
	}
}

func PostMessage(msg string) {
	for i := 0; i < retryCount; i++ {
		if _, _, err := client.client.PostMessage(client.channel, msg, slack.NewPostMessageParameters()); err != nil {
			time.Sleep(time.Second * 2)
			continue
		}
		break
	}
}
