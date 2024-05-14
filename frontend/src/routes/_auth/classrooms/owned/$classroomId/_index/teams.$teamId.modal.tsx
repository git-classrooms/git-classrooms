import { CreateOwnedTeamForm } from "@/components/createOwnedTeamForm";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { createFileRoute, useNavigate } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/owned/$classroomId/_index/teams/$teamId/modal")({
  component: TeamModal,
});
/* Needs to  be in index of classroomId to be able to make it a popup for this page */

function TeamModal() {
  const { classroomId, teamId } = Route.useParams();
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
        {teamId}
      </DialogContent>
    </Dialog>
  );
}
