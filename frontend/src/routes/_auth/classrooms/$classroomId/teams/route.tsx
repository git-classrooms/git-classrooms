import { classroomQueryOptions } from "@/api/classroom";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teams")({
  beforeLoad: async ({ context: { queryClient }, params }) => {
    const joinedClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));

    if (joinedClassroom.classroom.maxTeamSize === 1) {
      throw redirect({
        to: "/classrooms/$classroomId",
        params,
        replace: true,
      });
    }
  },
});
