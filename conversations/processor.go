package conversations

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/jsgoecke/go-wit"
	"github.com/pararthshah/srcbot/slack"
)

var (
	HelpResponse  = "To create a thread from your conversations, say `@srcbot: create thread {num_messages} {thread_title}`\nfor example, `@srcbot: create thread 2 Discussion about routers`"
	WhoResponse   = "I'm srcbot. Think of me as Siri for your source code. Except, I'm only as smart as a puppy."
	HelloResponse = "Hey! I'm doing good. What's the latest news about The Graph?"
	HALResponse   = "I'm sorry, Dave. I'm afraid I can't do that."
)

// history stores the last messageLimit messages for each channel
var history SlackHistory

func init() {
	history = make(map[string]*MessageHistory)
}

func OnHeardMessage(message *slack.Message) {
	channelHistory := history.GetChannel(message.Channel)
	channelHistory.Append(&Message{Author: message.From, Body: message.Text})
}

var createThreadCmd = regexp.MustCompile("(?i)create (?:(?i)ren)?(?i)thread ([0-9]+) (.*)")

func OnAskedMessage(message *slack.Message) {
	var response string

	log.Printf("%s-> %s", message.From, message.Text)

	if message.Text == "help" {
		response = HelpResponse
	}
	if response == "" {
		matches := createThreadCmd.FindStringSubmatch(message.Text)
		if len(matches) > 0 {
			response = createThread(message, matches[2], matches[1])
		}
	}

	if response == "" {
		witMsg, err := witParseMessage(message.Text)
		if err != nil {
			response = fmt.Sprintf("error parsing message: %v", err)
		} else {
			response = generateWitResponse(message, witMsg)
		}
	}

	if response == "" {
		response = "woof"
	}

	log.Printf("me-> %s", response)

	if err := message.Respond(response); err != nil {
		// gulp!
		log.Printf("Error responding to message: %s\nwith Message: '%s'", err, response)
	}
}

func createThread(message *slack.Message, title, numStr string) string {
	var response string
	num, err := strconv.Atoi(numStr)
	if err == nil {
		channelHistory := history.GetChannel(message.Channel)
		comments := channelHistory.GetLastN(num)

		info, err := postCreateThread(title, message.From, comments)
		if err != nil {
			response = fmt.Sprintf("error creating thread: %v", err)
		} else if info.Err != "" {
			response = fmt.Sprintf("error creating thread: %v", info.Err)
		} else {
			response = fmt.Sprintf("thread created (%s)", info.URL)
		}
	}
	return response
}

func generateWitResponse(message *slack.Message, witMsg *wit.Message) string {
	if len(witMsg.Outcomes) == 0 {
		return ""
	}

	log.Printf("num outcomes: %d", len(witMsg.Outcomes))

	for i, outcome := range witMsg.Outcomes {
		log.Printf("wit.ai intent %d: %s", i, outcome.Intent)
		switch outcome.Intent {
		case "create_thread":
			var title, numStr string
			if v, ok := outcome.Entities["thread_title"]; ok && len(v) > 0 {
				title = (*v[0].Value).(string)
			}
			if v, ok := outcome.Entities["thread_msg_count"]; ok && len(v) > 0 {
				numStr = (*v[0].Value).(string)
			}
			if title != "" && numStr != "" {
				log.Printf("title: %s, num: %s", title, numStr)
				return createThread(message, title, numStr)
			}
		case "how_are_you":
			return HelloResponse
		case "who_are_you":
			return WhoResponse
		case "help_me":
			return HelpResponse
		case "open_pod_bay_doors":
			return HALResponse
		}
	}

	log.Printf("no valid intent")

	return ""
}
