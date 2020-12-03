import { Options, Threshold } from 'k6/options';
import { utils, auth, defaults, playbook, types } from '../../../../../../lib';
import { times } from 'lodash';

// upload, download and delete of many files with several sizes and summary size of 500 MB in one directory

const files: {
    size: number;
    unit: types.AssetUnit;
}[] = [
    ...times(100, () => ({ size: 500, unit: 'KB' as types.AssetUnit })),
    ...times(50, () => ({ size: 5, unit: 'MB' as types.AssetUnit })),
    ...times(10, () => ({ size: 25, unit: 'MB' as types.AssetUnit })),
];
const authFactory = new auth.default(utils.buildAccount({ login: defaults.ACCOUNTS.EINSTEIN }));
const plays = {
    davUpload: new playbook.dav.Upload(),
    davDownload: new playbook.dav.Download(),
    davDelete: new playbook.dav.Delete(),
};
export const options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 3,
    vus: 1,
    thresholds: files.reduce((acc: { [name: string]: Threshold[] }, c) => {
        acc[`${plays.davUpload.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
        acc[`${plays.davDownload.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
        acc[`${plays.davDelete.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
        return acc;
    }, {}),
};
export default (): void => {
    const filesUploaded: { id: string; name: string }[] = [];
    const { account, credential } = authFactory;

    files.forEach((f) => {
        const id = f.unit + f.size.toString();

        const asset = utils.buildAsset({
            name: `${account.login}-dummy.zip`,
            unit: f.unit,
            size: f.size,
        });

        plays.davUpload.exec({
            credential,
            asset,
            userName: account.login,
            tags: { asset: id },
        });

        filesUploaded.push({ id, name: asset.name });
    });

    filesUploaded.forEach((f) => {
        plays.davDownload.exec({
            credential,
            userName: account.login,
            path: f.name,
            tags: { asset: f.id },
        });
    });

    filesUploaded.forEach((f) => {
        plays.davDelete.exec({
            credential,
            userName: account.login,
            path: f.name,
            tags: { asset: f.id },
        });
    });
};
