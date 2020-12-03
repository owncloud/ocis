import { Options, Threshold } from 'k6/options';
import { utils, auth, defaults, playbook, types } from '../../../../../../lib';
import { times } from 'lodash';

// Unpack standard data tarball, run PROPFIND on each dir

const files: {
    size: number;
    unit: types.AssetUnit;
}[] = times(1000, () => ({ size: 1, unit: 'KB' }));
const authFactory = new auth.default(utils.buildAccount({ login: defaults.ACCOUNTS.EINSTEIN }));
const plays = {
    davUpload: new playbook.dav.Upload(),
    davPropfind: new playbook.dav.Propfind(),
    davCreate: new playbook.dav.Create(),
    davDelete: new playbook.dav.Delete(),
};
export const options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 3,
    vus: 1,
    thresholds: files.reduce((acc: { [name: string]: Threshold[] }, c) => {
        acc[`${plays.davUpload.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
        acc[`${plays.davCreate.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
        acc[`${plays.davDelete.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
        return acc;
    }, {}),
};
export default (): void => {
    const filesUploaded: { id: string; name: string; folder: string }[] = [];
    const { account, credential } = authFactory;

    files.forEach((f) => {
        const id = f.unit + f.size.toString();

        const asset = utils.buildAsset({
            name: `${account.login}-dummy.zip`,
            unit: f.unit,
            size: f.size,
        });

        const folder = times(utils.randomNumber({ min: 1, max: 10 }), () => utils.randomString())
            .reduce((acc: string[], c) => {
                acc.push(c);

                plays.davCreate.exec({
                    credential,
                    path: acc.join('/'),
                    userName: account.login,
                    tags: { asset: id },
                });

                return acc;
            }, [])
            .join('/');

        plays.davUpload.exec({
            credential,
            asset,
            path: folder,
            userName: account.login,
            tags: { asset: id },
        });

        filesUploaded.push({ id, name: asset.name, folder });
    });

    plays.davPropfind.exec({
        credential,
        userName: account.login,
    });

    filesUploaded.forEach((f) => {
        plays.davDelete.exec({
            credential,
            userName: account.login,
            path: f.folder.split('/')[0],
            tags: { asset: f.id },
        });
    });
};
