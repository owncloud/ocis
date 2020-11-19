import {bytes} from "k6";
import http, {RefinedResponse, ResponseType} from "k6/http";

export const uploadFile = <RT extends ResponseType | undefined>(data: bytes, name: string, account: any): RefinedResponse<RT> => {
    const file = http.file(data, name);

    return http.put(
        `https://${account.login}:${account.password}@localhost:9200/remote.php/dav/files/${account.login}/${name}`,
        {file},
        {},
    );
}

export const userInfo = <RT extends ResponseType | undefined>(account: any): RefinedResponse<RT> => {
    return http.get(`https://${account.login}:${account.password}@localhost:9200/ocs/v1.php/cloud/users/${account.login}`);;
}