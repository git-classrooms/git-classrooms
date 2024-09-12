import { classroomQueryOptions } from "@/api/classroom";
import { Role } from "@/types/classroom";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments")({
  beforeLoad: async ({ context: { queryClient }, params: { classroomId } }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(classroomId));

    if (userClassroom.role === Role.Student && !userClassroom.classroom.studentsViewAllProjects) {
      throw redirect({
        to: "/classrooms/$classroomId",
        search: { tab: "assignments" },
        params: { classroomId },
        replace: true,
      });
    }

    return { userClassroom };
  },
  component: Outlet,
});
