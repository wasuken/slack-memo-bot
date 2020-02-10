package main

import (
	"database/sql"
	"github.com/BurntSushi/toml"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nlopes/slack"
	"github.com/wasuken/slack-memo-bot/dbio"
	"github.com/wasuken/slack-memo-bot/util"
	"log"
	"strings"
)

var DEFAULT_LOAD_PATH_LIST []string = []string{
	"./",
	"~/.config/slack-memo-bot/",
	"/etc/slack-memo-bot/"}

type Config struct {
	Slack SlackConfig
}

type SlackConfig struct {
	Apitoken     string
	WatchChannel string
}

func main() {
	var config Config
	_, err := toml.DecodeFile(util.LoadFiles(DEFAULT_LOAD_PATH_LIST, "config.tml"), &config)
	if err != nil {
		panic(err)
	}
	api := slack.New(
		config.Slack.Apitoken,
	)
	db, err := sql.Open("sqlite3", util.LoadFiles(DEFAULT_LOAD_PATH_LIST, "db.sqlite"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if ev.Channel == config.Slack.WatchChannel {
				text, tagList := util.ParseText(ev.Text)
				if strings.HasPrefix(strings.Trim(text, " 　"), "!!output!!") {
					message := dbio.OutputMemo(db, tagList[0], tagList[1:], ev.User)
					rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))
				} else if strings.HasPrefix(strings.Trim(text, " 　"), "!!delete!!") {
					dbio.DeleteMemoTags(db, tagList)
					message := dbio.OutputMemo(db, "markdown", tagList, ev.User)
					rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))
				} else {
					dbio.SaveMemo(db, text, tagList, ev.User)
				}
			}
		}
	}
}
