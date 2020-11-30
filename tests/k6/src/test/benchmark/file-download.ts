import {defaults, playbook} from '../../lib'
import {Options} from 'k6/options';
import {sleep} from "k6";
import auth from "../../lib/auth";

export const options: Options = {
    ...defaults.K6_OPTION_DEFAULTS,
    iterations: 1,
    vus: 1,
};
const authFactory = new auth(defaults.ACCOUNTS.for(defaults.ACCOUNTS.EINSTEIN));
const playbooks = {
    fileUpload: playbook.dav.fileUpload(),
    fileDownload: playbook.dav.fileDownload(),
    fileDelete: playbook.dav.fileDelete(),
}

export default () => {
    const {login: userName} = authFactory.account;
    const fileName = playbooks.fileUpload({
        credential: authFactory.credential,
        userName,
        asset: defaults.OC_TEST_FILE
    });

    sleep(1)

    playbooks.fileDownload({
        credential: authFactory.credential,
        userName,
        fileName,
    });

    sleep(1)

    playbooks.fileDelete({
        credential: authFactory.credential,
        userName,
        fileName,
    });

    sleep(1)
};