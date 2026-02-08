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
  Loader2,
  MoreHorizontal,
  Scale,
  Send,
  Settings,
  Users,
} from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { cn, createCloneScript, formatDate, formatDateWithTime, getDaysUntilDue, isModerator, isOwner } from "@/lib/utils.ts";
import { assignmentCloneUrlsQueryOptions, assignmentQueryOptions } from "@/api/assignment";
import { assignmentProjectsQueryOptions, useInviteToAssignment } from "@/api/project";
import { Assignment, ProjectResponse, ReportApiAxiosParamCreator, UserClassroomResponse } from "@/swagger-client";
import { classroomQueryOptions } from "@/api/classroom";
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
import { ApiProjectCloneUrlResponse } from "@/swagger-client";

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

    const urls = (
      await Promise.all(
        assignmentProjects.map(async (project) => ({
          url: (await ReportApiAxiosParamCreator().getClassroomTeamReport(classroomId, project.teamId)).url,
          projectId: project.id,
        })),
      )
    ).reduce((acc, { url, projectId }) => acc.set(projectId, url), new Map<string, string>());

    return { classroom, assignment, assignmentProjects, reportDownloadUrl, urls, cloneUrls };
  },
  component: AssignmentDetail,
  pendingComponent: Loader,
});

function AssignmentDetail() {
  const { classroomId, assignmentId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  const { data: assignmentProjects } = useSuspenseQuery(assignmentProjectsQueryOptions(classroomId, assignmentId));

  const [showStats, setShowStats] = useLocalStorage("assignment-stats-open", true);
  const { data: cloneUrls } = useSuspenseQuery(assignmentCloneUrlsQueryOptions(classroomId, assignmentId));

  const { mutateAsync, isError, isPending } = useInviteToAssignment(classroomId, assignmentId);

  const acceptedCount = assignmentProjects.filter((p) => p.projectStatus === "accepted").length;
  const pendingCount = assignmentProjects.filter((p) => p.projectStatus === "pending").length;
  const totalCount = assignmentProjects.length;
  const progressPercent = totalCount > 0 ? Math.round((acceptedCount / totalCount) * 100) : 0;

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
              <Link to="/classrooms">Classrooms</Link>
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
                {assignment.closed ? "Closed" : isOverdue ? "Overdue" : isUrgent ? "Due Soon" : "Open"}
              </StatusBadge>
            </div>
            <p className="text-muted-foreground">
              {classroom.classroom.maxTeamSize === 1 ? "Individual assignment" : "Team assignment"} ·{" "}
              {totalCount} project{totalCount !== 1 ? "s" : ""}
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
                Grading
              </Link>
            </Button>
          )}
          {isOwner(classroom) && (
            <Button variant="outline" size="sm" asChild>
              <Link
                to="/classrooms/$classroomId/assignments/$assignmentId/settings/"
                params={{ classroomId, assignmentId }}
              >
                <Settings className="w-4 h-4 mr-2" />
                Settings
              </Link>
            </Button>
          )}
        </div>
      </div>

      {/* Collapsible Stats Section */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">Overview</h2>
          <Button variant="ghost" size="sm" className="h-7 px-2" onClick={() => setShowStats(!showStats)}>
            {showStats ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
            <span className="ml-1 text-xs">{showStats ? "Hide" : "Show"}</span>
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
                    <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">Progress</p>
                    <p className="text-2xl font-mono font-bold">{progressPercent}%</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {acceptedCount} of {totalCount} accepted
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
                    <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">Due Date</p>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <p className={cn(
                          "text-2xl font-mono font-bold",
                          isOverdue && "text-destructive",
                          isUrgent && !isOverdue && "text-warning"
                        )}>
                          {assignment.dueDate ? formatDate(new Date(assignment.dueDate)) : "No deadline"}
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
                        {daysUntil === 0 ? "Due today" : daysUntil === 1 ? "Due tomorrow" :
                         isOverdue ? `${Math.abs(daysUntil!)} days overdue` :
                         `${daysUntil} days remaining`}
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
                    <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">Last Updated</p>
                    <p className="text-2xl font-mono font-bold">
                      {assignment.updatedAt ? formatDistanceToNow(new Date(assignment.updatedAt), { addSuffix: false }) : "—"}
                    </p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {assignment.updatedAt ? "ago" : "No activity yet"}
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
                    <p className="text-xs text-muted-foreground uppercase tracking-wide mb-1">Teams</p>
                    <div className="flex items-baseline gap-2">
                      <p className="text-2xl font-mono font-bold text-success">{acceptedCount}</p>
                      {pendingCount > 0 && (
                        <p className="text-lg font-mono text-warning">+{pendingCount}</p>
                      )}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      {pendingCount > 0 ? `${pendingCount} pending` : "All accepted"}
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
                <p className="text-xs text-muted-foreground uppercase tracking-wide mb-2">Description</p>
                <p className="text-sm leading-relaxed">{assignment.description}</p>
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
              <h2 className="font-semibold">Projects</h2>
              <p className="text-sm text-muted-foreground">
                {classroom.classroom.maxTeamSize === 1
                  ? "Individual student repositories"
                  : "Team repositories"}
              </p>
            </div>
          </div>

          {isModerator(classroom) && (
            <Tooltip delayDuration={0}>
              <TooltipTrigger asChild>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => mutateAsync()}
                  disabled={isPending || pendingCount === 0}
                >
                  {isPending ? (
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  ) : (
                    <Send className="w-4 h-4 mr-2" />
                  )}
                  Send Invites
                  {pendingCount > 0 && (
                    <span className="ml-2 px-1.5 py-0.5 text-xs bg-warning/20 text-warning rounded">
                      {pendingCount}
                    </span>
                  )}
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                {pendingCount > 0
                  ? `Send invitations to ${pendingCount} pending team${pendingCount !== 1 ? "s" : ""}`
                  : "All teams have accepted"}
              </TooltipContent>
            </Tooltip>
          )}
        </div>

        {assignmentProjects.length === 0 ? (
          <Card className="border-dashed border-border/50">
            <CardContent className="p-8 text-center">
              <GitFork className="w-10 h-10 mx-auto text-muted-foreground mb-3" />
              <h3 className="font-medium mb-1">No projects yet</h3>
              <p className="text-sm text-muted-foreground">
                Projects will appear here when teams are invited to this assignment
              </p>
            </CardContent>
          </Card>
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

      {isError && (
        <Alert variant="destructive" className="mt-6">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>The invitation could not be sent. Please try again.</AlertDescription>
        </Alert>
      )}
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
  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline" size="sm">
          <FolderGit2 className="w-4 h-4 mr-2" />
          Clone All
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-72" align="end">
        <div className="space-y-3">
          <div>
            <h4 className="font-medium text-sm">Clone all projects</h4>
            <p className="text-xs text-muted-foreground mt-1">
              Copy a shell script to clone all {assignmentProjects.length} project{assignmentProjects.length !== 1 ? "s" : ""}
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
                  toast.success("SSH clone script copied");
                }}
              >
                <ClipboardCopy className="w-4 h-4 mr-2" />
                Copy SSH script
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
                  toast.success("HTTPS clone script copied");
                }}
              >
                <ClipboardCopy className="w-4 h-4 mr-2" />
                Copy HTTPS script
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
  const { urls } = Route.useLoaderData();
  const reportUrl = urls.get(project.id)!;

  const isAccepted = project.projectStatus === "accepted";
  const isPending = project.projectStatus === "pending";

  const statusConfig = {
    accepted: { variant: "success" as const, label: "Accepted" },
    pending: { variant: "warning" as const, label: "Pending" },
    creating: { variant: "info" as const, label: "Creating" },
    failed: { variant: "destructive" as const, label: "Failed" },
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
                Invited {formatDistanceToNow(new Date(project.createdAt), { addSuffix: true })}
              </p>
            </div>
          </div>
          <StatusBadge variant={status.variant} size="sm" showDot>
            {status.label}
          </StatusBadge>
        </div>

        {/* Actions */}
        <div className="flex items-center justify-between pt-3 border-t border-border/50">
          <div className="flex items-center gap-1">
            {isAccepted && (
              <Button variant="ghost" size="sm" className="h-8 px-2" asChild>
                <a href={project.webUrl} target="_blank" rel="noopener noreferrer">
                  <ExternalLink className="w-3.5 h-3.5 mr-1" />
                  Open
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
                  Open in GitLab
                </a>
              </DropdownMenuItem>
              {isModerator(classroom) && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem disabled={!isAccepted} asChild>
                    <a href={reportUrl} target="_blank" rel="noopener noreferrer">
                      <Download className="w-4 h-4 mr-2" />
                      Download Report
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
