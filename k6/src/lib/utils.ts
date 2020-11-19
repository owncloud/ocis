import {sample} from "lodash";
import {defaults} from "./index";
import {bytes} from "k6";

export const randomString = (): string => {
    return Math.random().toString(36).slice(2)
}

export const getFile = (name?: string): { name: string, bytes: bytes } => {
    if (!name) {
        name = sample(Object.keys(defaults.files))
    }

    const selectedFile = defaults.files[name]

    return {
        name: `${name.split('.')[0]}-${randomString()}.${name.split('.')[1]}`,
        bytes: selectedFile,
    }
}