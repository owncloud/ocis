import {sample} from "lodash";
import {bytes} from "k6";
import * as defaults from "./defaults";

export const randomString = (): string => {
    return Math.random().toString(36).slice(2)
}

export const getAccount = (name?: string): { login: string, password: string } => {
    if (!name) {
        name = sample(Object.keys(defaults.accounts))
    }

    return defaults.accounts[name]
}

export const getFile = (name?: string): { name: string, nameRandom: string, bytes: bytes } => {
    if (!name) {
        name = sample(Object.keys(defaults.files))
    }

    const selectedFile = defaults.files[name]

    return {
        name,
        nameRandom: `${name.split('.')[0]}-${randomString()}.${name.split('.')[1]}`,
        bytes: selectedFile,
    }
}

export const utl = (name?: string): string => {
    return ""
}