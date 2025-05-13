package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	abstract "GoAPI/internal/app/core/domain/service/interface"
	"GoAPI/internal/app/core/dto"
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
		return true
	},
}

func (c *CreateVoiceController) HandleWebSocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()
	log.Println("WebSocket connection established")

	// 3秒バッファ用
	var voiceBuffer []byte
	ticker := time.NewTicker(6 * time.Second)
	defer ticker.Stop()

	// メッセージ受信用チャネル
	messageChan := make(chan []byte)

	// 非同期で受信処理
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				close(messageChan)
				return
			}
			messageChan <- msg
		}
	}()

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return
			}
			voiceBuffer = append(voiceBuffer, msg...)

		case <-ticker.C:
			if len(voiceBuffer) == 0 {
				continue
			}

			log.Println("Sending buffered audio data of length:", len(voiceBuffer))

			voiceDataDTO := dto.VoiceDataDTO{
				AudioData: voiceBuffer,
			}

			resAudioData, err := c.creatingVoiceUsecase.CreateVoice(voiceDataDTO)
			if err != nil {
				log.Println("Error processing voice data:", err)
				return
			}

			err = conn.WriteMessage(websocket.BinaryMessage, resAudioData.AudioData)
			if err != nil {
				log.Println("Error sending response:", err)
				return
			}

			// バッファをリセット
			voiceBuffer = nil
		}
	}
}
