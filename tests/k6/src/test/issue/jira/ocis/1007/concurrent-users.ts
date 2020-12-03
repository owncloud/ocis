import { Options, Threshold } from 'k6/options';
import { utils, auth, defaults, playbook, types } from '../../../../../lib';

const files: Array<{
    size: number;
    unit: types.AssetUnit;
}> = [
    { size: 50, unit: 'KB' },
    { size: 500, unit: 'KB' },
    { size: 5, unit: 'MB' },
    { size: 50, unit: 'MB' },
];
const adminAuthFactory = new auth.default(utils.buildAccount({ login: defaults.ACCOUNTS.ADMIN }));
const plays = {
    usersCreate: new playbook.users.Create(),
    usersDelete: new playbook.users.Delete(),
    davUpload: new playbook.dav.Upload(),
    davDelete: new playbook.dav.Delete(),
};
export const options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 10,
    vus: 10,
    thresholds: files.reduce((acc: { [name: string]: Threshold[] }, c) => {
        acc[`${plays.davUpload.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
        acc[`${plays.davDelete.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
        return acc;
    }, {}),
};

export default (): void => {
    const userName: string = utils.randomString();
    const password: string = utils.randomString();

    plays.usersCreate.exec({
        userName,
        password,
        email: `${userName}@owncloud.com`,
        credential: adminAuthFactory.credential,
    });
    const userAuthFactory = new auth.default({ login: userName, password });
    const filesUploaded: { id: string; name: string }[] = [];

    files.forEach((f) => {
        const id = f.unit + f.size.toString();

        const asset = utils.buildAsset({
            name: `${userName}-dummy.zip`,
            unit: f.unit,
            size: f.size,
        });

        plays.davUpload.exec({
            credential: userAuthFactory.credential,
            asset,
            userName,
            tags: { asset: id },
        });

        filesUploaded.push({ id, name: asset.name });
    });

    filesUploaded.forEach((f) => {
        plays.davDelete.exec({
            credential: userAuthFactory.credential,
            userName: userAuthFactory.account.login,
            path: f.name,
            tags: { asset: f.id },
        });
    });

    plays.usersDelete.exec({ userName: userName, credential: adminAuthFactory.credential });
};
