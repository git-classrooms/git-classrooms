import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";

import { Header } from "@/components/header";

import { Tabs, TabsContent, TabsList, TabsTrigger, } from "@/components/ui/tabs"
import { AssignmentListCard } from "@/components/classroomAssignments.tsx";
import { assignmentsQueryOptions } from "@/api/assignment.ts";
import { classroomQueryOptions } from "@/api/classroom.ts";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/")({
  component: Assignments,
  loader: async ({ context: { queryClient }, params }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    const assignments = await queryClient.ensureQueryData(assignmentsQueryOptions(params.classroomId));
    return {
      userClassroom,
      assignments
    };
  },
  pendingComponent: Loader,
});

function Assignments() {
  const { classroomId } = Route.useParams();
  const { data: assignments } = useSuspenseQuery(assignmentsQueryOptions(classroomId));
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  return (
    <div>
      <Header title="Assignments" />
      <Tabs defaultValue="managed" className="w-[400]">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="Active">Active Assignments</TabsTrigger>
          <TabsTrigger value="Past">Past Assignments</TabsTrigger>
        </TabsList>
        <TabsContent value="Active">
          <AssignmentListCard
            assignments={assignments}
            classroomId={classroomId}
            classroomName={userClassroom.classroom.name}
          /> {/* TODO: Change assignment list to properly work with user roles and active and past/overdue assignments */}
        </TabsContent>
        <TabsContent value="Past">
          <AssignmentListCard
            assignments={assignments}
            classroomId={classroomId}
            classroomName={userClassroom.classroom.name}
          />{/* TODO: Change assignment list to properly work with user roles and active and past/overdue assignments */}
        </TabsContent>
      </Tabs>
      <Outlet />
    </div>
  );
}
