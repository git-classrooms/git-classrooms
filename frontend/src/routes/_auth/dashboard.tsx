import { Button } from "@/components/ui/button";
import { useSuspenseQuery } from "@tanstack/react-query";
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
  const { data: auth } = useAuth();
  const { data: ownedClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Owned));
  const { data: moderatorClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Moderator));
  const { data: studentClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Student));
  const { data: activeAssignments } = useSuspenseQuery(activeAssignmentQueryOptions());

  const sortedAssignments = useMemo(() => {
    return [...activeAssignments].sort((a, b) => {
      if (a.dueDate === null) return 1;
      if (b.dueDate === null) return -1;
      return new Date(a.dueDate ?? 0).getTime() - new Date(b.dueDate ?? 0).getTime();
    });
  }, [activeAssignments]);

  const totalClassrooms = ownedClassrooms.length + moderatorClassrooms.length + studentClassrooms.length;
  const firstName = auth?.name?.split(" ")[0] ?? "there";

  return (
    <div className="space-y-8 pb-8">
      {/* Welcome Header */}
      <section className="animate-stagger-1">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              Welcome back, {firstName}
            </h1>
            <p className="text-muted-foreground mt-1">
              {activeAssignments.length === 0
                ? "You're all caught up!"
                : `You have ${activeAssignments.length} active assignment${activeAssignments.length !== 1 ? "s" : ""}`}
            </p>
          </div>
          <Button variant="glow" asChild>
            <Link to="/classrooms/create">
              <Plus className="w-4 h-4 mr-2" />
              Create Classroom
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
          <h2 className="text-xl font-semibold">Your Classrooms</h2>
          {totalClassrooms > 0 && (
            <Button variant="ghost" size="sm" asChild>
              <Link to="/classrooms">View all</Link>
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
                <h3 className="text-sm font-medium text-muted-foreground mb-4 uppercase tracking-wide">
                  Managed by you ({ownedClassrooms.length})
                </h3>
                <ClassroomCardGrid
                  classrooms={ownedClassrooms.slice(0, 4)}
                  role="owner"
                />
              </div>
            )}

            {/* Moderator Classrooms */}
            {moderatorClassrooms.length > 0 && (
              <div>
                <h3 className="text-sm font-medium text-muted-foreground mb-4 uppercase tracking-wide">
                  Moderating ({moderatorClassrooms.length})
                </h3>
                <ClassroomCardGrid
                  classrooms={moderatorClassrooms.slice(0, 4)}
                  role="moderator"
                />
              </div>
            )}

            {/* Joined Classrooms */}
            {studentClassrooms.length > 0 && (
              <div>
                <h3 className="text-sm font-medium text-muted-foreground mb-4 uppercase tracking-wide">
                  Joined ({studentClassrooms.length})
                </h3>
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
  return (
    <div className="border border-dashed border-border rounded-lg p-12 text-center">
      <h3 className="font-medium text-foreground mb-2">No classrooms yet</h3>
      <p className="text-sm text-muted-foreground mb-4">
        Create your first classroom to get started
      </p>
      <Button variant="outline" asChild>
        <Link to="/classrooms/create">
          <Plus className="w-4 h-4 mr-2" />
          Create Classroom
        </Link>
      </Button>
    </div>
  );
}
