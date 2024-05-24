import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { QueryClient } from "@tanstack/react-query";
import React, { Suspense } from "react";
import { authCsrfQueryOptions } from "@/api/auth.ts";
import { Loader } from "@/components/loader.tsx";
import { CsrfProvider } from "@/provider/csrfProvider";
import { Navbar } from "@/components/navbar.tsx";
import { GetMeResponse } from "@/swagger-client";

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
  auth: GetMeResponse | null;
}>()({
  component: RootComponent,
  loader: ({ context }) => context.queryClient.ensureQueryData(authCsrfQueryOptions),
  pendingComponent: Loader,
});

const TanStackRouterDevtools =
  process.env.NODE_ENV === "production"
    ? () => null // Render nothing in production
    : React.lazy(() =>
        // Lazy load in development
        import("@tanstack/router-devtools").then((res) => ({
          default: res.TanStackRouterDevtools,
          // For Embedded Mode
          // default: res.TanStackRouterDevtoolsPanel
        })),
      );

function RootComponent() {
  const { auth } = Route.useRouteContext();
  return (
    <CsrfProvider>
      <div className="min-w-screen min-h-screen">
        <Navbar auth={auth} />
        <div className="flex flex-col w-full items-center">
          <div className="w-full xl:max-w-[90rem]">
            <div className="mx-6 md:px-10">
              <Outlet />
            </div>
            <ReactQueryDevtools initialIsOpen={false} />
            <Suspense>
              <TanStackRouterDevtools />
            </Suspense>
          </div>
        </div>
      </div>
    </CsrfProvider>
  );
}
