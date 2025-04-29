package controllers

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"

	"GoAPI/internal/app/core/domain/service/interface"
	"GoAPI/internal/app/core/dto"
)

type CreateVoiceController struct {
	creatingVoiceService abstract.CreateVoiceUsecase
}

func NewCreateVoiceController(creatingVoiceService abstract.CreateVoiceUsecase) *CreateVoiceController {
	return &CreateVoiceController{
		creatingVoiceService: creatingVoiceService,
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

		voiceDataDTO := dto.VoiceDataDTO{
			AudioData: msg, // 受け取った音声データを設定
		}

		// 音声データをサービスで処理
		err = c.creatingVoiceService.CreatVoice(voiceDataDTO)
		if err != nil {
			log.Println("Error processing voice data:", err)
			return
		}

		// 成功した場合のレスポンス
		err = conn.WriteMessage(websocket.TextMessage, []byte("Voice processed successfully"))
		if err != nil {
			log.Println("Error sending response:", err)
			return
		}
	}
}

