import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { Header } from "@/components/header";
import { classroomsQueryOptions } from "@/api/classroom";
import { Filter } from "@/types/classroom";
import { useMemo } from "react";
import { activeAssignmentQueryOptions } from "@/api/assignment";
import { ActiveAssignmentListCard } from "@/components/activeAssignments";
import { JoinedClassroomTable, OwnedClassroomTable } from "@/components/classrooms";

export const Route = createFileRoute("/_auth/dashboard")({
  component: Classrooms,
  loader: async ({ context: { queryClient } }) => {
    const ownedClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Owned));
    const moderatorClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Moderator));
    const studentClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Student));
    const activeAssignments = await queryClient.ensureQueryData(activeAssignmentQueryOptions());

    return {
      ownedClassrooms,
      moderatorClassrooms,
      studentClassrooms,
      activeAssignments,
    };
  },
  pendingComponent: Loader,
});

function Classrooms() {
  const { data: ownedClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Owned));
  const { data: moderatorClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Moderator));
  const { data: studentClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Student));
  const { data: activeAssignments } = useSuspenseQuery(activeAssignmentQueryOptions());

  const joinedClassrooms = useMemo(
    () => [...moderatorClassrooms, ...studentClassrooms],
    [moderatorClassrooms, studentClassrooms],
  );

  const sortedAssignments = useMemo(() => {
    return [...activeAssignments].sort((a, b) => {
      if (a.dueDate === null) return 1;
      if (b.dueDate === null) return -1;
      return new Date(a.dueDate ?? 0).getTime() - new Date(b.dueDate ?? 0).getTime();
    });
  }, [activeAssignments]);

  return (
    <div>
      <div className="flex-1 space-y-4">
        <Header title="Dashboard" />
        <ActiveAssignmentListCard activeAssignments={sortedAssignments} />
        <div className="grid grid-cols-1 lg:grid-cols-2 justify-between gap-4">
          <OwnedClassroomTable showAll={false} classrooms={ownedClassrooms} />
          <JoinedClassroomTable classrooms={joinedClassrooms} />
          <Outlet />
        </div>
      </div>
    </div>
  );
}
