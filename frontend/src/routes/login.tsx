import { createFileRoute, redirect } from "@tanstack/react-router";
import GitlabLogo from "./../assets/gitlab_logo.svg";
import { Button } from "@/components/ui/button";
import { useCsrf } from "@/provider/csrfProvider";

export const Route = createFileRoute("/login")({
  validateSearch: (search: Record<string, unknown>) => {
    return {
      redirect: (search.redirect as string) || "",
    };
  },
  beforeLoad: async ({ context }) => {
    if (context.auth) {
      throw redirect({
        to: "/classrooms",
      });
    }
  },
  component: Login,
});

function Login() {
  const { csrfToken } = useCsrf();
  const { redirect } = Route.useSearch();

  return (
    <div className="flex flex-col items-center">
      <img src={GitlabLogo} className="h-96 w-96" />
      <form method="POST" action="/api/v1/auth/sign-in">
        <input type="hidden" name="redirect" value={redirect} />
        <input type="hidden" name="csrf_token" value={csrfToken} />
        <Button type="submit">Login with Gitlab</Button>
      </form>
    </div>
  );
}
