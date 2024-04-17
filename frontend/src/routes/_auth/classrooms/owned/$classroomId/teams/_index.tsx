import { Button } from "@/components/ui/button";
import { Link, Outlet, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/owned/$classroomId/teams/_index")({
  component: TeamsIndex,
});

function TeamsIndex() {
  const { classroomId } = Route.useParams();
  return (
    <div>
      <Button variant="default" asChild>
        <Link to="/classrooms/owned/$classroomId/teams/create/modal" replace params={{ classroomId }}>
          Create
        </Link>
      </Button>
      TeamsIndex
      <Outlet />
    </div>
  );
}
