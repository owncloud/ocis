import http, {RefinedResponse, ResponseType} from "k6/http";
import * as api from './api'
import * as defaults from "../defaults";
import * as types from "../types";

export const fileUpload = <RT extends ResponseType | undefined>(
    {credential, userName, asset}: { credential: types.Credential; userName: string; asset: types.Asset }
): RefinedResponse<RT> => {
    return http.put(
        `${defaults.ENV.HOST}/remote.php/dav/files/${userName}/${asset.fileName}`,
        asset.bytes as any,
        {
            headers: {
                ...api.headersDefault({credential})
            }
        }
    );
}

export const fileDownload = <RT extends ResponseType | undefined>(
    {credential, userName, fileName}: { credential: types.Credential; userName: string; fileName: string }
): RefinedResponse<RT> => {
    return http.get(
        `${defaults.ENV.HOST}/remote.php/dav/files/${userName}/${fileName}`,
        {
            headers: {
                ...api.headersDefault({credential})
            }
        }
    );
}

export const fileDelete = <RT extends ResponseType | undefined>(
    {credential, userName, fileName}: { credential: types.Credential; userName: string; fileName: string }
): RefinedResponse<RT> => {
    return http.del(
        `${defaults.ENV.HOST}/remote.php/dav/files/${userName}/${fileName}`,
        {},
        {
            headers: {
                ...api.headersDefault({credential})
            }
        }
    );
}