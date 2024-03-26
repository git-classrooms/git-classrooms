import { Outlet, createFileRoute, redirect } from "@tanstack/react-router";
import { isAuthenticated } from "@/lib/utils";
import { LogoutButton } from "@/components/logoutButton";
import { Loader } from "@/components/loader";

export const Route = createFileRoute("/_auth")({
  beforeLoad: async ({ location }) => {
    if (!(await isAuthenticated())) {
      throw redirect({
        to: "/login",
        search: {
          redirect: location.href,
        },
      });
    }
  },
  component: Index,
  pendingComponent: Loader,
});

function Index() {
  return (
    <>
      <LogoutButton />
      <Outlet />
    </>
  );
}
