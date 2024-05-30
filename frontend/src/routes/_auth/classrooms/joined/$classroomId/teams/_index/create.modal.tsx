import { classroomQueryOptions } from "@/api/classroom";
import { teamsQueryOptions } from "@/api/team";
import { CreateJoinedTeamForm } from "@/components/createJoinedTeamForm";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Role } from "@/types/classroom";
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/teams/_index/create/modal")({
  loader: async ({ context: { queryClient }, params }) => {
    const classroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    if (classroom.role === Role.Student || classroom.classroom.maxTeamSize <= teams.length) {
      throw redirect({
        to: "/classrooms/joined/$classroomId/teams",
        params,
        replace: true,
      });
    }
  },
  component: CreateTeamModal,
});

function CreateTeamModal() {
  const { classroomId } = Route.useParams();
  const navigate = useNavigate();
  return (
    <Dialog
      defaultOpen
      onOpenChange={(open) => {
        if (!open) {
          navigate({
            to: "/classrooms/joined/$classroomId/teams",
            params: { classroomId },
            replace: true,
          });
        }
      }}
    >
      <DialogContent>
        <CreateJoinedTeamForm classroomId={classroomId} />
      </DialogContent>
    </Dialog>
  );
}
