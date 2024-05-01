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

import { Role } from './role';
 /**
 * 
 *
 * @export
 * @interface ChangeOwnedClassroomMemberRequest
 */
export interface ChangeOwnedClassroomMemberRequest {

    /**
     * @type {Role}
     * @memberof ChangeOwnedClassroomMemberRequest
     */
    role?: Role;

    /**
     * @type {string}
     * @memberof ChangeOwnedClassroomMemberRequest
     */
    teamId?: string;
}
