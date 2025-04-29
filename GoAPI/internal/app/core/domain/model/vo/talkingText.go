package vo

import "errors"

type TalkingText struct {
	value string
}

func (t TalkingText) Value() string {
	return t.value
}

// コンストラクタ関数
func NewTalkingText(text string) (TalkingText, error) {
	if text == "" {
		return TalkingText{}, errors.New("talking text cannot be empty")
	}
	// 必要ならもっと制限を追加してもOK（長さ制限とか）
	return TalkingText{value: text}, nil
}
