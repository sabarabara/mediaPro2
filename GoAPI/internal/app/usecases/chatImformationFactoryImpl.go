package usecases

import (
	"GoAPI/internal/app/core/domain/model/vo"
	"GoAPI/internal/app/core/domain/model/factory"
)

type ChattingInformationFactoryImpl struct{}

func NewChattingInformationFactory() factory.ChattingInformationFactory {
	return &ChattingInformationFactoryImpl{}
}


func (c *ChattingInformationFactoryImpl) CreateChattingInformation(talkingText vo.TalkingText, emotionalParam vo.ImotionalParam) (vo.ChattingInformation, error) {

	// ImotionalParam のバリデーション
	emoParam, err := vo.NewImotionalParam(emotionalParam.Value())
	if err != nil {
		return vo.ChattingInformation{}, err
	}

	// TalkingTextのバリデーション(仮)
	talkingTextVo, err := vo.NewTalkingText(talkingText.Value())
    if err != nil {
    return vo.ChattingInformation{}, err
   }

	// ChattingInformationの生成
	return vo.ChattingInformation{
		TalkingText:    talkingTextVo,
		ImotionalParam: emoParam,
	}, nil
}