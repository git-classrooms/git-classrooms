import { CreateTeamForm } from "@/components/createTeamForm.tsx";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { classroomQueryOptions } from "@/api/classroom.ts";
import { teamsQueryOptions } from "@/api/team.ts";
import { Role } from "@/types/classroom.ts";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teams/create")({
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
  component: CreateTeam,
});
function CreateTeam() {
  const { classroomId } = Route.useParams();

  return (
    <div className="max-w-3xl mx-auto">
      <CreateTeamForm classroomId={classroomId} />{" "}
    </div>
  );
}
