import { Options } from 'k6/options';
import { times } from 'lodash';

import { auth, defaults, playbook, types, utils } from '../../../../../../lib';
import { default as upDownDelete, options as upDownDeleteOptions } from './shared.lib';

// upload, download and delete of many files with several sizes and summary size of 500 MB in one directory

const files: {
    size: number;
    unit: types.AssetUnit;
}[] = [
    ...times(100, () => ({ size: 500, unit: 'KB' as types.AssetUnit })),
    ...times(50, () => ({ size: 5, unit: 'MB' as types.AssetUnit })),
    ...times(10, () => ({ size: 25, unit: 'MB' as types.AssetUnit })),
];
const authFactory = new auth(utils.buildAccount({ login: defaults.ACCOUNTS.EINSTEIN }));
const plays = {
    davUpload: new playbook.dav.Upload(),
    davDownload: new playbook.dav.Download(),
    davDelete: new playbook.dav.Delete(),
};
export const options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 3,
    vus: 1,
    ...upDownDeleteOptions({ plays, files }),
};

export default (): void =>
    upDownDelete({ files, plays, credential: authFactory.credential, account: authFactory.account });
