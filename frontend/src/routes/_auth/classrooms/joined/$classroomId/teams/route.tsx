import { classroomQueryOptions } from "@/api/classroom";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/teams")({
  beforeLoad: async ({ context: { queryClient }, params }) => {
    const joinedClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));

    if (joinedClassroom.classroom.maxTeamSize === 1) {
      throw redirect({
        to: "/classrooms/joined/$classroomId",
        params,
        replace: true,
      });
    }
  },
});
