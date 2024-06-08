import { CreateTeamForm } from "@/components/createTeamForm.tsx";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/owned/$classroomId/teams/create")({
  component: CreateOwnedTeam,
});

function CreateOwnedTeam() {
  const { classroomId } = Route.useParams();

  return (
    <div className="max-w-3xl mx-auto">
      <CreateTeamForm classroomId={classroomId} />{" "}
    </div>
  );
}
