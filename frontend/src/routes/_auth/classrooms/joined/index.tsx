import { classroomsQueryOptions } from "@/api/classroom";
import { Loader } from "@/components/loader";
import { Filter } from "@/types/classroom";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useMemo } from "react";

export const Route = createFileRoute("/_auth/classrooms/joined/")({
  loader: async ({ context: { queryClient } }) => {
    const moderatorClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Moderator));
    const studentClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Student));

    return { moderatorClassrooms, studentClassrooms };
  },
  pendingComponent: Loader,
  component: JoinedClassrooms,
});

function JoinedClassrooms() {
  const { data: moderatorClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Moderator));
  const { data: studentClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Student));
  const joinedClassrooms = useMemo(
    () => [...moderatorClassrooms, ...studentClassrooms],
    [moderatorClassrooms, studentClassrooms],
  );

  joinedClassrooms;

  return <div>Joined Classrooms</div>;
}
