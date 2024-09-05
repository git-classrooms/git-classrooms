import { ClassroomEditForm } from "@/components/classroomsForm";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/settings/")({
  component: Index,
});

function Index() {
  const { classroomId } = Route.useParams();
  return <ClassroomEditForm classroomId={classroomId} />;
}
