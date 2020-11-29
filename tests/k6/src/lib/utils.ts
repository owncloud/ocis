import * as types from "./types";
import * as defaults from "./defaults";

export const randomString = (): string => {
    return Math.random().toString(36).slice(2)
}

export const extension = (p: string): string | undefined => {
    return (p.split('/').pop())!.split('.').pop()
}

export const getAccount = (key: string): types.Account => {
    if (defaults.OC_LOGIN && defaults.OC_PASSWORD) {
        return {
            login: defaults.OC_LOGIN,
            password: defaults.OC_PASSWORD,
        }
    }

    return defaults.knownAccounts[key];
}

