import { joinedClassroomQueryOptions } from "@/api/classrooms";
import { Loader } from "@/components/loader";
import { Outlet, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId")({
  loader: async ({ context, params }) => {
    context.queryClient.ensureQueryData(joinedClassroomQueryOptions(params.classroomId));
  },
  pendingComponent: Loader,
  component: Outlet,
});
