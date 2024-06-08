import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { MemberListCard } from "@/components/classroomMembers.tsx";
import { Role } from "@/types/classroom.ts";
import { TeamListCard } from "@/components/classroomTeams.tsx";
import { AssignmentListCard } from "@/components/classroomAssignments.tsx";
import { Header } from "@/components/header";
import { classroomQueryOptions } from "@/api/classroom";
import { assignmentsQueryOptions } from "@/api/assignment";
import { membersQueryOptions } from "@/api/member";
import { teamsQueryOptions } from "@/api/team";
import { UserClassroomResponse } from "@/swagger-client";

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
    const members = await queryClient.ensureQueryData(membersQueryOptions(params.classroomId));
    if(userClassroom.role !== Role.Student) {
      const assignments = await queryClient.ensureQueryData(assignmentsQueryOptions(params.classroomId));
      return { userClassroom, assignments, members, teams };
    }else{
      return { userClassroom, members, teams };
    }
  },
  pendingComponent: Loader,
});

function ClassroomDetail() {
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  if (userClassroom.role === Role.Student) {
    return <ClassroomStudentView />;
  } else {
    return <ClassroomOwnerView userClassroom={userClassroom} />;
  }
}

function ClassroomStudentView(){
  // Role.Student does not have access to assignments
  // This is a placeholder
  return (
    <div>
      <h1>Joined Classroom Info</h1>
    </div>
  );
}
function ClassroomOwnerView( {userClassroom}: {userClassroom: UserClassroomResponse}){
  const { classroomId } = Route.useParams();
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));
  const { data: assignments } = useSuspenseQuery(assignmentsQueryOptions(classroomId));
  return (
    <div>
      <Header title={`Classroom: ${userClassroom.classroom.name}`} subtitle={userClassroom.classroom.description} />
      <div className="grid grid-cols-1 lg:grid-cols-2 justify-between gap-10">
        <AssignmentListCard
          assignments={assignments}
          classroomId={classroomId}
          classroomName={userClassroom.classroom.name}
        />
        <MemberListCard
          classroomMembers={classroomMembers}
          classroomId={classroomId}
          userRole={Role.Owner}
          showTeams={userClassroom.classroom.maxTeamSize > 1}
        />
        {/* uses Role.Owner, as you can only be the owner, making a check if GetMe.id == OwnedClassroom.ownerId unnecessary*/}
        {userClassroom.classroom.maxTeamSize > 1 && (
          <TeamListCard teams={teams} classroomId={classroomId} userRole={Role.Owner} />
        )}
        <Outlet />
      </div>
    </div>
  );
}
