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
    fileUpload: playbook.dav.fileUpload({}),
    propfind: playbook.dav.propfind({}),
    fileDelete: playbook.dav.fileDelete({}),
}
export const options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 3,
    vus: 1,
    thresholds: files.reduce((acc: any, c) => {
        acc[`${plays.fileUpload.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
        acc[`${plays.propfind.metricTrendName}`] = []
        acc[`${plays.fileDelete.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
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

        plays.fileUpload.exec({
            credential,
            asset,
            userName: account.login,
            tags: {asset: id},
        });

        filesUploaded.push({id, name: asset.name})
    })

    plays.propfind.exec({
        credential,
        userName: account.login,
    })

    filesUploaded.forEach(f => {
        plays.fileDelete.exec({
            credential,
            userName: account.login,
            path: f.name,
            tags: {asset: f.id},
        });
    })
}