import {sleep, check} from 'k6';
import {Options} from "k6/options";
import {sample} from 'lodash';
import {defaults, tasks} from "./lib";
import {getFile} from "./lib/utils";

export let options: Options = {
    insecureSkipTLSVerify: true,
    vus: 10,
    duration: '10s',
};

export const setup = () => {
    return {}
}


export default () => {
    const file = getFile();
    const account = sample(defaults.accounts);

    const userInfoResponse = tasks.userInfo(account)
    check(userInfoResponse, {
        'status is 200': () => userInfoResponse.status === 200,
    });

    const uploadResponse = tasks.uploadFile(file.bytes, file.name, account)
    check(uploadResponse, {
        'status is 201': () => uploadResponse.status === 201,
    });

    sleep(1);
};