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
    const fileName = `kb_50-${__VU}-${__ITER}.jpg`
    const res = api.uploadFile(defaults.accounts.einstein, files['kb_50.jpg'], fileName)

    check(res, {
        'status is 204': () => res.status === 204,
    });

    sleep(1);
};

export const teardown = (): void => {
  console.log("teardown")
  for (let vu=1; vu <= options.vus; vu++) {
    for (let iter=0; iter < options.iterations/options.vus; iter++) {
      const res = api.deleteFile(defaults.accounts.einstein, `kb_50-${vu}-${iter}.jpg`)
      check(res, {
          // status could be either 204(if file was created) of 404(if file was not uploaded)
          'status is 204 or 404': () => [204, 404].includes(res.status),
      });
    }
  }
}
