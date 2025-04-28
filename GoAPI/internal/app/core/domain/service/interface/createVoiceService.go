package abstract

import(
	"GoAPI/internal/app/core/domain/model/vo"
)

type CreateVoiceService interface {
	CreateVoice(talkingText vo.TalkingText, emotionalParam vo.ImotionalParam) (vo.ChattingInformation, error)
}