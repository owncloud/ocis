import {Options} from "k6/options";
import {utils, auth, defaults, playbook} from '../../../../../../lib'
import {times} from 'lodash'

// put 1000 files into one dir and run a 'PROPFIND' through API

const files: {
    size: number;
    unit: any;
}[] = times(1000, () => ({size: 1, unit: 'KB'}))
const authFactory = new auth.default(utils.buildAccount({login: defaults.ACCOUNTS.EINSTEIN}));
const plays = {
    davUpload: new playbook.dav.Upload({}),
    davPropfind: new playbook.dav.Propfind({}),
    davDelete: new playbook.dav.Delete({}),
}
export const options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 3,
    vus: 1,
    thresholds: files.reduce((acc: any, c) => {
        acc[`${plays.davUpload.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
        acc[`${plays.davPropfind.metricTrendName}`] = []
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

    plays.davPropfind.exec({
        credential,
        userName: account.login,
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