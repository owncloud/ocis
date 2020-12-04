import { Options, Threshold } from 'k6/options';
import { times } from 'lodash';

import { playbook, types, utils } from '../../../../../../lib';

interface File {
    size: number;
    unit: types.AssetUnit;
}

interface Plays {
    davUpload: playbook.dav.Upload;
    davPropfind: playbook.dav.Propfind;
    davCreate: playbook.dav.Create;
    davDelete: playbook.dav.Delete;
}

export const options = ({ files, plays }: { files: File[]; plays: Plays }): Options => {
    return {
        thresholds: files.reduce((acc: { [name: string]: Threshold[] }, c) => {
            acc[`${plays.davUpload.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
            acc[`${plays.davCreate.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
            acc[`${plays.davDelete.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
            return acc;
        }, {}),
    };
};

export default ({
    files,
    account,
    credential,
    plays,
}: {
    plays: Plays;
    files: File[];
    account: types.Account;
    credential: types.Credential;
}): void => {
    const filesUploaded: { id: string; name: string; folder: string }[] = [];

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
