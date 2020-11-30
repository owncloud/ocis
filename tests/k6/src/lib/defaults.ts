import * as types from './types';
import {Options} from "k6/options";

const ocTestFile = '../_files/' + (__ENV.OC_TEST_FILE || 'kb_50.jpg').split('/').pop()
export const OC_HOST = __ENV.OC_HOST || 'https://localhost:9200'
export const OC_LOGIN = __ENV.OC_LOGIN
export const OC_PASSWORD = __ENV.OC_PASSWORD
export const OC_OIDC_HOST = __ENV.OC_OIDC_HOST || OC_HOST
export const OC_OIDC = __ENV.OC_OIDC === 'true' || false
export const OC_TEST_FILE = {
    fileName: ocTestFile,
    bytes: open(ocTestFile, 'b'),
}
export const K6_OPTION_DEFAULTS: Options = {
    insecureSkipTLSVerify: true,
};

export class ACCOUNTS {
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
        if (OC_LOGIN && OC_PASSWORD) {
            return {
                login: OC_LOGIN,
                password: OC_PASSWORD,
            }
        }

        return this.list[key];
    }
}
