import {bytes} from 'k6';

export interface Asset {
    bytes: bytes;
    name: string;
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

export type AssetUnit = 'KB' | 'MB' | 'GB'

export type Tags = { [key: string]: string }