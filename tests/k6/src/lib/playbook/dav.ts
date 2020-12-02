import {Gauge, Trend} from 'k6/metrics';
import * as api from '../api';
import {check} from 'k6';
import * as types from '../types';
import {RefinedResponse, ResponseType} from 'k6/http';

export const fileUpload = ({name, metricID = 'default'}: { name?: string; metricID?: string; }) => {
    const playName = name || `oc_${metricID}_play_dav_file_upload`;
    const metricTrendName = `${playName}_trend`;
    const metricTrend = new Trend(metricTrendName, true);
    const metricErrorRateName = `${playName}_error_rate`;
    const metricErrorRate = new Gauge(metricErrorRateName);

    return {
        playName,
        metricTrendName,
        metricErrorRateName,
        exec: (
            {
                credential,
                userName,
                path,
                asset,
                tags,
            }: {
                credential: types.Credential;
                path?: string;
                userName: string;
                asset: types.Asset;
                tags?: { [key: string]: string };
            }
        ): {
            response: RefinedResponse<ResponseType>;
            tags: { [key: string]: string };
        } => {
            tags = {play: playName, ...tags};

            const response = api.dav.fileUpload({
                credential: credential as types.Credential,
                asset,
                userName,
                tags,
                path,
            });

            check(response, {
                'file upload status is 201': () => response.status === 201,
            }, tags) || metricErrorRate.add(1, tags);

            metricTrend.add(response.timings.duration, tags)

            return {
                response,
                tags,
            }
        }
    }
};

export const fileDelete = ({name, metricID = 'default'}: { name?: string; metricID?: string; }) => {
    const playName = name || `oc_${metricID}_play_dav_file_delete`;
    const metricTrendName = `${playName}_trend`;
    const metricTrend = new Trend(metricTrendName, true);
    const metricErrorRateName = `${playName}_error_rate`;
    const metricErrorRate = new Gauge(metricErrorRateName);

    return {
        playName,
        metricTrendName,
        metricErrorRateName,
        exec: (
            {
                credential,
                userName,
                path,
                tags,
            }: {
                credential: types.Credential;
                userName: string;
                path: string;
                tags?: { [key: string]: string };
            }
        ): {
            response: RefinedResponse<ResponseType>;
            tags: { [key: string]: string };
        } => {
            tags = {play: playName, ...tags};

            const response = api.dav.fileDelete({
                credential: credential as types.Credential,
                path,
                userName,
                tags,
            });

            check(response, {
                'file delete status is 204': () => response.status === 204,
            }, tags) || metricErrorRate.add(1, tags);

            metricTrend.add(response.timings.duration, tags)

            return {
                response,
                tags,
            }
        }
    }
};

export const fileDownload = ({name, metricID = 'default'}: { name?: string; metricID?: string; }) => {
    const playName = name || `oc_${metricID}_play_dav_file_download`;
    const metricTrendName = `${playName}_trend`;
    const metricTrend = new Trend(metricTrendName, true);
    const metricErrorRateName = `${playName}_error_rate`;
    const metricErrorRate = new Gauge(metricErrorRateName);

    return {
        playName,
        metricTrendName,
        metricErrorRateName,
        exec: (
            {
                credential,
                userName,
                path,
                tags,
            }: {
                credential: types.Credential;
                userName: string;
                path: string;
                tags?: { [key: string]: string };
            }
        ): {
            response: RefinedResponse<ResponseType>;
            tags: { [key: string]: string };
        } => {
            tags = {play: playName, ...tags};

            const response = api.dav.fileDownload({
                credential: credential as types.Credential,
                path,
                userName,
                tags,
            });

            check(response, {
                'file download status is 200': () => response.status === 200,
            }, tags) || metricErrorRate.add(1, tags);

            metricTrend.add(response.timings.duration, tags)

            return {
                response,
                tags,
            }
        }
    }
};

export const folderCreate = ({name, metricID = 'default'}: { name?: string; metricID?: string; }) => {
    const playName = name || `oc_${metricID}_play_dav_folder_create`;
    const metricTrendName = `${playName}_trend`;
    const metricTrend = new Trend(metricTrendName, true);
    const metricErrorRateName = `${playName}_error_rate`;
    const metricErrorRate = new Gauge(metricErrorRateName);

    return {
        playName,
        metricTrendName,
        metricErrorRateName,
        exec: (
            {
                credential,
                userName,
                path,
                tags,
            }: {
                credential: types.Credential;
                userName: string;
                path: string;
                tags?: { [key: string]: string };
            }
        ): {
            response: RefinedResponse<ResponseType>;
            tags: { [key: string]: string };
        } => {
            tags = {play: playName, ...tags};

            const response = api.dav.folderCreate({
                credential: credential as types.Credential,
                path,
                userName,
                tags,
            });

            check(response, {
                'folder create status is 201': () => response.status === 201,
            }, tags) || metricErrorRate.add(1, tags);

            metricTrend.add(response.timings.duration, tags)

            return {
                response,
                tags,
            }
        }
    }
};

export const folderDelete = ({name, metricID = 'default'}: { name?: string; metricID?: string; }) => {
    const playName = name || `oc_${metricID}_play_dav_folder_delete`;
    const metricTrendName = `${playName}_trend`;
    const metricTrend = new Trend(metricTrendName, true);
    const metricErrorRateName = `${playName}_error_rate`;
    const metricErrorRate = new Gauge(metricErrorRateName);

    return {
        playName,
        metricTrendName,
        metricErrorRateName,
        exec: (
            {
                credential,
                userName,
                path,
                tags,
            }: {
                credential: types.Credential;
                userName: string;
                path: string;
                tags?: { [key: string]: string };
            }
        ): {
            response: RefinedResponse<ResponseType>;
            tags: { [key: string]: string };
        } => {
            tags = {play: playName, ...tags};

            const response = api.dav.folderDelete({
                credential: credential as types.Credential,
                path,
                userName,
                tags,
            });

            check(response, {
                'folder delete status is 204': () => response.status === 204,
            }, tags) || metricErrorRate.add(1, tags);

            metricTrend.add(response.timings.duration, tags)

            return {
                response,
                tags,
            }
        }
    }
};

export const propfind = ({name, metricID = 'default'}: { name?: string; metricID?: string; }) => {
    const playName = name || `oc_${metricID}_play_dav_propfind`;
    const metricTrendName = `${playName}_trend`;
    const metricTrend = new Trend(metricTrendName, true);
    const metricErrorRateName = `${playName}_error_rate`;
    const metricErrorRate = new Gauge(metricErrorRateName);

    return {
        playName,
        metricTrendName,
        metricErrorRateName,
        exec: (
            {
                credential,
                userName,
                path,
                tags,
            }: {
                credential: types.Credential;
                userName: string;
                path?: string;
                tags?: { [key: string]: string };
            }
        ): {
            response: RefinedResponse<ResponseType>;
            tags: { [key: string]: string };
        } => {
            tags = {play: playName, ...tags};

            const response = api.dav.propfind({
                credential: credential as types.Credential,
                path,
                userName,
                tags,
            });

            check(response, {
                'propfind status is 207': () => response.status === 207,
            }, tags) || metricErrorRate.add(1, tags);

            metricTrend.add(response.timings.duration, tags)

            return {
                response,
                tags,
            }
        }
    }
};