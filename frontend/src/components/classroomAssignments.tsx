import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { ArrowRight, Loader2 } from "lucide-react";
import { formatDate, formatDateWithTime } from "@/lib/utils.ts";
import { Link } from "@tanstack/react-router";
import { Assignment } from "@/swagger-client";
import { useMemo, useState } from "react";
import { useSuspenseQuery } from "@tanstack/react-query";
import { assignmentsQueryOptions } from "@/api/assignment";
import { DataTable, DataTableColumnHeader, ColumnDef } from "@/components/ui/data-table";

/**
 * AssignmentListSection is a React component that displays a list of assignments in a classroom.
 * It includes a table of assignments and a button to show more assignments.
 *
 * @param {Object} props - The properties passed to the component.
 * @param {Array} props.assignments - An array of Assignment objects representing the assignments in the classroom.
 * @param {string} props.classroomId - The ID of the classroom.
 * @param {string} props.classroomName - The name of the classroom.
 * @param {boolean} props.deactivateInteraction - A boolean indicating whether the user can interact with the assignments.
 * @returns {JSX.Element} A React component that displays a card with the list of assignments in a classroom.
 * @constructor
 */
export function AssignmentListSection({
  classroomId,
  deactivateInteraction,
}: {
  classroomId: string;
  deactivateInteraction: boolean;
}): JSX.Element {
  const { data: assignments } = useSuspenseQuery(assignmentsQueryOptions(classroomId));
  const [isLoading, setIsLoading] = useState(false);
  return (
    <>
      <Card className="p-2">
        <CardHeader className="md:flex md:flex-row md:items-center justify-between space-y-0 pb-2 mb-4">
          <div className="mb-4 md:mb-0">
            <CardTitle className="mb-1">Assignments</CardTitle>
            <CardDescription>Assignments managed by this classroom</CardDescription>
          </div>
          {!deactivateInteraction && (
            <Button variant="outline" asChild>
              <Link
                to="/classrooms/$classroomId/assignments/create"
                onClick={() => setIsLoading(true)}
                params={{ classroomId }}
              >
                {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : "Create assignment"}
              </Link>
            </Button>
          )}
        </CardHeader>
        <CardContent>
          <AssignmentTable
            assignments={assignments}
            classroomId={classroomId}
            deactivateInteraction={deactivateInteraction}
          />
        </CardContent>
      </Card>
    </>
  );
}

function AssignmentTable({
  assignments,
  classroomId,
  deactivateInteraction,
}: {
  assignments: Assignment[];
  classroomId: string;
  deactivateInteraction: boolean;
}) {
  const columns = useMemo<ColumnDef<Assignment, unknown>[]>(
    () => [
      {
        accessorKey: "name",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Name" />,
        cell: ({ row }) => (
          <div className="cursor-default">
            <Link
              to="/classrooms/$classroomId/assignments/$assignmentId"
              params={{ classroomId, assignmentId: row.original.id }}
            >
              <div className="font-medium">{row.original.name}</div>
              <div className="text-sm text-muted-foreground">{row.original.description}</div>
            </Link>
          </div>
        ),
      },
      {
        accessorKey: "createdAt",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Creation date" />,
        cell: ({ row }) => <span className="hidden md:inline">{formatDate(row.original.createdAt)}</span>,
        meta: { className: "hidden md:table-cell" },
      },
      {
        accessorKey: "dueDate",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Due date" />,
        cell: ({ row }) => (
          <span className="hidden md:inline">
            {row.original.dueDate ? formatDateWithTime(row.original.dueDate) : "-"}
          </span>
        ),
        meta: { className: "hidden md:table-cell" },
      },
      {
        id: "actions",
        header: () => <span className="sr-only">Actions</span>,
        cell: ({ row }) =>
          !deactivateInteraction && (
            <div className="text-right">
              <Button variant="ghost" size="icon" asChild>
                <Link
                  to="/classrooms/$classroomId/assignments/$assignmentId"
                  params={{ classroomId, assignmentId: row.original.id }}
                >
                  <ArrowRight className="text-gray-600 dark:text-white h-6 w-6" />
                </Link>
              </Button>
            </div>
          ),
      },
    ],
    [classroomId, deactivateInteraction],
  );

  return (
    <DataTable
      columns={columns}
      data={assignments}
      searchKey="name"
      searchPlaceholder="Search assignments..."
      showPagination={assignments.length > 10}
      emptyMessage="No assignments found."
    />
  );
}
