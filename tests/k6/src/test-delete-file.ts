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
    console.log("setup for delete")
    for (let vu=1; vu <= options.vus; vu++) {
      for (let iter=0; iter < options.iterations/options.vus; iter++) {
        const res = api.uploadFile(defaults.accounts.einstein, files['kb_50.jpg'], `deletefile_${vu}_${iter}.jpg`)
        check(res, {
          'status is 201': () => res.status === 201,
        })
      }
    }
    console.log("uploaded test files")
    sleep(1)
}

export default () => {
    const res = api.deleteFile(defaults.accounts.einstein, `deletefile_${__VU}_${__ITER}.jpg`)
    check(res, {
        'status is 204': () => res.status === 204,
    });
    sleep(1);
};

// this needs to be done because some delete requests fail and the file still remain in the server
export const teardown = (): void => {
    console.log("teardown for delete")
    for (let vu=1; vu <= options.vus; vu++) {
      for (let iter=0; iter < options.iterations/options.vus; iter++) {
        const res = api.deleteFile(defaults.accounts.einstein, `deletefile_${vu}_${iter}.jpg`)
        check(res, {
          'status is 204 or 404': () => [204, 404].includes(res.status),
        })
      }
    }
    console.log("cleaned up the files")
}
