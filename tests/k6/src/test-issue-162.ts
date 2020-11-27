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
    const fileName = `kb_50-${(__ITER * options.vus) + __VU  - 1}.jpg`
    const res = api.uploadFile(defaults.accounts.einstein, files['kb_50.jpg'], fileName)

    check(res, {
        'status is 204': () => res.status === 204,
    });

    sleep(1);
};

export const teardown = (): void => {
  console.log("teardown")
  for (let iter=0; iter < options.iterations; iter++) {
    const res = api.deleteFile(defaults.accounts.einstein, `kb_50-${iter}.jpg`)
    check(res, {
      // status could be either 204(if file was uploaded) of 404(if file was not uploaded)
      'status is 204 or 404': () => [204, 404].includes(res.status),
    });
  }
}
