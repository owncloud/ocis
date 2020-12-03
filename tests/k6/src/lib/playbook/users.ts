import * as api from '../api';
import { check } from 'k6';
import * as types from '../types';
import { RefinedResponse, ResponseType } from 'k6/http';
import { Play } from './playbook';

export class Create extends Play {
    constructor({ name, metricID = 'default' }: { name?: string; metricID?: string } = {}) {
        super({ name: name || `oc_${metricID}_play_users_create` });
    }

    public exec({
        credential,
        userName,
        password,
        email,
        tags,
    }: {
        credential: types.Credential;
        userName: string;
        password: string;
        email: string;
        tags?: types.Tags;
    }): { response: RefinedResponse<ResponseType>; tags: types.Tags } {
        tags = { ...this.tags, ...tags };

        const response = api.users.Create.exec({ credential: credential, userName, password, tags, email });

        check(
            response,
            {
                'users create status is 200': () => response.status === 200,
            },
            tags,
        ) || this.metricErrorRate.add(1, tags);

        this.metricTrend.add(response.timings.duration, tags);

        return { response, tags };
    }
}

export class Delete extends Play {
    constructor({ name, metricID = 'default' }: { name?: string; metricID?: string } = {}) {
        super({ name: name || `oc_${metricID}_play_users_delete` });
    }

    public exec({
        credential,
        userName,
        tags,
    }: {
        credential: types.Credential;
        userName: string;
        tags?: types.Tags;
    }): { response: RefinedResponse<ResponseType>; tags: types.Tags } {
        tags = { ...this.tags, ...tags };

        const response = api.users.Delete.exec({ credential: credential, userName, tags });

        check(
            response,
            {
                'users delete status is 200': () => response.status === 200,
            },
            tags,
        ) || this.metricErrorRate.add(1, tags);

        this.metricTrend.add(response.timings.duration, tags);

        return { response, tags };
    }
}
