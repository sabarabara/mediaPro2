package frameworks

import (
    "bytes"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "os"

    "GoAPI/internal/app/core/domain/model/vo"
		"GoAPI/internal/app/core/domain/service/interface"

		"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
			panic("Error loading .env file")
	}

	if os.Getenv("Gemini_API_URL") == "" || os.Getenv("GEMINI_API_KEY") == "" {
			panic("Gemini API URL or API Key is not set")
	}
}


var geminiURL = os.Getenv("Gemini_API_URL")
var geminiAPIKey = os.Getenv("GEMINI_API_KEY")


type GeminiRequester struct{}

type RequestPayload struct {
    Text    string `json:"text"`
    Emotion string `json:"emotion"`
}

type GeminiRequest struct {
    Contents []Content `json:"contents"`
}

type Content struct {
    Parts []Part `json:"parts"`
}

type Part struct {
    Text string `json:"text"`
}

type GeminiResponse struct {
    Candidates []struct {
        Content struct {
            Parts []struct {
                Text string `json:"text"`
            } `json:"parts"`
        } `json:"content"`
    } `json:"candidates"`
}

// GeminiRequesterがCreateChattingInformationインターフェースを満たす
var _ abstract.CreateChattingInformation = (*GeminiRequester)(nil)

func NewGeminiRequester() abstract.CreateChattingInformation {
    return &GeminiRequester{}
}

func (g *GeminiRequester) CreateChattingInformation(talkingText vo.TalkingText, emotionalParam vo.ImotionalParam) (vo.ChattingInformation, error) {
    prompt := "次のテキストと感情に基づいて、彼女のような話し方で、似た雰囲気のテキストを生成してください。\n\n" +
        "テキスト: " + talkingText.Value() + "\n" +
        "感情: " + emotionalParam.Value()

    geminiReq := GeminiRequest{
        Contents: []Content{
            {Parts: []Part{{Text: prompt}}},
        },
    }

    body, err := json.Marshal(geminiReq)
    if err != nil {
        return vo.ChattingInformation{}, err
    }

    req, err := http.NewRequest("POST", geminiURL+"?key="+geminiAPIKey, bytes.NewBuffer(body))
    if err != nil {
        return vo.ChattingInformation{}, err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return vo.ChattingInformation{}, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return vo.ChattingInformation{}, err
    }

    var geminiResp GeminiResponse
    if err := json.Unmarshal(respBody, &geminiResp); err != nil {
        return vo.ChattingInformation{}, err
    }

    if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
        return vo.ChattingInformation{}, errors.New("no response from Gemini API")
    }

    responseText := geminiResp.Candidates[0].Content.Parts[0].Text

		talkingTextVo, err := vo.NewTalkingText(responseText)
    if err != nil {
    return vo.ChattingInformation{}, err
   }

		return vo.ChattingInformation{
			TalkingText:    talkingTextVo,
			ImotionalParam: emotionalParam,
		}, nil
}