import {sleep, check} from 'k6';
import {Options} from "k6/options";
import {api, users} from "./lib";

const files = {
    'mb_100.jpg': open('./_files/mb_100.zip', 'b'),
}

export let options: Options = {
    insecureSkipTLSVerify: true,
    iterations: 20,
    vus: 20,
};

export const setup = (): void => {
    users.createBatchUsers(options.vus)
}

export default () => {
    const fileName = `mb_100-${__ITER}.jpg`
    const res = api.uploadFile(users.k6VirtualUser(__VU), files['mb_100.zip'], fileName)

    check(res, {
        'status is 201': () => res.status === 201,
    });
    sleep(1);
};

export const teardown = (): void => {
  console.log("teardown")
  for (let vu=1; vu <= options.vus; vu++) {
    for (let iter=0; iter < Math.floor(options.iterations/options.vus); iter++) {
        const res = api.deleteFile(users.k6VirtualUser(vu), `mb_100-${iter}.zip`)
        check(res, {
            // status could be either 204(if file was uploaded) of 404(if file was not uploaded)
            'status is 204': () => [204, 404].includes(res.status),
        });
    }
  } 
  users.deleteBatchUsers(options.vus)
}
