import { createFileRoute, redirect, useNavigate, Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import {
  AlertCircle,
  ArrowLeft,
  Calendar,
  CheckCircle2,
  Clock,
  ExternalLink,
  GitBranch,
  Loader2,
  Play,
  Sparkles,
} from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { useSuspenseQuery } from "@tanstack/react-query";
import { projectQueryOptions, useAcceptAssignment } from "@/api/project";
import { classroomQueryOptions } from "@/api/classroom";
import { cn, formatDate, formatDateWithTime, getDaysUntilDue, isStudent } from "@/lib/utils";
import { Card, CardContent } from "@/components/ui/card";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { useTranslation } from "react-i18next";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/projects/$projectId/accept")({
  loader: async ({ context: { queryClient }, params }) => {
    const project = await queryClient.ensureQueryData(projectQueryOptions(params.classroomId, params.projectId));
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    if (!isStudent(userClassroom) || userClassroom.team?.id !== project.team.id) {
      throw redirect({
        to: "/classrooms/$classroomId",
        search: { tab: "assignments" },
        replace: true,
        params,
      });
    }
    return { project };
  },
  component: AcceptAssignment,
});

function AcceptAssignment() {
  const { t } = useTranslation(["assignment", "classroom"]);
  const navigate = useNavigate({
    from: "/_auth/classrooms/$classroomId/projects/$projectId/accept/",
  });
  const { classroomId, projectId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: project } = useSuspenseQuery(projectQueryOptions(classroomId, projectId, 10000));
  const { mutateAsync, isError, isPending } = useAcceptAssignment(classroomId, projectId);

  const dueDate = project.assignment.dueDate;
  const daysUntil = getDaysUntilDue(dueDate);
  const isUrgent = daysUntil !== null && daysUntil <= 3;
  const isVeryUrgent = daysUntil !== null && daysUntil <= 1;

  const onClick = async () => {
    await mutateAsync();
    await navigate({ to: "/classrooms/$classroomId", params: { classroomId }, search: { tab: "assignments" } });
  };

  return (
    <div className="min-h-[80vh] flex flex-col items-center justify-center px-4">
      <div className="w-full max-w-xl">
        {/* Breadcrumb */}
        <Breadcrumb className="mb-6">
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link to="/classrooms">{t("classroom:title")}</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId }}>
                  {classroom.classroom.name}
                </Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>{t("acceptPage.title")}</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>

        {/* Main Card */}
        <Card className="relative overflow-hidden border-2 border-primary/20">
          {/* Decorative gradient background */}
          <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-transparent pointer-events-none" />
          <div className="absolute top-0 right-0 w-64 h-64 opacity-[0.03] pointer-events-none">
            <Sparkles className="w-full h-full" />
          </div>

          <CardContent className="relative p-8">
            {/* Icon */}
            <div className="flex justify-center mb-6">
              <div className="relative">
                <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center shadow-lg shadow-primary/25">
                  <Play className="w-10 h-10 text-primary-foreground" />
                </div>
                <div className="absolute -bottom-1 -right-1 w-8 h-8 rounded-full bg-success flex items-center justify-center border-4 border-background">
                  <CheckCircle2 className="w-4 h-4 text-success-foreground" />
                </div>
              </div>
            </div>

            {/* Title */}
            <div className="text-center mb-8">
              <h1 className="text-2xl font-bold tracking-tight mb-2">{t("acceptPage.title")}</h1>
              <p className="text-muted-foreground">
                {t("acceptPage.subtitle")}
              </p>
            </div>

            {/* Assignment Details Card */}
            <div className="bg-muted/30 rounded-xl p-5 mb-6 border border-border/50">
              <div className="flex items-start gap-4">
                <div className="w-12 h-12 rounded-lg bg-background border border-border flex items-center justify-center shrink-0">
                  <GitBranch className="w-6 h-6 text-muted-foreground" />
                </div>
                <div className="flex-1 min-w-0">
                  <h2 className="font-semibold text-lg truncate">{project.assignment.name}</h2>
                  <p className="text-sm text-muted-foreground">{classroom.classroom.name}</p>

                  {project.assignment.description && (
                    <p className="text-sm text-muted-foreground mt-2 line-clamp-2">
                      {project.assignment.description}
                    </p>
                  )}

                  {/* Meta info */}
                  <div className="flex flex-wrap items-center gap-x-4 gap-y-2 mt-3 text-xs text-muted-foreground">
                    <span className="flex items-center gap-1.5">
                      <Calendar className="w-3.5 h-3.5" />
                      {t("acceptPage.created")} {formatDate(project.assignment.createdAt)}
                    </span>
                    {dueDate && (
                      <span
                        className={cn(
                          "flex items-center gap-1.5",
                          isVeryUrgent && "text-destructive font-medium",
                          isUrgent && !isVeryUrgent && "text-warning font-medium"
                        )}
                      >
                        <Clock className="w-3.5 h-3.5" />
                        {isVeryUrgent
                          ? t("dueDate.dueToday")
                          : daysUntil === 1
                            ? t("dueDate.dueTomorrow")
                            : `${t("acceptPage.due")} ${formatDateWithTime(dueDate)}`}
                      </span>
                    )}
                  </div>
                </div>
              </div>
            </div>

            {/* Info text */}
            <div className="flex items-start gap-3 p-4 rounded-lg bg-primary/5 border border-primary/10 mb-6">
              <ExternalLink className="w-5 h-5 text-primary shrink-0 mt-0.5" />
              <div className="text-sm">
                <p className="font-medium text-foreground">{t("acceptPage.whatHappensNext")}</p>
                <p className="text-muted-foreground mt-1">
                  {t("acceptPage.whatHappensNextDescription")}
                </p>
              </div>
            </div>

            {/* Error Alert */}
            {isError && (
              <Alert variant="destructive" className="mb-6">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>{t("acceptPage.error")}</AlertTitle>
                <AlertDescription>
                  {t("acceptPage.errorDescription")}
                </AlertDescription>
              </Alert>
            )}

            {/* Actions */}
            <div className="flex flex-col-reverse sm:flex-row items-center gap-3">
              <Button variant="outline" className="w-full sm:w-auto" asChild>
                <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId }}>
                  <ArrowLeft className="w-4 h-4 mr-2" />
                  {t("acceptPage.goBack")}
                </Link>
              </Button>
              <Button
                variant="glow"
                size="lg"
                className="w-full sm:flex-1 gap-2 font-semibold"
                onClick={onClick}
                disabled={isPending}
              >
                {isPending ? (
                  <>
                    <Loader2 className="w-4 h-4 animate-spin" />
                    {t("acceptPage.settingUp")}
                  </>
                ) : (
                  <>
                    <CheckCircle2 className="w-4 h-4" />
                    {t("acceptPage.acceptAndStart")}
                  </>
                )}
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Footer hint */}
        <p className="text-center text-xs text-muted-foreground mt-4">
          {t("acceptPage.classroomNote")}
        </p>
      </div>
    </div>
  );
}
