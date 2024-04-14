// This file is auto-generated by TanStack Router

// Import Routes

import { Route as rootRoute } from "./routes/__root"
import { Route as LoginImport } from "./routes/login"
import { Route as AuthImport } from "./routes/_auth"
import { Route as IndexImport } from "./routes/index"
import { Route as AuthClassroomsRouteImport } from "./routes/_auth/classrooms/route"
import { Route as AuthClassroomsCreateImport } from "./routes/_auth/classrooms_/create"
import { Route as AuthClassroomsOwnedIndexImport } from "./routes/_auth/classrooms_/owned/index"
import { Route as AuthClassroomsJoinedIndexImport } from "./routes/_auth/classrooms_/joined/index"
import { Route as AuthClassroomsCreateModalImport } from "./routes/_auth/classrooms/create.modal"
import { Route as AuthClassroomsJoinedClassroomIdRouteImport } from "./routes/_auth/classrooms_/joined/$classroomId/route"
import { Route as AuthClassroomsOwnedClassroomIdIndexImport } from "./routes/_auth/classrooms_/owned/$classroomId/index"
import { Route as AuthClassroomsJoinedClassroomIdIndexImport } from "./routes/_auth/classrooms_/joined/$classroomId/index"
import { Route as AuthClassroomsOwnedClassroomIdInviteImport } from "./routes/_auth/classrooms_/owned/$classroomId/invite"
import { Route as AuthClassroomsJoinedClassroomIdTeamsRouteImport } from "./routes/_auth/classrooms_/joined/$classroomId/teams/route"
import { Route as AuthClassroomsOwnedClassroomIdAssignmentsCreateImport } from "./routes/_auth/classrooms_/owned/$classroomId/assignments/create"
import { Route as AuthClassroomsJoinedClassroomIdInvitationsInvitationIdImport } from "./routes/_auth/classrooms_/joined/$classroomId_/invitations/$invitationId"
import { Route as AuthClassroomsOwnedClassroomIdAssignmentsAssignmentIdIndexImport } from "./routes/_auth/classrooms_/owned/$classroomId/assignments/$assignmentId/index"
import { Route as AuthClassroomsJoinedClassroomIdTeamsJoinIndexImport } from "./routes/_auth/classrooms_/joined/$classroomId/teams/join.index"
import { Route as AuthClassroomsJoinedClassroomIdAssignmentsAssignmentIdAcceptImport } from "./routes/_auth/classrooms_/joined/$classroomId/assignments/$assignmentId/accept"

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

const AuthClassroomsRouteRoute = AuthClassroomsRouteImport.update({
  path: "/classrooms",
  getParentRoute: () => AuthRoute,
} as any)

const AuthClassroomsCreateRoute = AuthClassroomsCreateImport.update({
  path: "/classrooms/create",
  getParentRoute: () => AuthRoute,
} as any)

const AuthClassroomsOwnedIndexRoute = AuthClassroomsOwnedIndexImport.update({
  path: "/classrooms/owned/",
  getParentRoute: () => AuthRoute,
} as any)

const AuthClassroomsJoinedIndexRoute = AuthClassroomsJoinedIndexImport.update({
  path: "/classrooms/joined/",
  getParentRoute: () => AuthRoute,
} as any)

const AuthClassroomsCreateModalRoute = AuthClassroomsCreateModalImport.update({
  path: "/create/modal",
  getParentRoute: () => AuthClassroomsRouteRoute,
} as any)

const AuthClassroomsJoinedClassroomIdRouteRoute =
  AuthClassroomsJoinedClassroomIdRouteImport.update({
    path: "/classrooms/joined/$classroomId",
    getParentRoute: () => AuthRoute,
  } as any)

const AuthClassroomsOwnedClassroomIdIndexRoute =
  AuthClassroomsOwnedClassroomIdIndexImport.update({
    path: "/classrooms/owned/$classroomId/",
    getParentRoute: () => AuthRoute,
  } as any)

const AuthClassroomsJoinedClassroomIdIndexRoute =
  AuthClassroomsJoinedClassroomIdIndexImport.update({
    path: "/",
    getParentRoute: () => AuthClassroomsJoinedClassroomIdRouteRoute,
  } as any)

const AuthClassroomsOwnedClassroomIdInviteRoute =
  AuthClassroomsOwnedClassroomIdInviteImport.update({
    path: "/classrooms/owned/$classroomId/invite",
    getParentRoute: () => AuthRoute,
  } as any)

const AuthClassroomsJoinedClassroomIdTeamsRouteRoute =
  AuthClassroomsJoinedClassroomIdTeamsRouteImport.update({
    path: "/teams",
    getParentRoute: () => AuthClassroomsJoinedClassroomIdRouteRoute,
  } as any)

const AuthClassroomsOwnedClassroomIdAssignmentsCreateRoute =
  AuthClassroomsOwnedClassroomIdAssignmentsCreateImport.update({
    path: "/classrooms/owned/$classroomId/assignments/create",
    getParentRoute: () => AuthRoute,
  } as any)

const AuthClassroomsJoinedClassroomIdInvitationsInvitationIdRoute =
  AuthClassroomsJoinedClassroomIdInvitationsInvitationIdImport.update({
    path: "/classrooms/joined/$classroomId/invitations/$invitationId",
    getParentRoute: () => AuthRoute,
  } as any)

const AuthClassroomsOwnedClassroomIdAssignmentsAssignmentIdIndexRoute =
  AuthClassroomsOwnedClassroomIdAssignmentsAssignmentIdIndexImport.update({
    path: "/classrooms/owned/$classroomId/assignments/$assignmentId/",
    getParentRoute: () => AuthRoute,
  } as any)

const AuthClassroomsJoinedClassroomIdTeamsJoinIndexRoute =
  AuthClassroomsJoinedClassroomIdTeamsJoinIndexImport.update({
    path: "/join/",
    getParentRoute: () => AuthClassroomsJoinedClassroomIdTeamsRouteRoute,
  } as any)

const AuthClassroomsJoinedClassroomIdAssignmentsAssignmentIdAcceptRoute =
  AuthClassroomsJoinedClassroomIdAssignmentsAssignmentIdAcceptImport.update({
    path: "/assignments/$assignmentId/accept",
    getParentRoute: () => AuthClassroomsJoinedClassroomIdRouteRoute,
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
      preLoaderRoute: typeof AuthClassroomsRouteImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/create": {
      preLoaderRoute: typeof AuthClassroomsCreateImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/joined/$classroomId": {
      preLoaderRoute: typeof AuthClassroomsJoinedClassroomIdRouteImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/create/modal": {
      preLoaderRoute: typeof AuthClassroomsCreateModalImport
      parentRoute: typeof AuthClassroomsRouteImport
    }
    "/_auth/classrooms/joined/": {
      preLoaderRoute: typeof AuthClassroomsJoinedIndexImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/owned/": {
      preLoaderRoute: typeof AuthClassroomsOwnedIndexImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/joined/$classroomId/teams": {
      preLoaderRoute: typeof AuthClassroomsJoinedClassroomIdTeamsRouteImport
      parentRoute: typeof AuthClassroomsJoinedClassroomIdRouteImport
    }
    "/_auth/classrooms/owned/$classroomId/invite": {
      preLoaderRoute: typeof AuthClassroomsOwnedClassroomIdInviteImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/joined/$classroomId/": {
      preLoaderRoute: typeof AuthClassroomsJoinedClassroomIdIndexImport
      parentRoute: typeof AuthClassroomsJoinedClassroomIdRouteImport
    }
    "/_auth/classrooms/owned/$classroomId/": {
      preLoaderRoute: typeof AuthClassroomsOwnedClassroomIdIndexImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/joined/$classroomId/invitations/$invitationId": {
      preLoaderRoute: typeof AuthClassroomsJoinedClassroomIdInvitationsInvitationIdImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/owned/$classroomId/assignments/create": {
      preLoaderRoute: typeof AuthClassroomsOwnedClassroomIdAssignmentsCreateImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/joined/$classroomId/assignments/$assignmentId/accept": {
      preLoaderRoute: typeof AuthClassroomsJoinedClassroomIdAssignmentsAssignmentIdAcceptImport
      parentRoute: typeof AuthClassroomsJoinedClassroomIdRouteImport
    }
    "/_auth/classrooms/joined/$classroomId/teams/join/": {
      preLoaderRoute: typeof AuthClassroomsJoinedClassroomIdTeamsJoinIndexImport
      parentRoute: typeof AuthClassroomsJoinedClassroomIdTeamsRouteImport
    }
    "/_auth/classrooms/owned/$classroomId/assignments/$assignmentId/": {
      preLoaderRoute: typeof AuthClassroomsOwnedClassroomIdAssignmentsAssignmentIdIndexImport
      parentRoute: typeof AuthImport
    }
  }
}

// Create and export the route tree

export const routeTree = rootRoute.addChildren([
  IndexRoute,
  AuthRoute.addChildren([
    AuthClassroomsRouteRoute.addChildren([AuthClassroomsCreateModalRoute]),
    AuthClassroomsCreateRoute,
    AuthClassroomsJoinedClassroomIdRouteRoute.addChildren([
      AuthClassroomsJoinedClassroomIdTeamsRouteRoute.addChildren([
        AuthClassroomsJoinedClassroomIdTeamsJoinIndexRoute,
      ]),
      AuthClassroomsJoinedClassroomIdIndexRoute,
      AuthClassroomsJoinedClassroomIdAssignmentsAssignmentIdAcceptRoute,
    ]),
    AuthClassroomsJoinedIndexRoute,
    AuthClassroomsOwnedIndexRoute,
    AuthClassroomsOwnedClassroomIdInviteRoute,
    AuthClassroomsOwnedClassroomIdIndexRoute,
    AuthClassroomsJoinedClassroomIdInvitationsInvitationIdRoute,
    AuthClassroomsOwnedClassroomIdAssignmentsCreateRoute,
    AuthClassroomsOwnedClassroomIdAssignmentsAssignmentIdIndexRoute,
  ]),
  LoginRoute,
])
