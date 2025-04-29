from flask import Flask, request, jsonify, send_file
import pyttsx3
import audioread
import numpy as np
from io import BytesIO
import tempfile
import os

app = Flask(__name__)

def create_voice(talking_text: str, emotional_param: str) -> BytesIO:
    # TTSエンジン初期化
    engine = pyttsx3.init()

    # 話速設定（emotional_paramに応じた速度）
    if emotional_param == "happy":
        engine.setProperty('rate', 150)
    elif emotional_param == "sad":
        engine.setProperty('rate', 100)
    else:
        engine.setProperty('rate', 125)

    # 一時ファイルに出力（pyttsx3はファイル保存しかできない）
    with tempfile.NamedTemporaryFile(delete=False, suffix=".wav") as tmpfile:
        tmp_path = tmpfile.name
    
    try:
        engine.save_to_file(talking_text, tmp_path)
        engine.runAndWait()

        # 音声データをメモリに読み込む
        audio_output = BytesIO()
        with open(tmp_path, 'rb') as f:
            audio_output.write(f.read())
        audio_output.seek(0)
        return audio_output
    finally:
        if os.path.exists(tmp_path):
            os.remove(tmp_path)

@app.route('/createVoice', methods=['POST'])
def handle_create_voice():
    try:
        data = request.get_json()
        if not data or 'talkingText' not in data or 'emotionalParam' not in data:
            return jsonify({"error": "Invalid input: talkingText and emotionalParam are required"}), 400

        talking_text = data['talkingText']
        emotional_param = data['emotionalParam']

        audio_output = create_voice(talking_text, emotional_param)

        return send_file(audio_output, mimetype='audio/wav', as_attachment=True, download_name="output.wav")
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True, port=5001)
