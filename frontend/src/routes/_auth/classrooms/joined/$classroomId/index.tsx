import { classroomQueryOptions } from "@/api/classroom";
import { Role } from "@/types/classroom";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/")({
  loader: async ({ context, params }) => {
    const classroom = await context.queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));

    if (classroom.role === Role.Student && !classroom.team) {
      throw redirect({
        to: "/classrooms/joined/$classroomId/teams/join",
        params,
      });
    }

    return { classroom };
  },
  component: JoinedClassroom,
});

function JoinedClassroom() {
  const { classroomId } = Route.useParams();
  const { data: joinedClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  joinedClassroom;

  return <div>Joined Classroom</div>;
}
