import {Gauge, Trend} from "k6/metrics";
import * as api from "../api";
import * as utils from "../utils";
import {bytes, check} from "k6";
import * as types from "../types";

export const fileUpload = () => {
    const fileUploadTrend = new Trend('occ_file_upload_trend', true);
    const fileUploadErrorRate = new Gauge('occ_file_upload_error_rate');

    return ({credential, userName, asset}: { credential: types.Account | types.Token; userName: string; asset: types.Asset }): string => {
        const fileName = `upload-${userName}-${__VU}-${__ITER}.${utils.extension(asset.fileName)}`;
        const uploadResponse = api.dav.fileUpload({
            credential: credential as any,
            asset: {
                fileName,
                bytes: asset.bytes,
            },
            userName,
        });

        check(uploadResponse, {
            'file upload status is 201': () => uploadResponse.status === 201,
        }) || fileUploadErrorRate.add(1);

        fileUploadTrend.add(uploadResponse.timings.duration)

        return fileName
    }
};

export const fileDelete = () => {
    const fileDeleteTrend = new Trend('occ_file_delete_trend', true);
    const fileDeleteErrorRate = new Gauge('occ_file_delete_error_rate');

    return ({credential, userName, fileName}: { credential: types.Account | types.Token, userName: string; fileName: string }) => {
        const deleteResponse = api.dav.fileDelete({
            credential: credential as any,
            fileName,
            userName,
        });

        check(deleteResponse, {
            'file delete status is 204': () => deleteResponse.status === 204,
        }) || fileDeleteErrorRate.add(1);

        fileDeleteTrend.add(deleteResponse.timings.duration)
    }
};

export const fileDownload = () => {
    const fileDownloadTrend = new Trend('occ_file_download_trend', true);
    const fileDownloadErrorRate = new Gauge('occ_file_download_error_rate');

    return ({credential, userName, fileName}: { credential: types.Account | types.Token, userName: string; fileName: string }): bytes => {
        const downloadResponse = api.dav.fileDownload({
            credential: credential as any,
            fileName,
            userName,
        });

        check(downloadResponse, {
            'file download status is 200': () => downloadResponse.status === 200,
        }) || fileDownloadErrorRate.add(1);

        fileDownloadTrend.add(downloadResponse.timings.duration)

        return downloadResponse.body as bytes
    }
};
