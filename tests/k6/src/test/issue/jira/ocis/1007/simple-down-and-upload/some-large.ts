import {Options} from "k6/options";
import {utils, auth, defaults, playbook, types} from '../../../../../../lib'

// upload, download and delete of one file with sizes 50kb, 500kb, 5MB, 50MB, 500MB, 1GB

const files: {
    size: number;
    unit: types.AssetUnit;
}[] = [
    {size: 50, unit: 'KB'},
    {size: 500, unit: 'KB'},
    {size: 5, unit: 'MB'},
    {size: 50, unit: 'MB'},
    {size: 500, unit: 'MB'},
    {size: 1, unit: 'GB'},
]
const authFactory = new auth.default(utils.buildAccount({login: defaults.ACCOUNTS.EINSTEIN}));
const plays = {
    davUpload: new playbook.dav.Upload({}),
    davDownload: new playbook.dav.Download({}),
    davDelete: new playbook.dav.Delete({}),
}
export const options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 3,
    vus: 1,
    thresholds: files.reduce((acc: any, c) => {
        acc[`${plays.davUpload.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
        acc[`${plays.davDownload.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
        acc[`${plays.davDelete.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
        return acc
    }, {}),
};

export default (): void => {
    const filesUploaded: { id: string, name: string, }[] = []
    const {account, credential} = authFactory;

    files.forEach(f => {
        const id = f.unit + f.size.toString();

        const asset = utils.buildAsset({
            name: `${account.login}-dummy.zip`,
            unit: f.unit as any,
            size: f.size,
        })

        plays.davUpload.exec({
            credential,
            asset,
            userName: account.login,
            tags: {asset: id},
        });

        filesUploaded.push({id, name: asset.name})
    })

    filesUploaded.forEach(f => {
        plays.davDownload.exec({
            credential,
            userName: account.login,
            path: f.name,
            tags: {asset: f.id},
        });
    })

    filesUploaded.forEach(f => {
        plays.davDelete.exec({
            credential,
            userName: account.login,
            path: f.name,
            tags: {asset: f.id},
        });
    })
}