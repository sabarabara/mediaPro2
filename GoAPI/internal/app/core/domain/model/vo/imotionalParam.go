package vo

import (
	"errors"
	"strings"
)

type ImotionalParam string

// バリデーション関数と感情の返却
func (i *ImotionalParam) NewImotionalParam() (ImotionalParam, error) {
	// Google Colabで帰ってくる感情の種類（例）
	allowedEmotions := []string{
		"joy", "sadness", "anger", "fear", "surprise", "disgust", "neutral",
	}

	// 感情が許可されたリストに含まれているかチェック
	for _, emotion := range allowedEmotions {
		if strings.ToLower(string(*i)) == emotion {
			return *i, nil // 有効な感情を返す
		}
	}

	// 無効な感情の場合はエラーを返す
	return "", errors.New("invalid emotion: " + string(*i))
}
