import {Yopa, type ConfigurationInput} from "./Yopa";

type GlobalConfiguration = Partial<{
    [Property in keyof ConfigurationInput as Uppercase<`_YOPA_${Property}`>]: ConfigurationInput[Property];
}>

declare global {
    interface Window extends GlobalConfiguration { Yopa: Yopa; }
}

const config: Partial<ConfigurationInput> = {
    debug: window._YOPA_DEBUG,
    pixel: window._YOPA_PIXEL,
    scheme: window._YOPA_SCHEME,
    domain: window._YOPA_DOMAIN,
    site: window._YOPA_SITE,
};

window.Yopa = new Yopa(config);
Object.freeze(window.Yopa);
