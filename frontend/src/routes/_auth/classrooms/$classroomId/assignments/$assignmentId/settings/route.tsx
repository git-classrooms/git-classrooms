import { Separator } from "@/components/ui/separator";
import { createFileRoute, Link, Outlet, redirect } from "@tanstack/react-router";
import { buttonVariants } from "@/components/ui/button";
import { cn, isModerator } from "@/lib/utils";
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbSeparator,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { useSuspenseQuery } from "@tanstack/react-query";
import { assignmentQueryOptions } from "@/api/assignment";
import { classroomQueryOptions } from "@/api/classroom";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/$assignmentId/settings")({
  beforeLoad: async ({ context: { queryClient }, params: { classroomId, assignmentId } }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(classroomId));
    if (!isModerator(userClassroom)) {
      throw redirect({
        to: "/classrooms/$classroomId/assignments/$assignmentId",
        params: { classroomId, assignmentId },
        replace: true,
      });
    }
  },
  component: Settings,
});

function Settings() {
  const { classroomId, assignmentId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  return (
    <div>
      <Breadcrumb className="mb-5">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId }}>
                {classroom.classroom.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms/$classroomId/assignments/$assignmentId" params={{ classroomId, assignmentId }}>
                {assignment.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Settings</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <div className="space-y-0.5">
        <h1 className="text-4xl font-bold tracking-tight">Settings</h1>
        <p className="text-muted-foreground">Manage the settings of this assignment.</p>
      </div>
      <Separator className="my-6" />
      <div className="flex flex-col space-y-8 lg:flex-row lg:space-x-12 lg:space-y-0">
        <aside className="-mx-4 lg:w-1/5">
          <nav className={cn("flex space-x-2 lg:flex-col lg:space-x-0 lg:space-y-1")}>
            <Link
              to={"/classrooms/$classroomId/assignments/$assignmentId/settings/"}
              params={{ classroomId, assignmentId }}
              activeOptions={{ exact: true }}
              activeProps={{ className: "bg-muted hover:bg-muted" }}
              inactiveProps={{ className: "hover:bg-transparent hover:underline" }}
              className={cn(buttonVariants({ variant: "ghost" }), "justify-start")}
            >
              General
            </Link>
            <Link
              to={"/classrooms/$classroomId/assignments/$assignmentId/settings/grading"}
              activeOptions={{ exact: true }}
              params={{ classroomId, assignmentId }}
              activeProps={{ className: "bg-muted hover:bg-muted" }}
              inactiveProps={{ className: "hover:bg-transparent hover:underline" }}
              className={cn(buttonVariants({ variant: "ghost" }), "justify-start")}
            >
              Grading
            </Link>
          </nav>
        </aside>
        <div className="flex-1 flex w-full">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
