import * as defaults from "./defaults";
import http from "k6/http";
import queryString from "query-string";
import * as types from "./types";
import {fail} from 'k6';
import {get} from 'lodash'

export const oidc = (account: types.Account): types.Token => {
    const redirectUri = `${defaults.OC_OIDC_HOST}/oidc-callback.html`;

    const logonUri = `${defaults.OC_OIDC_HOST}/signin/v1/identifier/_/logon`;
    const logonResponse = http.post(
        logonUri,
        JSON.stringify(
            {
                params: [account.login, account.password, '1'],
                hello: {
                    scope: 'openid profile email',
                    client_id: 'phoenix',
                    redirect_uri: redirectUri,
                    flow: 'oidc'
                },
                'state': 'vp42cf'
            },
        ),
        {
            headers: {
                'Kopano-Konnect-XSRF': '1',
                Referer: defaults.OC_OIDC_HOST,
                'Content-Type': 'application/json',
            },
        },
    );
    const authorizeURI = get(logonResponse.json(), 'hello.continue_uri');

    if (logonResponse.status != 200 || !authorizeURI) {
        fail(logonUri);
    }

    const authorizeUri = `${authorizeURI}?${
        queryString.stringify(
            {
                client_id: 'phoenix',
                prompt: 'none',
                redirect_uri: redirectUri,
                response_mode: 'query',
                response_type: 'code',
                scope: 'openid profile email',
            },
        )
    }`;
    const authorizeResponse = http.get(
        authorizeUri,
        {
            redirects: 0,
        },
    )
    const authCode = get(queryString.parseUrl(authorizeResponse.headers.Location), 'query.code')

    if (authorizeResponse.status != 302 || !authCode) {
        fail(authorizeURI);
    }

    const tokenUrl = `${defaults.OC_OIDC_HOST}/konnect/v1/token`;
    const tokenResponse = http.post(
        tokenUrl,
        {
            client_id: 'phoenix',
            code: authCode,
            redirect_uri: redirectUri,
            grant_type: 'authorization_code'
        }
    )

    const token = {
        accessToken: get(tokenResponse.json(), 'access_token'),
        tokenType: get(tokenResponse.json(), 'token_type'),
        idToken: get(tokenResponse.json(), 'id_token'),
        expiresIn: get(tokenResponse.json(), 'expires_in'),
    }

    if (tokenResponse.status != 200 || !token.accessToken || !token.tokenType || !token.idToken || !token.expiresIn) {
        fail(authorizeURI);
    }

    return token
}