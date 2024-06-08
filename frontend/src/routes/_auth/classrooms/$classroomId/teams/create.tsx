import { CreateOwnedTeamForm } from "@/components/createOwnedTeamForm";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teams/create")({
  component: CreateOwnedTeam,
});

function CreateOwnedTeam() {
  const { classroomId } = Route.useParams();

  return (
    <div className="max-w-3xl mx-auto">
      <CreateOwnedTeamForm classroomId={classroomId} />{" "}
    </div>
  );
}
