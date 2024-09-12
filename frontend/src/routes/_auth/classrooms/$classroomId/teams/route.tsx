import { classroomQueryOptions } from "@/api/classroom";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teams")({
  beforeLoad: async ({ context: { queryClient }, params }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));

    if (userClassroom.classroom.maxTeamSize === 1) {
      throw redirect({
        to: "/classrooms/$classroomId",
        search: { tab: "assignments" },
        params,
        replace: true,
      });
    }
  },
});
