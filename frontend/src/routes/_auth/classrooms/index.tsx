import { Button } from "@/components/ui/button";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { Loader } from "@/components/loader";
import { Archive, Plus } from "lucide-react";
import { classroomsQueryOptions } from "@/api/classroom";
import { Filter } from "@/types/classroom";
import { Role as SwaggerRole } from "@/swagger-client";
import { useMemo } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ClassroomCard, ClassroomCardGrid } from "@/components/classroom-card";
import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";

export const Route = createFileRoute("/_auth/classrooms/")({
  component: Classrooms,
  validateSearch: (search: Record<string, unknown>): { view?: "managed" | "joined" | "archived" } => {
    const view = search.view;
    if (view === "managed" || view === "joined" || view === "archived") {
      return { view };
    }
    return {};
  },
  loader: async ({ context: { queryClient } }) => {
    const ownedClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Owned));
    const moderatorClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Moderator));
    const studentClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Student));
    const archivedClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(undefined, true));

    return {
      ownedClassrooms,
      moderatorClassrooms,
      studentClassrooms,
      archivedClassrooms,
    };
  },
  pendingComponent: Loader,
});

const roleMap: Record<SwaggerRole, "owner" | "moderator" | "student"> = {
  [SwaggerRole.NUMBER_0]: "owner",
  [SwaggerRole.NUMBER_1]: "moderator",
  [SwaggerRole.NUMBER_2]: "student",
};

function swaggerRoleToCardRole(role: SwaggerRole): "owner" | "moderator" | "student" {
  return roleMap[role] ?? "student";
}

function Classrooms() {
  const { t } = useTranslation("classroom");
  const { view } = Route.useSearch();
  const navigate = useNavigate();
  const { data: ownedClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Owned));
  const { data: moderatorClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Moderator));
  const { data: studentClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Student));
  const { data: archivedClassrooms } = useSuspenseQuery(classroomsQueryOptions(undefined, true));

  const joinedClassrooms = useMemo(
    () => [...moderatorClassrooms, ...studentClassrooms],
    [moderatorClassrooms, studentClassrooms],
  );

  const totalManaged = ownedClassrooms.length;
  const totalJoined = joinedClassrooms.length;
  const totalArchived = archivedClassrooms.length;

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
      search: { view: value as "managed" | "joined" | "archived" },
      replace: true,
    });
  };

  return (
    <div className="space-y-6 pb-8">
      {/* Header */}
      <section>
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">{t("title")}</h1>
            <p className="text-muted-foreground mt-1">
              {t("list.subtitle")}
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

      {/* Tabs */}
      <section>
        <Tabs value={defaultTab} onValueChange={handleTabChange} className="w-full">
          <TabsList>
            <TabsTrigger value="managed">
              {t("list.managed")}
              {totalManaged > 0 && (
                <span className="ml-2 px-1.5 py-0.5 text-xs bg-muted rounded-md font-mono">
                  {totalManaged}
                </span>
              )}
            </TabsTrigger>
            <TabsTrigger value="joined">
              {t("list.joined")}
              {totalJoined > 0 && (
                <span className="ml-2 px-1.5 py-0.5 text-xs bg-muted rounded-md font-mono">
                  {totalJoined}
                </span>
              )}
            </TabsTrigger>
            <TabsTrigger value="archived">
              <Archive className="w-3.5 h-3.5 mr-1.5" />
              {t("list.archived")}
              {totalArchived > 0 && (
                <span className="ml-2 px-1.5 py-0.5 text-xs bg-muted rounded-md font-mono">
                  {totalArchived}
                </span>
              )}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="managed">
            {ownedClassrooms.length === 0 ? (
              <EmptyState showCreate />
            ) : (
              <ClassroomCardGrid classrooms={ownedClassrooms} role="owner" />
            )}
          </TabsContent>

          <TabsContent value="joined">
            {joinedClassrooms.length === 0 ? (
              <EmptyState isJoined />
            ) : (
              <div className="space-y-8">
                {moderatorClassrooms.length > 0 && (
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground mb-4 uppercase tracking-wide">
                      {t("list.moderating")} ({moderatorClassrooms.length})
                    </h3>
                    <ClassroomCardGrid classrooms={moderatorClassrooms} role="moderator" />
                  </div>
                )}
                {studentClassrooms.length > 0 && (
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground mb-4 uppercase tracking-wide">
                      {t("list.asStudent")} ({studentClassrooms.length})
                    </h3>
                    <ClassroomCardGrid classrooms={studentClassrooms} role="student" />
                  </div>
                )}
              </div>
            )}
          </TabsContent>

          <TabsContent value="archived">
            {archivedClassrooms.length === 0 ? (
              <EmptyState isArchived />
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {archivedClassrooms.map((c) => (
                  <ClassroomCard
                    key={c.classroom.id}
                    classroom={c}
                    role={swaggerRoleToCardRole(c.role)}
                  />
                ))}
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
  showCreate = false,
  isJoined = false,
  isArchived = false,
}: {
  showCreate?: boolean;
  isJoined?: boolean;
  isArchived?: boolean;
}) {
  const { t } = useTranslation("classroom");

  const title = isArchived
    ? t("list.empty.archivedTitle")
    : isJoined
      ? t("list.empty.joinedTitle")
      : t("list.empty.ownedTitle");
  const description = isArchived
    ? t("list.empty.archivedDescription")
    : isJoined
      ? t("list.empty.joinedDescription")
      : t("list.empty.ownedDescription");

  return (
    <div className="border border-dashed border-border rounded-lg p-12 text-center">
      <h3 className="font-medium text-foreground mb-2">{title}</h3>
      <p className="text-sm text-muted-foreground mb-4">{description}</p>
      {showCreate && (
        <Button variant="outline" asChild>
          <Link to="/classrooms/create">
            <Plus className="w-4 h-4 mr-2" />
            {t("create.button")}
          </Link>
        </Button>
      )}
    </div>
  );
}
