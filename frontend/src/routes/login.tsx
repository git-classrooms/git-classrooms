import { createFileRoute, redirect } from "@tanstack/react-router";
import GitlabLogo from "./../assets/gitlab_logo.svg";
import { Button } from "@/components/ui/button";
import { useCsrf } from "@/provider/csrfProvider";
import { useSuspenseQuery } from "@tanstack/react-query";
import { gitlabInfoQueryOptions } from "@/api/info.ts";
import {
  ArrowRight,
  CheckCircle2,
  ExternalLink,
  GitBranch,
  GraduationCap,
  Lock,
  Mail,
  Shield,
  Sparkles,
} from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { useTranslation } from "react-i18next";

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
        replace: true,
      });
    }
  },
  loader: async ({ context: { queryClient } }) => {
    const gitlabInfo = await queryClient.ensureQueryData(gitlabInfoQueryOptions);
    return { gitlabInfo };
  },
  component: Login,
});

function Login() {
  const { t } = useTranslation("auth");
  const { csrfToken } = useCsrf();
  const { redirect: redirectUrl } = Route.useSearch();
  const { data } = useSuspenseQuery(gitlabInfoQueryOptions);

  const permissions = [
    {
      icon: Mail,
      label: t("permissions.email.title"),
      description: t("permissions.email.description"),
    },
    {
      icon: GitBranch,
      label: t("permissions.writeRepository.title"),
      description: t("permissions.writeRepository.description"),
    },
  ];

  return (
    <div className="min-h-[90vh] flex items-center justify-center px-4">
      <div className="w-full max-w-md">
        {/* Logo & Branding */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-gradient-to-br from-primary to-primary/80 shadow-lg shadow-primary/25 mb-6">
            <GraduationCap className="w-10 h-10 text-primary-foreground" />
          </div>
          <h1 className="text-3xl font-bold tracking-tight mb-2">
            {t("login.title")}
          </h1>
          <p className="text-muted-foreground">
            {t("login.subtitle")}
          </p>
        </div>

        {/* Main Card */}
        <Card className="relative overflow-hidden border-2 border-border/50">
          {/* Decorative elements */}
          <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-transparent pointer-events-none" />
          <div className="absolute top-0 right-0 w-48 h-48 opacity-[0.03] pointer-events-none">
            <Sparkles className="w-full h-full" />
          </div>

          <CardContent className="relative p-6 space-y-6">
            {/* GitLab Connection Info */}
            <div className="flex items-center gap-3 p-4 rounded-xl bg-muted/50 border border-border/50">
              <div className="w-12 h-12 rounded-lg bg-[#FC6D26]/10 flex items-center justify-center shrink-0">
                <img src={GitlabLogo} className="w-7 h-7" alt="GitLab" />
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium">{t("login.button")}</p>
                <a
                  href={data.gitlabUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-xs text-muted-foreground hover:text-primary transition-colors inline-flex items-center gap-1 truncate max-w-full"
                >
                  {data.gitlabUrl}
                  <ExternalLink className="w-3 h-3 shrink-0" />
                </a>
              </div>
            </div>

            {/* Permissions */}
            <div>
              <div className="flex items-center gap-2 mb-3">
                <Shield className="w-4 h-4 text-muted-foreground" />
                <p className="text-sm font-medium text-muted-foreground">
                  {t("permissions.title")}
                </p>
              </div>
              <div className="space-y-2">
                {permissions.map((permission) => (
                  <div
                    key={permission.label}
                    className={cn(
                      "flex items-center gap-3 p-3 rounded-lg",
                      "bg-background border border-border/50",
                      "transition-colors hover:border-border"
                    )}
                  >
                    <div className="w-8 h-8 rounded-md bg-primary/10 flex items-center justify-center shrink-0">
                      <permission.icon className="w-4 h-4 text-primary" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium">{permission.label}</p>
                      <p className="text-xs text-muted-foreground">
                        {permission.description}
                      </p>
                    </div>
                    <CheckCircle2 className="w-4 h-4 text-success shrink-0" />
                  </div>
                ))}
              </div>
            </div>

            {/* Login Button */}
            <form method="POST" action="/api/v1/auth/sign-in" className="space-y-4">
              <input type="hidden" name="redirect" value={redirectUrl} />
              <input type="hidden" name="csrf_token" value={csrfToken} />
              <Button
                type="submit"
                variant="glow"
                size="lg"
                className="w-full gap-2 font-semibold group"
              >
                <img src={GitlabLogo} className="w-5 h-5" alt="" />
                {t("login.button")}
                <ArrowRight className="w-4 h-4 transition-transform group-hover:translate-x-1" />
              </Button>
            </form>

            {/* Security note */}
            <div className="flex items-start gap-2 text-xs text-muted-foreground">
              <Lock className="w-3.5 h-3.5 mt-0.5 shrink-0" />
              <p>
                {t("login.securityNote")}
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Footer */}
        <p className="text-center text-xs text-muted-foreground mt-6">
          {t("login.termsNote")}
        </p>
      </div>
    </div>
  );
}
