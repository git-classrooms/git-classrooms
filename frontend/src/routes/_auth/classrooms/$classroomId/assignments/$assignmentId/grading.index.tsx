import { assignmentReportQueryOptions } from "@/api/report";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/$assignmentId/grading/")({
  loader: async ({ context: { queryClient }, params: { classroomId, assignmentId } }) => {
    const report = await queryClient.ensureQueryData(assignmentReportQueryOptions(classroomId, assignmentId));

    return { report };
  },
  component: GradingResult,
});

function GradingResult() {
  const { classroomId, assignmentId } = Route.useParams();
  const { data } = useSuspenseQuery(assignmentReportQueryOptions(classroomId, assignmentId));
  return (
    <div>
      Grading Results
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}
