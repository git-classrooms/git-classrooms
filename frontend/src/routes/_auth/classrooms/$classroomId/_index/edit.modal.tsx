import { ClassroomEditForm } from '@/components/classroomsForm';
import { Dialog, DialogContent } from '@/components/ui/dialog';
import { createFileRoute, useNavigate } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/classrooms/$classroomId/_index/edit/modal')({
  component: EditModal,
})

function EditModal() {
  const { classroomId } = Route.useParams();
  const navigate = useNavigate();
  return (
    <Dialog
      defaultOpen
      onOpenChange={(open) => {
        if (!open) {
          navigate({
            to: "/classrooms/$classroomId",
            params: { classroomId },
            replace: true,
          });
        }
      }}
    >
      <DialogContent>
        <ClassroomEditForm classroomId={classroomId} />
      </DialogContent>
    </Dialog>
  );
}
