# ========================
# 仮想環境の構築手順（メモ）
# ========================
# python3 -m venv new_venv
# source new_venv/bin/activate
# pip install -r requirements.txt
# ========================



from flask import Flask, request, jsonify
import whisper
import numpy as np
from io import BytesIO
import wave
import io

app = Flask(__name__)

# Whisperモデルをロード
model = whisper.load_model("base")  # 必要に応じて tiny / small などに変えてOK

# 感情推定用の関数
def estimate_emotion(audio_data):
    # 音声データを読み取るために一時ファイルに保存
    audio_file = io.BytesIO(audio_data)
    
    # サンプルデータを抽出する処理（波形を取得して音量やピッチを分析）
    with wave.open(audio_file, 'rb') as f:
        samples = np.frombuffer(f.readframes(f.getnframes()), dtype=np.int16)
        volume = np.mean(np.abs(samples))  # 音量を簡易的に求める

        # ピッチの簡易推定
        zero_crossings = np.mean(np.abs(np.diff(np.sign(samples))))

        # 判定基準
        if volume > 1000 and zero_crossings > 0.1:
            return "happy"
        elif volume < 500 and zero_crossings < 0.05:
            return "sad"
        else:
            return "neutral"

@app.route('/analyzeVoice', methods=['POST'])
def analyze_voice():
    if 'file' not in request.files:
        return jsonify({"error": "No file part"}), 400

    file = request.files['file']
    audio_data = file.read()

    print(len(audio_data))
    # 感情推定
    emotion = estimate_emotion(audio_data)

    # Whisperで文字起こし
    # 一時ファイルを使用して音声データを保存
    audio_file = io.BytesIO(audio_data)
    with open("temp_audio.wav", "wb") as f:
        f.write(audio_data)

    # Whisperによる文字起こし
    result = model.transcribe("temp_audio.wav")
    text = result['text']

    print(f"Transcribed text: {text}")

    return jsonify({
        "text": text,
        "emotion": emotion
    })


if __name__ == '__main__':
    app.run(debug=True)

