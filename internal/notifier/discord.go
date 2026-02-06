package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type payload struct {
	Content string `json:"content"`
}

func Send(webhookURL, content string) error {
	body, err := json.Marshal(payload{
		Content: content,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		webhookURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
