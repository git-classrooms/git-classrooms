import {
  createRootRouteWithContext,
  Outlet,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { QueryClient } from "@tanstack/react-query";

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
}>()({
  component: RootComponent,
});

function RootComponent() {
  return (
    <>
      <div className="w-screen h-screen overflow-scroll">
        <div className="max-w-2xl m-auto">
          <Outlet />
          <ReactQueryDevtools initialIsOpen={false} />
          <TanStackRouterDevtools />
        </div>
      </div>
    </>
  );
}
