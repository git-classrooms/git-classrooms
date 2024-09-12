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
import { HTTPError } from '../models';
import { UpdateMemberRoleRequest } from '../models';
import { UpdateMemberTeamRequest } from '../models';
import { UserClassroomResponse } from '../models';
/**
 * MemberApi - axios parameter creator
 * @export
 */
export const MemberApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * GetClassroomMember
         * @summary GetClassroomMember
         * @param {string} classroomId Classroom ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getClassroomMember: async (classroomId: string, memberId: number, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling getClassroomMember.');
            }
            // verify required parameter 'memberId' is not null or undefined
            if (memberId === null || memberId === undefined) {
                throw new RequiredError('memberId','Required parameter memberId was null or undefined when calling getClassroomMember.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/members/{memberId}`
                .replace(`{${"classroomId"}}`, encodeURIComponent(String(classroomId)))
                .replace(`{${"memberId"}}`, encodeURIComponent(String(memberId)));
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
         * GetClassroomMembers
         * @summary GetClassroomMembers
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getClassroomMembers: async (classroomId: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling getClassroomMembers.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/members`
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
         * GetClassroomTeamMember
         * @summary GetClassroomTeamMember
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getClassroomTeamMember: async (classroomId: string, teamId: string, memberId: number, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling getClassroomTeamMember.');
            }
            // verify required parameter 'teamId' is not null or undefined
            if (teamId === null || teamId === undefined) {
                throw new RequiredError('teamId','Required parameter teamId was null or undefined when calling getClassroomTeamMember.');
            }
            // verify required parameter 'memberId' is not null or undefined
            if (memberId === null || memberId === undefined) {
                throw new RequiredError('memberId','Required parameter memberId was null or undefined when calling getClassroomTeamMember.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/teams/{teamId}/members/{memberId}`
                .replace(`{${"classroomId"}}`, encodeURIComponent(String(classroomId)))
                .replace(`{${"teamId"}}`, encodeURIComponent(String(teamId)))
                .replace(`{${"memberId"}}`, encodeURIComponent(String(memberId)));
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
         * GetClassroomTeamMembers
         * @summary GetClassroomTeamMembers
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getClassroomTeamMembers: async (classroomId: string, teamId: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling getClassroomTeamMembers.');
            }
            // verify required parameter 'teamId' is not null or undefined
            if (teamId === null || teamId === undefined) {
                throw new RequiredError('teamId','Required parameter teamId was null or undefined when calling getClassroomTeamMembers.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/teams/{teamId}/members`
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
         * Remove Member from the team
         * @summary Remove Member from the team
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {number} memberId Member ID
         * @param {string} xCsrfToken Csrf-Token
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        removeMemberFromTeamV2: async (classroomId: string, teamId: string, memberId: number, xCsrfToken: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling removeMemberFromTeamV2.');
            }
            // verify required parameter 'teamId' is not null or undefined
            if (teamId === null || teamId === undefined) {
                throw new RequiredError('teamId','Required parameter teamId was null or undefined when calling removeMemberFromTeamV2.');
            }
            // verify required parameter 'memberId' is not null or undefined
            if (memberId === null || memberId === undefined) {
                throw new RequiredError('memberId','Required parameter memberId was null or undefined when calling removeMemberFromTeamV2.');
            }
            // verify required parameter 'xCsrfToken' is not null or undefined
            if (xCsrfToken === null || xCsrfToken === undefined) {
                throw new RequiredError('xCsrfToken','Required parameter xCsrfToken was null or undefined when calling removeMemberFromTeamV2.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/teams/{teamId}/members/{memberId}`
                .replace(`{${"classroomId"}}`, encodeURIComponent(String(classroomId)))
                .replace(`{${"teamId"}}`, encodeURIComponent(String(teamId)))
                .replace(`{${"memberId"}}`, encodeURIComponent(String(memberId)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }
            const localVarRequestOptions :AxiosRequestConfig = { method: 'DELETE', ...baseOptions, ...options};
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
         * Update Classroom Members role
         * @summary Update Classroom Members role
         * @param {UpdateMemberRoleRequest} body Update Member Role
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        updateMemberRole: async (body: UpdateMemberRoleRequest, xCsrfToken: string, classroomId: string, memberId: number, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'body' is not null or undefined
            if (body === null || body === undefined) {
                throw new RequiredError('body','Required parameter body was null or undefined when calling updateMemberRole.');
            }
            // verify required parameter 'xCsrfToken' is not null or undefined
            if (xCsrfToken === null || xCsrfToken === undefined) {
                throw new RequiredError('xCsrfToken','Required parameter xCsrfToken was null or undefined when calling updateMemberRole.');
            }
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling updateMemberRole.');
            }
            // verify required parameter 'memberId' is not null or undefined
            if (memberId === null || memberId === undefined) {
                throw new RequiredError('memberId','Required parameter memberId was null or undefined when calling updateMemberRole.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/members/{memberId}/role`
                .replace(`{${"classroomId"}}`, encodeURIComponent(String(classroomId)))
                .replace(`{${"memberId"}}`, encodeURIComponent(String(memberId)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }
            const localVarRequestOptions :AxiosRequestConfig = { method: 'PATCH', ...baseOptions, ...options};
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
         * Update Classroom Members team
         * @summary Update Classroom Members team
         * @param {UpdateMemberTeamRequest} body Update Member Team
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        updateMemberTeam: async (body: UpdateMemberTeamRequest, xCsrfToken: string, classroomId: string, memberId: number, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'body' is not null or undefined
            if (body === null || body === undefined) {
                throw new RequiredError('body','Required parameter body was null or undefined when calling updateMemberTeam.');
            }
            // verify required parameter 'xCsrfToken' is not null or undefined
            if (xCsrfToken === null || xCsrfToken === undefined) {
                throw new RequiredError('xCsrfToken','Required parameter xCsrfToken was null or undefined when calling updateMemberTeam.');
            }
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling updateMemberTeam.');
            }
            // verify required parameter 'memberId' is not null or undefined
            if (memberId === null || memberId === undefined) {
                throw new RequiredError('memberId','Required parameter memberId was null or undefined when calling updateMemberTeam.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/members/{memberId}/team`
                .replace(`{${"classroomId"}}`, encodeURIComponent(String(classroomId)))
                .replace(`{${"memberId"}}`, encodeURIComponent(String(memberId)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }
            const localVarRequestOptions :AxiosRequestConfig = { method: 'PATCH', ...baseOptions, ...options};
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
 * MemberApi - functional programming interface
 * @export
 */
export const MemberApiFp = function(configuration?: Configuration) {
    return {
        /**
         * GetClassroomMember
         * @summary GetClassroomMember
         * @param {string} classroomId Classroom ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomMember(classroomId: string, memberId: number, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<UserClassroomResponse>>> {
            const localVarAxiosArgs = await MemberApiAxiosParamCreator(configuration).getClassroomMember(classroomId, memberId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * GetClassroomMembers
         * @summary GetClassroomMembers
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomMembers(classroomId: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<Array<UserClassroomResponse>>>> {
            const localVarAxiosArgs = await MemberApiAxiosParamCreator(configuration).getClassroomMembers(classroomId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * GetClassroomTeamMember
         * @summary GetClassroomTeamMember
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomTeamMember(classroomId: string, teamId: string, memberId: number, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<UserClassroomResponse>>> {
            const localVarAxiosArgs = await MemberApiAxiosParamCreator(configuration).getClassroomTeamMember(classroomId, teamId, memberId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * GetClassroomTeamMembers
         * @summary GetClassroomTeamMembers
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomTeamMembers(classroomId: string, teamId: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<Array<UserClassroomResponse>>>> {
            const localVarAxiosArgs = await MemberApiAxiosParamCreator(configuration).getClassroomTeamMembers(classroomId, teamId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * Remove Member from the team
         * @summary Remove Member from the team
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {number} memberId Member ID
         * @param {string} xCsrfToken Csrf-Token
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async removeMemberFromTeamV2(classroomId: string, teamId: string, memberId: number, xCsrfToken: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<void>>> {
            const localVarAxiosArgs = await MemberApiAxiosParamCreator(configuration).removeMemberFromTeamV2(classroomId, teamId, memberId, xCsrfToken, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * Update Classroom Members role
         * @summary Update Classroom Members role
         * @param {UpdateMemberRoleRequest} body Update Member Role
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async updateMemberRole(body: UpdateMemberRoleRequest, xCsrfToken: string, classroomId: string, memberId: number, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<void>>> {
            const localVarAxiosArgs = await MemberApiAxiosParamCreator(configuration).updateMemberRole(body, xCsrfToken, classroomId, memberId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * Update Classroom Members team
         * @summary Update Classroom Members team
         * @param {UpdateMemberTeamRequest} body Update Member Team
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async updateMemberTeam(body: UpdateMemberTeamRequest, xCsrfToken: string, classroomId: string, memberId: number, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<void>>> {
            const localVarAxiosArgs = await MemberApiAxiosParamCreator(configuration).updateMemberTeam(body, xCsrfToken, classroomId, memberId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
    }
};

/**
 * MemberApi - factory interface
 * @export
 */
export const MemberApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    return {
        /**
         * GetClassroomMember
         * @summary GetClassroomMember
         * @param {string} classroomId Classroom ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomMember(classroomId: string, memberId: number, options?: AxiosRequestConfig): Promise<AxiosResponse<UserClassroomResponse>> {
            return MemberApiFp(configuration).getClassroomMember(classroomId, memberId, options).then((request) => request(axios, basePath));
        },
        /**
         * GetClassroomMembers
         * @summary GetClassroomMembers
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomMembers(classroomId: string, options?: AxiosRequestConfig): Promise<AxiosResponse<Array<UserClassroomResponse>>> {
            return MemberApiFp(configuration).getClassroomMembers(classroomId, options).then((request) => request(axios, basePath));
        },
        /**
         * GetClassroomTeamMember
         * @summary GetClassroomTeamMember
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomTeamMember(classroomId: string, teamId: string, memberId: number, options?: AxiosRequestConfig): Promise<AxiosResponse<UserClassroomResponse>> {
            return MemberApiFp(configuration).getClassroomTeamMember(classroomId, teamId, memberId, options).then((request) => request(axios, basePath));
        },
        /**
         * GetClassroomTeamMembers
         * @summary GetClassroomTeamMembers
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomTeamMembers(classroomId: string, teamId: string, options?: AxiosRequestConfig): Promise<AxiosResponse<Array<UserClassroomResponse>>> {
            return MemberApiFp(configuration).getClassroomTeamMembers(classroomId, teamId, options).then((request) => request(axios, basePath));
        },
        /**
         * Remove Member from the team
         * @summary Remove Member from the team
         * @param {string} classroomId Classroom ID
         * @param {string} teamId Team ID
         * @param {number} memberId Member ID
         * @param {string} xCsrfToken Csrf-Token
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async removeMemberFromTeamV2(classroomId: string, teamId: string, memberId: number, xCsrfToken: string, options?: AxiosRequestConfig): Promise<AxiosResponse<void>> {
            return MemberApiFp(configuration).removeMemberFromTeamV2(classroomId, teamId, memberId, xCsrfToken, options).then((request) => request(axios, basePath));
        },
        /**
         * Update Classroom Members role
         * @summary Update Classroom Members role
         * @param {UpdateMemberRoleRequest} body Update Member Role
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async updateMemberRole(body: UpdateMemberRoleRequest, xCsrfToken: string, classroomId: string, memberId: number, options?: AxiosRequestConfig): Promise<AxiosResponse<void>> {
            return MemberApiFp(configuration).updateMemberRole(body, xCsrfToken, classroomId, memberId, options).then((request) => request(axios, basePath));
        },
        /**
         * Update Classroom Members team
         * @summary Update Classroom Members team
         * @param {UpdateMemberTeamRequest} body Update Member Team
         * @param {string} xCsrfToken Csrf-Token
         * @param {string} classroomId Classroom ID
         * @param {number} memberId Member ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async updateMemberTeam(body: UpdateMemberTeamRequest, xCsrfToken: string, classroomId: string, memberId: number, options?: AxiosRequestConfig): Promise<AxiosResponse<void>> {
            return MemberApiFp(configuration).updateMemberTeam(body, xCsrfToken, classroomId, memberId, options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * MemberApi - object-oriented interface
 * @export
 * @class MemberApi
 * @extends {BaseAPI}
 */
export class MemberApi extends BaseAPI {
    /**
     * GetClassroomMember
     * @summary GetClassroomMember
     * @param {string} classroomId Classroom ID
     * @param {number} memberId Member ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MemberApi
     */
    public async getClassroomMember(classroomId: string, memberId: number, options?: AxiosRequestConfig) : Promise<AxiosResponse<UserClassroomResponse>> {
        return MemberApiFp(this.configuration).getClassroomMember(classroomId, memberId, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * GetClassroomMembers
     * @summary GetClassroomMembers
     * @param {string} classroomId Classroom ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MemberApi
     */
    public async getClassroomMembers(classroomId: string, options?: AxiosRequestConfig) : Promise<AxiosResponse<Array<UserClassroomResponse>>> {
        return MemberApiFp(this.configuration).getClassroomMembers(classroomId, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * GetClassroomTeamMember
     * @summary GetClassroomTeamMember
     * @param {string} classroomId Classroom ID
     * @param {string} teamId Team ID
     * @param {number} memberId Member ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MemberApi
     */
    public async getClassroomTeamMember(classroomId: string, teamId: string, memberId: number, options?: AxiosRequestConfig) : Promise<AxiosResponse<UserClassroomResponse>> {
        return MemberApiFp(this.configuration).getClassroomTeamMember(classroomId, teamId, memberId, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * GetClassroomTeamMembers
     * @summary GetClassroomTeamMembers
     * @param {string} classroomId Classroom ID
     * @param {string} teamId Team ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MemberApi
     */
    public async getClassroomTeamMembers(classroomId: string, teamId: string, options?: AxiosRequestConfig) : Promise<AxiosResponse<Array<UserClassroomResponse>>> {
        return MemberApiFp(this.configuration).getClassroomTeamMembers(classroomId, teamId, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * Remove Member from the team
     * @summary Remove Member from the team
     * @param {string} classroomId Classroom ID
     * @param {string} teamId Team ID
     * @param {number} memberId Member ID
     * @param {string} xCsrfToken Csrf-Token
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MemberApi
     */
    public async removeMemberFromTeamV2(classroomId: string, teamId: string, memberId: number, xCsrfToken: string, options?: AxiosRequestConfig) : Promise<AxiosResponse<void>> {
        return MemberApiFp(this.configuration).removeMemberFromTeamV2(classroomId, teamId, memberId, xCsrfToken, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * Update Classroom Members role
     * @summary Update Classroom Members role
     * @param {UpdateMemberRoleRequest} body Update Member Role
     * @param {string} xCsrfToken Csrf-Token
     * @param {string} classroomId Classroom ID
     * @param {number} memberId Member ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MemberApi
     */
    public async updateMemberRole(body: UpdateMemberRoleRequest, xCsrfToken: string, classroomId: string, memberId: number, options?: AxiosRequestConfig) : Promise<AxiosResponse<void>> {
        return MemberApiFp(this.configuration).updateMemberRole(body, xCsrfToken, classroomId, memberId, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * Update Classroom Members team
     * @summary Update Classroom Members team
     * @param {UpdateMemberTeamRequest} body Update Member Team
     * @param {string} xCsrfToken Csrf-Token
     * @param {string} classroomId Classroom ID
     * @param {number} memberId Member ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof MemberApi
     */
    public async updateMemberTeam(body: UpdateMemberTeamRequest, xCsrfToken: string, classroomId: string, memberId: number, options?: AxiosRequestConfig) : Promise<AxiosResponse<void>> {
        return MemberApiFp(this.configuration).updateMemberTeam(body, xCsrfToken, classroomId, memberId, options).then((request) => request(this.axios, this.basePath));
    }
}
