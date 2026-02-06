package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type message struct {
	Content string `json:"content"`
}

func Send(webhookURL, content string) error {
	body, _ := json.Marshal(message{
		Content: content,
	})

	req, err := http.NewRequest(
		http.MethodPost,
		webhookURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	return err
}
