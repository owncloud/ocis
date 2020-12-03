import * as defaults from './defaults';
import http from 'k6/http';
import queryString from 'query-string';
import * as types from './types';
import { fail } from 'k6';
import { get } from 'lodash';

export default class Factory {
    private provider!: types.AuthProvider;
    public account!: types.Account;

    constructor(account: types.Account) {
        this.account = account;

        if (defaults.ENV.OIDC_ENABLED) {
            this.provider = new OIDCProvider(account);
            return;
        }

        this.provider = new AccountProvider(account);
    }

    public get credential(): types.Credential {
        return this.provider.credential;
    }
}

class AccountProvider implements types.AuthProvider {
    private account: types.Account;

    constructor(account: types.Account) {
        this.account = account;
    }

    public get credential(): types.Account {
        return this.account;
    }
}

class OIDCProvider implements types.AuthProvider {
    private account: types.Account;
    private redirectUri = `${defaults.ENV.OIDC_HOST}/oidc-callback.html`;
    private logonUri = `${defaults.ENV.OIDC_HOST}/signin/v1/identifier/_/logon`;
    private tokenUrl = `${defaults.ENV.OIDC_HOST}/konnect/v1/token`;
    private cache!: {
        validTo: Date;
        token: types.Token;
    };

    constructor(account: types.Account) {
        this.account = account;
    }

    public get credential(): types.Token {
        if (!this.cache || this.cache.validTo <= new Date()) {
            const continueURI = this.getContinueURI();
            const code = this.getCode(continueURI);
            const token = this.getToken(code);

            this.cache = {
                validTo: ((): Date => {
                    const offset = 5;
                    const d = new Date();

                    d.setSeconds(d.getSeconds() + token.expiresIn - offset);

                    return d;
                })(),
                token,
            };
        }

        return this.cache.token;
    }

    private getContinueURI(): string {
        const logonResponse = http.post(
            this.logonUri,
            JSON.stringify({
                params: [this.account.login, this.account.password, '1'],
                hello: {
                    scope: 'openid profile email',
                    client_id: 'phoenix',
                    redirect_uri: this.redirectUri,
                    flow: 'oidc',
                },
                state: 'vp42cf',
            }),
            {
                headers: {
                    'Kopano-Konnect-XSRF': '1',
                    Referer: defaults.ENV.OIDC_HOST,
                    'Content-Type': 'application/json',
                },
            },
        );
        const continueURI = get(logonResponse.json(), 'hello.continue_uri');

        if (logonResponse.status != 200 || !continueURI) {
            fail(this.logonUri);
        }

        return continueURI;
    }

    private getCode(continueURI: string): string {
        const authorizeUri = `${continueURI}?${queryString.stringify({
            client_id: 'phoenix',
            prompt: 'none',
            redirect_uri: this.redirectUri,
            response_mode: 'query',
            response_type: 'code',
            scope: 'openid profile email',
        })}`;
        const authorizeResponse = http.get(authorizeUri, {
            redirects: 0,
        });

        const code = get(queryString.parseUrl(authorizeResponse.headers.Location), 'query.code');

        if (authorizeResponse.status != 302 || !code) {
            fail(continueURI);
        }

        return code;
    }

    private getToken(code: string): types.Token {
        const tokenResponse = http.post(this.tokenUrl, {
            client_id: 'phoenix',
            code,
            redirect_uri: this.redirectUri,
            grant_type: 'authorization_code',
        });

        const token = {
            accessToken: get(tokenResponse.json(), 'access_token'),
            tokenType: get(tokenResponse.json(), 'token_type'),
            idToken: get(tokenResponse.json(), 'id_token'),
            expiresIn: get(tokenResponse.json(), 'expires_in'),
        };

        if (
            tokenResponse.status != 200 ||
            !token.accessToken ||
            !token.tokenType ||
            !token.idToken ||
            !token.expiresIn
        ) {
            fail(this.tokenUrl);
        }

        return token;
    }
}
