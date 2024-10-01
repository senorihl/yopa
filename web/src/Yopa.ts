export type ConfigurationInput = {
    site: number,
    scheme: string,
    domain: string,
    pixel: string,
    cookie_domain: string,
    cookie_path: string,
    cookie_secure: true,
    debug: true,
}

const str = () => ('00000000000000000' + (Math.random() * 0xffffffffffffffff).toString(16)).slice(-16);

const uuid = () => {
    const a = str();
    const b = str();
    return a.slice(0, 8) + '-' + a.slice(8, 12) + '-4' + a.slice(13) + '-a' + b.slice(1, 4) + '-' + b.slice(4);
};

export class Yopa {
    private readonly _site: number;
    private readonly _scheme: string ;
    private readonly _domain: string;
    private readonly _pixel: string;
    private readonly _visitor: string;
    private readonly _hooks = new Hooks();

    constructor(config: Partial<ConfigurationInput>) {
        if (!config.site) {
            throw new Error('Missing config site or _YOPA_SITE');
        }
        this._site = config.site;
        this._scheme = config.scheme || 'https';
        this._domain = config.domain || 'www.yopa.io';
        this._pixel = config.pixel || '/pixel.gif';
        const visitorFromCookie = document.cookie
            .split("; ")
            .find((row) => row.startsWith("_yovi="))
            ?.split("=")[1];
        this._visitor = visitorFromCookie || uuid();

        if (this._visitor != visitorFromCookie) {
            const expires = new Date();
            expires.setMonth(expires.getMonth() + 13);
            document.cookie = "_yovi=" + this._visitor
                + ";expires=" + expires
                + ";path=" + (config.cookie_path || '/')
                + (config.cookie_domain ? `;domain=${config.cookie_domain}` : '')
            ;
            console.log('Setted a new visitor id', this._visitor)
        }

        if (config.debug) {
            this.Hooks().add("build:before", (properties) => {
                console.log('Yopa', "build:before", properties)
            })
            this.Hooks().add("build:after", (properties) => {
                console.log('Yopa', "build:after", properties)
            })
            this.Hooks().add("send:before", (url) => {
                console.log('Yopa', "send:before", url)
            })
            this.Hooks().add("send:after", (url) => {
                console.log('Yopa', "send:after", url)
            })
            this.Hooks().add("send:error", (url, err) => {
                console.warn('Yopa', "send:errorr", url, err)
            })
        }
    }

    public sendEvent(name: string) {
        this._hooks.trigger("build:before", { event_name: name });
        const event = { visitor: this._visitor, event_name: name, ts: Date.now() };
        this._hooks.trigger("build:after",  event);
        const url = `${this._scheme}://${this._domain}${this._pixel}?s=${this._site}&p=${encodeURIComponent(JSON.stringify(event))}`;
        this._hooks.trigger("send:before",  url);
        fetch(url)
            .then(() => {
                this._hooks.trigger("send:after",  url);
            })
            .catch((err) => {
                this._hooks.trigger("send:error", url, err);
            })
    }

    public Hooks() {
        return this._hooks;
    }
}

type Arrayable<T> = T | Array<T>;

type HooksList = {
    "send:before":  Array<(url: string) => void>,
    "send:after":   Array<(url: string) => void>,
    "send:error":   Array<(url: string, err: any) => void>,
    "build:before": Array<(properties: {[k: string]: number | string | boolean | Date }) => void>,
    "build:after":  Array<(properties: {[k: string]: number | string | boolean }) => void>,
}

class Hooks {
    private readonly _hooks: HooksList = {
        "send:before":  [],
        "send:after":   [],
        "send:error":   [],
        "build:before": [],
        "build:after":  [],
    }

    public add<K extends keyof HooksList>(name: K, callback: HooksList[K][number]) {
        this._hooks[name].push(callback as any);
    }

    public trigger<K extends keyof HooksList>(name: K, ...data: Parameters<HooksList[K][number]>) {
        for (const callback of this._hooks[name]) {
            callback.apply(null, data)
        }
    }
}