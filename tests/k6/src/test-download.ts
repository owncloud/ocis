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

export const setup = (): void => {
    console.log("setup for download")
    const res = api.uploadFile(defaults.accounts.einstein, files['kb_50.jpg'], 'downloadfile.jpg')
    check(res, {
      'status is 201': () => res.status === 201,
    })
    console.log("uploaded file")
    sleep(1)
}

export default () => {
    const res = api.downloadFile(defaults.accounts.einstein, 'downloadfile.jpg')
    check(res, {
        'status is 200': () => res.status === 200,
    });
    sleep(1);
};

export const teardown = (): void => {
    console.log("teardown for download")
    const res = api.deleteFile(defaults.accounts.einstein, 'downloadfile.jpg')
    check(res, {
      'status is 204': () => res.status === 204,
    })
    console.log("deleted file")
}

