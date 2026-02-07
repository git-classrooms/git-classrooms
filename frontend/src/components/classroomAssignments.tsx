import { Button } from "@/components/ui/button";
import { ArrowRight, Calendar, ClipboardList, Loader2, Plus } from "lucide-react";
import { formatDate, formatRelativeTime, getDaysUntilDue } from "@/lib/utils";
import { Link } from "@tanstack/react-router";
import { Assignment } from "@/swagger-client";
import { useState } from "react";
import { useSuspenseQuery } from "@tanstack/react-query";
import { assignmentsQueryOptions } from "@/api/assignment";
import { Card, CardContent } from "@/components/ui/card";
import { StatusBadge } from "@/components/ui/status-badge";
import { cn } from "@/lib/utils";

export function AssignmentListSection({
  classroomId,
  deactivateInteraction,
}: {
  classroomId: string;
  deactivateInteraction: boolean;
}) {
  const { data: assignments } = useSuspenseQuery(assignmentsQueryOptions(classroomId));
  const [isLoading, setIsLoading] = useState(false);

  const sortedAssignments = [...assignments].sort((a, b) => {
    if (!a.dueDate && !b.dueDate) return 0;
    if (!a.dueDate) return 1;
    if (!b.dueDate) return -1;
    return new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime();
  });

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
            <ClipboardList className="w-5 h-5 text-primary" />
          </div>
          <div>
            <h3 className="font-semibold">Assignments</h3>
            <p className="text-sm text-muted-foreground">
              {assignments.length} assignment{assignments.length !== 1 ? "s" : ""} in this classroom
            </p>
          </div>
        </div>

        {!deactivateInteraction && (
          <Button variant="glow" size="sm" asChild>
            <Link
              to="/classrooms/$classroomId/assignments/create"
              onClick={() => setIsLoading(true)}
              params={{ classroomId }}
            >
              {isLoading ? (
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              ) : (
                <Plus className="w-4 h-4 mr-2" />
              )}
              Create Assignment
            </Link>
          </Button>
        )}
      </div>

      {/* Assignment Grid */}
      {assignments.length === 0 ? (
        <EmptyAssignmentsState canCreate={!deactivateInteraction} classroomId={classroomId} />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {sortedAssignments.map((assignment) => (
            <AssignmentCard
              key={assignment.id}
              assignment={assignment}
              classroomId={classroomId}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function EmptyAssignmentsState({
  canCreate,
  classroomId,
}: {
  canCreate: boolean;
  classroomId: string;
}) {
  return (
    <div className="border border-dashed border-border rounded-lg p-8 text-center">
      <ClipboardList className="w-10 h-10 text-muted-foreground mx-auto mb-3" />
      <h3 className="font-medium text-foreground mb-1">No assignments yet</h3>
      <p className="text-sm text-muted-foreground mb-4">
        {canCreate ? "Create your first assignment to get started" : "No assignments have been created yet"}
      </p>
      {canCreate && (
        <Button variant="outline" asChild>
          <Link to="/classrooms/$classroomId/assignments/create" params={{ classroomId }}>
            <Plus className="w-4 h-4 mr-2" />
            Create Assignment
          </Link>
        </Button>
      )}
    </div>
  );
}

function AssignmentCard({
  assignment,
  classroomId,
}: {
  assignment: Assignment;
  classroomId: string;
}) {
  const daysUntil = getDaysUntilDue(assignment.dueDate);
  const isOverdue = daysUntil !== null && daysUntil < 0;
  const isUrgent = daysUntil !== null && daysUntil >= 0 && daysUntil <= 3;

  const getStatusVariant = () => {
    if (assignment.closed) return "neutral";
    if (isOverdue) return "destructive";
    if (isUrgent) return "warning";
    return "success";
  };

  const getStatusLabel = () => {
    if (assignment.closed) return "Closed";
    if (isOverdue) return "Overdue";
    if (isUrgent) return "Due Soon";
    return "Open";
  };

  return (
    <Link
      to="/classrooms/$classroomId/assignments/$assignmentId"
      params={{ classroomId, assignmentId: assignment.id }}
      className="group block"
    >
      <Card className="h-full transition-all duration-200 hover:border-primary/30 hover:shadow-lg hover:shadow-background/50">
        <CardContent className="p-5">
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-2">
                <h4 className="font-semibold truncate group-hover:text-primary transition-colors">
                  {assignment.name}
                </h4>
                <StatusBadge variant={getStatusVariant()} size="sm" showDot>
                  {getStatusLabel()}
                </StatusBadge>
              </div>

              {assignment.description && (
                <p className="text-sm text-muted-foreground line-clamp-2 mb-3">
                  {assignment.description}
                </p>
              )}

              <div className="flex items-center gap-4 text-xs text-muted-foreground">
                <span className="flex items-center gap-1">
                  <Calendar className="w-3 h-3" />
                  Created {formatRelativeTime(assignment.createdAt)}
                </span>
              </div>
            </div>

            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
              asChild
            >
              <span>
                <ArrowRight className="w-4 h-4" />
              </span>
            </Button>
          </div>

          {/* Due Date Section */}
          {assignment.dueDate && (
            <div
              className={cn(
                "mt-4 pt-3 border-t border-border flex items-center justify-between",
                isOverdue && "text-destructive",
                isUrgent && !isOverdue && "text-warning"
              )}
            >
              <span className="text-xs uppercase tracking-wide font-medium">Due Date</span>
              <span className="text-sm font-mono">
                {daysUntil === 0
                  ? "Today"
                  : daysUntil === 1
                    ? "Tomorrow"
                    : isOverdue
                      ? `${Math.abs(daysUntil)} days ago`
                      : formatDate(assignment.dueDate)}
              </span>
            </div>
          )}
        </CardContent>
      </Card>
    </Link>
  );
}
