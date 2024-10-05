import {Yopa} from "./Yopa";

declare global {
    const BUILD_BROWSER: boolean;
    interface Window {
        _yopa: Yopa
    }
}

export const yopa: Yopa = (() => {
    const _instance  = new Yopa();
    if (BUILD_BROWSER) {
        if (window && !window[_instance.getConfiguration('global_var_name')]) {
            window[_instance.getConfiguration('global_var_name')] = _instance;
        }
    }
    return _instance;
})();
