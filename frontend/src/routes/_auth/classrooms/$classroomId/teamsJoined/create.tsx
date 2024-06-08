import { classroomQueryOptions } from "@/api/classroom";
import { teamsQueryOptions } from "@/api/team";
import { CreateJoinedTeamForm } from "@/components/createJoinedTeamForm";
import { Role } from "@/types/classroom";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teamsJoined/create")({
  loader: async ({ context: { queryClient }, params }) => {
    const classroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    if (classroom.role === Role.Student || classroom.classroom.maxTeamSize <= teams.length) {
      throw redirect({
        to: "/classrooms/$classroomId/teams",
        params,
        replace: true,
      });
    }
  },
  component: CreateJoinedTeam,
});

function CreateJoinedTeam() {
  const { classroomId } = Route.useParams();

  return <CreateJoinedTeamForm classroomId={classroomId} />;
}
