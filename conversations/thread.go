package conversations

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var srcToken, srcRepo, srcUrl string

func init() {
	srcToken = os.Getenv("SRC_ACCESS_TOKEN")
	if srcToken == "" {
		log.Fatalf("no SRC_ACCESS_TOKEN found in env")
	}

	srcUrl = os.Getenv("SRC_THREAD_URL")
	if srcUrl == "" {
		log.Fatalf("no SRC_THREAD_URL found in env")
	}

	srcRepo = os.Getenv("SRC_REPO_NAME")
	if srcRepo == "" {
		log.Fatalf("no SRC_REPO_NAME found in env")
	}
}

type NewThread struct {
	Token    string
	Author   string
	Title    string
	Comments []*Message
}

type ThreadInfo struct {
	ID  int
	URL string

	Status int
	Err    string
}

func postCreateThread(title, author string, comments []*Message) (*ThreadInfo, error) {
	data, err := json.Marshal(NewThread{
		Token:    srcToken,
		Author:   author,
		Title:    title,
		Comments: comments,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", srcUrl, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var newThread ThreadInfo
	if err := json.NewDecoder(resp.Body).Decode(&newThread); err != nil {
		return nil, err
	}

	return &newThread, nil
}
