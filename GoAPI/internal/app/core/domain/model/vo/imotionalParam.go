package vo

import (
	"errors"
	"strings"
)

type ImotionalParam struct {
	value string
}

func (e ImotionalParam) Value() string {
	return e.value
}

// 正しいコンストラクタ
func NewImotionalParam(value string) (ImotionalParam, error) {
	allowedEmotions := []string{
		"joy", "sadness", "anger", "fear", "surprise", "disgust", "neutral","happy",
		"angry", "sad", "fearful", "surprised", "disgusted", "neutral",
	}

	// 感情が許可されたリストに含まれているかチェック
	for _, emotion := range allowedEmotions {
		if strings.ToLower(value) == emotion {
			return ImotionalParam{value: value}, nil
		}
	}

	return ImotionalParam{}, errors.New("invalid emotion: " + value)
}
