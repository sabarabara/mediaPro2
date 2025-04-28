package factory

import (
	"GoAPI/internal/app/core/domain/model/vo"
)


type ChattingInformationFactory interface {
	CreateChattingInformation(talkingText vo.TalkingText, emotionalParam vo.ImotionalParam) (vo.ChattingInformation,error)
}
