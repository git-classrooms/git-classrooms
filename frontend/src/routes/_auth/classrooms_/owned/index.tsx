import { ownedClassroomsQueryOptions } from "@/api/classrooms";
import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/owned/")({
  loader: async ({ context }) => {
    const ownedClassrooms = await context.queryClient.ensureQueryData(ownedClassroomsQueryOptions);

    return { ownedClassrooms };
  },
  pendingComponent: Loader,
  component: OwnedClassrooms,
});

function OwnedClassrooms() {
  const { data: ownedClassrooms } = useSuspenseQuery(ownedClassroomsQueryOptions);
  ownedClassrooms;

  return <div>Owned Classrooms</div>;
}
