import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ArrowRight, Calendar, CheckCircle2, Clock, FileCode2, GraduationCap } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card";
import { ActiveAssignmentResponse } from "@/swagger-client";
import { formatDate, formatRelativeTime, getDaysUntilDue } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { useTranslation } from "react-i18next";

type UrgencyLevel = "high" | "medium" | "none";

function getUrgencyLevel(dueDate: string | null | undefined): UrgencyLevel {
  const days = getDaysUntilDue(dueDate);
  if (days === null) return "none";
  if (days <= 3) return "high";
  return "medium";
}

const urgencyConfig = {
  high: {
    indicator: "bg-warning shadow-[0_0_8px_hsl(var(--glow-warning))]",
    row: "hover:border-l-warning",
  },
  medium: {
    indicator: "bg-success",
    row: "hover:border-l-success",
  },
  none: {
    indicator: "bg-muted-foreground",
    row: "hover:border-l-muted-foreground",
  },
};

export function ActiveAssignmentListCard({
  activeAssignments,
}: {
  activeAssignments: ActiveAssignmentResponse[];
}) {
  const { t: tc } = useTranslation("common");

  return (
    <Card>
      <CardHeader className="pb-3">
        <CardTitle>{tc("dashboard.activeAssignmentsTitle")}</CardTitle>
        <CardDescription>{tc("dashboard.activeAssignmentsSubtitle")}</CardDescription>
      </CardHeader>
      <CardContent>
        {activeAssignments.length === 0 ? (
          <EmptyState />
        ) : (
          <div className="space-y-1">
            {activeAssignments.map((assignment) => (
              <AssignmentRow key={assignment.id} assignment={assignment} />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function EmptyState() {
  const { t: tc } = useTranslation("common");

  return (
    <div className="flex flex-col items-center justify-center py-8 text-center">
      <div className="w-12 h-12 rounded-full bg-success/15 flex items-center justify-center mb-3">
        <CheckCircle2 className="w-6 h-6 text-success" />
      </div>
      <p className="font-medium text-foreground">{tc("dashboard.allCaughtUp")}</p>
      <p className="text-sm text-muted-foreground">{tc("dashboard.noActiveAssignments")}</p>
    </div>
  );
}

function AssignmentHoverContent({
  assignment,
  urgency,
}: {
  assignment: ActiveAssignmentResponse;
  urgency: UrgencyLevel;
}) {
  const { t } = useTranslation("assignment");
  const { t: tc } = useTranslation("common");
  const daysUntil = getDaysUntilDue(assignment.dueDate);

  return (
    <div className="relative">
      {/* Header with gradient */}
      <div
        className={cn(
          "px-4 py-3 border-b",
          urgency === "high"
            ? "bg-gradient-to-r from-warning/10 to-transparent border-warning/20"
            : "bg-gradient-to-r from-primary/10 to-transparent border-border/50"
        )}
      >
        <div className="flex items-start gap-3">
          <div
            className={cn(
              "w-9 h-9 rounded-lg flex items-center justify-center shrink-0",
              urgency === "high"
                ? "bg-warning/20 text-warning"
                : "bg-primary/20 text-primary"
            )}
          >
            <FileCode2 className="w-4 h-4" />
          </div>
          <div className="flex-1 min-w-0">
            <h4 className="font-semibold text-sm leading-tight truncate">
              {assignment.name}
            </h4>
            <div className="flex items-center gap-1.5 mt-1 text-xs text-muted-foreground">
              <GraduationCap className="w-3 h-3" />
              <span className="truncate">{assignment.classroom.name}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="px-4 py-3 space-y-3">
        {/* Description */}
        {assignment.description && (
          <p className="text-sm text-muted-foreground leading-relaxed line-clamp-2">
            {assignment.description}
          </p>
        )}

        {/* Meta info */}
        <div className="flex items-center justify-between text-xs">
          <div className="flex items-center gap-1.5 text-muted-foreground">
            <Calendar className="w-3.5 h-3.5" />
            <span>{tc("time.createdAt")} {formatDate(assignment.createdAt)}</span>
          </div>

          {assignment.dueDate && (
            <div
              className={cn(
                "flex items-center gap-1.5 font-medium",
                urgency === "high" ? "text-warning" : "text-muted-foreground"
              )}
            >
              <Clock className="w-3.5 h-3.5" />
              <span>
                {daysUntil !== null && daysUntil <= 0
                  ? t("dueDate.dueToday")
                  : daysUntil === 1
                    ? t("dueDate.dueTomorrow")
                    : `${t("dueDate.due")} ${formatRelativeTime(assignment.dueDate)}`}
              </span>
            </div>
          )}
        </div>
      </div>

      {/* Footer */}
      <div className="px-4 py-2.5 bg-muted/30 border-t border-border/50 cursor-pointer">
        <div className="flex items-center justify-between">
          <span className="text-xs text-muted-foreground">
            {t("dueDate.clickToView")}
          </span>
          <ArrowRight className="w-3.5 h-3.5 text-muted-foreground" />
        </div>
      </div>
    </div>
  );
}

function AssignmentRow({ assignment }: { assignment: ActiveAssignmentResponse }) {
  const { t } = useTranslation("assignment");
  const urgency = getUrgencyLevel(assignment.dueDate);
  const config = urgencyConfig[urgency];
  const daysUntil = getDaysUntilDue(assignment.dueDate);

  return (
    <Link
      to="/classrooms/$classroomId"
      search={{ tab: "assignments" }}
      params={{ classroomId: assignment.classroomId }}
      className="group block"
    >
      <div
        className={cn(
          "relative flex items-center gap-4 px-4 py-3 rounded-lg",
          "border-l-2 border-l-transparent",
          "hover:bg-muted/50 transition-all duration-200",
          config.row
        )}
      >
        <div className={cn("w-2 h-2 rounded-full shrink-0", config.indicator)} />

        <HoverCard>
          <HoverCardTrigger asChild>
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <span className="font-mono font-medium truncate">
                  {assignment.name}
                </span>
                <Badge variant="outline" className="shrink-0 text-xs">
                  {assignment.classroom.name}
                </Badge>
              </div>
            </div>
          </HoverCardTrigger>
          <HoverCardContent className="w-80 p-0 overflow-hidden cursor-pointer">
            <AssignmentHoverContent assignment={assignment} urgency={urgency} />
          </HoverCardContent>
        </HoverCard>

        <div className="flex items-center gap-4 shrink-0">
          <div className="text-right hidden sm:block">
            {assignment.dueDate ? (
              <>
                <div className="text-xs text-muted-foreground">{t("dueDate.due")}</div>
                <div
                  className={cn(
                    "text-sm font-medium",
                    urgency === "high" && "text-warning"
                  )}
                >
                  {daysUntil !== null && daysUntil <= 0
                    ? t("dueDate.today")
                    : daysUntil === 1
                      ? t("dueDate.tomorrow")
                      : formatRelativeTime(assignment.dueDate)}
                </div>
              </>
            ) : (
              <span className="text-sm text-muted-foreground">{t("dueDate.noDueDate")}</span>
            )}
          </div>

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
    </Link>
  );
}
