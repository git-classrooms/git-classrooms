import { classroomReportQueryOptions } from "@/api/report";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/grading/")({
  loader: async ({ context: { queryClient }, params: { classroomId } }) => {
    const report = await queryClient.ensureQueryData(classroomReportQueryOptions(classroomId));

    return { report };
  },
  component: GradingResult,
});

function GradingResult() {
  const { classroomId } = Route.useParams();
  const { data } = useSuspenseQuery(classroomReportQueryOptions(classroomId));
  return (
    <div>
      Grading Results
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}
