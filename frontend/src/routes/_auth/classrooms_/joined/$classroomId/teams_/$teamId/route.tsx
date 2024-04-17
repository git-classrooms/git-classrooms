import { joinedClassroomTeamQueryOptions } from "@/api/teams";
import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Outlet, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/teams/$teamId")({
  loader: async ({ context, params }) => {
    const team = await context.queryClient.ensureQueryData(
      joinedClassroomTeamQueryOptions(params.classroomId, params.teamId),
    );

    return { team };
  },
  pendingComponent: Loader,
  component: Team,
});

function Team() {
  const { classroomId, teamId } = Route.useParams();
  const { data: team } = useSuspenseQuery(joinedClassroomTeamQueryOptions(classroomId, teamId));
  team;
  return (
    <div>
      Team Page
      <Outlet />
    </div>
  );
}
