package di

import (
	handlers "GoAPI/cmd/handler"
	"GoAPI/internal/app/controllers"
	"GoAPI/internal/app/core/domain/service/interface"
	"GoAPI/internal/app/frameworks"
	"GoAPI/internal/app/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var OKcreateVoiceUsecaseSet = wire.NewSet(
	usecases.NewCreateVoiceUsecaseImpl,
	wire.Bind(new(abstract.CreateVoiceUsecase), new(*usecases.CreateVoiceUsecaseImpl)),
)

var OKfactorySet = wire.NewSet(
	usecases.NewChattingInformationFactory,
)


func OKInitializeRouter() *gin.Engine {
	wire.Build(
		handlers.SetupRouter,
		controllers.NewCreateVoiceController,
		frameworks.NewGeminiRequester,
		frameworks.NewCreateVoiceService,
		frameworks.NewAnalyzingVoiceService,
		createVoiceUsecaseSet,
		factorySet, 
	)

	return handlers.SetupRouter(nil)
}
