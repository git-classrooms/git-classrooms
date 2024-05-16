import { Dialog, DialogContent } from "@/components/ui/dialog";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import {
  ownedAssignmentsQueryOptions,
  ownedClassroomTeamProjectsQueryOptions,
} from "@/api/assignments.ts";
import { ownedClassroomTeamQueryOptions } from "@/api/teams.ts";
import { Loader } from "@/components/loader.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { ClassroomTeamModal } from "@/components/classroomTeam.tsx";

export const Route = createFileRoute("/_auth/classrooms/owned/$classroomId/_index/teams/$teamId/modal")({
  component: TeamModal,
  loader: async ({ context, params }) => {
    const team = await context.queryClient.ensureQueryData(ownedClassroomTeamQueryOptions(params.classroomId, params.teamId));
    const assignments = await context.queryClient.ensureQueryData(ownedAssignmentsQueryOptions(params.classroomId));
    const projects = await context.queryClient.ensureQueryData(ownedClassroomTeamProjectsQueryOptions(params.classroomId, params.teamId));
    return { team, assignments, projects };
  },
  pendingComponent: Loader,
});
/* Needs to  be in index of classroomId to be able to make it a popup for this page */

function TeamModal() {
  const { classroomId, teamId } = Route.useParams();
  const { data: team } = useSuspenseQuery(ownedClassroomTeamQueryOptions(classroomId, teamId));
  const { data: assignments } = useSuspenseQuery(ownedAssignmentsQueryOptions(classroomId)); //should normally also be in projects.assignment but isn't filled in yet so we need to do it like this
  const { data: projects } = useSuspenseQuery(ownedClassroomTeamProjectsQueryOptions(classroomId, teamId));
  const navigate = useNavigate();
  //fill in the assignment data in the project object, normally this should be done in the query but it isn't filled in yet
  for (const project of projects) {
    for (const assignment of assignments) {
      if (project.assignmentId === assignment.id) {
        if(assignment.dueDate == null) assignment.dueDate = ""
        project.assignment.id = assignment.id
        project.assignment.name = assignment.name
        project.assignment.dueDate = assignment.dueDate
      }
    }
  }
  return (
    <Dialog
      defaultOpen
      onOpenChange={(open) => {
        if (!open) {
          navigate({
            to: "/classrooms/owned/$classroomId",
            params: { classroomId },
            replace: true,
          });
        }
      }}
    >
      <DialogContent>
        <ClassroomTeamModal classroomId={classroomId} team={team} projects={projects}/>
      </DialogContent>
    </Dialog>
  );
}
