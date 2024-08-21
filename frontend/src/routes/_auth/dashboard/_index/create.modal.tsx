import { ClassroomCreateForm } from "@/components/classroomsForm";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { createFileRoute, useNavigate } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/dashboard/_index/create/modal")({
  component: CreateModal,
});

function CreateModal() {
  const navigate = useNavigate();
  return (
    <Dialog
      defaultOpen
      onOpenChange={(open) => {
        if (!open) {
          navigate({
            to: "/dashboard",
            replace: true,
          });
        }
      }}
    >
      <DialogContent>
        <ClassroomCreateForm />
      </DialogContent>
    </Dialog>
  );
}
