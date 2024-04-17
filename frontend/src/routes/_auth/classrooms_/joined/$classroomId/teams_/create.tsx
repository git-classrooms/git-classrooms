import { joinedClassroomQueryOptions } from "@/api/classrooms";
import { CreateJoinedTeamForm } from "@/components/createJoinedTeamForm";
import { Role } from "@/types/classroom";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/teams/create")({
  loader: async ({ context, params }) => {
    const classroom = await context.queryClient.ensureQueryData(joinedClassroomQueryOptions(params.classroomId));
    if (!classroom.classroom.createTeams) {
      throw new Error("This classroom does not allow creating teams");
    }
    if (classroom.role !== Role.Moderator) {
      throw redirect({
        to: "/classrooms/joined/$classroomId/teams",
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
