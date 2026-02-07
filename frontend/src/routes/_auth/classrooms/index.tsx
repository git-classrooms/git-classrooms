import { Button } from "@/components/ui/button";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { Loader } from "@/components/loader";
import { Plus } from "lucide-react";
import { classroomsQueryOptions } from "@/api/classroom";
import { Filter } from "@/types/classroom";
import { useMemo } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ClassroomCardGrid } from "@/components/classroom-card";
import { useNavigate } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/")({
  component: Classrooms,
  validateSearch: (search: Record<string, unknown>): { view?: "managed" | "joined" } => {
    const view = search.view;
    if (view === "managed" || view === "joined") {
      return { view };
    }
    return {};
  },
  loader: async ({ context: { queryClient } }) => {
    const ownedClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Owned));
    const moderatorClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Moderator));
    const studentClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Student));

    return {
      ownedClassrooms,
      moderatorClassrooms,
      studentClassrooms,
    };
  },
  pendingComponent: Loader,
});

function Classrooms() {
  const { view } = Route.useSearch();
  const navigate = useNavigate();
  const { data: ownedClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Owned));
  const { data: moderatorClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Moderator));
  const { data: studentClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Student));

  const joinedClassrooms = useMemo(
    () => [...moderatorClassrooms, ...studentClassrooms],
    [moderatorClassrooms, studentClassrooms],
  );

  const totalManaged = ownedClassrooms.length;
  const totalJoined = joinedClassrooms.length;

  // Smart default: if view is specified use it, otherwise pick tab with content
  const defaultTab = useMemo(() => {
    if (view) return view;
    // If managed has content, show managed (default for teachers)
    if (totalManaged > 0) return "managed";
    // If only joined has content, show joined
    if (totalJoined > 0) return "joined";
    // Fallback to managed
    return "managed";
  }, [view, totalManaged, totalJoined]);

  const handleTabChange = (value: string) => {
    navigate({
      to: "/classrooms",
      search: { view: value as "managed" | "joined" },
      replace: true,
    });
  };

  return (
    <div className="space-y-6 pb-8">
      {/* Header */}
      <section className="animate-stagger-1">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Classrooms</h1>
            <p className="text-muted-foreground mt-1">
              Manage and access all your classrooms
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

      {/* Tabs */}
      <section className="animate-stagger-2">
        <Tabs value={defaultTab} onValueChange={handleTabChange} className="w-full">
          <TabsList>
            <TabsTrigger value="managed">
              Managed
              {totalManaged > 0 && (
                <span className="ml-2 px-1.5 py-0.5 text-xs bg-muted rounded-md font-mono">
                  {totalManaged}
                </span>
              )}
            </TabsTrigger>
            <TabsTrigger value="joined">
              Joined
              {totalJoined > 0 && (
                <span className="ml-2 px-1.5 py-0.5 text-xs bg-muted rounded-md font-mono">
                  {totalJoined}
                </span>
              )}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="managed">
            {ownedClassrooms.length === 0 ? (
              <EmptyState
                title="No managed classrooms"
                description="Create your first classroom to get started teaching"
                showCreate
              />
            ) : (
              <ClassroomCardGrid classrooms={ownedClassrooms} role="owner" />
            )}
          </TabsContent>

          <TabsContent value="joined">
            {joinedClassrooms.length === 0 ? (
              <EmptyState
                title="No joined classrooms"
                description="Join a classroom using an invite link from your instructor"
              />
            ) : (
              <div className="space-y-8">
                {moderatorClassrooms.length > 0 && (
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground mb-4 uppercase tracking-wide">
                      Moderating ({moderatorClassrooms.length})
                    </h3>
                    <ClassroomCardGrid classrooms={moderatorClassrooms} role="moderator" />
                  </div>
                )}
                {studentClassrooms.length > 0 && (
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground mb-4 uppercase tracking-wide">
                      As Student ({studentClassrooms.length})
                    </h3>
                    <ClassroomCardGrid classrooms={studentClassrooms} role="student" />
                  </div>
                )}
              </div>
            )}
          </TabsContent>
        </Tabs>
      </section>

      <Outlet />
    </div>
  );
}

function EmptyState({
  title,
  description,
  showCreate = false,
}: {
  title: string;
  description: string;
  showCreate?: boolean;
}) {
  return (
    <div className="border border-dashed border-border rounded-lg p-12 text-center">
      <h3 className="font-medium text-foreground mb-2">{title}</h3>
      <p className="text-sm text-muted-foreground mb-4">{description}</p>
      {showCreate && (
        <Button variant="outline" asChild>
          <Link to="/classrooms/create">
            <Plus className="w-4 h-4 mr-2" />
            Create Classroom
          </Link>
        </Button>
      )}
    </div>
  );
}
