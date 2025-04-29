package di

import (
	handlers "GoAPI/cmd/handler"
	"GoAPI/internal/app/controllers"
	"GoAPI/internal/app/frameworks"
	"GoAPI/internal/app/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)


func InitializeRouter() *gin.Engine {

	wire.Build(

		handlers.SetupRouter,
		controllers.NewCreateVoiceController,
		usecases.NewCreateVoiceUsecaseImpl,
		frameworks.NewGeminiRequester,
		frameworks.NewCreateVoiceService,
		frameworks.NewAnalyzingVoiceService,

	)

	return handlers.SetupRouter(nil)
}