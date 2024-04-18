import { joinedClassroomsQueryOptions } from "@/api/classrooms";
import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/")({
  loader: async ({ context }) => {
    const joinedClassrooms = await context.queryClient.ensureQueryData(joinedClassroomsQueryOptions);

    return { joinedClassrooms };
  },
  pendingComponent: Loader,
  component: JoinedClassrooms,
});

function JoinedClassrooms() {
  const { data: joinedClassrooms } = useSuspenseQuery(joinedClassroomsQueryOptions);
  joinedClassrooms;

  return <div>Joined Classrooms</div>;
}
