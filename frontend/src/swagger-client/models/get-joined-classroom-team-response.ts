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

import { User } from './user';
 /**
 * 
 *
 * @export
 * @interface GetJoinedClassroomTeamResponse
 */
export interface GetJoinedClassroomTeamResponse {

    /**
     * @type {string}
     * @memberof GetJoinedClassroomTeamResponse
     */
    createdAt?: string;

    /**
     * @type {string}
     * @memberof GetJoinedClassroomTeamResponse
     */
    gitlabUrl?: string;

    /**
     * @type {number}
     * @memberof GetJoinedClassroomTeamResponse
     */
    groupId?: number;

    /**
     * @type {string}
     * @memberof GetJoinedClassroomTeamResponse
     */
    id?: string;

    /**
     * @type {Array<User>}
     * @memberof GetJoinedClassroomTeamResponse
     */
    members?: Array<User>;

    /**
     * @type {string}
     * @memberof GetJoinedClassroomTeamResponse
     */
    name?: string;

    /**
     * @type {string}
     * @memberof GetJoinedClassroomTeamResponse
     */
    updatedAt?: string;
}