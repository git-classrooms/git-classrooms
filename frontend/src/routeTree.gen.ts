/* prettier-ignore-start */

/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file is auto-generated by TanStack Router

import { createFileRoute } from "@tanstack/react-router"

// Import Routes

import { Route as rootRoute } from "./routes/__root"
import { Route as LoginImport } from "./routes/login"
import { Route as AuthImport } from "./routes/_auth"
import { Route as IndexImport } from "./routes/index"
import { Route as AuthDashboardCreateImport } from "./routes/_auth/dashboard/create"
import { Route as AuthDashboardIndexImport } from "./routes/_auth/dashboard/_index"
import { Route as AuthClassroomsCreateImport } from "./routes/_auth/classrooms/create"
import { Route as AuthClassroomsIndexImport } from "./routes/_auth/classrooms/_index"
import { Route as AuthDashboardIndexIndexImport } from "./routes/_auth/dashboard/_index/index"
import { Route as AuthClassroomsIndexIndexImport } from "./routes/_auth/classrooms/_index/index"
import { Route as AuthClassroomsClassroomIdInviteImport } from "./routes/_auth/classrooms/$classroomId/invite"
import { Route as AuthClassroomsClassroomIdEditImport } from "./routes/_auth/classrooms/$classroomId/edit"
import { Route as AuthClassroomsClassroomIdIndexImport } from "./routes/_auth/classrooms/$classroomId/_index"
import { Route as AuthClassroomsClassroomIdTeamsRouteImport } from "./routes/_auth/classrooms/$classroomId/teams/route"
import { Route as AuthClassroomsClassroomIdIndexIndexImport } from "./routes/_auth/classrooms/$classroomId/_index/index"
import { Route as AuthDashboardIndexCreateModalImport } from "./routes/_auth/dashboard/_index/create.modal"
import { Route as AuthClassroomsIndexCreateModalImport } from "./routes/_auth/classrooms/_index/create.modal"
import { Route as AuthClassroomsClassroomIdTeamsCreateImport } from "./routes/_auth/classrooms/$classroomId/teams/create"
import { Route as AuthClassroomsClassroomIdTeamsIndexImport } from "./routes/_auth/classrooms/$classroomId/teams/_index"
import { Route as AuthClassroomsClassroomIdInvitationsInvitationIdImport } from "./routes/_auth/classrooms/$classroomId/invitations/$invitationId"
import { Route as AuthClassroomsClassroomIdAssignmentsCreateImport } from "./routes/_auth/classrooms/$classroomId/assignments/create"
import { Route as AuthClassroomsClassroomIdTeamsTeamIdRouteImport } from "./routes/_auth/classrooms/$classroomId/teams/$teamId/route"
import { Route as AuthClassroomsClassroomIdTeamsJoinIndexImport } from "./routes/_auth/classrooms/$classroomId/teams/join.index"
import { Route as AuthClassroomsClassroomIdTeamsIndexIndexImport } from "./routes/_auth/classrooms/$classroomId/teams/_index/index"
import { Route as AuthClassroomsClassroomIdAssignmentsAssignmentIdIndexImport } from "./routes/_auth/classrooms/$classroomId/assignments/$assignmentId/index"
import { Route as AuthClassroomsClassroomIdProjectsProjectIdAcceptImport } from "./routes/_auth/classrooms/$classroomId/projects/$projectId/accept"
import { Route as AuthClassroomsClassroomIdTeamsIndexCreateModalImport } from "./routes/_auth/classrooms/$classroomId/teams/_index/create.modal"
import { Route as AuthClassroomsClassroomIdIndexTeamsTeamIdModalImport } from "./routes/_auth/classrooms/$classroomId/_index/teams.$teamId.modal"
import { Route as AuthClassroomsClassroomIdIndexTeamCreateModalImport } from "./routes/_auth/classrooms/$classroomId/_index/team.create.modal"

// Create Virtual Routes

const AuthDashboardImport = createFileRoute("/_auth/dashboard")()
const AuthClassroomsImport = createFileRoute("/_auth/classrooms")()
const AuthClassroomsClassroomIdImport = createFileRoute(
  "/_auth/classrooms/$classroomId",
)()

// Create/Update Routes

const LoginRoute = LoginImport.update({
  path: "/login",
  getParentRoute: () => rootRoute,
} as any)

const AuthRoute = AuthImport.update({
  id: "/_auth",
  getParentRoute: () => rootRoute,
} as any)

const IndexRoute = IndexImport.update({
  path: "/",
  getParentRoute: () => rootRoute,
} as any)

const AuthDashboardRoute = AuthDashboardImport.update({
  path: "/dashboard",
  getParentRoute: () => AuthRoute,
} as any)

const AuthClassroomsRoute = AuthClassroomsImport.update({
  path: "/classrooms",
  getParentRoute: () => AuthRoute,
} as any)

const AuthClassroomsClassroomIdRoute = AuthClassroomsClassroomIdImport.update({
  path: "/$classroomId",
  getParentRoute: () => AuthClassroomsRoute,
} as any)

const AuthDashboardCreateRoute = AuthDashboardCreateImport.update({
  path: "/create",
  getParentRoute: () => AuthDashboardRoute,
} as any)

const AuthDashboardIndexRoute = AuthDashboardIndexImport.update({
  id: "/_index",
  getParentRoute: () => AuthDashboardRoute,
} as any)

const AuthClassroomsCreateRoute = AuthClassroomsCreateImport.update({
  path: "/create",
  getParentRoute: () => AuthClassroomsRoute,
} as any)

const AuthClassroomsIndexRoute = AuthClassroomsIndexImport.update({
  id: "/_index",
  getParentRoute: () => AuthClassroomsRoute,
} as any)

const AuthDashboardIndexIndexRoute = AuthDashboardIndexIndexImport.update({
  path: "/",
  getParentRoute: () => AuthDashboardIndexRoute,
} as any)

const AuthClassroomsIndexIndexRoute = AuthClassroomsIndexIndexImport.update({
  path: "/",
  getParentRoute: () => AuthClassroomsIndexRoute,
} as any)

const AuthClassroomsClassroomIdInviteRoute =
  AuthClassroomsClassroomIdInviteImport.update({
    path: "/invite",
    getParentRoute: () => AuthClassroomsClassroomIdRoute,
  } as any)

const AuthClassroomsClassroomIdEditRoute =
  AuthClassroomsClassroomIdEditImport.update({
    path: "/edit",
    getParentRoute: () => AuthClassroomsClassroomIdRoute,
  } as any)

const AuthClassroomsClassroomIdIndexRoute =
  AuthClassroomsClassroomIdIndexImport.update({
    id: "/_index",
    getParentRoute: () => AuthClassroomsClassroomIdRoute,
  } as any)

const AuthClassroomsClassroomIdTeamsRouteRoute =
  AuthClassroomsClassroomIdTeamsRouteImport.update({
    path: "/$classroomId/teams",
    getParentRoute: () => AuthClassroomsRoute,
  } as any)

const AuthClassroomsClassroomIdIndexIndexRoute =
  AuthClassroomsClassroomIdIndexIndexImport.update({
    path: "/",
    getParentRoute: () => AuthClassroomsClassroomIdIndexRoute,
  } as any)

const AuthDashboardIndexCreateModalRoute =
  AuthDashboardIndexCreateModalImport.update({
    path: "/create/modal",
    getParentRoute: () => AuthDashboardIndexRoute,
  } as any)

const AuthClassroomsIndexCreateModalRoute =
  AuthClassroomsIndexCreateModalImport.update({
    path: "/create/modal",
    getParentRoute: () => AuthClassroomsIndexRoute,
  } as any)

const AuthClassroomsClassroomIdTeamsCreateRoute =
  AuthClassroomsClassroomIdTeamsCreateImport.update({
    path: "/create",
    getParentRoute: () => AuthClassroomsClassroomIdTeamsRouteRoute,
  } as any)

const AuthClassroomsClassroomIdTeamsIndexRoute =
  AuthClassroomsClassroomIdTeamsIndexImport.update({
    id: "/_index",
    getParentRoute: () => AuthClassroomsClassroomIdTeamsRouteRoute,
  } as any)

const AuthClassroomsClassroomIdInvitationsInvitationIdRoute =
  AuthClassroomsClassroomIdInvitationsInvitationIdImport.update({
    path: "/invitations/$invitationId",
    getParentRoute: () => AuthClassroomsClassroomIdRoute,
  } as any)

const AuthClassroomsClassroomIdAssignmentsCreateRoute =
  AuthClassroomsClassroomIdAssignmentsCreateImport.update({
    path: "/assignments/create",
    getParentRoute: () => AuthClassroomsClassroomIdRoute,
  } as any)

const AuthClassroomsClassroomIdTeamsTeamIdRouteRoute =
  AuthClassroomsClassroomIdTeamsTeamIdRouteImport.update({
    path: "/$teamId",
    getParentRoute: () => AuthClassroomsClassroomIdTeamsRouteRoute,
  } as any)

const AuthClassroomsClassroomIdTeamsJoinIndexRoute =
  AuthClassroomsClassroomIdTeamsJoinIndexImport.update({
    path: "/join/",
    getParentRoute: () => AuthClassroomsClassroomIdTeamsRouteRoute,
  } as any)

const AuthClassroomsClassroomIdTeamsIndexIndexRoute =
  AuthClassroomsClassroomIdTeamsIndexIndexImport.update({
    path: "/",
    getParentRoute: () => AuthClassroomsClassroomIdTeamsIndexRoute,
  } as any)

const AuthClassroomsClassroomIdAssignmentsAssignmentIdIndexRoute =
  AuthClassroomsClassroomIdAssignmentsAssignmentIdIndexImport.update({
    path: "/assignments/$assignmentId/",
    getParentRoute: () => AuthClassroomsClassroomIdRoute,
  } as any)

const AuthClassroomsClassroomIdProjectsProjectIdAcceptRoute =
  AuthClassroomsClassroomIdProjectsProjectIdAcceptImport.update({
    path: "/projects/$projectId/accept",
    getParentRoute: () => AuthClassroomsClassroomIdRoute,
  } as any)

const AuthClassroomsClassroomIdTeamsIndexCreateModalRoute =
  AuthClassroomsClassroomIdTeamsIndexCreateModalImport.update({
    path: "/create/modal",
    getParentRoute: () => AuthClassroomsClassroomIdTeamsIndexRoute,
  } as any)

const AuthClassroomsClassroomIdIndexTeamsTeamIdModalRoute =
  AuthClassroomsClassroomIdIndexTeamsTeamIdModalImport.update({
    path: "/teams/$teamId/modal",
    getParentRoute: () => AuthClassroomsClassroomIdIndexRoute,
  } as any)

const AuthClassroomsClassroomIdIndexTeamCreateModalRoute =
  AuthClassroomsClassroomIdIndexTeamCreateModalImport.update({
    path: "/team/create/modal",
    getParentRoute: () => AuthClassroomsClassroomIdIndexRoute,
  } as any)

// Populate the FileRoutesByPath interface

declare module "@tanstack/react-router" {
  interface FileRoutesByPath {
    "/": {
      preLoaderRoute: typeof IndexImport
      parentRoute: typeof rootRoute
    }
    "/_auth": {
      preLoaderRoute: typeof AuthImport
      parentRoute: typeof rootRoute
    }
    "/login": {
      preLoaderRoute: typeof LoginImport
      parentRoute: typeof rootRoute
    }
    "/_auth/classrooms": {
      preLoaderRoute: typeof AuthClassroomsImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/_index": {
      preLoaderRoute: typeof AuthClassroomsIndexImport
      parentRoute: typeof AuthClassroomsRoute
    }
    "/_auth/classrooms/create": {
      preLoaderRoute: typeof AuthClassroomsCreateImport
      parentRoute: typeof AuthClassroomsImport
    }
    "/_auth/dashboard": {
      preLoaderRoute: typeof AuthDashboardImport
      parentRoute: typeof AuthImport
    }
    "/_auth/dashboard/_index": {
      preLoaderRoute: typeof AuthDashboardIndexImport
      parentRoute: typeof AuthDashboardRoute
    }
    "/_auth/dashboard/create": {
      preLoaderRoute: typeof AuthDashboardCreateImport
      parentRoute: typeof AuthDashboardImport
    }
    "/_auth/classrooms/$classroomId/teams": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdTeamsRouteImport
      parentRoute: typeof AuthClassroomsImport
    }
    "/_auth/classrooms/$classroomId": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdImport
      parentRoute: typeof AuthClassroomsImport
    }
    "/_auth/classrooms/$classroomId/_index": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdIndexImport
      parentRoute: typeof AuthClassroomsClassroomIdRoute
    }
    "/_auth/classrooms/$classroomId/edit": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdEditImport
      parentRoute: typeof AuthClassroomsClassroomIdImport
    }
    "/_auth/classrooms/$classroomId/invite": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdInviteImport
      parentRoute: typeof AuthClassroomsClassroomIdImport
    }
    "/_auth/classrooms/_index/": {
      preLoaderRoute: typeof AuthClassroomsIndexIndexImport
      parentRoute: typeof AuthClassroomsIndexImport
    }
    "/_auth/dashboard/_index/": {
      preLoaderRoute: typeof AuthDashboardIndexIndexImport
      parentRoute: typeof AuthDashboardIndexImport
    }
    "/_auth/classrooms/$classroomId/teams/$teamId": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdTeamsTeamIdRouteImport
      parentRoute: typeof AuthClassroomsClassroomIdTeamsRouteImport
    }
    "/_auth/classrooms/$classroomId/assignments/create": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdAssignmentsCreateImport
      parentRoute: typeof AuthClassroomsClassroomIdImport
    }
    "/_auth/classrooms/$classroomId/invitations/$invitationId": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdInvitationsInvitationIdImport
      parentRoute: typeof AuthClassroomsClassroomIdImport
    }
    "/_auth/classrooms/$classroomId/teams/_index": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdTeamsIndexImport
      parentRoute: typeof AuthClassroomsClassroomIdTeamsRouteImport
    }
    "/_auth/classrooms/$classroomId/teams/create": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdTeamsCreateImport
      parentRoute: typeof AuthClassroomsClassroomIdTeamsRouteImport
    }
    "/_auth/classrooms/_index/create/modal": {
      preLoaderRoute: typeof AuthClassroomsIndexCreateModalImport
      parentRoute: typeof AuthClassroomsIndexImport
    }
    "/_auth/dashboard/_index/create/modal": {
      preLoaderRoute: typeof AuthDashboardIndexCreateModalImport
      parentRoute: typeof AuthDashboardIndexImport
    }
    "/_auth/classrooms/$classroomId/_index/": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdIndexIndexImport
      parentRoute: typeof AuthClassroomsClassroomIdIndexImport
    }
    "/_auth/classrooms/$classroomId/projects/$projectId/accept": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdProjectsProjectIdAcceptImport
      parentRoute: typeof AuthClassroomsClassroomIdImport
    }
    "/_auth/classrooms/$classroomId/assignments/$assignmentId/": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdAssignmentsAssignmentIdIndexImport
      parentRoute: typeof AuthClassroomsClassroomIdImport
    }
    "/_auth/classrooms/$classroomId/teams/_index/": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdTeamsIndexIndexImport
      parentRoute: typeof AuthClassroomsClassroomIdTeamsIndexImport
    }
    "/_auth/classrooms/$classroomId/teams/join/": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdTeamsJoinIndexImport
      parentRoute: typeof AuthClassroomsClassroomIdTeamsRouteImport
    }
    "/_auth/classrooms/$classroomId/_index/team/create/modal": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdIndexTeamCreateModalImport
      parentRoute: typeof AuthClassroomsClassroomIdIndexImport
    }
    "/_auth/classrooms/$classroomId/_index/teams/$teamId/modal": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdIndexTeamsTeamIdModalImport
      parentRoute: typeof AuthClassroomsClassroomIdIndexImport
    }
    "/_auth/classrooms/$classroomId/teams/_index/create/modal": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdTeamsIndexCreateModalImport
      parentRoute: typeof AuthClassroomsClassroomIdTeamsIndexImport
    }
  }
}

// Create and export the route tree

export const routeTree = rootRoute.addChildren([
  IndexRoute,
  AuthRoute.addChildren([
    AuthClassroomsRoute.addChildren([
      AuthClassroomsIndexRoute.addChildren([
        AuthClassroomsIndexIndexRoute,
        AuthClassroomsIndexCreateModalRoute,
      ]),
      AuthClassroomsCreateRoute,
      AuthClassroomsClassroomIdTeamsRouteRoute.addChildren([
        AuthClassroomsClassroomIdTeamsTeamIdRouteRoute,
        AuthClassroomsClassroomIdTeamsIndexRoute.addChildren([
          AuthClassroomsClassroomIdTeamsIndexIndexRoute,
          AuthClassroomsClassroomIdTeamsIndexCreateModalRoute,
        ]),
        AuthClassroomsClassroomIdTeamsCreateRoute,
        AuthClassroomsClassroomIdTeamsJoinIndexRoute,
      ]),
      AuthClassroomsClassroomIdRoute.addChildren([
        AuthClassroomsClassroomIdIndexRoute.addChildren([
          AuthClassroomsClassroomIdIndexIndexRoute,
          AuthClassroomsClassroomIdIndexTeamCreateModalRoute,
          AuthClassroomsClassroomIdIndexTeamsTeamIdModalRoute,
        ]),
        AuthClassroomsClassroomIdEditRoute,
        AuthClassroomsClassroomIdInviteRoute,
        AuthClassroomsClassroomIdAssignmentsCreateRoute,
        AuthClassroomsClassroomIdInvitationsInvitationIdRoute,
        AuthClassroomsClassroomIdProjectsProjectIdAcceptRoute,
        AuthClassroomsClassroomIdAssignmentsAssignmentIdIndexRoute,
      ]),
    ]),
    AuthDashboardRoute.addChildren([
      AuthDashboardIndexRoute.addChildren([
        AuthDashboardIndexIndexRoute,
        AuthDashboardIndexCreateModalRoute,
      ]),
      AuthDashboardCreateRoute,
    ]),
  ]),
  LoginRoute,
])

/* prettier-ignore-end */
