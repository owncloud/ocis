import {Options} from "k6/options";
import {utils, auth, defaults, playbook} from '../../../../../../lib'
import {times} from 'lodash'

// upload, download and delete of many files with several sizes and summary size of 500 MB in one directory

const files: {
    size: number;
    unit: any;
}[] = [
    ...times(100, () => ({size: 500, unit: 'KB'})),
    ...times(50, () => ({size: 5, unit: 'MB'})),
    ...times(10, () => ({size: 25, unit: 'MB'})),
]
const authFactory = new auth.default(utils.buildAccount({login: defaults.ACCOUNTS.EINSTEIN}));
const plays = {
    fileUpload: playbook.dav.fileUpload({}),
    fileDownload: playbook.dav.fileDownload({}),
    fileDelete: playbook.dav.fileDelete({}),
}
export const options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 3,
    vus: 1,
    thresholds: files.reduce((acc: any, c) => {
        acc[`${plays.fileUpload.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
        acc[`${plays.fileDownload.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
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

    filesUploaded.forEach(f => {
        plays.fileDownload.exec({
            credential,
            userName: account.login,
            path: f.name,
            tags: {asset: f.id},
        });
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