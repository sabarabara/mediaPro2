"use client";

import React, { useState, useRef } from 'react';

const Home = () => {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [isRecording, setIsRecording] = useState(false);
  const audioContextRef = useRef<AudioContext | null>(null);
  const processorRef = useRef<ScriptProcessorNode | null>(null);
  const streamRef = useRef<MediaStream | null>(null);

  const connectWebSocket = () => {
    const socket = new WebSocket('ws://localhost:8080/ws');
    socket.binaryType = "arraybuffer";

    socket.onopen = () => {
      console.log('WebSocket connection established');
      setWs(socket);
    };

    socket.onmessage = (event) => {
      console.log('音声データ受信');
      const audioBlob = new Blob([event.data], { type: 'audio/wav' });
      const audioUrl = URL.createObjectURL(audioBlob);
      const audio = new Audio(audioUrl);
      audio.play();
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    socket.onclose = () => {
      console.log('WebSocket connection closed');
    };
  };

  const startRecording = async () => {
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      console.error('音声録音非対応のブラウザです');
      return;
    }

    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    streamRef.current = stream;

    const audioContext = new AudioContext({ sampleRate: 16000 });
    audioContextRef.current = audioContext;

    const source = audioContext.createMediaStreamSource(stream);
    const processor = audioContext.createScriptProcessor(4096, 1, 1);
    processorRef.current = processor;

    processor.onaudioprocess = (e) => {
      const floatData = e.inputBuffer.getChannelData(0);
      const int16Data = convertFloat32ToInt16(floatData);
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(int16Data.buffer);
      }
    };

    source.connect(processor);
    processor.connect(audioContext.destination);

    setIsRecording(true);
  };

  const stopRecording = () => {
    processorRef.current?.disconnect();
    streamRef.current?.getTracks().forEach((track) => track.stop());
    audioContextRef.current?.close();

    setIsRecording(false);
  };

  const convertFloat32ToInt16 = (buffer: Float32Array): Int16Array => {
    const int16 = new Int16Array(buffer.length);
    for (let i = 0; i < buffer.length; i++) {
      const s = Math.max(-1, Math.min(1, buffer[i]));
      int16[i] = s < 0 ? s * 0x8000 : s * 0x7FFF;
    }
    return int16;
  };

  return (
    <div>
      <h1>音声録音と送信（PCM形式）</h1>
      <button onClick={connectWebSocket}>WebSocket接続</button>
      <button onClick={startRecording} disabled={isRecording}>録音開始</button>
      <button onClick={stopRecording} disabled={!isRecording}>録音停止</button>
    </div>
  );
};

export default Home;
