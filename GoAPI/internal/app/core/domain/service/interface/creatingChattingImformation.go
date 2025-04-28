package abstract

import(
	"GoAPI/internal/app/core/domain/model/vo"
)

type CreateChattingInformation interface {
	CreateChattingInformation(talkingText vo.TalkingText, emotionalParam vo.ImotionalParam) (vo.ChattingInformation, error)
}