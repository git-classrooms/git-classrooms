// @ts-nocheck
/* tslint:disable */
/* eslint-disable */
/**
 * GitLab Classrooms – Backend API
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
import { ClassroomRunnerResponse } from '../models';
import { HTTPError } from '../models';
/**
 * RunnersApi - axios parameter creator
 * @export
 */
export const RunnersApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * GetClassroomRunners
         * @summary GetClassroomRunners
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getClassroomRunners: async (classroomId: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling getClassroomRunners.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/runners`
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
         * GetClassroomRunnersAreAvailable
         * @summary GetClassroomRunnersAreAvailable
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getClassroomRunnersAreAvailable: async (classroomId: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'classroomId' is not null or undefined
            if (classroomId === null || classroomId === undefined) {
                throw new RequiredError('classroomId','Required parameter classroomId was null or undefined when calling getClassroomRunnersAreAvailable.');
            }
            const localVarPath = `/api/v2/classrooms/{classroomId}/runners/available`
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
    }
};

/**
 * RunnersApi - functional programming interface
 * @export
 */
export const RunnersApiFp = function(configuration?: Configuration) {
    return {
        /**
         * GetClassroomRunners
         * @summary GetClassroomRunners
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomRunners(classroomId: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<Array<ClassroomRunnerResponse>>>> {
            const localVarAxiosArgs = await RunnersApiAxiosParamCreator(configuration).getClassroomRunners(classroomId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * GetClassroomRunnersAreAvailable
         * @summary GetClassroomRunnersAreAvailable
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomRunnersAreAvailable(classroomId: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => Promise<AxiosResponse<boolean>>> {
            const localVarAxiosArgs = await RunnersApiAxiosParamCreator(configuration).getClassroomRunnersAreAvailable(classroomId, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs :AxiosRequestConfig = {...localVarAxiosArgs.options, url: basePath + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
    }
};

/**
 * RunnersApi - factory interface
 * @export
 */
export const RunnersApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    return {
        /**
         * GetClassroomRunners
         * @summary GetClassroomRunners
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomRunners(classroomId: string, options?: AxiosRequestConfig): Promise<AxiosResponse<Array<ClassroomRunnerResponse>>> {
            return RunnersApiFp(configuration).getClassroomRunners(classroomId, options).then((request) => request(axios, basePath));
        },
        /**
         * GetClassroomRunnersAreAvailable
         * @summary GetClassroomRunnersAreAvailable
         * @param {string} classroomId Classroom ID
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getClassroomRunnersAreAvailable(classroomId: string, options?: AxiosRequestConfig): Promise<AxiosResponse<boolean>> {
            return RunnersApiFp(configuration).getClassroomRunnersAreAvailable(classroomId, options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * RunnersApi - object-oriented interface
 * @export
 * @class RunnersApi
 * @extends {BaseAPI}
 */
export class RunnersApi extends BaseAPI {
    /**
     * GetClassroomRunners
     * @summary GetClassroomRunners
     * @param {string} classroomId Classroom ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof RunnersApi
     */
    public async getClassroomRunners(classroomId: string, options?: AxiosRequestConfig) : Promise<AxiosResponse<Array<ClassroomRunnerResponse>>> {
        return RunnersApiFp(this.configuration).getClassroomRunners(classroomId, options).then((request) => request(this.axios, this.basePath));
    }
    /**
     * GetClassroomRunnersAreAvailable
     * @summary GetClassroomRunnersAreAvailable
     * @param {string} classroomId Classroom ID
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof RunnersApi
     */
    public async getClassroomRunnersAreAvailable(classroomId: string, options?: AxiosRequestConfig) : Promise<AxiosResponse<boolean>> {
        return RunnersApiFp(this.configuration).getClassroomRunnersAreAvailable(classroomId, options).then((request) => request(this.axios, this.basePath));
    }
}
