import https from "https";
import type http from "node:http";

export default {
    post(url: string, data: string, callback?: (url: string, data: string, res?: http.IncomingMessage) => void) {
        const bodyContent = data;
        const _url = new URL(url);
        const options = {
            hostname: _url.hostname,
            port: 443,
            path: _url.pathname + _url.search,
            method: 'POST',
            headers: {
                'Content-Type': 'text/plain;charset=UTF-8'
            }
        };
        const req = https.request(options, res => {
            callback && callback(url, data, res);
        });
        req.write(bodyContent);
        req.end();
    }
}