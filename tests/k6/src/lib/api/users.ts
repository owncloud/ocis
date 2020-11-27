import http, {RefinedResponse, ResponseType} from "k6/http";
import * as api from './api'
import * as defaults from "../defaults";
import * as types from "../types";

export const userInfo = <RT extends ResponseType | undefined>(
    {credential, userName}: { credential: types.Account | types.Token; userName: string; }
): RefinedResponse<RT> => {
    return http.get(
        `${defaults.OC_OCIS_HOST}/ocs/v1.php/cloud/users/${userName}`,
        {
            headers: {
                ...api.headersDefault({credential})
            },
        }
    );
}