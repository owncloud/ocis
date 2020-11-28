import encoding from 'k6/encoding';
import {bytes} from "k6";
import http, {RefinedResponse, ResponseType} from "k6/http";
import * as defaults from "./defaults";
import * as types from "./types";

export const uploadFile = <RT extends ResponseType | undefined>(account: types.Account, data: bytes, name: string): RefinedResponse<RT> => {
    return http.put(
        `https://${defaults.host.name}/remote.php/dav/files/${account.login}/${name}`,
        data as any,
        {
            headers: {
                Authorization: `Basic ${encoding.b64encode(`${account.login}:${account.password}`)}`,
            }
        }
    );
}

export const downloadFile = <RT extends ResponseType | undefined>(account: types.Account, name: string): RefinedResponse<RT> => {
    return http.get(
        `https://${defaults.host.name}/remote.php/dav/files/${account.login}/${name}`,
        {
            headers: {
                Authorization: `Basic ${encoding.b64encode(`${account.login}:${account.password}`)}`,
            }
        }
    );
}

export const userInfo = <RT extends ResponseType | undefined>(account: any): RefinedResponse<RT> => {
    return http.get(
        `https://${defaults.host.name}/ocs/v1.php/cloud/users/${account.login}`,
        {
            headers: {
                Authorization: `Basic ${encoding.b64encode(`${account.login}:${account.password}`)}`,
            },
        }
    );
}

export const deleteFile = <RT extends ResponseType | undefined>(account: types.Account, name: string): RefinedResponse<RT> => {
    return http.del(
        `https://${defaults.host.name}/remote.php/dav/files/${account.login}/${name}`,
        {},
        {
            headers: {
                Authorization: `Basic ${encoding.b64encode(`${account.login}:${account.password}`)}`,
            }
        }
    );
}