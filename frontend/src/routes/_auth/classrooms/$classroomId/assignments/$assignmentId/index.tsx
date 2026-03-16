import { createFileRoute, Link } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Card, CardContent } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  AlertCircle,
  Calendar,
  ChevronDown,
  ChevronUp,
  ClipboardCopy,
  Clock,
  Download,
  ExternalLink,
  FolderGit2,
  GitFork,
  Info,
  Loader2,
  MoreHorizontal,
  Rocket,
  Scale,
  Settings,
  Sparkles,
  Users,
} from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { cn, createCloneScript, formatDate, formatDateWithTime, formatRelativeTime, getDaysUntilDue, getDateLocale, isModerator, isOwner } from "@/lib/utils.ts";
import { assignmentCloneUrlsQueryOptions, assignmentQueryOptions } from "@/api/assignment";
import { assignmentProjectsQueryOptions, useInviteToAssignment } from "@/api/project";
import { Assignment, ProjectResponse, ReportApiAxiosParamCreator, UserClassroomResponse } from "@/swagger-client";
import { classroomQueryOptions } from "@/api/classroom";
import { teamsQueryOptions } from "@/api/team";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { formatDistanceToNow } from "date-fns";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useLocalStorage } from "@/hooks/useLocalStorage";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { toast } from "sonner";
import { PopoverClose } from "@radix-ui/react-popover";
import { StatusBadge } from "@/components/ui/status-badge";
import { Markdown } from "@/components/ui/markdown";
import { ApiProjectCloneUrlResponse } from "@/swagger-client";
import { useTranslation } from "react-i18next";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/$assignmentId/")({
  loader: async ({ context: { queryClient }, params: { classroomId, assignmentId } }) => {
    const classroom = await queryClient.ensureQueryData(classroomQueryOptions(classroomId));
    const assignment = await queryClient.ensureQueryData(assignmentQueryOptions(classroomId, assignmentId));
    const assignmentProjects = await queryClient.ensureQueryData(
      assignmentProjectsQueryOptions(classroomId, assignmentId),
    );

    const { url: reportDownloadUrl } = await ReportApiAxiosParamCreator().getClassroomAssignmentReport(
      classroomId,
      assignmentId,
    );
    const cloneUrls = await queryClient.ensureQueryData(assignmentCloneUrlsQueryOptions(classroomId, assignmentId));
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(classroomId));

    const urls = (
      await Promise.all(
        assignmentProjects.map(async (project) => ({
          url: (await ReportApiAxiosParamCreator().getClassroomTeamReport(classroomId, project.teamId)).url,
          projectId: project.id,
        })),
      )
    ).reduce((acc, { url, projectId }) => acc.set(projectId, url), new Map<string, string>());

    return { classroom, assignment, assignmentProjects, reportDownloadUrl, urls, cloneUrls, teams };
  },
  component: AssignmentDetail,
  pendingComponent: Loader,
});

function AssignmentDetail() {
  const { t } = useTranslation("assignment");
  const { t: tc } = useTranslation("classroom");
  const { t: tco } = useTranslation("common");
  const { classroomId, assignmentId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  const { data: assignmentProjects } = useSuspenseQuery(assignmentProjectsQueryOptions(classroomId, assignmentId));

  const [showStats, setShowStats] = useLocalStorage("assignment-stats-open", true);
  const { data: cloneUrls } = useSuspenseQuery(assignmentCloneUrlsQueryOptions(classroomId, assignmentId));

  const { mutateAsync, isError, isPending } = useInviteToAssignment(classroomId, assignmentId);
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));

  const acceptedCount = assignmentProjects.filter((p) => p.projectStatus === "accepted").length;
  const pendingCount = assignmentProjects.filter((p) => p.projectStatus === "pending").length;
  const totalCount = assignmentProjects.length;
  const progressPercent = totalCount > 0 ? Math.round((acceptedCount / totalCount) * 100) : 0;

  const isReleased = totalCount > 0;
  const needsRelease = !isReleased && isOwner(classroom);

  const daysUntil = getDaysUntilDue(assignment.dueDate);
  const isOverdue = daysUntil !== null && daysUntil < 0;
  const isUrgent = daysUntil !== null && daysUntil >= 0 && daysUntil <= 3;

  return (
    <div>
      {/* Breadcrumb */}
      <Breadcrumb className="mb-6">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms">{tc("title")}</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link
                to="/classrooms/$classroomId"
                search={{ tab: "assignments" }}
                params={{ classroomId }}
              >
                {classroom.classroom.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{assignment.name}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {/* Header */}
      <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-4 mb-6">
        <div className="flex items-start gap-4">
          <div className="w-14 h-14 rounded-xl bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center shrink-0">
            <FolderGit2 className="w-7 h-7 text-primary" />
          </div>
          <div>
            <div className="flex items-center gap-3 mb-1">
              <h1 className="text-3xl font-bold font-mono tracking-tight">{assignment.name}</h1>
              <StatusBadge
                variant={assignment.closed ? "neutral" : isOverdue ? "destructive" : isUrgent ? "warning" : "success"}
                size="sm"
                showDot
              >
                {assignment.closed ? t("status.closed") : isOverdue ? t("status.overdue") : isUrgent ? t("status.dueSoon") : t("status.open")}
              </StatusBadge>
            </div>
            <p className="text-muted-foreground">
              {classroom.classroom.maxTeamSize === 1 ? t("detail.individual") : t("detail.team")} ·{" "}
              {t("projects.count", { count: totalCount })}
            </p>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-wrap items-center gap-2">
          <CloneProjectsPopover
            assignment={assignment}
            cloneUrls={cloneUrls}
            assignmentProjects={assignmentProjects}
          />
          {isModerator(classroom) && (
            <Button variant="outline" size="sm" asChild>
              <Link
                to="/classrooms/$classroomId/assignments/$assignmentId/grading"
                params={{ classroomId, assignmentId }}
              >
                <Scale className="w-4 h-4 mr-2" />
                {t("grading.title")}
              </Link>
            </Button>
          )}
          {isOwner(classroom) && (
            <Button variant="outline" size="sm" asChild>
              <Link
                to="/classrooms/$classroomId/assignments/$assignmentId/settings"
                params={{ classroomId, assignmentId }}
              >
                <Settings className="w-4 h-4 mr-2" />
                {t("settings.title")}
              </Link>
            </Button>
          )}
        </div>
      </div>

      {/* Release Banner */}
      {isError && (
        <Alert variant="destructive" className="mb-6">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>{tco("status.error")}</AlertTitle>
          <AlertDescription>{t("detail.releaseError")}</AlertDescription>
        </Alert>
      )}
      {needsRelease && (
        <div className="relative overflow-hidden rounded-xl border-2 border-primary/50 bg-gradient-to-br from-primary/10 via-primary/5 to-background p-6 mb-6 transition-all duration-300">
          {/* Decorative elements */}
          <div className="absolute top-0 right-0 w-64 h-64 opacity-[0.03] pointer-events-none">
            <div className="absolute inset-0 rounded-full blur-3xl bg-primary" />
          </div>
          <div className="absolute -bottom-8 -left-8 w-32 h-32 opacity-[0.02] pointer-events-none">
            <Sparkles className="w-full h-full" />
          </div>

          <div className="relative flex flex-col lg:flex-row lg:items-center gap-6">
            {/* Icon */}
            <div className="shrink-0 w-16 h-16 rounded-2xl flex items-center justify-center shadow-lg bg-gradient-to-br from-primary to-primary/80 shadow-primary/25">
              <Rocket className="w-8 h-8 text-primary-foreground" />
            </div>

            {/* Content */}
            <div className="flex-1 min-w-0 space-y-2">
              <h2 className="text-lg font-bold tracking-tight text-primary">
                {teams.length === 0 ? t("detail.release.noTeams") : t("detail.release.title")}
              </h2>
              <p className="text-sm text-muted-foreground">
                {teams.length === 0
                  ? t("detail.release.noTeamsDescription")
                  : t("detail.release.description", { count: teams.length })}
              </p>
            </div>

            {/* CTA */}
            {teams.length > 0 && (
              <div className="flex items-center gap-3 shrink-0">
                <Button
                  variant="glow"
                  size="lg"
                  className="group gap-2 font-semibold shadow-lg"
                  onClick={() => mutateAsync().catch(() => {})}
                  disabled={isPending}
                >
                  {isPending ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    <Rocket className="w-4 h-4" />
                  )}
                  {isPending ? t("detail.release.releasing") : t("detail.release.button")}
                  <span className="ml-1 px-1.5 py-0.5 text-xs bg-white/20 rounded">
                    {teams.length}
                  </span>
                </Button>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Collapsible Stats Section */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">{t("detail.overview")}</h2>
          <Button variant="ghost" size="sm" className="h-7 px-2" onClick={() => setShowStats(!showStats)}>
            {showStats ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
            <span className="ml-1 text-xs">{showStats ? t("detail.hide") : t("detail.show")}</span>
          </Button>
        </div>

        {showStats && (
          <>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {/* Progress Card */}
            <Card className="border-border/50 bg-gradient-to-br from-primary/5 to-transparent">
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">{t("detail.progress")}</p>
                    <p className="text-2xl font-mono font-bold">{progressPercent}%</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {t("detail.ofTotal", { accepted: acceptedCount, total: totalCount })}
                    </p>
                  </div>
                  <ProgressRing percent={progressPercent} size={56} />
                </div>
              </CardContent>
            </Card>

            {/* Due Date Card */}
            <Card className={cn(
              "border-border/50",
              isOverdue && "border-destructive/30 bg-destructive/5",
              isUrgent && !isOverdue && "border-warning/30 bg-warning/5"
            )}>
              <CardContent className="p-4">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">{t("dueDate.label")}</p>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <p className={cn(
                          "text-2xl font-mono font-bold",
                          isOverdue && "text-destructive",
                          isUrgent && !isOverdue && "text-warning"
                        )}>
                          {assignment.dueDate ? formatDate(new Date(assignment.dueDate)) : t("detail.noDeadline")}
                        </p>
                      </TooltipTrigger>
                      {assignment.dueDate && (
                        <TooltipContent>{formatDateWithTime(new Date(assignment.dueDate))}</TooltipContent>
                      )}
                    </Tooltip>
                    {assignment.dueDate && (
                      <p className={cn(
                        "text-xs mt-1",
                        isOverdue ? "text-destructive" : isUrgent ? "text-warning" : "text-muted-foreground"
                      )}>
                        {daysUntil === 0 ? t("detail.dueToday") : daysUntil === 1 ? t("detail.dueTomorrow") :
                         isOverdue ? t("detail.daysOverdue", { count: Math.abs(daysUntil!) }) :
                         t("detail.daysRemaining", { count: daysUntil! })}
                      </p>
                    )}
                  </div>
                  <div className={cn(
                    "w-10 h-10 rounded-lg flex items-center justify-center",
                    isOverdue ? "bg-destructive/10 text-destructive" :
                    isUrgent ? "bg-warning/10 text-warning" :
                    "bg-muted/50 text-muted-foreground"
                  )}>
                    <Calendar className="w-5 h-5" />
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Last Activity Card */}
            <Card className="border-border/50">
              <CardContent className="p-4">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">{t("detail.lastUpdated")}</p>
                    <p className="text-2xl font-mono font-bold">
                      {assignment.updatedAt ? formatDistanceToNow(new Date(assignment.updatedAt), { addSuffix: false, locale: getDateLocale() }) : "—"}
                    </p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {assignment.updatedAt ? t("detail.ago") : t("detail.noActivity")}
                    </p>
                  </div>
                  <div className="w-10 h-10 rounded-lg bg-muted/50 flex items-center justify-center text-muted-foreground">
                    <Clock className="w-5 h-5" />
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Teams Card */}
            <Card className="border-border/50">
              <CardContent className="p-4">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">{classroom.classroom.maxTeamSize === 1 ? t("grading.matrix.student") : t("grading.matrix.team")}</p>
                    <div className="flex items-baseline gap-2">
                      <p className="text-2xl font-mono font-bold text-success">{acceptedCount}</p>
                      {pendingCount > 0 && (
                        <p className="text-lg font-mono text-warning">+{pendingCount}</p>
                      )}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      {pendingCount > 0 ? t("detail.pending", { count: pendingCount }) : t("detail.allAccepted")}
                    </p>
                  </div>
                  <div className="w-10 h-10 rounded-lg bg-muted/50 flex items-center justify-center text-muted-foreground">
                    <Users className="w-5 h-5" />
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Description */}
          {assignment.description && (
            <Card className="mt-4 border-border/50">
              <CardContent className="p-4">
                <p className="text-xs text-muted-foreground uppercase tracking-wide mb-2">{t("form.description")}</p>
                <div className="text-sm leading-relaxed">
                  <Markdown>{assignment.description}</Markdown>
                </div>
              </CardContent>
            </Card>
          )}
        </>
        )}
      </div>

      {/* Projects Section */}
      <section>
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
              <GitFork className="w-5 h-5 text-primary" />
            </div>
            <div>
              <h2 className="font-semibold">{t("projects.title")}</h2>
              <p className="text-sm text-muted-foreground">
                {classroom.classroom.maxTeamSize === 1
                  ? t("detail.individualRepos")
                  : t("detail.teamRepos")}
              </p>
            </div>
          </div>

        </div>

        {assignmentProjects.length === 0 ? (
          needsRelease ? null : !isReleased && !isOwner(classroom) ? (
            <Card className="border-dashed border-border/50">
              <CardContent className="p-8 text-center">
                <Info className="w-10 h-10 mx-auto text-muted-foreground mb-3" />
                <h3 className="font-medium mb-1">{t("detail.release.notReleasedInfo")}</h3>
                <p className="text-sm text-muted-foreground">
                  {t("detail.release.notReleasedInfoDescription")}
                </p>
              </CardContent>
            </Card>
          ) : (
            <Card className="border-dashed border-border/50">
              <CardContent className="p-8 text-center">
                <GitFork className="w-10 h-10 mx-auto text-muted-foreground mb-3" />
                <h3 className="font-medium mb-1">{t("detail.noProjects")}</h3>
                <p className="text-sm text-muted-foreground">
                  {t("detail.noProjectsDescription")}
                </p>
              </CardContent>
            </Card>
          )
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {assignmentProjects.map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                classroom={classroom}
              />
            ))}
          </div>
        )}
      </section>

    </div>
  );
}

function ProgressRing({ percent, size = 56 }: { percent: number; size?: number }) {
  const strokeWidth = 4;
  const radius = (size - strokeWidth) / 2;
  const circumference = radius * 2 * Math.PI;
  const offset = circumference - (percent / 100) * circumference;

  return (
    <div className="relative" style={{ width: size, height: size }}>
      <svg className="transform -rotate-90" width={size} height={size}>
        <circle
          className="text-muted/30"
          strokeWidth={strokeWidth}
          stroke="currentColor"
          fill="transparent"
          r={radius}
          cx={size / 2}
          cy={size / 2}
        />
        <circle
          className="text-primary transition-all duration-500 ease-out"
          strokeWidth={strokeWidth}
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          strokeLinecap="round"
          stroke="currentColor"
          fill="transparent"
          r={radius}
          cx={size / 2}
          cy={size / 2}
        />
      </svg>
      <div className="absolute inset-0 flex items-center justify-center">
        <span className="text-xs font-mono font-semibold">{percent}%</span>
      </div>
    </div>
  );
}

function CloneProjectsPopover({
  assignment,
  cloneUrls,
  assignmentProjects,
}: {
  assignment: Assignment;
  cloneUrls: ApiProjectCloneUrlResponse[];
  assignmentProjects: ProjectResponse[];
}) {
  const { t } = useTranslation("assignment");

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline" size="sm">
          <FolderGit2 className="w-4 h-4 mr-2" />
          {t("detail.cloneAll")}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-72" align="end">
        <div className="space-y-3">
          <div>
            <h4 className="font-medium text-sm">{t("detail.cloneAllTitle")}</h4>
            <p className="text-xs text-muted-foreground mt-1">
              {t("detail.cloneAllDescription", { count: assignmentProjects.length })}
            </p>
          </div>
          <div className="grid gap-2">
            <PopoverClose asChild>
              <Button
                variant="outline"
                size="sm"
                className="w-full justify-start"
                onClick={() => {
                  navigator.clipboard.writeText(
                    createCloneScript("ssh", assignment, cloneUrls, assignmentProjects),
                  );
                  toast.success(t("detail.sshCopied"));
                }}
              >
                <ClipboardCopy className="w-4 h-4 mr-2" />
                {t("detail.copySSH")}
              </Button>
            </PopoverClose>
            <PopoverClose asChild>
              <Button
                variant="outline"
                size="sm"
                className="w-full justify-start"
                onClick={() => {
                  navigator.clipboard.writeText(
                    createCloneScript("https", assignment, cloneUrls, assignmentProjects),
                  );
                  toast.success(t("detail.httpsCopied"));
                }}
              >
                <ClipboardCopy className="w-4 h-4 mr-2" />
                {t("detail.copyHTTPS")}
              </Button>
            </PopoverClose>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
}

function ProjectCard({
  project,
  classroom,
}: {
  project: ProjectResponse;
  classroom: UserClassroomResponse;
}) {
  const { t } = useTranslation("assignment");
  const { urls } = Route.useLoaderData();
  const reportUrl = urls.get(project.id)!;

  const isAccepted = project.projectStatus === "accepted";
  const isPending = project.projectStatus === "pending";

  const statusConfig = {
    accepted: { variant: "success" as const, labelKey: "accepted" },
    pending: { variant: "warning" as const, labelKey: "pending" },
    creating: { variant: "info" as const, labelKey: "creating" },
    failed: { variant: "destructive" as const, labelKey: "failed" },
  };

  const status = statusConfig[project.projectStatus as keyof typeof statusConfig] || statusConfig.pending;

  return (
    <Card
      className={cn(
        "group transition-all duration-200 hover:shadow-lg hover:shadow-background/50",
        isAccepted && "hover:border-success/30",
        isPending && "hover:border-warning/30 opacity-80"
      )}
    >
      <CardContent className="p-4">
        <div className="flex items-start justify-between gap-3 mb-3">
          <div className="flex items-center gap-3 min-w-0">
            <TeamAvatar name={project.team.name} />
            <div className="min-w-0">
              <h3 className="font-semibold truncate">{project.team.name}</h3>
              <p className="text-xs text-muted-foreground">
                {t("detail.invited")} {formatRelativeTime(new Date(project.createdAt))}
              </p>
            </div>
          </div>
          <StatusBadge variant={status.variant} size="sm" showDot>
            {/* @ts-expect-error - dynamic key lookup */}
            {t(`projects.status.${status.labelKey}`)}
          </StatusBadge>
        </div>

        {/* Actions */}
        <div className="flex items-center justify-between pt-3 border-t border-border/50">
          <div className="flex items-center gap-1">
            {isAccepted && (
              <Button variant="ghost" size="sm" className="h-8 px-2" asChild>
                <a href={project.webUrl} target="_blank" rel="noopener noreferrer">
                  <ExternalLink className="w-3.5 h-3.5 mr-1" />
                  {t("detail.open")}
                </a>
              </Button>
            )}
          </div>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <MoreHorizontal className="w-4 h-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem disabled={!isAccepted} asChild>
                <a href={project.webUrl} target="_blank" rel="noopener noreferrer">
                  <ExternalLink className="w-4 h-4 mr-2" />
                  {t("detail.openInGitlab")}
                </a>
              </DropdownMenuItem>
              {isModerator(classroom) && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem disabled={!isAccepted} asChild>
                    <a href={reportUrl} target="_blank" rel="noopener noreferrer">
                      <Download className="w-4 h-4 mr-2" />
                      {t("detail.downloadReport")}
                    </a>
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardContent>
    </Card>
  );
}

function TeamAvatar({ name }: { name: string }) {
  const gradients = [
    "from-primary to-[hsl(280,100%,60%)]",
    "from-[hsl(142,71%,45%)] to-[hsl(185,100%,50%)]",
    "from-[hsl(38,92%,55%)] to-[hsl(0,72%,55%)]",
    "from-[hsl(280,65%,60%)] to-[hsl(210,100%,60%)]",
  ];
  const hash = name.split("").reduce((acc, char) => acc + char.charCodeAt(0), 0);
  const gradient = gradients[hash % gradients.length];

  return (
    <div
      className={cn(
        "w-10 h-10 rounded-lg flex items-center justify-center bg-gradient-to-br shrink-0",
        gradient
      )}
    >
      <span className="font-mono font-semibold text-primary-foreground">
        {name.charAt(0).toUpperCase()}
      </span>
    </div>
  );
}
