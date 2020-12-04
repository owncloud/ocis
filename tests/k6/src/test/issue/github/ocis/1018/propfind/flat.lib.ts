import { Options, Threshold } from 'k6/options';

import { playbook, types, utils } from '../../../../../../lib';

interface File {
    size: number;
    unit: types.AssetUnit;
}

interface Plays {
    davUpload: playbook.dav.Upload;
    davPropfind: playbook.dav.Propfind;
    davDelete: playbook.dav.Delete;
}

export const options = ({ files, plays }: { files: File[]; plays: Plays }): Options => {
    return {
        thresholds: files.reduce((acc: { [name: string]: Threshold[] }, c) => {
            acc[`${plays.davUpload.metricTrendName}{asset:${c.unit + c.size.toString()}}`] = [];
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
    const filesUploaded: { id: string; name: string }[] = [];

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

    plays.davPropfind.exec({
        credential,
        userName: account.login,
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
