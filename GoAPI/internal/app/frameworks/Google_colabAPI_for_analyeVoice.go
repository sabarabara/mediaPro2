package frameworks

import (
    "bytes"
    "mime/multipart"
    "net/http"
    "os"
    "io"
    "encoding/json"
    "errors"
    "fmt"
    "encoding/binary"
    "github.com/go-audio/audio"
    "github.com/go-audio/wav"
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

type WriteSeekerBuffer struct {
    buf *bytes.Buffer
    pos int64
}

// Ensure we implement io.WriteSeeker
var _ io.WriteSeeker = (*WriteSeekerBuffer)(nil)

// NewWriteSeekerBuffer は空の WriteSeekerBuffer を返します
func NewWriteSeekerBuffer() *WriteSeekerBuffer {
    return &WriteSeekerBuffer{buf: &bytes.Buffer{}}
}

// Write は現在のシーク位置にバイト列を書き込みます
func (w *WriteSeekerBuffer) Write(p []byte) (int, error) {
    // バッファ末尾にいる場合は単純に追加
    if w.pos == int64(w.buf.Len()) {
        n, err := w.buf.Write(p)
        w.pos += int64(n)
        return n, err
    }
    // それ以外は overwrite or extend
    data := w.buf.Bytes()
    end := w.pos + int64(len(p))
    if end > int64(len(data)) {
        // 拡張が必要
        newData := make([]byte, end)
        copy(newData, data)
        copy(newData[w.pos:], p)
        w.buf = bytes.NewBuffer(newData)
    } else {
        // 既存領域を上書き
        copy(data[w.pos:end], p)
        w.buf = bytes.NewBuffer(data)
    }
    w.pos = end
    return len(p), nil
}

// Seek は現在のシーク位置を更新します
func (w *WriteSeekerBuffer) Seek(offset int64, whence int) (int64, error) {
    var newPos int64
    switch whence {
    case io.SeekStart:
        newPos = offset
    case io.SeekCurrent:
        newPos = w.pos + offset
    case io.SeekEnd:
        newPos = int64(w.buf.Len()) + offset
    default:
        return 0, fmt.Errorf("invalid whence: %d", whence)
    }
    if newPos < 0 {
        return 0, errors.New("negative position")
    }
    w.pos = newPos
    return newPos, nil
}




func pcmToWav(pcm []byte, sampleRate, bitDepth, numChannels int) ([]byte, error) {
    // PCMバイト列を int スライスに変換
    n := len(pcm) / 2
    samples := make([]int, n)
    for i := 0; i < n; i++ {
        u := binary.LittleEndian.Uint16(pcm[2*i : 2*i+2])
        samples[i] = int(int16(u))
    }

    // ← ここで bytes.Buffer ではなく WriteSeekerBuffer を使う
    wsb := NewWriteSeekerBuffer()
    enc := wav.NewEncoder(
        wsb,          // ← *WriteSeekerBuffer を渡す
        sampleRate,   // 例: 16000
        bitDepth,     // 例: 16
        numChannels,  // 例: 1
        1,            // PCM フォーマット
    )

    audioBuf := &audio.IntBuffer{
        Data:           samples,
        Format:         &audio.Format{SampleRate: sampleRate, NumChannels: numChannels},
        SourceBitDepth: bitDepth,
    }
    if err := enc.Write(audioBuf); err != nil {
        return nil, err
    }
    if err := enc.Close(); err != nil {
        return nil, err
    }

    // WAV ヘッダー付きバイト列を返す
    return wsb.buf.Bytes(), nil
}

// VoiceAnalyzer implements abstract.AnalyzingVoiceService
var _ abstract.AnalyzingVoiceService = (*VoiceAnalyzer)(nil)

type VoiceAnalyzer struct{}

func NewAnalyzingVoiceService() abstract.AnalyzingVoiceService {
	return &VoiceAnalyzer{}
}

// AnalyzeVoice sends a WAV-wrapped PCM payload to Flask and parses the response
func (v *VoiceAnalyzer) AnalyzeVoice(voiceDataDTO dto.VoiceDataDTO) (vo.ChattingInformation, error) {
	// 1. Wrap raw PCM into WAV
	wavBytes, err := pcmToWav(voiceDataDTO.AudioData, 16000, 16, 1)
	if err != nil {
		return vo.ChattingInformation{}, err
	}

	// 2. Prepare multipart/form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "audio.wav")
	if err != nil {
		return vo.ChattingInformation{}, err
	}
	if _, err := part.Write(wavBytes); err != nil {
		return vo.ChattingInformation{}, err
	}
	writer.Close()

	// 3. Send HTTP request to Flask
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

	// 4. Read and parse response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return vo.ChattingInformation{}, err
	}
	var result struct {
		Text    string `json:"text"`
		Emotion string `json:"emotion"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return vo.ChattingInformation{}, err
	}

    println("Received from Flask:", result.Text, result.Emotion)

	// 5. Convert to vo.ChattingInformation
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
