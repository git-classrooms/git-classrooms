import { ClassroomsForm } from "@/components/classroomsForm";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/dashboard/create")({
  component: () => (
    <div className="max-w-3xl mx-auto">
      <ClassroomsForm />
    </div>
  ),
});
