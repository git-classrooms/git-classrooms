import { projectsQueryOptions } from "@/api/project";
import { ProjectResponse, UserClassroomResponse } from "@/swagger-client";
import { useQueries } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { AlertTriangle, ArrowRight, Clock, Play, Sparkles } from "lucide-react";
import { Button } from "./ui/button";
import { cn, getDaysUntilDue } from "@/lib/utils";
import { Status } from "@/types/projects";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";

interface PendingProject extends ProjectResponse {
  classroomId: string;
  classroomName: string;
}

export function PendingAssignmentsBanner({
  studentClassrooms,
}: {
  studentClassrooms: UserClassroomResponse[];
}) {
  const { t } = useTranslation("assignment");
  const projectQueries = useQueries({
    queries: studentClassrooms.map((userClassroom) => ({
      ...projectsQueryOptions(userClassroom.classroom.id),
      select: (projects: ProjectResponse[]) =>
        projects
          .filter(
            (p) =>
              p.projectStatus === Status.Pending ||
              p.projectStatus === Status.Failed
          )
          .map((p) => ({
            ...p,
            classroomId: userClassroom.classroom.id,
            classroomName: userClassroom.classroom.name,
          })),
    })),
    combine: (results) => {
      const allPending: PendingProject[] = [];
      results.forEach((result) => {
        if (result.data) {
          allPending.push(...result.data);
        }
      });
      return {
        data: allPending,
        isPending: results.some((r) => r.isPending),
        isError: results.some((r) => r.isError),
      };
    },
  });

  const sortedPending = useMemo(() => {
    return [...projectQueries.data].sort((a, b) => {
      const dueDateA = a.assignment.dueDate;
      const dueDateB = b.assignment.dueDate;
      if (!dueDateA && !dueDateB) return 0;
      if (!dueDateA) return 1;
      if (!dueDateB) return -1;
      return new Date(dueDateA).getTime() - new Date(dueDateB).getTime();
    });
  }, [projectQueries.data]);

  if (projectQueries.isPending || sortedPending.length === 0) {
    return null;
  }

  const mostUrgent = sortedPending[0];
  const daysUntil = getDaysUntilDue(mostUrgent.assignment.dueDate);
  const isUrgent = daysUntil !== null && daysUntil <= 3;
  const isVeryUrgent = daysUntil !== null && daysUntil <= 1;
  const isFailed = mostUrgent.projectStatus === Status.Failed;

  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-xl border-2 p-6 transition-all duration-300",
        isFailed
          ? "border-destructive/50 bg-gradient-to-br from-destructive/10 via-destructive/5 to-background"
          : isVeryUrgent
            ? "border-warning/50 bg-gradient-to-br from-warning/10 via-warning/5 to-background"
            : isUrgent
              ? "border-warning/30 bg-gradient-to-br from-warning/5 via-background to-background"
              : "border-primary/50 bg-gradient-to-br from-primary/10 via-primary/5 to-background"
      )}
    >
      {/* Decorative elements */}
      <div className="absolute top-0 right-0 w-64 h-64 opacity-[0.03] pointer-events-none">
        <div
          className={cn(
            "absolute inset-0 rounded-full blur-3xl",
            isFailed
              ? "bg-destructive"
              : isUrgent
                ? "bg-warning"
                : "bg-primary"
          )}
        />
      </div>
      <div className="absolute -bottom-8 -left-8 w-32 h-32 opacity-[0.02] pointer-events-none">
        <Sparkles className="w-full h-full" />
      </div>

      <div className="relative flex flex-col lg:flex-row lg:items-center gap-6">
        {/* Icon */}
        <div
          className={cn(
            "shrink-0 w-16 h-16 rounded-2xl flex items-center justify-center",
            "shadow-lg",
            isFailed
              ? "bg-gradient-to-br from-destructive to-destructive/80 shadow-destructive/25"
              : isUrgent
                ? "bg-gradient-to-br from-warning to-warning/80 shadow-warning/25"
                : "bg-gradient-to-br from-primary to-primary/80 shadow-primary/25"
          )}
        >
          {isFailed ? (
            <AlertTriangle className="w-8 h-8 text-destructive-foreground" />
          ) : (
            <Play className="w-8 h-8 text-primary-foreground" />
          )}
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0 space-y-2">
          <div className="flex items-center gap-2 flex-wrap">
            <h2
              className={cn(
                "text-lg font-bold tracking-tight",
                isFailed
                  ? "text-destructive"
                  : isUrgent
                    ? "text-warning"
                    : "text-primary"
              )}
            >
              {isFailed
                ? t("pending.setupFailed")
                : sortedPending.length === 1
                  ? t("pending.actionRequired")
                  : t("pending.assignmentsWaiting", { count: sortedPending.length })}
            </h2>
            {sortedPending.length > 1 && (
              <span className="text-xs px-2 py-0.5 rounded-full bg-muted text-muted-foreground font-medium">
                {t("pending.more", { count: sortedPending.length - 1 })}
              </span>
            )}
          </div>

          <p className="text-foreground font-medium">
            {isFailed ? t("pending.retryAccepting") : t("pending.accept")}{" "}
            <span className="font-bold">{mostUrgent.assignment.name}</span>
            <span className="text-muted-foreground font-normal">
              {" "}
              {t("pending.in")} {mostUrgent.classroomName}
            </span>
          </p>

          {mostUrgent.assignment.dueDate && (
            <div
              className={cn(
                "flex items-center gap-1.5 text-sm",
                isVeryUrgent
                  ? "text-destructive font-semibold"
                  : isUrgent
                    ? "text-warning font-medium"
                    : "text-muted-foreground"
              )}
            >
              <Clock className="w-4 h-4" />
              {isVeryUrgent ? (
                <span>{t("pending.dueToday")}</span>
              ) : daysUntil === 1 ? (
                <span>{t("pending.dueTomorrow")}</span>
              ) : (
                <span>{t("dueDate.dueIn", { count: daysUntil! })}</span>
              )}
            </div>
          )}
        </div>

        {/* CTA */}
        <div className="flex items-center gap-3 shrink-0">
          <Button
            variant="glow"
            size="lg"
            asChild
            className={cn(
              "group gap-2 font-semibold shadow-lg",
              isFailed && "bg-destructive hover:bg-destructive/90",
              isUrgent && !isFailed && "bg-warning hover:bg-warning/90 text-warning-foreground"
            )}
          >
            <Link
              to="/classrooms/$classroomId/projects/$projectId/accept"
              params={{
                classroomId: mostUrgent.classroomId,
                projectId: mostUrgent.id,
              }}
            >
              <Play className="w-4 h-4" />
              {isFailed ? t("pending.retryNow") : t("pending.acceptNow")}
              <ArrowRight className="w-4 h-4 transition-transform group-hover:translate-x-1" />
            </Link>
          </Button>
        </div>
      </div>

      {/* Additional pending assignments preview */}
      {sortedPending.length > 1 && (
        <div className="relative mt-4 pt-4 border-t border-border/50">
          <p className="text-xs text-muted-foreground mb-2 uppercase tracking-wide font-medium">
            {t("pending.alsoWaiting")}
          </p>
          <div className="flex flex-wrap gap-2">
            {sortedPending.slice(1, 4).map((project) => (
              <Link
                key={project.id}
                to="/classrooms/$classroomId/projects/$projectId/accept"
                params={{
                  classroomId: project.classroomId,
                  projectId: project.id,
                }}
                className={cn(
                  "inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm",
                  "bg-muted/50 hover:bg-muted transition-colors",
                  "border border-transparent hover:border-border"
                )}
              >
                <span className="font-medium truncate max-w-[150px]">
                  {project.assignment.name}
                </span>
                {project.projectStatus === Status.Failed && (
                  <AlertTriangle className="w-3.5 h-3.5 text-destructive shrink-0" />
                )}
                <ArrowRight className="w-3.5 h-3.5 text-muted-foreground shrink-0" />
              </Link>
            ))}
            {sortedPending.length > 4 && (
              <span className="inline-flex items-center px-3 py-1.5 text-sm text-muted-foreground">
                {t("pending.more", { count: sortedPending.length - 4 })}
              </span>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
