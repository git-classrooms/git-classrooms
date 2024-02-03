import { Outlet, createFileRoute, redirect } from "@tanstack/react-router";
import { isAuthenticated } from '@/lib/utils'

export const Route = createFileRoute("/_auth")({
  beforeLoad: async ({ location }) => {
    if (!await isAuthenticated()) {
      throw redirect({
        to: '/login',
        search: {
          redirect: location.href
        }
      })

    }
  },
  component: Index,
});

function Index() {
  return (
    <Outlet />
  );
}
