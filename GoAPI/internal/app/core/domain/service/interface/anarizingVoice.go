package abstract

import (
    "GoAPI/internal/app/core/dto"
    "GoAPI/internal/app/core/domain/model/vo"
)

type AnalyzingVoiceService interface {
    AnalyzeVoice(voiceDataDTO dto.VoiceDataDTO) (vo.ChattingInformation, error)
}
