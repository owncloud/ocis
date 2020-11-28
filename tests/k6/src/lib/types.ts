export interface Account {
    login: string
    password: string
}
export interface UserRequestData {
    displayname: string,
    email: string,
    password: string,
    userid: string
}

export enum HostType {
    Ocis = 1,
    Oc10
}

export interface Host {
    name: string,
    type: HostType,
}