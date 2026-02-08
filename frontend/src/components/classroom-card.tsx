import { Link } from "@tanstack/react-router";
import { ArrowRight, ClipboardList } from "lucide-react";
import { StatusBadge } from "@/components/ui/status-badge";
import { Button } from "@/components/ui/button";
import { UserClassroomResponse } from "@/swagger-client";
import { formatRelativeTime } from "@/lib/utils";
import { useTranslation } from "react-i18next";

interface ClassroomCardProps {
  classroom: UserClassroomResponse;
  role: "owner" | "moderator" | "student";
}

const roleVariant = {
  owner: "success" as const,
  moderator: "info" as const,
  student: "neutral" as const,
};

const avatarGradients = [
  "from-primary to-[hsl(280,100%,60%)]",
  "from-[hsl(142,71%,45%)] to-[hsl(185,100%,50%)]",
  "from-[hsl(38,92%,55%)] to-[hsl(0,72%,55%)]",
  "from-[hsl(280,65%,60%)] to-[hsl(210,100%,60%)]",
  "from-[hsl(0,72%,55%)] to-[hsl(38,92%,55%)]",
];

function getGradientForName(name: string): string {
  const hash = name.split("").reduce((acc, char) => acc + char.charCodeAt(0), 0);
  return avatarGradients[hash % avatarGradients.length];
}

export function ClassroomCard({ classroom, role }: ClassroomCardProps) {
  const { t } = useTranslation("classroom");
  const variant = roleVariant[role];
  const gradient = getGradientForName(classroom.classroom.name);

  const roleLabel = {
    owner: t("card.owner"),
    moderator: t("card.moderator"),
    student: t("card.member"),
  }[role];

  return (
    <Link
      to="/classrooms/$classroomId"
      params={{ classroomId: classroom.classroom.id }}
      search={{ tab: "assignments" }}
      className="group block"
    >
      <div className="bg-card border border-border rounded-lg p-5 transition-all duration-200 hover:-translate-y-0.5 hover:border-primary/30 hover:shadow-lg hover:shadow-background/50">
        <div className="flex items-start gap-4">
          <div
            className={`w-10 h-10 rounded-lg bg-gradient-to-br ${gradient} flex items-center justify-center shrink-0`}
          >
            <span className="font-mono font-semibold text-primary-foreground">
              {classroom.classroom.name.charAt(0).toUpperCase()}
            </span>
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <h3 className="font-semibold truncate group-hover:text-primary transition-colors">
                {classroom.classroom.name}
              </h3>
              <StatusBadge variant={variant} size="sm">
                {roleLabel}
              </StatusBadge>
            </div>
            {classroom.classroom.description && (
              <p className="text-sm text-muted-foreground line-clamp-2 mb-3">
                {classroom.classroom.description}
              </p>
            )}
          </div>
        </div>

        <div className="flex items-center justify-between mt-4 pt-4 border-t border-border">
          <div className="flex items-center gap-4 text-sm text-muted-foreground">
            <span className="flex items-center gap-1.5">
              <ClipboardList className="w-4 h-4" />
              {t("card.assignments", { count: classroom.assignmentsCount })}
            </span>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">
              {formatRelativeTime(classroom.createdAt)}
            </span>
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 opacity-0 group-hover:opacity-100 transition-opacity"
              asChild
            >
              <span>
                <ArrowRight className="w-4 h-4" />
              </span>
            </Button>
          </div>
        </div>
      </div>
    </Link>
  );
}

interface ClassroomCardGridProps {
  classrooms: UserClassroomResponse[];
  role: "owner" | "moderator" | "student";
  emptyMessage?: string;
}

export function ClassroomCardGrid({
  classrooms,
  role,
  emptyMessage,
}: ClassroomCardGridProps) {
  const { t } = useTranslation("common");

  if (classrooms.length === 0) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        {emptyMessage ?? t("empty.noResults")}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {classrooms.map((classroom) => (
        <ClassroomCard
          key={classroom.classroom.id}
          classroom={classroom}
          role={role}
        />
      ))}
    </div>
  );
}
