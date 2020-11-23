import {sleep, check} from 'k6';
import {Options} from "k6/options";
import {defaults, api, utils} from "./lib";

export let options: Options = {
    insecureSkipTLSVerify: true,
};

export default () => {
    Object.keys(defaults.files).forEach(name => {
        const file = utils.getFile(name);
        const uploadResponse = api.uploadFile(utils.getAccount(), file.bytes, file.nameRandom)

        check(uploadResponse, {
            'status is 201': () => uploadResponse.status === 201,
        });
    })

    sleep(1);
};