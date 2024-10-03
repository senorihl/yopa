export default {
    post(url: string, data: string, callback?: (url: string, data: string, res?: never) => void) {
        let queued = false;
        if (window.navigator && typeof window.navigator.sendBeacon === 'function') {
            queued = window.navigator.sendBeacon(url, data);
        }
        if (!queued && window.fetch) {
            window.fetch(url, {
                method: 'POST',
                body: data,
                headers: {
                    'Content-Type': 'text/plain;charset=UTF-8'
                }
            }).then((res) => {
                callback && callback(url, data);
            });
        }
    }
}