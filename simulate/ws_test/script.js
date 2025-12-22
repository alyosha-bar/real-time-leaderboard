import ws from 'k6/ws';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 },  // Ramp up
    { duration: '1m', target: 200 },   // Stay at 200
    { duration: '30s', target: 0 },    // Ramp down
  ],
};

export default function () {
  const url = 'ws://host.docker.internal:8080/ws';

  const res = ws.connect(url, {}, function (socket) {
    // k6 uses .on('event', callback) NOT .onopen = ...

    socket.on('open', () => {
      console.log('Connected');
      socket.send(JSON.stringify({ message: "Hello from k6!" }));
    });

    socket.on('message', (data) => {
      console.log('Message received: ' + data);
    });

    socket.on('error', (e) => {
      console.log('An error occurred: ', e.error());
    });

    socket.on('close', () => {
      console.log('Disconnected');
    });

    // CRITICAL: Keep the connection open for 10 seconds 
    // before finishing the iteration.
    socket.setTimeout(() => {
      socket.close();
    }, 10000);
  });

  check(res, { 'status is 101': (r) => r && r.status === 101 });
}