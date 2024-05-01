import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { QueryClient } from "@tanstack/react-query";
import React, { Suspense } from "react";
import { ThemeProvider } from "@/provider/themeProvider.tsx";
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
      <ThemeProvider defaultTheme="system" storageKey="gitlab-classrooms-theme">
        <div className="w-screen h-screen overflow-scroll">
          <Navbar auth={auth} />
          <div className="max-w-2xl m-auto">
            <Outlet />
            <ReactQueryDevtools initialIsOpen={false} />
            <Suspense>
              <TanStackRouterDevtools />
            </Suspense>
          </div>
        </div>
      </ThemeProvider>
    </CsrfProvider>
  );
}
