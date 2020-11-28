import { types } from '@babel/core';
import {sleep, check} from 'k6';
import {Options} from "k6/options";
import {api, users, defaults} from "./lib";

const files = {
    'kb_50.jpg': open('./_files/kb_50.jpg', 'b'),
}

export let options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 100,
    vus: 100,
};

export const setup = (): void => {
    console.log("hello world")
    users.createBatchUsers(options.vus)
}

export default () => {
    const fileName = `kb_50-${__ITER}.jpg`
    const res = api.uploadFile(users.k6VirtualUser(__VU), files['kb_50.jpg'], fileName)

    check(res, {
        'status is 201': () => res.status === 201,
    });
    sleep(1);
};

export const teardown = (): void => {
  console.log("teardown")
  for (let vu=1; vu <= options.vus; vu++) {
    for (let iter=0; iter < Math.floor(options.iterations/options.vus); iter++) {
        const res = api.deleteFile(users.k6VirtualUser(vu), `kb_50-${iter}.jpg`)
        check(res, {
            // status could be either 204(if file was uploaded) of 404(if file was not uploaded)
            'status is 204': () => [204, 404].includes(res.status),
        });
    }
  } 
  users.deleteBatchUsers(options.vus)
}
