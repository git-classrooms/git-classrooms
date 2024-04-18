import { joinedClassroomQueryOptions } from "@/api/classrooms";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/teams")({
  beforeLoad: async ({ context, params }) => {
    const joinedClassroom = await context.queryClient.ensureQueryData(joinedClassroomQueryOptions(params.classroomId));

    if (joinedClassroom.classroom.maxTeamSize === 1) {
      throw redirect({
        to: "/classrooms/joined/$classroomId",
        params,
        replace: true,
      });
    }
  },
});
