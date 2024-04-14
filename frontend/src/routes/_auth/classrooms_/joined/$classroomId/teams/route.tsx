import { joinedClassroomTeamsQueryOptions } from "@/api/teams";
import { Loader } from "@/components/loader";
import { Outlet, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/teams")({
  loader: async ({ context, params }) => {
    const teams = await context.queryClient.ensureQueryData(joinedClassroomTeamsQueryOptions(params.classroomId));
    return { teams };
  },
  pendingComponent: Loader,
  component: Outlet,
});
