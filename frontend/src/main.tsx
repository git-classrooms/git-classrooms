import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { RouterProvider, createRouteMask, createRouter } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

// Import the generated route tree
import { routeTree } from "./routeTree.gen";

const classroomCreateModalToClassroomCreateMask = createRouteMask({
  routeTree,
  from: "/classrooms/create/modal",
  to: "/classrooms/create",
  params: true,
});

// Create a new router instance
const router = createRouter({
  routeTree,
  context: {
    queryClient,
  },
  routeMasks: [classroomCreateModalToClassroomCreateMask],
  defaultPreload: "intent",
  unmaskOnReload: true,
  // Since we're using React Query, we don't want loader calls to ever be stale
  // This will ensure that the loader is always called when the route is preloaded or visited
  defaultPreloadStaleTime: 0,
});

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </React.StrictMode>,
);
