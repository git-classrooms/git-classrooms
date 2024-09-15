import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, Link } from "@tanstack/react-router";
import { MemberListCard } from "@/components/classroomMembers.tsx";
import { TeamListCard } from "@/components/classroomTeams.tsx";
import { AssignmentListSection } from "@/components/classroomAssignments.tsx";
import { Header } from "@/components/header";
import { classroomQueryOptions } from "@/api/classroom";
import { assignmentsQueryOptions } from "@/api/assignment";
import { membersQueryOptions } from "@/api/member";
import { teamsQueryOptions } from "@/api/team";
import { ReportApiAxiosParamCreator, UserClassroomResponse } from "@/swagger-client";
import { Button } from "@/components/ui/button.tsx";
import {
  Archive,
  CalendarCheck2,
  CalendarClock,
  Download,
  ExternalLink,
  Eye,
  EyeOff,
  Info,
  Settings,
  Users,
} from "lucide-react";
import { useArchiveClassroom } from "@/api/classroom";
import { Text } from "lucide-react";
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
import { formatDate, isModerator, isStudent } from "@/lib/utils";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { formatDistanceToNow } from "date-fns";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { z } from "zod";
import { useLocalStorage } from "@/hooks/useLocalStorage";
import { projectsQueryOptions } from "@/api/project";
import { Breadcrumb, BreadcrumbItem, BreadcrumbList, BreadcrumbPage } from "@/components/ui/breadcrumb";
import { ProjectListSection } from "@/components/classroomProjects";

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

  const [showHeaderCards, setShowHeaderCards] = useLocalStorage("classroom-header", true);
  const toggleHeaderCards = () => setShowHeaderCards((old) => !old);
  const { teamsReportUrls } = Route.useLoaderData();

  const handleConfirmArchive = () => {
    mutate();
  };

  return (
    <>
      <Breadcrumb className="mb-5">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbPage>{userClassroom.classroom.name}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <div className="lg:flex justify-between gap-1 mb-4">
        <Header
          title={
            <a
              className="flex items-center"
              href={userClassroom.webUrl}
              target="_blank"
              referrerPolicy="no-referrer"
              title="Go to classroom"
            >
              {userClassroom.classroom.archived && "Archived "}
              {userClassroom.classroom.name}
              <ExternalLink className="h-4 w-4 ml-2" />
            </a>
          }
          subtitle="Classroom overview"
        />
        <div className="flex flex-col lg:flex-row gap-3">
          <Button
            variant="secondary"
            className="min-w-[137px]"
            onClick={toggleHeaderCards}
            size="sm"
            title="Toggle details"
          >
            {showHeaderCards ? (
              <>
                <EyeOff className="mr-2 w-4 h-4" /> Hide
              </>
            ) : (
              <>
                <Eye className="mr-2 w-4 h-4" /> Show
              </>
            )}{" "}
            details
          </Button>
          {!userClassroom.classroom.archived && isModerator(userClassroom) && (
            <>
              <Button variant="secondary" asChild size="sm" title="Download report">
                <a href={reportDownloadUrl} target="_blank" referrerPolicy="no-referrer">
                  <Download className="mr-2 h-4 w-4" />
                  Download report
                </a>
              </Button>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="secondary" size="sm" title="Archive classroom">
                    <Archive className="mr-2 h-4 w-4" /> Archive
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                    <AlertDialogDescription>
                      Are you sure that you wanna archive this classroom? This action can not be undone!
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction onClick={handleConfirmArchive} variant="destructive">
                      Confirm
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>

              <Button variant="secondary" asChild size="sm" title="Settings">
                <Link to="/classrooms/$classroomId/settings/" params={{ classroomId }}>
                  <Settings className="mr-2 h-4 w-4" />
                  Settings
                </Link>
              </Button>
            </>
          )}
        </div>
      </div>

      {showHeaderCards && (
        <ClassroomHeaderCards userClassroom={userClassroom} classroomMemberLength={classroomMembers.length} />
      )}

      <Tabs value={tab} className="w-full">
        <TabsList className="w-full">
          <TabsTrigger asChild value="assignments" className="w-full">
            <Link search={{ tab: "assignments" }}>Assignments</Link>
          </TabsTrigger>
          <TabsTrigger asChild value="members" className="w-full">
            <Link search={{ tab: "members" }}>Members</Link>
          </TabsTrigger>
          {userClassroom.classroom.maxTeamSize > 1 && (
            <TabsTrigger asChild value="teams" className="w-full">
              <Link search={{ tab: "teams" }}>Teams</Link>
            </TabsTrigger>
          )}
        </TabsList>
        <TabsContent value="assignments" className="pt-2">
          {isModerator(userClassroom) && (
            <AssignmentListSection classroomId={classroomId} deactivateInteraction={userClassroom.classroom.archived} />
          )}
          {isStudent(userClassroom) && <ProjectListSection classroomId={classroomId} />}
        </TabsContent>
        <TabsContent value="members" className="pt-2">
          <div className="grid grid-cols-1 justify-between gap-4">
            <MemberListCard
              teamsReportUrls={teamsReportUrls}
              classroomMembers={classroomMembers}
              classroomId={classroomId}
              userClassroom={userClassroom}
              showTeams={userClassroom.classroom.maxTeamSize > 1}
              deactivateInteraction={userClassroom.classroom.archived}
            />
            {/* uses Role.Owner, as you can only be the owner, making a check if GetMe.id == OwnedClassroom.ownerId unnecessary*/}
          </div>
        </TabsContent>
        {userClassroom.classroom.maxTeamSize > 1 && (
          <TabsContent value="teams" className="pt-2">
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
      <Outlet />
    </>
  );
}

const ClassroomHeaderCards = ({
  userClassroom,
  classroomMemberLength,
}: {
  userClassroom: UserClassroomResponse;
  classroomMemberLength: number;
}) => {
  return (
    <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-4 mb-8">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Creation date</CardTitle>
          <CalendarClock className="mr-2 h-4 w-4" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{formatDate(userClassroom.classroom.createdAt)}</div>
          <p className="text-xs text-muted-foreground">
            {formatDistanceToNow(new Date(userClassroom.classroom.createdAt)) + " ago"}
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Members</CardTitle>
          <Users className="mr-2 h-4 w-4" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{classroomMemberLength}</div>
          <p className="text-xs text-muted-foreground">{classroomMemberLength == 1 ? "member" : "members"}</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Assignments</CardTitle>
          <CalendarCheck2 className="mr-2 h-4 w-4" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{userClassroom.assignmentsCount}</div>
          <p className="text-xs text-muted-foreground">
            {userClassroom.assignmentsCount == 1 ? "assignment" : "assignments"}
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Status</CardTitle>
          <Info className="mr-2 h-4 w-4" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold">{userClassroom.classroom.archived === true ? "Archived" : "Active"}</div>
        </CardContent>
      </Card>

      <Card className="col-span-1 md:col-span-2 lg:col-span-4">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Description</CardTitle>
          <Text className="mr-2 h-4 w-4" />
        </CardHeader>
        <CardContent>
          <p>{userClassroom.classroom.description ?? <i>No description available</i>}</p>
        </CardContent>
      </Card>
    </div>
  );
};
