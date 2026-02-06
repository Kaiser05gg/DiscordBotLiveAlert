package platform

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"
)

type twitchTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type twitchStreamsResponse struct {
	Data []struct {
		ID        string `json:"id"`
		UserLogin string `json:"user_login"`
		Title     string `json:"title"`
	} `json:"data"`
}

type TwitchClient struct {
	clientID     string
	clientSecret string

	accessToken string
	expiresAt   time.Time
}

func NewTwitchClient() (*TwitchClient, error) {
	id := os.Getenv("TWITCH_CLIENT_ID")
	sec := os.Getenv("TWITCH_CLIENT_SECRET")
	if id == "" || sec == "" {
		return nil, errors.New("TWITCH_CLIENT_ID or TWITCH_CLIENT_SECRET is empty")
	}
	return &TwitchClient{
		clientID:     id,
		clientSecret: sec,
	}, nil
}

func (t *TwitchClient) ensureToken() error {
	// 期限が残ってたら再利用（30秒前に更新）
	if t.accessToken != "" && time.Now().Before(t.expiresAt.Add(-30*time.Second)) {
		return nil
	}

	endpoint := "https://id.twitch.tv/oauth2/token"
	values := url.Values{}
	values.Set("client_id", t.clientID)
	values.Set("client_secret", t.clientSecret)
	values.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, endpoint+"?"+values.Encode(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("twitch token request failed: " + resp.Status)
	}

	var tr twitchTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return err
	}

	t.accessToken = tr.AccessToken
	t.expiresAt = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return nil
}

// IsLive は user_login（例: "shroud"）を渡すと配信中かどうかを返す
func (t *TwitchClient) IsLive(userLogin string) (bool, error) {
	if err := t.ensureToken(); err != nil {
		return false, err
	}

	u := "https://api.twitch.tv/helix/streams?user_login=" + url.QueryEscape(userLogin)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Client-Id", t.clientID)
	req.Header.Set("Authorization", "Bearer "+t.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, errors.New("twitch streams request failed: " + resp.Status)
	}

	var sr twitchStreamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return false, err
	}

	return len(sr.Data) > 0, nil
}
