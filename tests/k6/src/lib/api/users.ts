import http, {RefinedResponse, ResponseType} from "k6/http";
import * as api from './api'
import * as defaults from "../defaults";
import * as types from "../types";

export const userInfo = <RT extends ResponseType | undefined>(
    {credential, userName}: { credential: types.Credential; userName: string; }
): RefinedResponse<RT> => {
    return http.get(
        `${defaults.ENV.HOST}/ocs/v1.php/cloud/users/${userName}`,
        {
            headers: {
                ...api.headersDefault({credential})
            },
        }
    );
}