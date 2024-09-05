import { Separator } from "@/components/ui/separator";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbSeparator,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { classroomQueryOptions } from "@/api/classroom";
import { useSuspenseQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/settings")({
  component: Settings,
});

function Settings() {
  const { classroomId } = Route.useParams();
  const { data } = useSuspenseQuery(classroomQueryOptions(classroomId));
  return (
    <div>
      <Breadcrumb className="mb-5">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms/$classroomId" params={{ classroomId }}>
                {data.classroom.name}
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
        <p className="text-muted-foreground">Manage the settings of your classroom.</p>
      </div>
      <Separator className="my-6" />
      <div className="flex flex-col space-y-8 lg:flex-row lg:space-x-12 lg:space-y-0">
        <aside className="-mx-4 lg:w-1/5">
          <nav className={cn("flex space-x-2 lg:flex-col lg:space-x-0 lg:space-y-1")}>
            <Link
              to={"/classrooms/$classroomId/settings/"}
              params={{ classroomId }}
              activeOptions={{ exact: true }}
              activeProps={{ className: "bg-muted hover:bg-muted" }}
              inactiveProps={{ className: "hover:bg-transparent hover:underline" }}
              className={cn(buttonVariants({ variant: "ghost" }), "justify-start")}
            >
              General
            </Link>
            <Link
              to={"/classrooms/$classroomId/settings/grading"}
              activeOptions={{ exact: true }}
              params={{ classroomId }}
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
