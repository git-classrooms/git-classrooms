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
 * @interface GetOwnedClassroomTeamResponse
 */
export interface GetOwnedClassroomTeamResponse {

    /**
     * @type {string}
     * @memberof GetOwnedClassroomTeamResponse
     */
    createdAt: string;

    /**
     * @type {string}
     * @memberof GetOwnedClassroomTeamResponse
     */
    gitlabUrl: string;

    /**
     * @type {number}
     * @memberof GetOwnedClassroomTeamResponse
     */
    groupId: number;

    /**
     * @type {string}
     * @memberof GetOwnedClassroomTeamResponse
     */
    id: string;

    /**
     * @type {Array<User>}
     * @memberof GetOwnedClassroomTeamResponse
     */
    members: Array<User>;

    /**
     * @type {string}
     * @memberof GetOwnedClassroomTeamResponse
     */
    name: string;

    /**
     * @type {string}
     * @memberof GetOwnedClassroomTeamResponse
     */
    updatedAt: string;
}
