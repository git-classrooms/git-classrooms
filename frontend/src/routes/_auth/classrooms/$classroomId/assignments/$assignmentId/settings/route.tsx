import { createFileRoute, Link, Outlet, redirect } from "@tanstack/react-router";
import { cn, isOwner } from "@/lib/utils";
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
import { ClipboardList, GraduationCap, Settings2, Sliders } from "lucide-react";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/$assignmentId/settings")({
  beforeLoad: async ({ context: { queryClient }, params: { classroomId, assignmentId } }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(classroomId));
    if (!isOwner(userClassroom)) {
      throw redirect({
        to: "/classrooms/$classroomId/assignments/$assignmentId",
        params: { classroomId, assignmentId },
        replace: true,
      });
    }
  },
  component: Settings,
});

const navItems = [
  {
    to: "/classrooms/$classroomId/assignments/$assignmentId/settings/" as const,
    label: "General",
    icon: Sliders,
    description: "Basic configuration",
  },
  {
    to: "/classrooms/$classroomId/assignments/$assignmentId/settings/grading" as const,
    label: "Grading",
    icon: GraduationCap,
    description: "Assessment settings",
  },
];

function Settings() {
  const { classroomId, assignmentId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));

  return (
    <div>
      {/* Breadcrumb */}
      <Breadcrumb className="mb-6">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms">Classrooms</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link
                to="/classrooms/$classroomId"
                search={{ tab: "assignments" }}
                params={{ classroomId }}
              >
                {classroom.classroom.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link
                to="/classrooms/$classroomId/assignments/$assignmentId"
                params={{ classroomId, assignmentId }}
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                {assignment.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage className="text-foreground">Settings</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {/* Header */}
      <div className="flex items-center gap-4 mb-8">
        <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center">
          <Settings2 className="w-6 h-6 text-primary" />
        </div>
        <div>
          <h1 className="text-3xl font-bold font-mono tracking-tight">Assignment Settings</h1>
          <p className="text-muted-foreground flex items-center gap-2">
            <ClipboardList className="w-4 h-4" />
            {assignment.name}
          </p>
        </div>
      </div>

      {/* Main Layout */}
      <div className="flex flex-col lg:flex-row gap-8">
        {/* Sidebar Navigation */}
        <aside className="lg:w-56 shrink-0">
          <nav className="lg:sticky lg:top-20">
            {/* Glass card container */}
            <div className="p-1.5 rounded-xl bg-card/50 backdrop-blur-sm border border-border/50">
              <div className="space-y-1">
                {navItems.map((item, index) => (
                  <Link
                    key={item.to}
                    to={item.to}
                    params={{ classroomId, assignmentId }}
                    activeOptions={{ exact: true }}
                    className="group block"
                    style={{ animationDelay: `${index * 50}ms` }}
                  >
                    {({ isActive }) => (
                      <div
                        className={cn(
                          "relative flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all duration-200",
                          isActive
                            ? "bg-primary/10 text-foreground"
                            : "text-muted-foreground hover:text-foreground hover:bg-muted/50"
                        )}
                      >
                        {/* Active indicator */}
                        {isActive && (
                          <div className="absolute left-0 top-1/2 -translate-y-1/2 w-0.5 h-6 bg-primary rounded-full shadow-[0_0_8px_hsl(var(--primary))]" />
                        )}

                        <div
                          className={cn(
                            "w-8 h-8 rounded-lg flex items-center justify-center transition-all duration-200",
                            isActive
                              ? "bg-primary/20 text-primary"
                              : "bg-muted/50 text-muted-foreground group-hover:bg-muted group-hover:text-foreground"
                          )}
                        >
                          <item.icon className="w-4 h-4" />
                        </div>

                        <div className="flex-1 min-w-0">
                          <span className="block text-sm font-medium">{item.label}</span>
                          <span className="block text-xs text-muted-foreground truncate">
                            {item.description}
                          </span>
                        </div>
                      </div>
                    )}
                  </Link>
                ))}
              </div>
            </div>

            {/* Info card */}
            <div className="mt-4 p-3 rounded-lg border border-dashed border-border/50 bg-muted/20">
              <p className="text-xs text-muted-foreground">
                Some settings cannot be changed once students have accepted the assignment.
              </p>
            </div>
          </nav>
        </aside>

        {/* Content Area */}
        <div className="flex-1 min-w-0">
          <div className="p-6 rounded-xl bg-card/30 border border-border/50 backdrop-blur-sm">
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  );
}
