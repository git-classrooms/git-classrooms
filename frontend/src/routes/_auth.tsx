import { Outlet, createFileRoute, redirect } from "@tanstack/react-router";
import { isAuthenticated } from '@/lib/utils'

export const Route = createFileRoute("/_auth")({
  beforeLoad: async () => {
    if (!await isAuthenticated()) {
      throw redirect({
        to: '/login',
        search: {
          // Use the current location to power a redirect after login
          // (Do not use `router.state.resolvedLocation` as it can
          // potentially lag behind the actual current location)
          redirect: location.href,
        },
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
