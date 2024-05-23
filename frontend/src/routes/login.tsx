import { createFileRoute, redirect } from "@tanstack/react-router";
import GitlabLogo from "./../assets/gitlab_logo.svg";
import { Button } from "@/components/ui/button";
import { useCsrf } from "@/provider/csrfProvider";
import { Separator } from "@/components/ui/separator.tsx";
import { Mail as MailIcon } from "lucide-react";
import { GitBranch as GitBranchIcon } from "lucide-react";

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
    <div>
      <img src={GitlabLogo} className="h-96 w-96" />
      <div className="p-6 rounded-lg border flex flex-col gap-5">
        <h1 className="text-5xl font-bold text-center mb-5">Login</h1>
        <p className="text-slate-500 text-lg">
          Authenticate Classrooms with your GitLab account at <span
          className="font-bold text-slate-900">gitlab.hs-flensburg.de</span>.
        </p>
        <Separator />
        <p className="text-slate-500">
          Classrooms will access to the following resources
        </p>
        <div className="flex flex-col gap-4">
          <div className="flex space-x-2 items-center">
            <MailIcon className="size-4"></MailIcon>
            <p className="text-sm">Email Address</p>
          </div>
          <div className="flex space-x-2 items-center">
            <GitBranchIcon className="size-4 align-middle"></GitBranchIcon>
            <p className="text-sm">Repository information</p>
          </div>
        </div>


        <Separator />
        <div className="flex justify-center">
          <form method="POST" action="/api/v1/auth/sign-in">
            <input type="hidden" name="redirect" value={redirect} />
            <input type="hidden" name="csrf_token" value={csrfToken} />
            <Button type="submit">
              <span className="mr-2 w-4 overflow-hidden">
                <img src={GitlabLogo} className="object-cover h-8" />
              </span>
              Login with Gitlab
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
}
