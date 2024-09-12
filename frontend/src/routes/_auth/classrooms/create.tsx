import { ClassroomCreateForm } from "@/components/classroomsForm";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/create")({
  component: () => (
    <div className="max-w-3xl mx-auto">
      <ClassroomCreateForm />
    </div>
  ),
});
