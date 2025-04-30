package frameworks

import (
    "bytes"
    "mime/multipart"
    "net/http"
    "os"
    "io"
    "encoding/json"
    "errors"
    "github.com/joho/godotenv"

    "GoAPI/internal/app/core/dto"
    "GoAPI/internal/app/core/domain/model/vo"
    "GoAPI/internal/app/core/domain/service/interface"
)


var flaskServerURL string

func init() {
    if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }
    flaskServerURL = os.Getenv("Python_API_PORT_FOR_ANALIZE")
}




type VoiceAnalyzer struct{}

var _ abstract.AnalyzingVoiceService = (*VoiceAnalyzer)(nil)

func NewAnalyzingVoiceService() abstract.AnalyzingVoiceService {
    return &VoiceAnalyzer{}
}

func (v *VoiceAnalyzer) AnalyzeVoice(voiceDataDTO dto.VoiceDataDTO) (vo.ChattingInformation, error) {
    // Flaskサーバー叩いて、vo.ChattingInformationを作る処理
    // 1. ファイルをmultipart/form-dataで送る
    println("0-1")
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("file", "audio.wav")
    if err != nil {
        return vo.ChattingInformation{}, err
    }
    _, err = part.Write(voiceDataDTO.AudioData)
    if err != nil {
        return vo.ChattingInformation{}, err
    }
    writer.Close()
    println("0-2")
    req, err := http.NewRequest("POST", flaskServerURL, body)
    if err != nil {
        return vo.ChattingInformation{}, err
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return vo.ChattingInformation{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return vo.ChattingInformation{}, errors.New("Flask server returned non-200")
    }

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return vo.ChattingInformation{}, err
    }

    // 2. レスポンスをパース
    var result struct {
        Text    string `json:"text"`
        Emotion string `json:"emotion"`
    }
    if err := json.Unmarshal(respBody, &result); err != nil {
        return vo.ChattingInformation{}, err
    }

    // 3. voに変換
    talkingText, err := vo.NewTalkingText(result.Text)
    if err != nil {
        return vo.ChattingInformation{}, err
    }
    emotionalParam, err := vo.NewImotionalParam(result.Emotion)
    if err != nil {
        return vo.ChattingInformation{}, err
    }

    return vo.ChattingInformation{
        TalkingText:    talkingText,
        ImotionalParam: emotionalParam,
    }, nil
}
