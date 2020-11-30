import {defaults, playbook} from '../../lib'
import {Options} from 'k6/options';
import {sleep} from "k6";
import auth from "../../lib/auth";

export const options: Options = {
    ...defaults.K6.OPTIONS,
};
const authFactory = new auth(defaults.ACCOUNT.for(defaults.ACCOUNT.EINSTEIN));
const plays = {
    fileUpload: playbook.dav.fileUpload(),
    fileDownload: playbook.dav.fileDownload(),
    fileDelete: playbook.dav.fileDelete(),
}
export default () => {
    const {login: userName} = authFactory.account;
    const fileName = plays.fileUpload({
        credential: authFactory.credential,
        userName,
        asset: defaults.FILE,
    });

    sleep(1)

    plays.fileDownload({
        credential: authFactory.credential,
        userName,
        fileName,
    });

    sleep(1)

    plays.fileDelete({
        credential: authFactory.credential,
        userName,
        fileName,
    });

    sleep(1)
};