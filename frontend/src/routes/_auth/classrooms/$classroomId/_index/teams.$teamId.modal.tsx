import { Dialog, DialogContent } from "@/components/ui/dialog";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { ClassroomTeamModal } from "@/components/classroomTeam.tsx";
import { teamQueryOptions } from "@/api/team";
import { teamProjectsQueryOptions } from "@/api/project";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/_index/teams/$teamId/modal")({
  component: TeamModal,
  loader: async ({ context: { queryClient }, params }) => {
    const team = await queryClient.ensureQueryData(teamQueryOptions(params.classroomId, params.teamId));
    const projects = await queryClient.ensureQueryData(teamProjectsQueryOptions(params.classroomId, params.teamId));
    return { team, projects };
  },
  pendingComponent: Loader,
});
/* Needs to  be in index of classroomId to be able to make it a popup for this page */

function TeamModal() {
  const { classroomId, teamId } = Route.useParams();
  const { data: team } = useSuspenseQuery(teamQueryOptions(classroomId, teamId));
  const { data: projects } = useSuspenseQuery(teamProjectsQueryOptions(classroomId, teamId));
  const navigate = useNavigate();

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
        <ClassroomTeamModal classroomId={classroomId} team={team} projects={projects} />
      </DialogContent>
    </Dialog>
  );
}
