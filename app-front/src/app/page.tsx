"use client";

import React, { useState, useRef } from 'react';

const Home = () => {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [isRecording, setIsRecording] = useState(false);
  const [audioChunks, setAudioChunks] = useState<Blob[]>([]);
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);

  // WebSocket接続を確立
  const connectWebSocket = () => {
    const socket = new WebSocket('//localhost:8080/ws'); // WebSocketのURLを調整
    socket.binaryType = "arraybuffer"; // 重要！！ バイナリデータを受信できるようにする

    socket.onopen = () => {
      console.log('WebSocket connection established');
      setWs(socket);
    };

    socket.onmessage = (event) => {
      console.log('音声データ受信');

      // 受け取ったArrayBufferをBlobに変換
      const audioBlob = new Blob([event.data], { type: 'audio/wav' });

      // BlobをURLにしてAudioタグで再生
      const audioUrl = URL.createObjectURL(audioBlob);
      const audio = new Audio(audioUrl);
      audio.play();
    };

    socket.onerror = (error) => {
      console.log('WebSocket error:', error);
    };

    socket.onclose = () => {
      console.log('WebSocket connection closed');
    };
  };

  // 録音開始
  const startRecording = () => {
    if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
      navigator.mediaDevices.getUserMedia({ audio: true })
        .then((stream) => {
          mediaRecorderRef.current = new MediaRecorder(stream);
          mediaRecorderRef.current.ondataavailable = (event) => {
            setAudioChunks((prevChunks) => [...prevChunks, event.data]);
          };
          mediaRecorderRef.current.start();
          setIsRecording(true);
        })
        .catch((err) => {
          console.error('録音開始エラー:', err);
        });
    } else {
      console.error('音声録音機能はこのブラウザではサポートされていません。');
    }
  };

  // 録音停止
  const stopRecording = () => {
    if (mediaRecorderRef.current && isRecording) {
      mediaRecorderRef.current.stop();
      setIsRecording(false);

      // 録音データをBlobとして統合
      const audioBlob = new Blob(audioChunks, { type: 'audio/wav' });

      // 録音した音声をWebSocketで送信
      if (ws) {
        audioBlob.arrayBuffer().then(buffer => {
          ws.send(buffer); // ArrayBufferで送信
          console.log('音声ファイルを送信しました');
        });
      }

      // 録音データをリセット
      setAudioChunks([]);
    }
  };

  return (
    <div>
      <h1>音声録音と送信（リアルタイム再生）</h1>
      <button onClick={connectWebSocket}>WebSocket接続</button>
      <button onClick={startRecording} disabled={isRecording}>録音開始</button>
      <button onClick={stopRecording} disabled={!isRecording}>録音停止</button>
    </div>
  );
};

export default Home;
