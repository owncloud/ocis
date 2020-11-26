import {sleep, check} from 'k6';
import {Options} from "k6/options";
import {api} from "./lib";
import {users} from "./lib"

const files = {
    'kb_50.jpg': open('./_files/kb_50.jpg', 'b'),
}

export let options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 100,
    vus: 100,
};

export const setup = (): void => {
    users.createBatchUsers(options.vus)   
}

export default () => {
    // Upload 6 50kb files
    for (let i = 1; i <= 6; i++) {
        const res = api.uploadFile(users.k6VirtualUser(__VU), files['kb_50.jpg'], `downloadfile_${i}.jpg`)
        check(res, {
            'status is 200': () => res.status === 200,
        });
    }

    // Download all the uploaded files
    for (let i = 1; i <= 6; i++) {
        const res = api.uploadFile(users.k6VirtualUser(__VU), files['kb_50.jpg'], `downloadfile_${i}.jpg`)
        check(res, {
            'status is 200': () => res.status === 200,
        });
    }

    const res = api.downloadFile(users.k6VirtualUser(__VU), 'downloadfile.jpg')
    check(res, {
        'status is 200': () => res.status === 200,
    });
    sleep(1);
};

export const teardown = (): void => {
    console.log("teardown for download")
    for (let vu=1; vu <= options.vus; vu++) {
        const res = api.deleteFile(users.k6VirtualUser(vu), 'downloadfile.jpg')
        if (res.status !== 204) {
            fail("Status is not 204 while deleting the file")
        }
    }
    console.log("deleted files")
    users.deleteBatchUsers(options.vus)
}

