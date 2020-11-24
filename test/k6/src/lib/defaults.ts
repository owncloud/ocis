import * as types from './types';

export const host = {
    name: __ENV.OC_HOST_NAME || 'localhost:9200'
}

export const accounts: { [key: string]: types.Account; } = {
    einstein: {
        login: 'einstein',
        password: 'relativity',
    },
    richard: {
        login: 'richard',
        password: 'superfluidity',
    },
}