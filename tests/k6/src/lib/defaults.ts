import * as types from './types';

export const host: types.Host = {
    name: __ENV.OC_HOST_NAME || 'localhost:9200',
    type: !!__ENV.TEST_OC10 ? types.HostType.Oc10 : types.HostType.Ocis
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
    admin: {
        login: 'admin',
        password: 'admin',
    },
}