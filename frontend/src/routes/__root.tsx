import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { QueryClient } from "@tanstack/react-query";
import React, { Suspense } from "react";
import { ThemeProvider } from "@/provider/themeProvider.tsx";
import { ModeToggle } from "@/components/modeToggle.tsx";
import { authCsrfQueryOptions } from "@/api/auth.ts";
import { Loader } from "@/components/loader.tsx";
import { CsrfProvider } from "@/provider/csrfProvider";

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
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
  return (
    <CsrfProvider>
      <ThemeProvider defaultTheme="system" storageKey="gitlab-classrooms-theme">
        <div className="w-screen h-screen overflow-scroll">
          <ModeToggle />
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