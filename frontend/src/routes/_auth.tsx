import { Outlet, createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth")({
  beforeLoad: async ({ location, context }) => {
    if (!context.auth) {
      throw redirect({
        to: "/login",
        replace: true,
        search: {
          redirect: location.href,
        },
      });
    }
  },
  component: Index,
});

function Index() {
  return <Outlet />;
}
