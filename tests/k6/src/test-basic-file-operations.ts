import {sleep, check, fail} from 'k6';
import {Options} from "k6/options";
import {api} from "./lib";
import {users} from "./lib"

const files = {
    'kb_50.jpg': open('./_files/kb_50.jpg', 'b'),
}

export let options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 300,
    vus: 100,
    setupTimeout: '200s',
};

export const setup = (): void => {
    users.createBatchUsers(options.vus)   
}

export default () => {
    // Upload 6 50kb files
    for (let i = 1; i <= 6; i++) {
        const res = api.uploadFile(users.k6VirtualUser(__VU), files['kb_50.jpg'], `downloadfile_${i}.jpg`)
        check(res, {
            'status is 201': () => res.status === 201,
        });
    }

    // Download all the uploaded files
    for (let i = 1; i <= 6; i++) {
        const res = api.downloadFile(users.k6VirtualUser(__VU), `downloadfile_${i}.jpg`)
        check(res, {
            'status is 200': () => res.status === 200,
        });
    }

    // Delete all the uploaded files
    for (let i = 1; i <= 6; i++) {
        const res = api.deleteFile(users.k6VirtualUser(__VU), `downloadfile_${i}.jpg`)
        check(res, {
            'status is 204': () => res.status === 204  ,
        });
    }
    sleep(1);
};

export const teardown = (): void => {
    console.log("teardown step")
    users.deleteBatchUsers(options.vus)
}

