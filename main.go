package main

import (
	"database/sql"
	"github.com/BurntSushi/toml"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nlopes/slack"
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
	"~/.config/slackbot-rss/config.tml",
	"/etc/slackbot/config.tml"}

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
				rtm.SendMessage(rtm.NewOutgoingMessage("test", ev.Channel))
			}
		}
	}
}

func writeDB(text string) (responseText string) {

}

func parseText(text string) {

}
