package main

import (
	"github.com/BurntSushi/toml"
	"github.com/nlopes/slack"
	"github.com/wasuken/slack-memo-bot/dbio"
	"os"
)

func loadFiles(filepaths []string) string {
	_, err := os.Stat(filepaths[0])
	if err != nil {
		return loadFiles(filepaths[1:])
	} else {
		return filepaths[0]
	}
}

var DEFAULT_LOAD_FILES []string = []string{
	"config.tml",
	"~/.config/slack-memo-bot/config.tml",
	"/etc/slack-memo-bot/config.tml"}

type Config struct {
	Slack SlackConfig
}

type SlackConfig struct {
	Apitoken     string
	WatchChannel string
}

func main() {
	var config Config
	_, err := toml.DecodeFile(loadFiles(DEFAULT_LOAD_FILES), &config)
	if err != nil {
		panic(err)
	}
	api := slack.New(
		config.Slack.Apitoken,
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if ev.Channel == config.Slack.WatchChannel {
				dbio.SaveMemo(ev.Text)
				rtm.SendMessage(rtm.NewOutgoingMessage("test", ev.Channel))
			}
		}
	}
}
