package frameworks

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/joho/godotenv"

	"GoAPI/internal/app/core/domain/model/vo"
	abstract "GoAPI/internal/app/core/domain/service/interface"
	"GoAPI/internal/app/core/dto"
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
	// PCMバイト列を int16 スライスに変換

	fmt.Println("PCM Dump (first 10 samples):")
	for i := 0; i < 10 && i < len(pcm)/2; i++ {
		u := binary.LittleEndian.Uint16(pcm[2*i : 2*i+2])
		fmt.Printf("Raw bytes: %x %x → Int16: %d\n", pcm[2*i], pcm[2*i+1], int16(u))
	}

	n := len(pcm) / 2
	samples16 := make([]int16, n)
	for i := 0; i < n; i++ {
		samples16[i] = int16(binary.LittleEndian.Uint16(pcm[2*i : 2*i+2]))
	}

	// int16 → int に変換（go-audioが int 期待する）
	samples := make([]int, n)
	for i := 0; i < n; i++ {
		samples[i] = int(samples16[i])
	}

	// WriteSeekerBuffer に書き込む
	wsb := NewWriteSeekerBuffer()
	enc := wav.NewEncoder(
		wsb,
		sampleRate,
		bitDepth,
		numChannels,
		1, // PCMフォーマット
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

	// デバッグ用に保存
	err := os.WriteFile("debug_output.wav", wsb.buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("保存失敗:", err)
	} else {
		fmt.Println("保存完了: debug_output.wav")
	}

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
