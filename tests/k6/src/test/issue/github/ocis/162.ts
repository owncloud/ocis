import {Options} from 'k6/options';
import * as uploadFilesBenchmark from '../../../benchmark/file-upload'

export const options: Options = {
    ...uploadFilesBenchmark.options,
    iterations: 200,
    vus: 50,
};
export const {setup} = uploadFilesBenchmark;
export default uploadFilesBenchmark.default;