package abstract

import(
	"GoAPI/internal/app/core/dto"
)
type AnalyzingVoiceService interface {
	AnalyzingVoiceService(voiceDataDTO dto.VoiceDataDTO) (string, error)
}