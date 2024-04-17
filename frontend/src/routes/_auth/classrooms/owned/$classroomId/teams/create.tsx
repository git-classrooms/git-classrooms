import { CreateOwnedTeamForm } from "@/components/createOwnedTeamForm";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/owned/$classroomId/teams/create")({
  component: CreateOwnedTeam,
});

function CreateOwnedTeam() {
  const { classroomId } = Route.useParams();

  return <CreateOwnedTeamForm classroomId={classroomId} />;
}
