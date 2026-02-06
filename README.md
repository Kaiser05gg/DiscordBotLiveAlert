# DiscordBotLiveAlert

discordnot/
├── cmd/
│ └── api/
│ └── main.go # Cloud Run エントリーポイント
│
├── internal/
│ ├── config/
│ │ └── config.go # 環境変数読み取り
│ │
│ ├── handler/
│ │ └── check.go # /check エンドポイント
│ │
│ ├── notifier/
│ │ └── discord.go # Discord Webhook 送信
│ │
│ ├── platform/
│ │ ├── twitch.go # Twitch 判定
│ │ └── tiktok.go # TikTok 判定
│ │
│ └── state/
│ └── memory.go # 状態管理（最初はインメモリ）
│
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
