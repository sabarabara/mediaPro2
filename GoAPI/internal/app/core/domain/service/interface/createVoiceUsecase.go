package abstract

import(
	"GoAPI/internal/app/core/dto"
)

type CreateVoiceUsecase interface {
	CreateVoice(voiceDataDTO dto.VoiceDataDTO)(*dto.VoiceDataDTO , error)
}