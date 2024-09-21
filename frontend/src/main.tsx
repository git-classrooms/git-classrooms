import React, { Suspense } from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient();

// Import the generated route tree
import { routeTree } from "./routeTree.gen";
import { useAuth } from "./api/auth";
import { Loader } from "./components/loader";
import { ThemeProvider } from "./provider/themeProvider";
import { Toaster } from "./components/ui/sonner";
import { TooltipProvider } from "./components/ui/tooltip";
import { NotFound } from "./components/not-found";

// Create a new router instance
const router = createRouter({
  routeTree,
  context: {
    queryClient,
    auth: undefined!,
  },
  defaultPreload: "intent",
  // Since we're using React Query, we don't want loader calls to ever be stale
  // This will ensure that the loader is always called when the route is preloaded or visited
  defaultPreloadStaleTime: 0,
  defaultNotFoundComponent: NotFound,
});

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

function InnerApp() {
  const { data: auth } = useAuth();
  return <RouterProvider router={router} context={{ auth }} />;
}

function App() {
  return (
    <React.StrictMode>
      <ThemeProvider defaultTheme="system" storageKey="git-classrooms-theme">
        <TooltipProvider>
          <QueryClientProvider client={queryClient}>
            <Suspense fallback={<Loader />}>
              <InnerApp />
            </Suspense>
            <Toaster position="bottom-center" />
          </QueryClientProvider>
        </TooltipProvider>
      </ThemeProvider>
    </React.StrictMode>
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(<App />);
