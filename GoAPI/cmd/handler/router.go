package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"GoAPI/internal/app/controllers"
)

func SetupRouter(controller *controllers.CreateVoiceController) *gin.Engine {
	r := gin.Default()

	// テスト用ルート
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "API is up and running!")
	})

	// WebSocket用のルート
	r.GET("/ws", controller.HandleWebSocket)

	return r
}
