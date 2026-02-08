import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, Link, useRouter } from "@tanstack/react-router";
import { MemberListCard } from "@/components/classroomMembers";
import { TeamListCard } from "@/components/classroomTeams";
import { AssignmentListSection } from "@/components/classroomAssignments";
import { classroomQueryOptions } from "@/api/classroom";
import { assignmentsQueryOptions } from "@/api/assignment";
import { membersQueryOptions } from "@/api/member";
import { teamsQueryOptions } from "@/api/team";
import { ReportApiAxiosParamCreator, TeamResponse, UserClassroomResponse } from "@/swagger-client";
import { Button } from "@/components/ui/button";
import {
  CalendarClock,
  CheckCircle2,
  ChevronRight,
  Clipboard,
  ClipboardList,
  Clock,
  Download,
  ExternalLink,
  FileText,
  Play,
  Settings,
  Users,
} from "lucide-react";
import { formatRelativeTime, isModerator, isOwner, isStudent } from "@/lib/utils";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { z } from "zod";
import { projectsQueryOptions } from "@/api/project";
import { useMemo } from "react";
import { ProjectListSection } from "@/components/classroomProjects";
import { toast } from "sonner";
import { StatusBadge } from "@/components/ui/status-badge";
import { Markdown } from "@/components/ui/markdown";
import { cn } from "@/lib/utils";
import { useTranslation } from "react-i18next";

const tabs = ["assignments", "members", "teams"] as const;
const tabSchema = z.enum(tabs);

export const Route = createFileRoute("/_auth/classrooms/$classroomId/")({
  validateSearch: z.object({ tab: tabSchema.catch("assignments") }),
  component: ClassroomDetail,
  loader: async ({ context: { queryClient }, params }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));

    // Students can only see members/teams if "Mutual Code View" is enabled
    const canViewMembersAndTeams = isModerator(userClassroom) || userClassroom.classroom.studentsViewAllProjects;

    let teams: TeamResponse[] = [];
    let members: UserClassroomResponse[] = [];
    let teamsReportUrls = new Map<string, string>();
    let reportDownloadUrl = "";

    if (canViewMembersAndTeams) {
      teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
      members = await queryClient.ensureQueryData(membersQueryOptions(params.classroomId));

      if (isModerator(userClassroom)) {
        const reportResult = await ReportApiAxiosParamCreator().getClassroomReport(params.classroomId);
        reportDownloadUrl = reportResult.url;
        teamsReportUrls = (
          await Promise.all(
            teams.map(async (team) => ({
              teamId: team.id,
              url: (await ReportApiAxiosParamCreator().getClassroomTeamReport(params.classroomId, team.id)).url,
            })),
          )
        ).reduce((acc, { url, teamId }) => acc.set(teamId, url), new Map<string, string>());
      }
    }

    if (isModerator(userClassroom)) {
      const assignments = await queryClient.ensureQueryData(assignmentsQueryOptions(params.classroomId));
      return { userClassroom, assignments, members, teams, reportDownloadUrl, teamsReportUrls };
    } else {
      const projects = await queryClient.ensureQueryData(projectsQueryOptions(params.classroomId));
      return { userClassroom, projects, members, teams, reportDownloadUrl, teamsReportUrls };
    }
  },
  pendingComponent: Loader,
});

function ClassroomDetail() {
  const { t } = useTranslation("classroom");
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { tab } = Route.useSearch();
  const { reportDownloadUrl } = Route.useLoaderData();
  const router = useRouter();

  // Students can only see members/teams if "Mutual Code View" is enabled
  const canViewMembersAndTeams = isModerator(userClassroom) || userClassroom.classroom.studentsViewAllProjects;

  // Redirect to assignments tab if user doesn't have access to current tab
  const effectiveTab = useMemo(() => {
    if (!canViewMembersAndTeams && (tab === "members" || tab === "teams")) {
      return "assignments";
    }
    return tab;
  }, [canViewMembersAndTeams, tab]);

  const { teamsReportUrls } = Route.useLoaderData();

  const handleCopyInviteLink = () => {
    const path = router.buildLocation({
      to: "/classrooms/$classroomId/invitations/$invitationId",
      params: { classroomId: userClassroom.classroom.id, invitationId: userClassroom.inviteCode },
      search: { groupLink: true }
    });
    navigator.clipboard.writeText(`${location.origin}${path.href}`);
    toast.success(t("detail.inviteCopied"));
  };

  return (
    <div className="space-y-6 pb-8">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-2 text-sm text-muted-foreground">
        <Link to="/classrooms" className="hover:text-foreground transition-colors">
          {t("title")}
        </Link>
        <ChevronRight className="w-4 h-4" />
        <span className="text-foreground font-medium">{userClassroom.classroom.name}</span>
      </nav>

      {/* Header Section */}
      <header>
        <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-4">
          <div className="space-y-3">
            <div className="flex items-center gap-3">
              <ClassroomAvatar name={userClassroom.classroom.name} />
              <div>
                <div className="flex items-center gap-3">
                  <h1 className="text-2xl font-bold tracking-tight">
                    {userClassroom.classroom.name}
                  </h1>
                  <StatusBadge
                    variant={userClassroom.classroom.archived ? "neutral" : "success"}
                    showDot
                  >
                    {userClassroom.classroom.archived ? t("status.archived") : t("status.active")}
                  </StatusBadge>
                </div>
                <a
                  href={userClassroom.webUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  {t("detail.viewOnGitlab")}
                  <ExternalLink className="w-3 h-3" />
                </a>
              </div>
            </div>

            {userClassroom.classroom.description && (
              <div className="text-muted-foreground max-w-2xl">
                <Markdown>{userClassroom.classroom.description}</Markdown>
              </div>
            )}
          </div>

          {/* Action Buttons */}
          {!userClassroom.classroom.archived && isModerator(userClassroom) && (
            <div className="flex flex-wrap gap-2">
              <Button variant="outline" size="sm" onClick={handleCopyInviteLink}>
                <Clipboard className="w-4 h-4 mr-2" />
                {t("detail.copyInvite")}
              </Button>
              {isOwner(userClassroom) && (
                <Button variant="outline" size="sm" asChild>
                  <a href={reportDownloadUrl} target="_blank" rel="noopener noreferrer">
                    <Download className="w-4 h-4 mr-2" />
                    {t("detail.report")}
                  </a>
                </Button>
              )}
              {isOwner(userClassroom) && (
                <Button variant="outline" size="sm" asChild>
                  <Link to="/classrooms/$classroomId/settings/" params={{ classroomId }}>
                    <Settings className="w-4 h-4 mr-2" />
                    {t("settings.title")}
                  </Link>
                </Button>
              )}
            </div>
          )}
        </div>
      </header>

      {/* Stats Cards - Different for Moderators and Students */}
      <section>
        {isModerator(userClassroom) ? (
          <ModeratorStatsCards classroomId={classroomId} userClassroom={userClassroom} />
        ) : (
          <StudentStatsCards
            classroomId={classroomId}
            userClassroom={userClassroom}
          />
        )}
      </section>

      {/* Tabs Section */}
      <section>
        <Tabs value={effectiveTab} className="w-full">
          <TabsList>
            <TabsTrigger asChild value="assignments">
              <Link search={{ tab: "assignments" }}>
                <FileText className="w-4 h-4 mr-2" />
                {t("tabs.assignments")}
              </Link>
            </TabsTrigger>
            {canViewMembersAndTeams && (
              <TabsTrigger asChild value="members">
                <Link search={{ tab: "members" }}>
                  <Users className="w-4 h-4 mr-2" />
                  {t("tabs.members")}
                </Link>
              </TabsTrigger>
            )}
            {canViewMembersAndTeams && userClassroom.classroom.maxTeamSize > 1 && (
              <TabsTrigger asChild value="teams">
                <Link search={{ tab: "teams" }}>
                  <Users className="w-4 h-4 mr-2" />
                  {t("tabs.teams")}
                </Link>
              </TabsTrigger>
            )}
          </TabsList>

          <TabsContent value="assignments">
            {isModerator(userClassroom) && (
              <AssignmentListSection
                classroomId={classroomId}
                deactivateInteraction={userClassroom.classroom.archived || !isOwner(userClassroom)}
              />
            )}
            {isStudent(userClassroom) && <ProjectListSection classroomId={classroomId} />}
          </TabsContent>

          {canViewMembersAndTeams && (
            <TabsContent value="members">
              <MembersTabContent
                teamsReportUrls={teamsReportUrls}
                classroomId={classroomId}
                userClassroom={userClassroom}
              />
            </TabsContent>
          )}

          {canViewMembersAndTeams && userClassroom.classroom.maxTeamSize > 1 && (
            <TabsContent value="teams">
              <TeamsTabContent
                classroomId={classroomId}
                userClassroom={userClassroom}
                teamsReportUrls={teamsReportUrls}
              />
            </TabsContent>
          )}
        </Tabs>
      </section>

      <Outlet />
    </div>
  );
}

function ClassroomAvatar({ name }: { name: string }) {
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
        "w-12 h-12 rounded-xl flex items-center justify-center",
        "bg-gradient-to-br",
        gradient
      )}
    >
      <span className="font-mono font-bold text-xl text-primary-foreground">
        {name.charAt(0).toUpperCase()}
      </span>
    </div>
  );
}

function StatCard({
  icon,
  label,
  value,
  isText = false,
}: {
  icon: React.ReactNode;
  label: string;
  value: string | number;
  isText?: boolean;
}) {
  return (
    <div className="bg-card border border-border rounded-lg p-4 transition-all duration-200 hover:border-primary/30">
      <div className="flex items-center gap-2 text-muted-foreground mb-2">
        {icon}
        <span className="text-xs font-medium uppercase tracking-wide">{label}</span>
      </div>
      <div className={cn("font-bold", isText ? "text-lg" : "text-2xl font-mono")}>
        {value}
      </div>
    </div>
  );
}

// Stats cards for moderators - fetches members/teams for reactive counts
function ModeratorStatsCards({
  classroomId,
  userClassroom,
}: {
  classroomId: string;
  userClassroom: UserClassroomResponse;
}) {
  const { t } = useTranslation("classroom");
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      <StatCard
        icon={<Users className="w-4 h-4" />}
        label={t("stats.members")}
        value={classroomMembers.length}
      />
      <StatCard
        icon={<ClipboardList className="w-4 h-4" />}
        label={t("stats.assignments")}
        value={userClassroom.assignmentsCount}
      />
      <StatCard
        icon={<Users className="w-4 h-4" />}
        label={t("stats.teams")}
        value={teams.length}
      />
      <StatCard
        icon={<CalendarClock className="w-4 h-4" />}
        label={t("stats.created")}
        value={formatRelativeTime(userClassroom.classroom.createdAt)}
        isText
      />
    </div>
  );
}

// Wrapper component for Members tab - fetches its own data for reactivity
function MembersTabContent({
  teamsReportUrls,
  classroomId,
  userClassroom,
}: {
  teamsReportUrls: Map<string, string>;
  classroomId: string;
  userClassroom: UserClassroomResponse;
}) {
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));

  return (
    <MemberListCard
      teamsReportUrls={teamsReportUrls}
      classroomMembers={classroomMembers}
      classroomId={classroomId}
      userClassroom={userClassroom}
      showTeams={userClassroom.classroom.maxTeamSize > 1}
      deactivateInteraction={userClassroom.classroom.archived}
    />
  );
}

// Wrapper component for Teams tab - fetches its own data for reactivity
function TeamsTabContent({
  classroomId,
  userClassroom,
  teamsReportUrls,
}: {
  classroomId: string;
  userClassroom: UserClassroomResponse;
  teamsReportUrls: Map<string, string>;
}) {
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));

  return (
    <TeamListCard
      teams={teams}
      studentsCanCreateTeams={userClassroom.classroom.createTeams}
      classroomId={classroomId}
      userClassroom={userClassroom}
      maxTeamSize={userClassroom.classroom.maxTeamSize}
      numInvitedMembers={classroomMembers.filter(isStudent).length}
      deactivateInteraction={userClassroom.classroom.archived}
      teamsReportUrls={teamsReportUrls}
    />
  );
}

function StudentStatsCards({
  classroomId,
  userClassroom,
}: {
  classroomId: string;
  userClassroom: UserClassroomResponse;
}) {
  const { t } = useTranslation("classroom");
  const { t: ta } = useTranslation("assignment");
  const { data: projects } = useSuspenseQuery(projectsQueryOptions(classroomId));

  const acceptedCount = projects.filter((p) => p.projectStatus === "accepted").length;
  const pendingCount = projects.filter((p) => p.projectStatus === "pending").length;

  const nextDueProject = projects
    .filter((p) => p.assignment.dueDate && new Date(p.assignment.dueDate) > new Date())
    .sort((a, b) => new Date(a.assignment.dueDate!).getTime() - new Date(b.assignment.dueDate!).getTime())[0];

  const daysUntilDue = nextDueProject?.assignment.dueDate
    ? Math.ceil((new Date(nextDueProject.assignment.dueDate).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24))
    : null;

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      <StatCard
        icon={<CheckCircle2 className="w-4 h-4" />}
        label={t("stats.completed")}
        value={acceptedCount}
      />
      <StatCard
        icon={<Play className="w-4 h-4" />}
        label={t("stats.pending")}
        value={pendingCount}
      />
      {userClassroom.classroom.maxTeamSize > 1 && userClassroom.team && (
        <StatCard
          icon={<Users className="w-4 h-4" />}
          label={t("stats.myTeam")}
          value={userClassroom.team.name}
          isText
        />
      )}
      <StatCard
        icon={<Clock className="w-4 h-4" />}
        label={t("stats.nextDue")}
        value={daysUntilDue !== null ? (daysUntilDue <= 0 ? ta("dueDate.today") : `${daysUntilDue}d`) : "-"}
        isText
      />
    </div>
  );
}
