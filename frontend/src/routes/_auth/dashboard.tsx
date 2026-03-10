import { Button } from "@/components/ui/button";
import { useQueries, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { Loader } from "@/components/loader";
import { Plus } from "lucide-react";
import { classroomsQueryOptions } from "@/api/classroom";
import { Filter } from "@/types/classroom";
import { useMemo, Suspense } from "react";
import { activeAssignmentQueryOptions } from "@/api/assignment";
import { ActiveAssignmentListCard } from "@/components/activeAssignments";
import { ClassroomCardGrid } from "@/components/classroom-card";
import { useAuth } from "@/api/auth";
import { PendingAssignmentsBanner } from "@/components/pendingAssignmentsBanner";
import { projectsQueryOptions } from "@/api/project";
import { ProjectResponse } from "@/swagger-client";
import { useTranslation } from "react-i18next";

export const Route = createFileRoute("/_auth/dashboard")({
  component: Dashboard,
  loader: async ({ context: { queryClient } }) => {
    const ownedClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Owned));
    const moderatorClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Moderator));
    const studentClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Student));
    const activeAssignments = await queryClient.ensureQueryData(activeAssignmentQueryOptions());

    return {
      ownedClassrooms,
      moderatorClassrooms,
      studentClassrooms,
      activeAssignments,
    };
  },
  pendingComponent: Loader,
});

function Dashboard() {
  const { t } = useTranslation("classroom");
  const { t: tc } = useTranslation("common");
  const { data: auth } = useAuth();
  const { data: ownedClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Owned));
  const { data: moderatorClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Moderator));
  const { data: studentClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Student));
  const { data: activeAssignments } = useSuspenseQuery(activeAssignmentQueryOptions());

  // Memoize queries array to prevent unnecessary re-renders
  const projectQueriesConfig = useMemo(
    () => studentClassrooms.map((userClassroom) => projectsQueryOptions(userClassroom.classroom.id)),
    [studentClassrooms]
  );

  // Fetch projects for all student classrooms to filter assignments
  const projectQueries = useQueries({
    queries: projectQueriesConfig,
    combine: (results) => {
      const allProjects: ProjectResponse[] = [];
      results.forEach((result) => {
        if (result.data) {
          allProjects.push(...result.data);
        }
      });
      return {
        data: allProjects,
        isPending: results.some((r) => r.isPending),
      };
    },
  });

  // Get set of classroom IDs where user is owner/moderator (can see all assignments)
  const managedClassroomIds = useMemo(() => {
    const ids = new Set<string>();
    ownedClassrooms.forEach((c) => ids.add(c.classroom.id));
    moderatorClassrooms.forEach((c) => ids.add(c.classroom.id));
    return ids;
  }, [ownedClassrooms, moderatorClassrooms]);

  // Get set of assignment IDs the user has projects for
  const projectAssignmentIds = useMemo(() => {
    return new Set(projectQueries.data.map((p) => p.assignment.id));
  }, [projectQueries.data]);

  // Filter and sort assignments
  const sortedAssignments = useMemo(() => {
    // Filter: show if user manages the classroom OR has a project for the assignment
    const filtered = activeAssignments.filter((assignment) => {
      if (managedClassroomIds.has(assignment.classroomId)) {
        return true; // Show all assignments in managed classrooms
      }
      return projectAssignmentIds.has(assignment.id); // Only show if has project
    });

    return filtered.sort((a, b) => {
      if (a.dueDate === null) return 1;
      if (b.dueDate === null) return -1;
      return new Date(a.dueDate ?? 0).getTime() - new Date(b.dueDate ?? 0).getTime();
    });
  }, [activeAssignments, managedClassroomIds, projectAssignmentIds]);

  const totalClassrooms = ownedClassrooms.length + moderatorClassrooms.length + studentClassrooms.length;
  const firstName = auth?.name?.split(" ")[0] ?? "there";

  return (
    <div className="space-y-8 pb-8">
      {/* Welcome Header */}
      <section className="animate-stagger-1">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {tc("dashboard.welcomeBack", { name: firstName })}
            </h1>
            <p className="text-muted-foreground mt-1">
              {sortedAssignments.length === 0
                ? tc("dashboard.allCaughtUp")
                : tc("dashboard.activeAssignments", { count: sortedAssignments.length })}
            </p>
          </div>
          <Button variant="glow" asChild>
            <Link to="/classrooms/create">
              <Plus className="w-4 h-4 mr-2" />
              {t("create.button")}
            </Link>
          </Button>
        </div>
      </section>

      {/* Pending Assignments Banner - Main CTA */}
      {studentClassrooms.length > 0 && (
        <section className="animate-stagger-2">
          <Suspense fallback={null}>
            <PendingAssignmentsBanner studentClassrooms={studentClassrooms} />
          </Suspense>
        </section>
      )}

      {/* Active Assignments */}
      <section className="animate-stagger-3">
        <ActiveAssignmentListCard activeAssignments={sortedAssignments} />
      </section>

      {/* Classrooms Grid */}
      <section className="space-y-6 animate-stagger-4">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-semibold">{tc("dashboard.yourClassrooms")}</h2>
          {totalClassrooms > 0 && (
            <Button variant="ghost" size="sm" asChild>
              <Link to="/classrooms">{tc("navigation.viewAll")}</Link>
            </Button>
          )}
        </div>

        {totalClassrooms === 0 ? (
          <EmptyClassroomsState />
        ) : (
          <div className="space-y-8">
            {/* Managed Classrooms */}
            {ownedClassrooms.length > 0 && (
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
                    {tc("dashboard.managedByYou")} ({ownedClassrooms.length})
                  </h3>
                  {ownedClassrooms.length > 4 && (
                    <Link
                      to="/classrooms"
                      search={{ view: "managed" }}
                      className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                    >
                      {tc("navigation.viewAll")}
                    </Link>
                  )}
                </div>
                <ClassroomCardGrid
                  classrooms={ownedClassrooms.slice(0, 4)}
                  role="owner"
                />
              </div>
            )}

            {/* Moderator Classrooms */}
            {moderatorClassrooms.length > 0 && (
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
                    {tc("dashboard.moderating")} ({moderatorClassrooms.length})
                  </h3>
                  {moderatorClassrooms.length > 4 && (
                    <Link
                      to="/classrooms"
                      search={{ view: "joined" }}
                      className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                    >
                      {tc("navigation.viewAll")}
                    </Link>
                  )}
                </div>
                <ClassroomCardGrid
                  classrooms={moderatorClassrooms.slice(0, 4)}
                  role="moderator"
                />
              </div>
            )}

            {/* Joined Classrooms */}
            {studentClassrooms.length > 0 && (
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
                    {tc("dashboard.joined")} ({studentClassrooms.length})
                  </h3>
                  {studentClassrooms.length > 4 && (
                    <Link
                      to="/classrooms"
                      search={{ view: "joined" }}
                      className="text-xs text-muted-foreground hover:text-foreground transition-colors"
                    >
                      {tc("navigation.viewAll")}
                    </Link>
                  )}
                </div>
                <ClassroomCardGrid
                  classrooms={studentClassrooms.slice(0, 4)}
                  role="student"
                />
              </div>
            )}
          </div>
        )}
      </section>

      <Outlet />
    </div>
  );
}

function EmptyClassroomsState() {
  const { t } = useTranslation("classroom");

  return (
    <div className="border border-dashed border-border rounded-lg p-12 text-center">
      <h3 className="font-medium text-foreground mb-2">{t("list.empty.title")}</h3>
      <p className="text-sm text-muted-foreground mb-4">
        {t("list.empty.description")}
      </p>
      <Button variant="outline" asChild>
        <Link to="/classrooms/create">
          <Plus className="w-4 h-4 mr-2" />
          {t("create.button")}
        </Link>
      </Button>
    </div>
  );
}
