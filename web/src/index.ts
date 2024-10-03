import {Yopa} from "./Yopa";

declare global {
    const BUILD_BROWSER: boolean;
}

export const yopa = (() => {
    const _instance  = new Yopa();
    if (BUILD_BROWSER) {
        // @ts-ignore
        if (window && !window[_instance.getConfiguration('global_var_name')]) {
            // @ts-ignore
            window[_instance.getConfiguration('global_var_name')] = _instance;
        }
    }
    return _instance;
})();
