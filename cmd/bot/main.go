package main

import (
	"log"
	"os"

	conv "github.com/pararthshah/srcbot/conversations"
	"github.com/pararthshah/srcbot/slack"
)

func main() {
	slackToken := os.Getenv("SLACK_ACCESS_TOKEN")
	if slackToken == "" {
		log.Fatalf("no SLACK_ACCESS_TOKEN found in env")
	}

	conn, err := slack.Connect(slackToken)
	if err != nil {
		log.Fatal(err)
	}

	slack.EventProcessor(conn, conv.OnAskedMessage, conv.OnHeardMessage)
}
