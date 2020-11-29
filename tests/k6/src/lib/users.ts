import encoding from 'k6/encoding';
import {fail} from "k6";
import http, { file } from "k6/http";
import * as defaults from "./defaults";
import * as types from "./types";
import { uploadFile } from './api';

const adminUser = defaults.accounts.admin

const files = {
    'kb_50.jpg': open('../_files/kb_50.jpg', 'b'),
}

export const k6VirtualUser = (userNumber: number): types.Account => {
    return {
        login: `k6users_${userNumber}`,
        password: `k6users_${userNumber}`,
    }
}

export const createBatchUsers = (numberOfUsers: number): void => {
    deleteBatchUsers(numberOfUsers, false)
    console.log("Cleared all existing users")
    console.log("Creating demo users...")
    for (let i=1; i <= numberOfUsers; i++ ) {
        const userName = `k6users_${i}`
        const userData: types.UserRequestData = {
            displayname: userName,
            email: `${userName}@example.com`,
            password: userName,
            userid: userName
        }
        const res = http.post(
            `${defaults.host.name}/ocs/v2.php/cloud/users`,
            userData as any,
            {
                headers: {
                    Authorization: `Basic ${encoding.b64encode(`${adminUser.login}:${adminUser.password}`)}`,
                }
            } 
        );
        if (res.status != 200) {
            fail("Failed while creating user")
        }
        console.log(`Created user ${userName}, initializing user now...`)

        // running on oc10 will also require users to be initialized
        uploadFile(k6VirtualUser(i), files['kb_50.jpg'], 'initialfile.txt')
    }
}

export const deleteBatchUsers = (numberOfUsers: number, stopOnFailure: boolean = true): void => {
    console.log("Clearing all demo users...")
    for (let i=1; i <= numberOfUsers; i++ ) {
        const userName = `k6users_${i}`
        const res = http.del(
            `${defaults.host.name}/ocs/v2.php/cloud/users/${userName}`,
            {},
            {
                headers: {
                    Authorization: `Basic ${encoding.b64encode(`${adminUser.login}:${adminUser.password}`)}`,
                }
            }
        );
        if (stopOnFailure && res.status != 200) {
            fail("Failed while deleting created user")
        }
    }
}