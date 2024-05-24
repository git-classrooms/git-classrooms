import { ownedClassroomMemberQueryOptions, ownedClassroomQueryOptions } from "@/api/classrooms";
import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { ownedAssignmentsQueryOptions } from "@/api/assignments.ts";
import { ownedClassroomTeamsQueryOptions } from "@/api/teams";
import { MemberListCard } from "@/components/classroomMembers.tsx";
import { Role } from "@/types/classroom.ts";
import { TeamListCard } from "@/components/classroomTeams.tsx";
import { AssignmentListCard } from "@/components/classroomAssignments.tsx";
import { Header } from "@/components/header";

export const Route = createFileRoute("/_auth/classrooms/owned/$classroomId/_index")({
  component: ClassroomDetail,
  loader: async ({ context, params }) => {
    const classroom = await context.queryClient.ensureQueryData(ownedClassroomQueryOptions(params.classroomId));
    const assignments = await context.queryClient.ensureQueryData(ownedAssignmentsQueryOptions(params.classroomId));
    const members = await context.queryClient.ensureQueryData(ownedClassroomMemberQueryOptions(params.classroomId));
    const teams = await context.queryClient.ensureQueryData(ownedClassroomTeamsQueryOptions(params.classroomId));
    return { classroom, assignments, members, teams };
  },
  pendingComponent: Loader,
});

function ClassroomDetail() {
  const { classroomId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(ownedClassroomQueryOptions(classroomId));
  const { data: assignments } = useSuspenseQuery(ownedAssignmentsQueryOptions(classroomId));
  const { data: classroomMembers } = useSuspenseQuery(ownedClassroomMemberQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(ownedClassroomTeamsQueryOptions(classroomId));
  return (
    <div>
      <Header title={`Classroom: ${classroom.name}`} subtitle={classroom.description} />
      <div className="grid grid-cols-1 lg:grid-cols-2 justify-between gap-10">
        <AssignmentListCard assignments={assignments} classroomId={classroomId} classroomName={classroom.name} />
        <MemberListCard
          classroomMembers={classroomMembers}
          classroomId={classroomId}
          userRole={Role.Owner}
          showTeams={classroom.maxTeamSize > 1}
        />
        {/* uses Role.Owner, as you can only be the owner, making a check if GetMe.id == OwnedClassroom.ownerId unnecessary*/}
        {classroom.maxTeamSize > 1 && <TeamListCard teams={teams} classroomId={classroomId} userRole={Role.Owner} />}
        <Outlet />
      </div>
    </div>
  );
}
