import { teamReportQueryOptions } from "@/api/report";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teams/$teamId/grading/")({
  loader: async ({ context: { queryClient }, params: { classroomId, teamId } }) => {
    const report = await queryClient.ensureQueryData(teamReportQueryOptions(classroomId, teamId));

    return { report };
  },
  component: GradingResult,
});

function GradingResult() {
  const { classroomId, teamId } = Route.useParams();
  const { data } = useSuspenseQuery(teamReportQueryOptions(classroomId, teamId));
  return (
    <div>
      Grading Results
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}
