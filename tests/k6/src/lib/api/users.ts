import { RefinedResponse, ResponseType } from 'k6/http';
import * as api from './api';
import * as types from '../types';

export class Create {
    public static exec({
        userName,
        password,
        email,
        credential,
        tags,
    }: {
        credential: types.Credential;
        userName: string;
        password: string;
        email: string;
        tags?: types.Tags;
    }): RefinedResponse<ResponseType> {
        return api.request({
            method: 'POST',
            credential,
            path: `/ocs/v1.php/cloud/users`,
            params: { tags },
            body: { userid: userName, password, email },
        });
    }
}

export class Delete {
    public static exec({
        userName,
        credential,
        tags,
    }: {
        credential: types.Credential;
        userName: string;
        tags?: types.Tags;
    }): RefinedResponse<ResponseType> {
        return api.request({
            method: 'DELETE',
            credential,
            path: `/ocs/v1.php/cloud/users/${userName}`,
            params: { tags },
        });
    }
}
