from flask import Flask, request, jsonify, send_file
from gtts import gTTS
import subprocess
from io import BytesIO
import traceback

app = Flask(__name__)

def create_voice(talking_text: str, emotional_param: str) -> BytesIO:
    # 1) gTTS で MP3 を生成
    mp3_buf = BytesIO()
    gTTS(text=talking_text, lang='ja').write_to_fp(mp3_buf)
    mp3_buf.seek(0)

    # 2) FFmpeg CLI で MP3 → WAV に変換
    #    stdin: mp3_buf、 stdout: wav_buf
    wav_buf = BytesIO()
    ffmpeg_cmd = [
        'ffmpeg', '-i', 'pipe:0',    # 標準入力から読む
        '-f', 'wav',                 # フォーマット WAV
        'pipe:1'                     # 標準出力へ書く
    ]
    proc = subprocess.run(
        ffmpeg_cmd,
        input=mp3_buf.read(),
        stdout=subprocess.PIPE,
        stderr=subprocess.DEVNULL  # 進捗ログは捨てる
    )
    wav_buf.write(proc.stdout)
    wav_buf.seek(0)

    return wav_buf

@app.route('/createVoice', methods=['POST'])
def handle_create_voice():
    try:
        data = request.get_json()
        if not data or 'talkingText' not in data or 'emotionalParam' not in data:
            return jsonify({"error": "Invalid input"}), 400

        wav_output = create_voice(
            talking_text=data['talkingText'],
            emotional_param=data['emotionalParam']
        )
        return send_file(
            wav_output,
            mimetype='audio/wav',
            as_attachment=True,
            download_name="output.wav"
        )
    except Exception as e:
        traceback.print_exc()
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True, port=5001)
