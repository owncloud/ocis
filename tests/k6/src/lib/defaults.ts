import * as types from './types';

export class ENV {
    public static readonly HOST = __ENV.OC_HOST || 'https://localhost:9200';
    public static readonly LOGIN = __ENV.OC_LOGIN;
    public static readonly PASSWORD = __ENV.OC_PASSWORD;
    public static readonly OIDC_HOST = __ENV.OC_OIDC_HOST || ENV.HOST;
    public static readonly OIDC_ENABLED = __ENV.OC_OIDC_ENABLED === 'true' || false;
}

export class ACCOUNTS {
    public static readonly EINSTEIN = 'einstein';
    public static readonly RICHARD = 'richard';
    public static readonly ALL: { [key: string]: types.Account; } = {
        einstein: {
            login: 'einstein',
            password: 'relativity',
        },
        richard: {
            login: 'richard',
            password: 'superfluidity',
        },
    }
}