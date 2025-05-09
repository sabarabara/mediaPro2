// server.js
const http = require('http');
const next = require('next');
const { WebSocketServer } = require('ws');
const fs = require('fs');
const path = require('path');
// const Speaker = require('speaker'); // Disable speaker to avoid ALSA errors

const dev = process.env.NODE_ENV !== 'production';
const app = next({ dev });
const handle = app.getRequestHandler();

app.prepare().then(() => {
  // Create HTTP server to handle Next.js pages
  const server = http.createServer((req, res) => handle(req, res));

  // Attach WebSocket server on /ws path
  const wss = new WebSocketServer({ server, path: '/ws' });
  wss.on('connection', (ws) => {
    console.log('WebSocket connected');
    ws.on('message', (data) => {
      try {
        const buffer = Buffer.isBuffer(data) ? data : Buffer.from(data);
        //const audioPath = path.join(process.cwd(), 'received_audio.wav');
        const audioPath = path.join(process.cwd(), 'received_audio.webm');
        fs.writeFileSync(audioPath, buffer);
        console.log('Saved audio:', audioPath);

        // Optionally play audio on server - commented out to prevent ALSA errors
        // try {
        //   const speaker = new Speaker({ channels: 1, bitDepth: 16, sampleRate: 44100 });
        //   fs.createReadStream(audioPath).pipe(speaker);
        // } catch (e) {
        //   console.warn('Speaker playback failed, skipping on server');
        // }

        // Send raw WAV back to client
        const out = fs.readFileSync(audioPath);
        ws.send(out);
        console.log('Sent raw WAV back to client');
      } catch (err) {
        console.error('Error processing message:', err);
        ws.send(JSON.stringify({ error: 'Server error' }));
      }
    });
  });

  server.listen(3000, () => {
    console.log('> Listening on http://localhost:3000');
  });
});
