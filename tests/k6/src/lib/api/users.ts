import http, {RefinedResponse, ResponseType} from 'k6/http';
import * as api from './api'
import * as defaults from '../defaults';
import * as types from '../types';

export const userInfo = (
    {
        credential,
        userName,
        tags,
    }: {
        credential: types.Credential;
        userName: string;
        tags?: { [name: string]: string };
    }
): RefinedResponse<ResponseType> => {
    return http.get(
        `${defaults.ENV.HOST}/ocs/v1.php/cloud/users/${userName}`,
        {
            tags,
            headers: {
                ...api.headersDefault({credential})
            },
        }
    );
}