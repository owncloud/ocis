import encoding from 'k6/encoding';
import * as types from '../types';
import * as defaults from "../defaults";
import {merge} from 'lodash';
import http, {RefinedParams, RefinedResponse, RequestBody, ResponseType} from "k6/http";

export const buildHeaders = ({credential}: { credential: types.Credential }): { [key: string]: string } => {
    const isOIDCGuard = (credential as types.Token).tokenType !== undefined;
    const authOIDC = credential as types.Token;
    const authBasic = credential as types.Account;

    return {
        Authorization: isOIDCGuard ? `${authOIDC.tokenType} ${authOIDC.accessToken}` : `Basic ${encoding.b64encode(`${authBasic.login}:${authBasic.password}`)}`,
    }
}

export const buildURL = ({path}: { path: string }): string => {
    return [
        defaults.ENV.HOST,
        ...path.split('/').filter(Boolean)
    ].join('/')
}

export const request = ({method, path, body = {}, params = {}, credential}: {
    method: 'PROPFIND' | 'PUT' | 'GET' | 'DELETE' | 'MKCOL',
    path: string,
    credential: types.Credential;
    body?: RequestBody | null,
    params?: RefinedParams<ResponseType> | null
}): RefinedResponse<ResponseType> => {
    return http.request(
        method,
        buildURL({path}),
        body,
        merge({
            headers: {
                ...buildHeaders({credential})
            }
        }, params)
    );
}