package abstract

import (
	"GoAPI/internal/app/core/domain/model/vo"
	"GoAPI/internal/app/core/dto"
)

type CreateVoiceService interface {
	CreateVoice(vo.ChattingInformation) (dto.VoiceDataDTO, error)
}