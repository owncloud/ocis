import {bytes} from "k6";

export interface Asset {
    bytes: bytes;
    fileName: string;
}

export interface Token {
    accessToken: string;
    tokenType: string;
    idToken: string;
    expiresIn: number;
}

export interface Account {
    login: string
    password: string
}

export type Credential = Token | Account

export interface AuthProvider {
    credential: Credential
}