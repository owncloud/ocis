import * as types from './types';
import {Options} from "k6/options";

export class K6 {
    public static readonly OPTIONS: Options = {
        insecureSkipTLSVerify: true,
        iterations: 1,
        vus: 1,
    };
}

export class ENV {
    public static readonly HOST = __ENV.OC_HOST || 'https://localhost:9200';
    public static readonly LOGIN = __ENV.OC_LOGIN;
    public static readonly PASSWORD = __ENV.OC_PASSWORD;
    public static readonly OIDC_HOST = __ENV.OC_OIDC_HOST || ENV.HOST;
    public static readonly OIDC_ENABLED = __ENV.OC_OIDC_ENABLED === 'true' || false;
    public static readonly FILE_NAME = '../_files/' + (__ENV.OC_TEST_FILE || 'kb_50.jpg').split('/').pop();
}

export const FILE = {
    fileName: ENV.FILE_NAME,
    bytes: open(ENV.FILE_NAME, 'b'),
};

export class ACCOUNT {
    public static readonly EINSTEIN = 'einstein';
    public static readonly RICHARD = 'richard';
    private static readonly list: { [key: string]: types.Account; } = {
        einstein: {
            login: 'einstein',
            password: 'relativity',
        },
        richard: {
            login: 'richard',
            password: 'superfluidity',
        },
    }

    public static for(key: string): types.Account {
        if (ENV.LOGIN && ENV.PASSWORD) {
            return {
                login: ENV.LOGIN,
                password: ENV.PASSWORD,
            }
        }

        return this.list[key];
    }
}
