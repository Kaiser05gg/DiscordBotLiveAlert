package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"discordbotlivealert/internal/handler"
)

func main() {
	// ローカル開発用（.env がなくてもエラーにしない）
	_ = godotenv.Load()

	mux := http.NewServeMux()

	// Cloud Scheduler から叩くエンドポイント
	mux.HandleFunc("/check", handler.Check)

	addr := ":8080" // Cloud Run は 8080 固定
	log.Println("server starting on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
