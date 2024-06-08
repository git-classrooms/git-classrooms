import { teamQueryOptions } from "@/api/team";
import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Outlet, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teams/$teamId")({
  loader: async ({ context: { queryClient }, params }) => {
    const team = await queryClient.ensureQueryData(teamQueryOptions(params.classroomId, params.teamId));

    return { team };
  },
  pendingComponent: Loader,
  component: Team,
});

function Team() {
  const { classroomId, teamId } = Route.useParams();
  const { data: team } = useSuspenseQuery(teamQueryOptions(classroomId, teamId));
  team;
  return (
    <div>
      Team Page
      <Outlet />
    </div>
  );
}
