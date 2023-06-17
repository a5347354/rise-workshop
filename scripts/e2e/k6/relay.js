import ws from 'k6/ws';
import { check } from 'k6';

export default function () {
    const url = `${__ENV.RELAY_HOSTNAME}`;
    const params = {};

    const res = ws.connect(url, params, function (socket) {
        socket.on('open', function open() {
            console.log('connected');

            const req = JSON.stringify([
                "REQ",
                "12j312n31knkajsndaksndas",
                { limit: 5 },
            ]);
            socket.send(req)
            socket.setInterval(function timeout() {
                socket.ping();
                console.log('Pinging every 1sec (setInterval test)');
            }, 1000);
        });

        socket.on('close', () => console.log('disconnected'));
    });

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}