package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"GoAPI/internal/app/di"
	"github.com/gin-contrib/cors"
)

func init() {
	// .env ファイルの読み込み
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 必要な環境変数がセットされているかチェック
	if os.Getenv("Gemini_API_URL") == "" || os.Getenv("Gemini_API_KEY") == "" {
		log.Fatal("Gemini API URL or API Key is not set")
	}
}

func main() {
	// wire で依存関係を解決してコントローラーを生成
	r := di.InitializeRouter()
	cors.Default()


	// サーバーを開始
	r.Run(":8080")
}
