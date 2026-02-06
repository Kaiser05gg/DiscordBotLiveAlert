package handler

import (
	"log"
	"net/http"
	"os"

	"discordbotlivealert/internal/config"
	"discordbotlivealert/internal/notifier"
	"discordbotlivealert/internal/platform"
	"discordbotlivealert/internal/state"
)

func Check(w http.ResponseWriter, r *http.Request) {
	log.Println("check start")
	cfg := config.Load()

	/* -------- Twitch -------- */

	if user := os.Getenv("TWITCH_USER_LOGIN"); user != "" {
		tw, _ := platform.NewTwitchClient()
		isLive, _ := tw.IsLive(user)
		wasLive := state.WasLive("twitch:" + user)

		log.Printf("twitch status user=%s wasLive=%v isLive=%v\n", user, wasLive, isLive)

		if !wasLive && isLive {
			_ = notifier.Send(
				cfg.DiscordWebhookURL,
				"🔴 Twitch配信開始！\nhttps://www.twitch.tv/"+user,
			)
		}
		state.SetLive("twitch:"+user, isLive)
	}

	/* -------- TikTok -------- */

	if user := os.Getenv("TIKTOK_USER"); user != "" {
		tt := platform.NewTikTokClient()
		isLive, err := tt.IsLive(user)
		if err != nil {
			log.Println("ERROR: tiktok api:", err)
		} else {
			wasLive := state.WasLive("tiktok:" + user)
			log.Printf("tiktok status user=%s wasLive=%v isLive=%v\n", user, wasLive, isLive)

			if !wasLive && isLive {
				_ = notifier.Send(
					cfg.DiscordWebhookURL,
					"🔴 TikTok配信開始！\nhttps://www.tiktok.com/@"+user,
				)
			}
			state.SetLive("tiktok:"+user, isLive)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("checked"))
}
