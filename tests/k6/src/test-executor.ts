import {sleep, check} from 'k6';
import {Options} from "k6/options";
import {defaults, api, utils} from "./lib";

export let options: Options = {
  insecureSkipTLSVerify: true,
  iterations: 2,
  vus: 1,
};

export const setup = (): string[] => {
  return ["hello", "world", "test"]
}

export default (data: string[]) => {
  data.pop()
};

// this needs to be done because some delete requests fail and the file still remain in the server
export const teardown = (files: string[]): void => {
  console.log(files)
}
