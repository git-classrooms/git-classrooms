import { assignmentQueryOptions } from "@/api/assignment";
import { assignmentProjectsQueryOptions } from "@/api/project";
import { assignmentReportQueryOptions } from "@/api/report";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useMemo } from "react";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/$assignmentId/grading/")({
  loader: async ({ context: { queryClient }, params: { classroomId, assignmentId } }) => {
    const assignment = await queryClient.ensureQueryData(assignmentQueryOptions(classroomId, assignmentId));
    const report = await queryClient.ensureQueryData(assignmentReportQueryOptions(classroomId, assignmentId));
    const projects = await queryClient.ensureQueryData(assignmentProjectsQueryOptions(classroomId, assignmentId));

    return { assignment, report, projects };
  },
  component: GradingResult,
});

function GradingResult() {
  const { classroomId, assignmentId } = Route.useParams();
  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  const { data: gradingResults } = useSuspenseQuery(assignmentReportQueryOptions(classroomId, assignmentId));
  const { data: projects } = useSuspenseQuery(assignmentProjectsQueryOptions(classroomId, assignmentId));

  const zippedProjects = useMemo(
    () =>
      projects.map((project) => ({
        ...project,
        gradingResult: gradingResults.find((result) => result.projectId === project.id),
      })),
    [projects, gradingResults],
  );

  return (
    <div>
      Grading Results
      <pre>{JSON.stringify(assignment, null, 2)}</pre>
      <pre>{JSON.stringify(zippedProjects, null, 2)}</pre>
    </div>
  );
}
