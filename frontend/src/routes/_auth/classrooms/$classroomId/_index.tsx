import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, redirect, Link } from "@tanstack/react-router";
import { MemberListCard } from "@/components/classroomMembers.tsx";
import { Role } from "@/types/classroom.ts";
import { TeamListCard } from "@/components/classroomTeams.tsx";
import { AssignmentListCard } from "@/components/classroomAssignments.tsx";
import { Header } from "@/components/header";
import { classroomQueryOptions } from "@/api/classroom";
import { assignmentsQueryOptions } from "@/api/assignment";
import { membersQueryOptions } from "@/api/member";
import { teamsQueryOptions } from "@/api/team";
import { ReportApiAxiosParamCreator, UserClassroomResponse } from "@/swagger-client";
import { Button } from "@/components/ui/button.tsx";
import { Archive, Settings } from "lucide-react";
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

export const Route = createFileRoute("/_auth/classrooms/$classroomId/_index")({
  component: ClassroomDetail,
  loader: async ({ context: { queryClient }, params }) => {
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    if (userClassroom.role === Role.Student && !userClassroom.team) {
      throw redirect({
        to: "/classrooms/$classroomId/teams/join",
        params,
      });
    }
    const { url: reportDownloadUrl } = await ReportApiAxiosParamCreator().getClassroomReport(params.classroomId);
    const members = await queryClient.ensureQueryData(membersQueryOptions(params.classroomId));
    if (userClassroom.role !== Role.Student) {
      const assignments = await queryClient.ensureQueryData(assignmentsQueryOptions(params.classroomId));
      return { userClassroom, assignments, members, teams, reportDownloadUrl };
    } else {
      return { userClassroom, members, teams, reportDownloadUrl };
    }
  },
  pendingComponent: Loader,
});

function ClassroomDetail() {
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  if (userClassroom.role === Role.Student) {
    return <ClassroomStudentView />;
  } else if (userClassroom.role === Role.Owner || userClassroom.role === Role.Moderator) {
    return <ClassroomSupervisorView userClassroom={userClassroom} />;
  }
}

function ClassroomStudentView() {
  // Role.Student does not have access to assignments
  // This is a placeholder
  return (
    <div>
      <h1>Joined Classroom Info</h1>
    </div>
  );
}

function ClassroomSupervisorView({ userClassroom }: { userClassroom: UserClassroomResponse }) {
  const { classroomId } = Route.useParams();
  const { reportDownloadUrl } = Route.useLoaderData();
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));
  const { data: assignments } = useSuspenseQuery(assignmentsQueryOptions(classroomId));

  const { mutate } = useArchiveClassroom(classroomId);

  const handleConfirmArchive = () => {
    mutate();
  };

  return (
    <div>
      <div className="grid grid-cols-[1fr,auto] justify-between gap-1">
        <Header
          title={`${userClassroom.classroom.archived ? "Archived " : ""}Classroom: ${userClassroom.classroom.name}`}
          subtitle={userClassroom.classroom.description}
        />
        <div className="grid grid-cols-2 gap-3">
          {!userClassroom.classroom.archived && (
            <>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button className="col-start-1" variant="secondary" size="sm" title="Archive classroom">
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

              <Button className="col-start-2" variant="secondary" asChild size="sm" title="Settings">
                <Link to="/classrooms/$classroomId/settings/" params={{ classroomId: classroomId }}>
                  <Settings className="mr-2 h-4 w-4" />
                  Settings
                </Link>
              </Button>
            </>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 justify-between gap-10">
        <Button asChild>
          <a href={reportDownloadUrl} target="_blank" referrerPolicy="no-referrer">
            Download Report
          </a>
        </Button>

        <AssignmentListCard
          assignments={assignments}
          classroomId={classroomId}
          classroomName={userClassroom.classroom.name}
          deactivateInteraction={userClassroom.classroom.archived}
        />
        <MemberListCard
          classroomMembers={classroomMembers}
          classroomId={classroomId}
          userRole={Role.Owner}
          showTeams={userClassroom.classroom.maxTeamSize > 1}
          deactivateInteraction={userClassroom.classroom.archived}
        />
        {/* uses Role.Owner, as you can only be the owner, making a check if GetMe.id == OwnedClassroom.ownerId unnecessary*/}
        {userClassroom.classroom.maxTeamSize > 1 && (
          <TeamListCard
            teams={teams}
            classroomId={classroomId}
            userRole={Role.Owner}
            maxTeamSize={userClassroom.classroom.maxTeamSize}
            numInvitedMembers={classroomMembers.length}
            deactivateInteraction={userClassroom.classroom.archived}
          />
        )}
        <Outlet />
      </div>
    </div>
  );
}
