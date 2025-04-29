package frameworks

import (
    "GoAPI/internal/app/core/domain/model/vo"
    "GoAPI/internal/app/core/domain/service/interface"
    "GoAPI/internal/app/core/dto"
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
			panic("Error loading .env file")
	}
}


// FlaskサーバーのURL
var flaskServerURL2 = os.Getenv("Python_API_PORT_FOR_CREATE")


type VoiceCreater struct{}

// インターフェースを満たす
var _ abstract.CreateVoiceService = (*VoiceCreater)(nil)

func NewCreateVoiceService() abstract.CreateVoiceService {
    return &VoiceCreater{}
}

func (v *VoiceCreater) CreateVoice(chattingInfo vo.ChattingInformation) (dto.VoiceDataDTO, error) {
    // 例として、ChattingInformationを使ってリクエストを生成
    text := chattingInfo.TalkingText.Value()  // ここでTalkingTextを取得
    emotion := chattingInfo.ImotionalParam.Value()  // ここでImotionalParamを取得

    // createVoiceエンドポイントのURL
    url := flaskServerURL2

    // リクエストデータをJSONに変換
    data := map[string]string{
        "text":    text,
        "emotion": emotion,
    }
    jsonData, err := json.Marshal(data)
    if err != nil {
        return dto.VoiceDataDTO{}, err
    }

    // POSTリクエストの送信
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return dto.VoiceDataDTO{}, err
    }
    defer resp.Body.Close()

    // レスポンスコードチェック
    if resp.StatusCode != http.StatusOK {
        return dto.VoiceDataDTO{}, fmt.Errorf("failed to create voice, status code: %d", resp.StatusCode)
    }

    // レスポンスデータを読み取り
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return dto.VoiceDataDTO{}, err
    }

    // 音声データをVoiceDataDTOに格納して返す
    return dto.VoiceDataDTO{
        AudioData: body, // ここでレスポンスの音声データを返す
    }, nil
}
