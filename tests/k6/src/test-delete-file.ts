import {sleep, check,fail} from 'k6';
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

export const setup = (): string[] => {
    const uploadedFiles: string[] = []
    console.log("setup for delete")
    for (let iter=0; iter < options.iterations; iter++) {
      const fileName = `deletefile_${iter}.jpg`
      const res = api.uploadFile(defaults.accounts.einstein, files['kb_50.jpg'], fileName)
      if (!check(res, {
        'status code MUST be 201': (res) => res.status == 201,
      })) {
        fail(`status code was *not* 201, expected: 201, actual: ${res.status}`);
      }
      uploadedFiles.push(fileName)
    }
    console.log("uploaded test files")
    return uploadedFiles
}

export default (files: string[]) => {
    const fileName: string = files[(__ITER * options.vus) + __VU  - 1]

    const res = api.deleteFile(defaults.accounts.einstein, fileName)
    check(res, {
        'status is 204': () => res.status === 204,
    })
    sleep(1);
};
