package usecases

import (
	abstract "GoAPI/internal/app/core/domain/service/interface"
	"GoAPI/internal/app/core/domain/model/factory"
	"GoAPI/internal/app/core/dto"
)

type CreateVoiceUsecaseImpl struct {
	analyzeVoiceService abstract.AnalyzingVoiceService
	createChattingInformationService abstract.CreateChattingInformation
	createVoiceService abstract.CreateVoiceService
	chattingInformationFactory factory.ChattingInformationFactory
}

func NewCreateVoiceUsecaseImpl(
	analyzeVoiceService abstract.AnalyzingVoiceService,
	createChattingInformationService abstract.CreateChattingInformation,
	createVoiceService abstract.CreateVoiceService,
	chattingInformationFactory factory.ChattingInformationFactory,
) *CreateVoiceUsecaseImpl {
	return &CreateVoiceUsecaseImpl{
		analyzeVoiceService: analyzeVoiceService,
		createChattingInformationService: createChattingInformationService,
		createVoiceService: createVoiceService,
		chattingInformationFactory: chattingInformationFactory,
	}
}

func (c *CreateVoiceUsecaseImpl) CreateVoice(voiceDataDTO dto.VoiceDataDTO) (*dto.VoiceDataDTO, error) {

	// 1. Analyze the voice data
	analyzedData, err := c.analyzeVoiceService.AnalyzeVoice(voiceDataDTO)
	if err != nil {
		return nil, err
	}

	// 2. Use Factory to create ChattingInformation domain object
	chattingInformation, err := c.chattingInformationFactory.CreateChattingInformation(
		analyzedData.TalkingText,
		analyzedData.ImotionalParam,
	)
	if err != nil {
		return nil, err
	}

	// 3. Save ChattingInformation
	responseChattingInformation, err := c.createChattingInformationService.CreateChattingInformation(chattingInformation.TalkingText, chattingInformation.ImotionalParam)
	if err != nil {
		return nil, err
	}

	// 4. Create the voice (WAV) from ChattingInformation
	audioData, err := c.createVoiceService.CreateVoice(responseChattingInformation)
	if err != nil {
		return nil, err
	}

	return &audioData, nil
}
