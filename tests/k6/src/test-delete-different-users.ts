import {sleep, check,fail} from 'k6';
import {Options} from "k6/options";
import {api, utils, users} from "./lib";

const files = {
    'kb_50.jpg': open('./_files/kb_50.jpg', 'b'),
}

export let options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 100,
    vus: 100,
};

export const setup = (): string[][] => {
    users.createBatchUsers(options.vus)
    const uploadedFiles: string[][] = []
    for (let vu=1; vu <= options.vus; vu++) {
        uploadedFiles[vu] = []
        const user = users.k6VirtualUser(vu)
        for (let iter=0; iter < options.iterations/options.vus; iter++) {
            const fileName = `deletefile_${utils.randomString()}.jpg`
            const res = api.uploadFile(user, files['kb_50.jpg'], fileName)
            if (res.status !== 201) {
                fail(`status code was *not* 201, expected: 201, actual: ${res.status}`);
            }
            uploadedFiles[vu].push(fileName)
        }
    }
    return uploadedFiles
}

export default (files: string[]) => {
    const fileName: string = files[__VU][__ITER]

    const res = api.deleteFile(users.k6VirtualUser(__VU), fileName)
    check(res, {
        'status is 204': () => res.status === 204,
    })
    sleep(1)
}

export const teardown = (): void => {
    users.deleteBatchUsers(options.vus)
}
