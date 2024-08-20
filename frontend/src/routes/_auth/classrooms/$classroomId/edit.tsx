import { classroomQueryOptions } from "@/api/classroom";
import { ClassroomsForm } from "@/components/classroomsForm";
import { Role } from "@/types/classroom";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/edit")({
  loader: async ({ context: { queryClient }, params }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    console.log(userClassroom);
    if (userClassroom.role !== Role.Owner) {
      throw redirect({
        to: "/classrooms",
        params,
      });
    }
  },
  pendingComponent: Loader,
  component: () => (
    <div className="max-w-3xl mx-auto">
      <ClassroomsForm />
    </div>
  ),
});
