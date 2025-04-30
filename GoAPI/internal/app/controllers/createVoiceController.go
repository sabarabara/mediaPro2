package controllers

import (
	abstract "GoAPI/internal/app/core/domain/service/interface"
	"GoAPI/internal/app/core/dto"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type CreateVoiceController struct {
	creatingVoiceUsecase abstract.CreateVoiceUsecase
}

func NewCreateVoiceController(creatingVoiceUsecase abstract.CreateVoiceUsecase) *CreateVoiceController {
	return &CreateVoiceController{
		creatingVoiceUsecase: creatingVoiceUsecase,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 任意でOriginをチェック（必要に応じて変更）
		return true
	},
}

// WebSocket接続を処理するハンドラ
func (c *CreateVoiceController) HandleWebSocket(ctx *gin.Context) {
	// WebSocket接続のアップグレード
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()
	log.Println("WebSocket connection established")

	// 音声データの受信と処理
	for {
		// メッセージを受信（音声データ）
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		log.Println("Received message:")

		voiceDataDTO := dto.VoiceDataDTO{
			AudioData: msg, //
			//  受け取った音声データを設定
		}

		log.Println("Received voice data:")

		// 音声データをサービスで処理
		resAudioData, err := c.creatingVoiceUsecase.CreateVoice(voiceDataDTO)
		if err != nil {
			log.Println("Error processing voice data:", err)
			return
		}

		// 構造体をJSONに変換
		jsonData, err := json.Marshal(resAudioData)
		if err != nil {
			log.Println("Error marshalling response:", err)
			return
		}

		// クライアントへ送信（JSON文字列）
		err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("Error sending response:", err)
			return
		}
	}
}
