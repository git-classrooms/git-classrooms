import { projectGradingResultsQueryOptions } from "@/api/grading";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/projects/$projectId/grading/")({
  loader: async ({ context: { queryClient }, params: { classroomId, projectId } }) => {
    const results = await queryClient.ensureQueryData(projectGradingResultsQueryOptions(classroomId, projectId));

    return { results };
  },
  component: GradingResult,
});

function GradingResult() {
  const { classroomId, projectId } = Route.useParams();
  const { data } = useSuspenseQuery(projectGradingResultsQueryOptions(classroomId, projectId));
  return (
    <div>
      Grading Results
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}
