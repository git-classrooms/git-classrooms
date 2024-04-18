import { joinedClassroomQueryOptions } from "@/api/classrooms";
import { Role } from "@/types/classroom";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/")({
  loader: async ({ context, params }) => {
    const joinedClassroom = await context.queryClient.ensureQueryData(joinedClassroomQueryOptions(params.classroomId));

    if (joinedClassroom.role === Role.Student && !joinedClassroom.team) {
      throw redirect({
        to: "/classrooms/joined/$classroomId/teams/join",
        params,
      });
    }

    return { joinedClassroom };
  },
  component: JoinedClassroom,
});

function JoinedClassroom() {
  const { classroomId } = Route.useParams();
  const { data: joinedClassroom } = useSuspenseQuery(joinedClassroomQueryOptions(classroomId));
  joinedClassroom;

  return <div>Joined Classroom</div>;
}
