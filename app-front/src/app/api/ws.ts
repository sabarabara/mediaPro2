import { WebSocketServer, WebSocket } from 'ws';
import fs from 'fs';
import path from 'path';
import Speaker from 'speaker';

// wssをグローバルで保持する
let wss: WebSocketServer | null = null;

export default function handler(req: any, res: any): void {
  if (!wss) {
    console.log('Initializing WebSocket server...');
    wss = new WebSocketServer({ noServer: true });

    wss.on('connection', (ws: WebSocket) => {
      console.log('WebSocket connection established.');
      
      //@ts-expect-error
      ws.on('message', (message: WebSocket.Data) => {
        const audioPath = path.join(process.cwd(), 'received_audio.wav');
        //@ts-expect-error
        fs.writeFile(audioPath, message, (err) => {
          if (err) {
            ws.send('Failed to save audio file');
            return;
          }

          playAudio(audioPath);

          ws.send('音声ファイルを受け取り、再生しました');
        });
      });
    });

    (res.socket.server as any).on('upgrade', (request: any, socket: any, head: any) => {
      wss?.handleUpgrade(request, socket, head, (ws: WebSocket) => {
        wss?.emit('connection', ws, request);
      });
    });
  }

  res.status(200).end();
}

function playAudio(filePath: string): void {
  const speaker = new Speaker({
    channels: 1,
    bitDepth: 16,
    sampleRate: 44100,
  });

  const audioStream = fs.createReadStream(filePath);
  audioStream.pipe(speaker);
}
