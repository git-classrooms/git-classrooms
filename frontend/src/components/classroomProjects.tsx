import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Calendar,
  ChevronRight,
  Clock,
  ExternalLink,
  FileCode2,
  GitBranch,
  Play,
  XCircle,
} from "lucide-react";
import { cn, formatDate, formatDateWithTime } from "@/lib/utils";
import { Link } from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";
import { projectsQueryOptions } from "@/api/project";
import { ProjectResponse, UserClassroomResponse } from "@/swagger-client";
import { Status } from "@/types/projects";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import { classroomQueryOptions } from "@/api/classroom";
import { StatusBadge } from "./ui/status-badge";
import { useTranslation } from "react-i18next";

export function ProjectListSection({ classroomId }: { classroomId: string }): JSX.Element {
  const { t } = useTranslation("assignment");
  const { data: projects } = useSuspenseQuery(projectsQueryOptions(classroomId));
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));

  return (
    <Card className="border-border/50">
      <CardHeader className="pb-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center">
            <FileCode2 className="w-5 h-5 text-primary" />
          </div>
          <div>
            <CardTitle className="text-lg font-mono">{t("projects.myAssignments")}</CardTitle>
            <CardDescription>{t("projects.yourAssignments")}</CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        {projects.length === 0 ? (
          <EmptyState />
        ) : (
          <div className="space-y-3">
            {projects.map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                userClassroom={userClassroom}
              />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function EmptyState() {
  const { t } = useTranslation("assignment");

  return (
    <div className="py-12 text-center">
      <div className="w-16 h-16 mx-auto mb-4 rounded-2xl bg-muted/50 flex items-center justify-center">
        <FileCode2 className="w-8 h-8 text-muted-foreground/50" />
      </div>
      <p className="text-muted-foreground">{t("projects.empty.title")}</p>
      <p className="text-sm text-muted-foreground/70 mt-1">
        {t("projects.empty.description")}
      </p>
    </div>
  );
}

function ProjectCard({
  project,
  userClassroom,
}: {
  project: ProjectResponse;
  userClassroom: UserClassroomResponse;
}) {
  const { t } = useTranslation("assignment");
  const isClosed = project.assignment.dueDate && new Date(project.assignment.dueDate) < new Date();
  const isPending = project.projectStatus === Status.Pending;
  const isFailed = project.projectStatus === Status.Failed;
  const isAccepted = project.projectStatus === Status.Accepted;
  const isCreating = project.projectStatus === Status.Creating;

  const dueDate = project.assignment.dueDate ? new Date(project.assignment.dueDate) : null;
  const now = new Date();
  const daysUntilDue = dueDate ? Math.ceil((dueDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)) : null;
  const isUrgent = daysUntilDue !== null && daysUntilDue <= 3 && daysUntilDue > 0;
  const isVeryUrgent = daysUntilDue !== null && daysUntilDue <= 1 && daysUntilDue > 0;

  const getStatusVariant = () => {
    if (isClosed) return "neutral";
    if (isAccepted) return "success";
    if (isPending) return "info";
    if (isFailed) return "destructive";
    if (isCreating) return "warning";
    return "neutral";
  };

  const getStatusLabel = () => {
    if (isClosed) return t("projects.status.closed");
    if (isAccepted) return t("projects.status.accepted");
    if (isPending) return t("projects.status.pending");
    if (isFailed) return t("projects.status.failed");
    if (isCreating) return t("projects.status.creating");
    return t("projects.status.unknown");
  };

  return (
    <div
      className={cn(
        "group relative p-4 rounded-lg border transition-all duration-200",
        "hover:bg-muted/30 hover:border-border",
        isPending || isFailed
          ? "border-primary/30 bg-primary/5"
          : "border-border/50 bg-card/50"
      )}
    >
      <div className="flex flex-col sm:flex-row sm:items-center gap-4">
        {/* Assignment Info */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start gap-3">
            <div className="w-10 h-10 rounded-lg bg-muted flex items-center justify-center shrink-0">
              <GitBranch className="w-5 h-5 text-muted-foreground" />
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 flex-wrap">
                <h3 className="font-medium truncate">{project.assignment.name}</h3>
                <StatusBadge variant={getStatusVariant()} size="sm" showDot>
                  {getStatusLabel()}
                </StatusBadge>
              </div>
              {project.assignment.description && (
                <p className="text-sm text-muted-foreground mt-0.5 line-clamp-1">
                  {project.assignment.description}
                </p>
              )}

              {/* Meta Info */}
              <div className="flex items-center gap-4 mt-2 text-xs text-muted-foreground">
                <span className="flex items-center gap-1.5">
                  <Calendar className="w-3.5 h-3.5" />
                  {formatDate(project.createdAt)}
                </span>
                {dueDate && (
                  <span
                    className={cn(
                      "flex items-center gap-1.5",
                      isVeryUrgent && !isClosed && "text-destructive font-medium",
                      isUrgent && !isVeryUrgent && !isClosed && "text-warning font-medium"
                    )}
                  >
                    <Clock className="w-3.5 h-3.5" />
                    {isClosed ? (
                      t("projects.status.closed")
                    ) : isVeryUrgent ? (
                      t("dueDate.dueToday")
                    ) : isUrgent ? (
                      t("dueDate.dueIn", { count: daysUntilDue })
                    ) : (
                      formatDateWithTime(project.assignment.dueDate!)
                    )}
                  </span>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-2 sm:ml-auto">
          {isAccepted ? (
            <AcceptedActions project={project} userClassroom={userClassroom} />
          ) : isPending || isFailed ? (
            <PendingActions project={project} userClassroom={userClassroom} isFailed={isFailed} />
          ) : isCreating ? (
            <CreatingState />
          ) : null}
        </div>
      </div>
    </div>
  );
}

function AcceptedActions({
  project,
  userClassroom,
}: {
  project: ProjectResponse;
  userClassroom: UserClassroomResponse;
}) {
  const { t } = useTranslation("assignment");

  return (
    <div className="flex items-center gap-2">
      <Tooltip>
        <TooltipTrigger asChild>
          <Button variant="outline" size="sm" asChild className="gap-2">
            <a href={project.webUrl} target="_blank" rel="noopener noreferrer">
              <ExternalLink className="h-4 w-4" />
              <span className="hidden sm:inline">{t("projects.viewCode")}</span>
            </a>
          </Button>
        </TooltipTrigger>
        <TooltipContent>{t("projects.viewCodeTooltip")}</TooltipContent>
      </Tooltip>

      {userClassroom.classroom.studentsViewAllProjects && (
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="secondary" size="sm" asChild className="gap-2">
              <Link
                to="/classrooms/$classroomId/assignments/$assignmentId"
                params={{
                  classroomId: userClassroom.classroom.id,
                  assignmentId: project.assignment.id,
                }}
              >
                <span className="hidden sm:inline">{t("projects.details")}</span>
                <ChevronRight className="h-4 w-4" />
              </Link>
            </Button>
          </TooltipTrigger>
          <TooltipContent>{t("projects.viewDetailsTooltip")}</TooltipContent>
        </Tooltip>
      )}
    </div>
  );
}

function PendingActions({
  project,
  userClassroom,
  isFailed,
}: {
  project: ProjectResponse;
  userClassroom: UserClassroomResponse;
  isFailed: boolean;
}) {
  const { t } = useTranslation("assignment");

  return (
    <div className="flex items-center gap-2">
      {isFailed && (
        <Tooltip>
          <TooltipTrigger>
            <div className="flex items-center gap-1.5 text-destructive text-xs">
              <XCircle className="w-4 h-4" />
              <span className="hidden sm:inline">{t("projects.setupFailed")}</span>
            </div>
          </TooltipTrigger>
          <TooltipContent>
            {t("projects.setupFailedTooltip")}
          </TooltipContent>
        </Tooltip>
      )}

      <Button variant="glow" size="sm" asChild className="gap-2">
        <Link
          to="/classrooms/$classroomId/projects/$projectId/accept"
          params={{
            classroomId: userClassroom.classroom.id,
            projectId: project.id,
          }}
        >
          <Play className="h-4 w-4" />
          <span>{isFailed ? t("projects.retry") : t("projects.accept")}</span>
          <ChevronRight className="h-4 w-4 -ml-1" />
        </Link>
      </Button>
    </div>
  );
}

function CreatingState() {
  const { t } = useTranslation("assignment");

  return (
    <div className="flex items-center gap-2 text-warning">
      <div className="relative w-4 h-4">
        <div className="absolute inset-0 rounded-full border-2 border-warning/30" />
        <div className="absolute inset-0 rounded-full border-2 border-warning border-t-transparent animate-spin" />
      </div>
      <span className="text-sm font-medium">{t("projects.status.settingUp")}</span>
    </div>
  );
}
