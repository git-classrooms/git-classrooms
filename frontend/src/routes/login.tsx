import { createFileRoute, redirect } from "@tanstack/react-router";
import GitlabLogo from "./../assets/gitlab_logo.svg";
import { Button } from "@/components/ui/button";
import { useCsrf } from "@/provider/csrfProvider";
import { Separator } from "@/components/ui/separator.tsx";
import { GitBranch as GitBranchIcon, Mail as MailIcon } from "lucide-react";
import { useSuspenseQuery } from "@tanstack/react-query";
import { gitlabInfoQueryOptions } from "@/api/info.ts";

export const Route = createFileRoute("/login")({
  validateSearch: (search: Record<string, unknown>) => {
    return {
      redirect: (search.redirect as string) || "",
    };
  },
  beforeLoad: async ({ context }) => {
    if (context.auth) {
      throw redirect({
        to: "/dashboard",
      });
    }
  },
  loader: async ({ context: { queryClient}}) => {
    const gitlabInfo = await queryClient.ensureQueryData(gitlabInfoQueryOptions)
    return { gitlabInfo }
  },
  component: Login,
});

function Login() {
  const { csrfToken } = useCsrf();
  const { redirect } = Route.useSearch();
  const { data } = useSuspenseQuery(gitlabInfoQueryOptions);

  return (
    <div className="m-auto max-w-lg ">
      <div className="flex justify-center">
        <img src={GitlabLogo} className="max-w-xs" />
      </div>

      <div className="p-6 rounded-lg border flex flex-col gap-5">
        <h1 className="text-5xl font-bold text-center mb-5">Login</h1>
        <p className="text-slate-500 text-lg">
          Authenticate Classrooms with your GitLab account at <span
          className="text-slate-900  dark:text-slate-300 font-bold">{data.gitlabUrl}</span>.
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
