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

 /**
 * 
 *
 * @export
 * @interface GradingManualRubricRequest
 */
export interface GradingManualRubricRequest {

    /**
     * @type {string}
     * @memberof GradingManualRubricRequest
     */
    description: string;

    /**
     * @type {string}
     * @memberof GradingManualRubricRequest
     */
    id?: string;

    /**
     * @type {number}
     * @memberof GradingManualRubricRequest
     */
    maxScore: number;

    /**
     * @type {string}
     * @memberof GradingManualRubricRequest
     */
    name: string;
}
