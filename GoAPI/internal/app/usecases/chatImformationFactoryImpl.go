package usecases

import (
	"GoAPI/internal/app/core/domain/model/vo"
	"GoAPI/internal/app/core/domain/model/factory"
)

type ChattingInformationFactoryImpl struct{}

func NewChattingInformationFactory() factory.ChattingInformationFactory {
	return &ChattingInformationFactoryImpl{}
}

// CreateChattingInformation メソッドを実装
func (c *ChattingInformationFactoryImpl) CreateChattingInformation(talkingText vo.TalkingText, emotionalParam vo.ImotionalParam) (vo.ChattingInformation, error) {
	// ImotionalParam のバリデーション
	emoParam, err := emotionalParam.NewImotionalParam()
	if err != nil {
		return vo.ChattingInformation{}, err
	}

	// TalkingTextのバリデーション(仮)
	//talkingText = talkingText;


	
	// ChattingInformationの生成
	return vo.ChattingInformation{
		TalkingText:    talkingText,
		ImotionalParam: emoParam,
	}, nil
}
