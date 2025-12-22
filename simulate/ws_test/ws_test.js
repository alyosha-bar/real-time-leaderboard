import ws from 'k6/ws';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 100 },  
    { duration: '1m', target: 200 },
    { duration: '3m', target: 200 }, 
    { duration: '1m', target: 0 },
    ],
};


export default function () {
    const url = 'ws://localhost:8080/ws';
    const res = ws.connect(url, {}, function (socket) {
        socket.on('open', function () {
            console.log('WebSocket connection established');
            socket.send(JSON.stringify({ action: 'ping' }));
        });
        socker.on('message', function (message) {
            console.log('Received message: ' + message);
            socket.close();
        });
        socket.on('close', function () {
            console.log('WebSocket connection closed');
        });
    });

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}