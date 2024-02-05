// This file is auto-generated by TanStack Router

// Import Routes

import { Route as rootRoute } from "./routes/__root"
import { Route as LoginImport } from "./routes/login"
import { Route as AuthImport } from "./routes/_auth"
import { Route as IndexImport } from "./routes/index"
import { Route as AuthClassroomsIndexImport } from "./routes/_auth/classrooms/index"
import { Route as AuthClassroomsCreateImport } from "./routes/_auth/classrooms/create"
import { Route as AuthClassroomsClassroomIdImport } from "./routes/_auth/classrooms/$classroomId"
import { Route as AuthClassroomsClassroomIdIndexImport } from "./routes/_auth/classrooms/$classroomId/index"
import { Route as AuthClassroomsClassroomIdInviteImport } from "./routes/_auth/classrooms/$classroomId/invite"

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

const AuthClassroomsIndexRoute = AuthClassroomsIndexImport.update({
  path: "/classrooms/",
  getParentRoute: () => AuthRoute,
} as any)

const AuthClassroomsCreateRoute = AuthClassroomsCreateImport.update({
  path: "/classrooms/create",
  getParentRoute: () => AuthRoute,
} as any)

const AuthClassroomsClassroomIdRoute = AuthClassroomsClassroomIdImport.update({
  path: "/classrooms/$classroomId",
  getParentRoute: () => AuthRoute,
} as any)

const AuthClassroomsClassroomIdIndexRoute =
  AuthClassroomsClassroomIdIndexImport.update({
    path: "/",
    getParentRoute: () => AuthClassroomsClassroomIdRoute,
  } as any)

const AuthClassroomsClassroomIdInviteRoute =
  AuthClassroomsClassroomIdInviteImport.update({
    path: "/invite",
    getParentRoute: () => AuthClassroomsClassroomIdRoute,
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
    "/_auth/classrooms/$classroomId": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/create": {
      preLoaderRoute: typeof AuthClassroomsCreateImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/": {
      preLoaderRoute: typeof AuthClassroomsIndexImport
      parentRoute: typeof AuthImport
    }
    "/_auth/classrooms/$classroomId/invite": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdInviteImport
      parentRoute: typeof AuthClassroomsClassroomIdImport
    }
    "/_auth/classrooms/$classroomId/": {
      preLoaderRoute: typeof AuthClassroomsClassroomIdIndexImport
      parentRoute: typeof AuthClassroomsClassroomIdImport
    }
  }
}

// Create and export the route tree

export const routeTree = rootRoute.addChildren([
  IndexRoute,
  AuthRoute.addChildren([
    AuthClassroomsClassroomIdRoute.addChildren([
      AuthClassroomsClassroomIdInviteRoute,
      AuthClassroomsClassroomIdIndexRoute,
    ]),
    AuthClassroomsCreateRoute,
    AuthClassroomsIndexRoute,
  ]),
  LoginRoute,
])
