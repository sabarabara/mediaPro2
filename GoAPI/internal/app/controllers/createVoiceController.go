package controllers

import (
	abstract "GoAPI/internal/app/core/domain/service/interface"
	"GoAPI/internal/app/core/dto"
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
    _, msg, err := conn.ReadMessage()
    if err != nil {
        log.Println("Error reading message:", err)
        return
    }

    voiceDataDTO := dto.VoiceDataDTO{
        AudioData: msg,
    }

    resAudioData, err := c.creatingVoiceUsecase.CreateVoice(voiceDataDTO)
    if err != nil {
        log.Println("Error processing voice data:", err)
        return
    }

    // ★ JSONじゃなくてバイナリ送信 ★
    err = conn.WriteMessage(websocket.BinaryMessage, resAudioData.AudioData)
    if err != nil {
        log.Println("Error sending response:", err)
        return
    }
}
}
