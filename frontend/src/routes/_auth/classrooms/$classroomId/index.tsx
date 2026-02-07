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
import { ReportApiAxiosParamCreator } from "@/swagger-client";
import { Button } from "@/components/ui/button";
import {
  Archive,
  CalendarClock,
  ChevronRight,
  Clipboard,
  ClipboardList,
  Download,
  ExternalLink,
  FileText,
  Settings,
  Users,
} from "lucide-react";
import { useArchiveClassroom } from "@/api/classroom";
import {
  AlertDialog,
  AlertDialogTrigger,
  AlertDialogContent,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogCancel,
  AlertDialogAction,
  AlertDialogHeader,
  AlertDialogFooter,
} from "@/components/ui/alert-dialog";
import { formatRelativeTime, isModerator, isStudent } from "@/lib/utils";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { z } from "zod";
import { projectsQueryOptions } from "@/api/project";
import { ProjectListSection } from "@/components/classroomProjects";
import { toast } from "sonner";
import { StatusBadge } from "@/components/ui/status-badge";
import { cn } from "@/lib/utils";

const tabs = ["assignments", "members", "teams"] as const;
const tabSchema = z.enum(tabs);

export const Route = createFileRoute("/_auth/classrooms/$classroomId/")({
  validateSearch: z.object({ tab: tabSchema.catch("assignments") }),
  component: ClassroomDetail,
  loader: async ({ context: { queryClient }, params }) => {
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));

    const { url: reportDownloadUrl } = await ReportApiAxiosParamCreator().getClassroomReport(params.classroomId);
    const teamsReportUrls = (
      await Promise.all(
        teams.map(async (team) => ({
          teamId: team.id,
          url: (await ReportApiAxiosParamCreator().getClassroomTeamReport(params.classroomId, team.id)).url,
        })),
      )
    ).reduce((acc, { url, teamId }) => acc.set(teamId, url), new Map<string, string>());

    const members = await queryClient.ensureQueryData(membersQueryOptions(params.classroomId));

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
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { tab } = Route.useSearch();
  const { reportDownloadUrl } = Route.useLoaderData();
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));
  const { mutate } = useArchiveClassroom(classroomId);

  const { teamsReportUrls } = Route.useLoaderData();
  const router = useRouter();

  const handleConfirmArchive = () => {
    mutate();
  };

  const handleCopyInviteLink = () => {
    const path = router.buildLocation({
      to: "/classrooms/$classroomId/invitations/$invitationId",
      params: { classroomId: userClassroom.classroom.id, invitationId: userClassroom.inviteCode },
      search: { groupLink: true }
    });
    navigator.clipboard.writeText(`${location.origin}${path.href}`);
    toast.success("Invite link copied to clipboard");
  };

  return (
    <div className="space-y-6 pb-8">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-2 text-sm text-muted-foreground animate-stagger-1">
        <Link to="/classrooms" className="hover:text-foreground transition-colors">
          Classrooms
        </Link>
        <ChevronRight className="w-4 h-4" />
        <span className="text-foreground font-medium">{userClassroom.classroom.name}</span>
      </nav>

      {/* Header Section */}
      <header className="animate-stagger-1">
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
                    {userClassroom.classroom.archived ? "Archived" : "Active"}
                  </StatusBadge>
                </div>
                <a
                  href={userClassroom.webUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  View on GitLab
                  <ExternalLink className="w-3 h-3" />
                </a>
              </div>
            </div>

            {userClassroom.classroom.description && (
              <p className="text-muted-foreground max-w-2xl">
                {userClassroom.classroom.description}
              </p>
            )}
          </div>

          {/* Action Buttons */}
          {!userClassroom.classroom.archived && isModerator(userClassroom) && (
            <div className="flex flex-wrap gap-2">
              <Button variant="outline" size="sm" onClick={handleCopyInviteLink}>
                <Clipboard className="w-4 h-4 mr-2" />
                Copy Invite
              </Button>
              <Button variant="outline" size="sm" asChild>
                <a href={reportDownloadUrl} target="_blank" rel="noopener noreferrer">
                  <Download className="w-4 h-4 mr-2" />
                  Report
                </a>
              </Button>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="outline" size="sm">
                    <Archive className="w-4 h-4 mr-2" />
                    Archive
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Archive this classroom?</AlertDialogTitle>
                    <AlertDialogDescription>
                      This action cannot be undone. The classroom will be marked as archived
                      and no new assignments can be created.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction onClick={handleConfirmArchive} variant="destructive">
                      Archive Classroom
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
              <Button variant="outline" size="sm" asChild>
                <Link to="/classrooms/$classroomId/settings/" params={{ classroomId }}>
                  <Settings className="w-4 h-4 mr-2" />
                  Settings
                </Link>
              </Button>
            </div>
          )}
        </div>
      </header>

      {/* Stats Cards */}
      <section className="animate-stagger-2">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <StatCard
            icon={<Users className="w-4 h-4" />}
            label="Members"
            value={classroomMembers.length}
          />
          <StatCard
            icon={<ClipboardList className="w-4 h-4" />}
            label="Assignments"
            value={userClassroom.assignmentsCount}
          />
          <StatCard
            icon={<Users className="w-4 h-4" />}
            label="Teams"
            value={teams.length}
          />
          <StatCard
            icon={<CalendarClock className="w-4 h-4" />}
            label="Created"
            value={formatRelativeTime(userClassroom.classroom.createdAt)}
            isText
          />
        </div>
      </section>

      {/* Tabs Section */}
      <section className="animate-stagger-3">
        <Tabs value={tab} className="w-full">
          <TabsList>
            <TabsTrigger asChild value="assignments">
              <Link search={{ tab: "assignments" }}>
                <FileText className="w-4 h-4 mr-2" />
                Assignments
              </Link>
            </TabsTrigger>
            <TabsTrigger asChild value="members">
              <Link search={{ tab: "members" }}>
                <Users className="w-4 h-4 mr-2" />
                Members
              </Link>
            </TabsTrigger>
            {userClassroom.classroom.maxTeamSize > 1 && (
              <TabsTrigger asChild value="teams">
                <Link search={{ tab: "teams" }}>
                  <Users className="w-4 h-4 mr-2" />
                  Teams
                </Link>
              </TabsTrigger>
            )}
          </TabsList>

          <TabsContent value="assignments">
            {isModerator(userClassroom) && (
              <AssignmentListSection
                classroomId={classroomId}
                deactivateInteraction={userClassroom.classroom.archived}
              />
            )}
            {isStudent(userClassroom) && <ProjectListSection classroomId={classroomId} />}
          </TabsContent>

          <TabsContent value="members">
            <MemberListCard
              teamsReportUrls={teamsReportUrls}
              classroomMembers={classroomMembers}
              classroomId={classroomId}
              userClassroom={userClassroom}
              showTeams={userClassroom.classroom.maxTeamSize > 1}
              deactivateInteraction={userClassroom.classroom.archived}
            />
          </TabsContent>

          {userClassroom.classroom.maxTeamSize > 1 && (
            <TabsContent value="teams">
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
