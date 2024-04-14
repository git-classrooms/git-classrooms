/* tslint:disable */
/* eslint-disable */
/**
 * Gitlab Classroom API
 * This is the API for our Gitlab Classroom Webapp
 *
 * OpenAPI spec version: 1.0
 * Contact: support@swagger.io
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 * Do not edit the class manually.
 */

 /**
 * 
 *
 * @export
 * @interface DefaultControllerCreateClassroomRequest
 */
export interface DefaultControllerCreateClassroomRequest {

    /**
     * @type {boolean}
     * @memberof DefaultControllerCreateClassroomRequest
     */
    createTeams?: boolean;

    /**
     * @type {string}
     * @memberof DefaultControllerCreateClassroomRequest
     */
    description?: string;

    /**
     * @type {number}
     * @memberof DefaultControllerCreateClassroomRequest
     */
    maxTeamSize?: number;

    /**
     * @type {number}
     * @memberof DefaultControllerCreateClassroomRequest
     */
    maxTeams?: number;

    /**
     * @type {string}
     * @memberof DefaultControllerCreateClassroomRequest
     */
    name?: string;
}
