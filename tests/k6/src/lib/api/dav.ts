import http, {RefinedResponse, ResponseType} from 'k6/http';
import * as api from './api'
import * as defaults from '../defaults';
import * as types from '../types';

export const fileUpload = (
    {
        credential,
        userName,
        path = '',
        asset,
        tags,
    }: {
        credential: types.Credential;
        userName: string;
        asset: types.Asset;
        path?: string;
        tags?: { [key: string]: string };
    }
): RefinedResponse<ResponseType> => {

    return http.put(
        [
            defaults.ENV.HOST,
            ...`/remote.php/dav/files/${userName}/${path}/${asset.name}`.split('/').filter(Boolean)
        ].join('/'),
        asset.bytes as any,
        {
            tags,
            headers: {
                ...api.headersDefault({credential})
            }
        }
    );
}

export const fileDownload = (
    {
        credential,
        userName,
        path,
        tags,
    }: {
        credential: types.Credential;
        userName: string;
        path: string;
        tags?: { [key: string]: string };
    }
): RefinedResponse<ResponseType> => {
    return http.get(
        [
            defaults.ENV.HOST,
            ...`/remote.php/dav/files/${userName}/${path}`.split('/').filter(Boolean)
        ].join('/'),
        {
            tags,
            headers: {
                ...api.headersDefault({credential})
            }
        }
    );
}

export const fileDelete = (
    {
        credential,
        userName,
        path,
        tags,
    }: {
        credential: types.Credential;
        userName: string;
        path: string;
        tags?: { [key: string]: string };
    }
): RefinedResponse<ResponseType> => {
    return http.del(
        [
            defaults.ENV.HOST,
            ...`/remote.php/dav/files/${userName}/${path}`.split('/').filter(Boolean)
        ].join('/'),
        {},
        {
            tags,
            headers: {
                ...api.headersDefault({credential})
            }
        }
    );
}

export const folderCreate = (
    {
        credential,
        userName,
        path,
        tags,
    }: {
        credential: types.Credential;
        userName: string;
        path: string;
        tags?: { [key: string]: string };
    }
): RefinedResponse<ResponseType> => {
    return http.request(
        'MKCOL',
        [
            defaults.ENV.HOST,
            ...`/remote.php/dav/files/${userName}/${path}`.split('/').filter(Boolean)
        ].join('/'),
        {},
        {
            tags,
            headers: {
                ...api.headersDefault({credential})
            }
        }
    );
}

export const folderDelete = fileDelete

export const propfind = (
    {
        credential,
        userName,
        path = '',
        tags,
    }: {
        credential: types.Credential;
        userName: string;
        path?: string;
        tags?: { [key: string]: string };
    }
): RefinedResponse<ResponseType> => {
    return http.request(
        'PROPFIND',
        [
            defaults.ENV.HOST,
            ...`/remote.php/dav/files/${userName}/${path}`.split('/').filter(Boolean)
        ].join('/'),
        {},
        {
            tags,
            headers: {
                ...api.headersDefault({credential})
            }
        }
    );
}