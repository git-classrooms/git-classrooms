import { classroomsQueryOptions } from "@/api/classroom";
import { Loader } from "@/components/loader";
import { Filter } from "@/types/classroom";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/owned/")({
  loader: async ({ context: { queryClient } }) => {
    const ownedClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Owned));

    return { ownedClassrooms };
  },
  pendingComponent: Loader,
  component: OwnedClassrooms,
});

function OwnedClassrooms() {
  const { data: ownedClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Owned));
  ownedClassrooms;

  return <div>Owned Classrooms</div>;
}
