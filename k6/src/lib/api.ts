import {bytes} from "k6";
import http, {RefinedResponse, ResponseType} from "k6/http";
import * as defaults from "./defaults";

export const uploadFile = <RT extends ResponseType | undefined>(account: any, data: bytes, name: string): RefinedResponse<RT> => {
    const file = http.file(data, name);

    return http.put(
        `https://${account.login}:${account.password}@${defaults.host.name}/remote.php/dav/files/${account.login}/${name}`,
        {file},
        {},
    );
}

export const userInfo = <RT extends ResponseType | undefined>(account: any): RefinedResponse<RT> => {
    return http.get(`https://${account.login}:${account.password}@${defaults.host.name}/ocs/v1.php/cloud/users/${account.login}`);
}