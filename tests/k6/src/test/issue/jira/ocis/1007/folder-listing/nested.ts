import {Options} from "k6/options";
import {utils, auth, defaults, playbook} from '../../../../../../lib'
import {times} from 'lodash'

// Unpack standard data tarball, run PROPFIND on each dir

const files: {
    size: number;
    unit: any;
}[] = times(1000, () => ({size: 1, unit: 'KB'}))
const authFactory = new auth.default(utils.buildAccount({login: defaults.ACCOUNTS.EINSTEIN}));
const plays = {
    fileUpload: playbook.dav.fileUpload({}),
    propfind: playbook.dav.propfind({}),
    folderCreate: playbook.dav.folderCreate({}),
    folderDelete: playbook.dav.folderDelete({}),
}
export const options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 3,
    vus: 1,
    thresholds: files.reduce((acc: any, c) => {
        acc[`${plays.fileUpload.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
        acc[`${plays.propfind.metricTrendName}`] = []
        acc[`${plays.folderCreate.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
        acc[`${plays.folderDelete.metricTrendName}{asset:${c.unit + c.size.toString()}`] = []
        return acc
    }, {}),
};
export default (): void => {
    const filesUploaded: { id: string, name: string, folder: string }[] = []
    const {account, credential} = authFactory;

    files.forEach(f => {
        const id = f.unit + f.size.toString();

        const asset = utils.buildAsset({
            name: `${account.login}-dummy.zip`,
            unit: f.unit as any,
            size: f.size,
        })

        const folder = times(utils.randomNumber({min: 1, max: 10}), () => utils.randomString()).reduce((acc: string[], c) => {
            acc.push(c)

            plays.folderCreate.exec({
                credential,
                path: acc.join('/'),
                userName: account.login,
                tags: {asset: id},
            });

            return acc
        }, []).join('/')


        plays.fileUpload.exec({
            credential,
            asset,
            path: folder,
            userName: account.login,
            tags: {asset: id},
        });

        filesUploaded.push({id, name: asset.name, folder})
    })

    plays.propfind.exec({
        credential,
        userName: account.login,
    })

    filesUploaded.forEach(f => {
        plays.folderDelete.exec({
            credential,
            userName: account.login,
            path: f.folder.split('/')[0],
            tags: {asset: f.id},
        });
    })
}