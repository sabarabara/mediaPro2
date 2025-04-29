package abstract

import(
	"GoAPI/internal/app/core/dto"
)

type CreateVoiceUsecase interface {
	CreatVoice(voiceDataDTO dto.VoiceDataDTO)(dto.VoiceDataDTO , error)
}