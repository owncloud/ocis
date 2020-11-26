import {sleep, check} from 'k6';
import {Options} from "k6/options";
import {defaults, api} from "./lib";

const files = {
    'kb_50.jpg': open('./_files/kb_50.jpg', 'b'),
}

export let options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 100,
    vus: 100,
};

export default () => {
    const res = api.uploadFile(defaults.accounts.einstein, files['kb_50.jpg'], `kb_50-${__VU}-${__ITER}.jpg`)

    check(res, {
        'status is 201': () => res.status === 201,
    });

    sleep(1);
};