// @ts-nocheck
/* tslint:disable */
/* eslint-disable */
/**
 * GitClassrooms – Backend API
 * This is the API for our Gitlab Classroom Webapp
 *
 * OpenAPI spec version: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 * Do not edit the class manually.
 */

import globalAxios, { AxiosResponse, AxiosInstance, AxiosRequestConfig } from 'axios';
import { Configuration } from '../configuration';
// Some imports not used depending on template conditions
// @ts-ignore
import { BASE_PATH, COLLECTION_FORMATS, RequestArgs, BaseAPI, RequiredError } from '../base';
import { CreateTeamRequest } from '../models';
import { HTTPError } from '../models';
import { TeamResponse } from '../models';
import { UpdateTeamRequest } from '../models';
/**
 * TeamApi - axios parameter creator
 * @export
 */
export const TeamApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * Create a new Team for the given classroom and join it if you are a student
         * @summary Create new Team
         * @param {CreateTeamRequest} body Classroom Info
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        createTeam: async (body: CreateTeamRequest, xCsrfToken: string, classroomId: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'body' is not null or undefined
            if (body === null || body === undefined) {
                throw new RequiredError('body','Required parameter body was null or undefined when calling createTeam.');
            }
            // verify required parameter 'xCsrfToken' is not null or undefined
            if (xCsrfToken === null || xCsrfToken === undefined) {
                throw new RequiredError('xCsrfToken','Required parameter xCsrfToken was null or undefined when calling createTeam.');
            }
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling createTeam.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/teams`
                .replace(`{${"classroomId"}}`, encodeURIComponent(String(classroomId)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }
            const localVarRequestOptions :AxiosRequestConfig = { method: 'POST', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            if (xCsrfToken !== undefined && xCsrfToken !== null) {
                localVarHeaderParameter['X-Csrf-Token'] = String(xCsrfToken);
            }

            localVarHeaderParameter['Content-Type'] = 'application/json';

            const query = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                query.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.params) {
                query.set(key, options.params[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(query)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            const needsSerialization = (typeof body !== "string") || localVarRequestOptions.headers['Content-Type'] === 'application/json';
            localVarRequestOptions.data =  needsSerialization ? JSON.stringify(body !== undefined ? body : {}) : (body || "");

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * GetClassroomTeam
         * @summary GetClassroomTeam
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getClassroomTeam: async (classroomId: string, teamId: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling getClassroomTeam.');
            }
            // verify required parameter 'teamId' is not null or undefined
            if (teamId === null || teamId === undefined) {
                throw new RequiredError('teamId','Required parameter teamId was null or undefined when calling getClassroomTeam.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/teams/{teamId}`
                .replace(`{${"classroomId"}}`, encodeURIComponent(String(classroomId)))
                .replace(`{${"teamId"}}`, encodeURIComponent(String(teamId)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }
            const localVarRequestOptions :AxiosRequestConfig = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            const query = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                query.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.params) {
                query.set(key, options.params[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(query)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * GetClassroomTeams
         * @summary GetClassroomTeams
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getClassroomTeams: async (classroomId: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling getClassroomTeams.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/teams`
                .replace(`{${"classroomId"}}`, encodeURIComponent(String(classroomId)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }
            const localVarRequestOptions :AxiosRequestConfig = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            const query = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                query.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.params) {
                query.set(key, options.params[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(query)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * Join the Team if we aren't in another team
         * @summary Join the team
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {string} xCsrfToken Csrf-Token
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        joinTeam: async (classroomId: string, teamId: string, xCsrfToken: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling joinTeam.');
            }
            // verify required parameter 'teamId' is not null or undefined
            if (teamId === null || teamId === undefined) {
                throw new RequiredError('teamId','Required parameter teamId was null or undefined when calling joinTeam.');
            }
            // verify required parameter 'xCsrfToken' is not null or undefined
            if (xCsrfToken === null || xCsrfToken === undefined) {
                throw new RequiredError('xCsrfToken','Required parameter xCsrfToken was null or undefined when calling joinTeam.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/teams/{teamId}/join`
                .replace(`{${"classroomId"}}`, encodeURIComponent(String(classroomId)))
                .replace(`{${"teamId"}}`, encodeURIComponent(String(teamId)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }
            const localVarRequestOptions :AxiosRequestConfig = { method: 'POST', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            if (xCsrfToken !== undefined && xCsrfToken !== null) {
                localVarHeaderParameter['X-Csrf-Token'] = String(xCsrfToken);
            }

            const query = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                query.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.params) {
                query.set(key, options.params[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(query)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * Update Team
         * @summary Update Team
         * @param {UpdateTeamRequest} body Update Team
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        updateTeam: async (body: UpdateTeamRequest, xCsrfToken: string, classroomId: string, teamId: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'body' is not null or undefined
            if (body === null || body === undefined) {
                throw new RequiredError('body','Required parameter body was null or undefined when calling updateTeam.');
            }
            // verify required parameter 'xCsrfToken' is not null or undefined
            if (xCsrfToken === null || xCsrfToken === undefined) {
                throw new RequiredError('xCsrfToken','Required parameter xCsrfToken was null or undefined when calling updateTeam.');
            }
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling updateTeam.');
            }
            // verify required parameter 'teamId' is not null or undefined
            if (teamId === null || teamId === undefined) {
                throw new RequiredError('teamId','Required parameter teamId was null or undefined when calling updateTeam.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/teams/{teamId}`
                .replace(`{${"classroomId"}}`, encodeURIComponent(String(classroomId)))
                .replace(`{${"teamId"}}`, encodeURIComponent(String(teamId)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }
            const localVarRequestOptions :AxiosRequestConfig = { method: 'PUT', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            if (xCsrfToken !== undefined && xCsrfToken !== null) {
                localVarHeaderParameter['X-Csrf-Token'] = String(xCsrfToken);
            }

            localVarHeaderParameter['Content-Type'] = 'application/json';

            const query = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                query.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.params) {
                query.set(key, options.params[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(query)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            const needsSerialization = (typeof body !== "string") || localVarRequestOptions.headers['Content-Type'] === 'application/json';
            localVarRequestOptions.data =  needsSerialization ? JSON.stringify(body !== undefined ? body : {}) : (body || "");

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
    }
};

/**
 * TeamApi - functional programming interface
 * @export
 */
export const TeamApiFp = function(configuration?: Configuration) {
    return {
        /**
         * Create a new Team for the given classroom and join it if you are a student
         * @summary Create new Team
         * @param {CreateTeamRequest} body Classroom Info
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async createTeam(body: CreateTeamRequest, xCsrfToken: string, classroomId: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<void>>> {
            const localVarAxiosArgs = await TeamApiAxiosParamCreator(configuration).createTeam(body, xCsrfToken, classroomId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * GetClassroomTeam
         * @summary GetClassroomTeam
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomTeam(classroomId: string, teamId: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<TeamResponse>>> {
            const localVarAxiosArgs = await TeamApiAxiosParamCreator(configuration).getClassroomTeam(classroomId, teamId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * GetClassroomTeams
         * @summary GetClassroomTeams
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomTeams(classroomId: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<Array<TeamResponse>>>> {
            const localVarAxiosArgs = await TeamApiAxiosParamCreator(configuration).getClassroomTeams(classroomId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * Join the Team if we aren't in another team
         * @summary Join the team
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {string} xCsrfToken Csrf-Token
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async joinTeam(classroomId: string, teamId: string, xCsrfToken: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<void>>> {
            const localVarAxiosArgs = await TeamApiAxiosParamCreator(configuration).joinTeam(classroomId, teamId, xCsrfToken, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * Update Team
         * @summary Update Team
         * @param {UpdateTeamRequest} body Update Team
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async updateTeam(body: UpdateTeamRequest, xCsrfToken: string, classroomId: string, teamId: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<void>>> {
            const localVarAxiosArgs = await TeamApiAxiosParamCreator(configuration).updateTeam(body, xCsrfToken, classroomId, teamId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
    }
};

/**
 * TeamApi - factory interface
 * @export
 */
export const TeamApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    return {
        /**
         * Create a new Team for the given classroom and join it if you are a student
         * @summary Create new Team
         * @param {CreateTeamRequest} body Classroom Info
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async createTeam(body: CreateTeamRequest, xCsrfToken: string, classroomId: string, options?: AxiosRequestConfig): Promise<AxiosResponse<void>> {
            return TeamApiFp(configuration).createTeam(body, xCsrfToken, classroomId, options).then((request) => request(axios, basePath));
        },
        /**
         * GetClassroomTeam
         * @summary GetClassroomTeam
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomTeam(classroomId: string, teamId: string, options?: AxiosRequestConfig): Promise<AxiosResponse<TeamResponse>> {
            return TeamApiFp(configuration).getClassroomTeam(classroomId, teamId, options).then((request) => request(axios, basePath));
        },
        /**
         * GetClassroomTeams
         * @summary GetClassroomTeams
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomTeams(classroomId: string, options?: AxiosRequestConfig): Promise<AxiosResponse<Array<TeamResponse>>> {
            return TeamApiFp(configuration).getClassroomTeams(classroomId, options).then((request) => request(axios, basePath));
        },
        /**
         * Join the Team if we aren't in another team
         * @summary Join the team
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {string} xCsrfToken Csrf-Token
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async joinTeam(classroomId: string, teamId: string, xCsrfToken: string, options?: AxiosRequestConfig): Promise<AxiosResponse<void>> {
            return TeamApiFp(configuration).joinTeam(classroomId, teamId, xCsrfToken, options).then((request) => request(axios, basePath));
        },
        /**
         * Update Team
         * @summary Update Team
         * @param {UpdateTeamRequest} body Update Team
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async updateTeam(body: UpdateTeamRequest, xCsrfToken: string, classroomId: string, teamId: string, options?: AxiosRequestConfig): Promise<AxiosResponse<void>> {
            return TeamApiFp(configuration).updateTeam(body, xCsrfToken, classroomId, teamId, options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * TeamApi - object-oriented interface
 * @export
 * @class TeamApi
 * @extends {BaseAPI}
 */
export class TeamApi extends BaseAPI {
    /**
     * Create a new Team for the given classroom and join it if you are a student
     * @summary Create new Team
     * @param {CreateTeamRequest} body Classroom Info
     * @param {string} xCsrfToken Csrf-Token
     * @param {string} classroomId Classroom ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof TeamApi
     */
    public async createTeam(body: CreateTeamRequest, xCsrfToken: string, classroomId: string, options?: AxiosRequestConfig) : Promise<AxiosResponse<void>> {
        return TeamApiFp(this.configuration).createTeam(body, xCsrfToken, classroomId, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * GetClassroomTeam
     * @summary GetClassroomTeam
     * @param {string} classroomId Classroom ID
     * @param {string} teamId Team ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof TeamApi
     */
    public async getClassroomTeam(classroomId: string, teamId: string, options?: AxiosRequestConfig) : Promise<AxiosResponse<TeamResponse>> {
        return TeamApiFp(this.configuration).getClassroomTeam(classroomId, teamId, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * GetClassroomTeams
     * @summary GetClassroomTeams
     * @param {string} classroomId Classroom ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof TeamApi
     */
    public async getClassroomTeams(classroomId: string, options?: AxiosRequestConfig) : Promise<AxiosResponse<Array<TeamResponse>>> {
        return TeamApiFp(this.configuration).getClassroomTeams(classroomId, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * Join the Team if we aren't in another team
     * @summary Join the team
     * @param {string} classroomId Classroom ID
     * @param {string} teamId Team ID
     * @param {string} xCsrfToken Csrf-Token
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof TeamApi
     */
    public async joinTeam(classroomId: string, teamId: string, xCsrfToken: string, options?: AxiosRequestConfig) : Promise<AxiosResponse<void>> {
        return TeamApiFp(this.configuration).joinTeam(classroomId, teamId, xCsrfToken, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * Update Team
     * @summary Update Team
     * @param {UpdateTeamRequest} body Update Team
     * @param {string} xCsrfToken Csrf-Token
     * @param {string} classroomId Classroom ID
     * @param {string} teamId Team ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof TeamApi
     */
    public async updateTeam(body: UpdateTeamRequest, xCsrfToken: string, classroomId: string, teamId: string, options?: AxiosRequestConfig) : Promise<AxiosResponse<void>> {
        return TeamApiFp(this.configuration).updateTeam(body, xCsrfToken, classroomId, teamId, options).then((request) => request(this.axios, this.basePath));
    }
}
