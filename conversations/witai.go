package conversations

import (
	"log"
	"os"

	"github.com/jsgoecke/go-wit"
)

var client *wit.Client

func init() {
	token := os.Getenv("WIT_ACCESS_TOKEN")
	if token == "" {
		log.Fatalf("no WIT_ACCESS_TOKEN found in env")
	}
	client = wit.NewClient(token)
}

// witParseMessage processes a text message
func witParseMessage(msg string) (*wit.Message, error) {
	request := &wit.MessageRequest{Query: msg, N: 3}
	result, err := client.Message(request)
	if err != nil {
		return nil, err
	}
	return result, nil
}
