import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { MemberListCard } from "@/components/classroomMembers.tsx";
import { Role } from "@/types/classroom.ts";
import { TeamListCard } from "@/components/classroomTeams.tsx";
import { AssignmentListCard } from "@/components/classroomAssignments.tsx";
import { Header } from "@/components/header";
import { classroomQueryOptions } from "@/api/classroom";
import { assignmentsQueryOptions } from "@/api/assignment";
import { membersQueryOptions } from "@/api/member";
import { teamsQueryOptions } from "@/api/team";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/_index")({
  component: ClassroomDetail,
  loader: async ({ context: { queryClient }, params }) => {
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    const classroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    const assignments = await queryClient.ensureQueryData(assignmentsQueryOptions(params.classroomId));
    const members = await queryClient.ensureQueryData(membersQueryOptions(params.classroomId));
    return { classroom, assignments, members, teams };
  },
  pendingComponent: Loader,
});

function ClassroomDetail() {
  const { classroomId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: assignments } = useSuspenseQuery(assignmentsQueryOptions(classroomId));
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));
  return (
    <div>
      <Header title={`Classroom: ${classroom.classroom.name}`} subtitle={classroom.classroom.description} />
      <div className="grid grid-cols-1 lg:grid-cols-2 justify-between gap-10">
        <AssignmentListCard
          assignments={assignments}
          classroomId={classroomId}
          classroomName={classroom.classroom.name}
        />
        <MemberListCard
          classroomMembers={classroomMembers}
          classroomId={classroomId}
          userRole={Role.Owner}
          showTeams={classroom.classroom.maxTeamSize > 1}
        />
        {/* uses Role.Owner, as you can only be the owner, making a check if GetMe.id == OwnedClassroom.ownerId unnecessary*/}
        {classroom.classroom.maxTeamSize > 1 && (
          <TeamListCard teams={teams} classroomId={classroomId} userRole={Role.Owner} />
        )}
        <Outlet />
      </div>
    </div>
  );
}
