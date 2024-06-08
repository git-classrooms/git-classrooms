import { CreateTeamForm } from "@/components/createTeamForm.tsx";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { classroomQueryOptions } from "@/api/classroom.ts";
import { teamsQueryOptions } from "@/api/team.ts";
import { Role } from "@/types/classroom.ts";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teams/_index/create/modal")({
  loader: async ({ context: { queryClient }, params }) => {
    const classroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    if (classroom.role === Role.Student || classroom.classroom.maxTeamSize <= teams.length) {
      throw redirect({
        to: "/classrooms/$classroomId/teams",
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
            to: "/classrooms/$classroomId/teams",
            params: { classroomId },
            replace: true,
          });
        }
      }}
    >
      <DialogContent>
        <CreateTeamForm classroomId={classroomId} />
      </DialogContent>
    </Dialog>
  );
}

