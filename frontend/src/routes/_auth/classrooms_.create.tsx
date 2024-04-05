import { ClassroomsForm } from "@/components/classroomsForm";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/create")({
  component: ClassroomsForm,
});
