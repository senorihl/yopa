export type ConfigurationInput = {
    scheme: string,
    domain: string,
    pixel: string,
    debug: true,
}

export class Yopa {
    private readonly _scheme: string ;
    private readonly _domain: string;
    private readonly _pixel: string;
    private readonly _hooks = new Hooks();

    constructor(config: Partial<ConfigurationInput>) {
        this._scheme = config.scheme || 'https';
        this._domain = config.domain || 'www.yopa.io';
        this._pixel = config.pixel || '/pixel.gif';

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
        }
    }

    public sendEvent(name: string) {
        this._hooks.trigger("build:before", { event_name: name });
        const event = { event_name: name };
        this._hooks.trigger("build:after",  event);
        this._hooks.trigger("send:before",  `${this._scheme}://${this._domain}${this._pixel}?p=${encodeURIComponent(JSON.stringify(event))}`);
        this._hooks.trigger("send:after",  ``);
    }

    public Hooks() {
        return this._hooks;
    }
}

type Arrayable<T> = T | Array<T>;

type HooksList = {
    "send:before":  Array<(url: string) => void>,
    "send:after":   Array<(url: string) => void>,
    "build:before": Array<(properties: {[k: string]: number | string | boolean | Date }) => void>,
    "build:after":  Array<(properties: {[k: string]: number | string | boolean }) => void>,
}

class Hooks {
    private readonly _hooks: HooksList = {
        "send:before":  [],
        "send:after":   [],
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