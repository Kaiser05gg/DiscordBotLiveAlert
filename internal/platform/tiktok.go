package platform

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type TikTokClient struct {
	client *http.Client
}

func NewTikTokClient() *TikTokClient {
	return &TikTokClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

/* ================================
   共通リクエスト（超重要）
================================ */

func (c *TikTokClient) doRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// TikTok が EOF を返さない最低限セット
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/120.0.0.0 Safari/537.36",
	)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://www.tiktok.com/")
	req.Header.Set("Connection", "keep-alive")

	return c.client.Do(req)
}

/* ================================
   STEP 1: user -> secUid
================================ */

func (c *TikTokClient) fetchSecUID(user string) (string, error) {
	url := "https://www.tiktok.com/api/user/detail/?uniqueId=" + user

	resp, err := c.doRequest(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		UserInfo struct {
			User struct {
				SecUID string `json:"secUid"`
			} `json:"user"`
		} `json:"userInfo"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if res.UserInfo.User.SecUID == "" {
		return "", errors.New("secUid not found")
	}

	return res.UserInfo.User.SecUID, nil
}

/* ================================
   STEP 2: secUid -> LIVE 判定
================================ */

func (c *TikTokClient) IsLive(user string) (bool, error) {
	secUid, err := c.fetchSecUID(user)
	if err != nil {
		return false, err
	}

	url := "https://www.tiktok.com/api/live/detail/?secUid=" + secUid

	resp, err := c.doRequest(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var res struct {
		RoomInfo *struct {
			Status int `json:"status"`
		} `json:"roomInfo"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}

	// status == 2 → LIVE
	return res.RoomInfo != nil && res.RoomInfo.Status == 2, nil
}
