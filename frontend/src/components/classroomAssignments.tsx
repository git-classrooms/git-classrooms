import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { ArrowRight, Loader2 } from "lucide-react";
import { formatDate, formatDateWithTime } from "@/lib/utils.ts";
import { Link } from "@tanstack/react-router";
import { Assignment } from "@/swagger-client";
import { useMemo, useState } from "react";
import { useSuspenseQuery } from "@tanstack/react-query";
import { assignmentsQueryOptions } from "@/api/assignment";
import {
  ColumnFiltersState,
  SortingState,
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { Input } from "./ui/input";

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

const createAssignmentColumns = (classroomId: string, deactivateInteraction: boolean) => {
  const assignmentColumnHelper = createColumnHelper<Assignment>();

  return [
    assignmentColumnHelper.accessor((row) => row.name, {
      id: "name",
      header: "Name",
      cell: ({ row: { original: assignment } }) => (
        <div className="cursor-default flex justify-between">
          <Link
            to="/classrooms/$classroomId/assignments/$assignmentId"
            params={{ classroomId, assignmentId: assignment.id }}
          >
            <div className="font-medium">{assignment.name}</div>
            <div className="text-sm text-muted-foreground md:inline">{assignment.description}</div>
          </Link>
        </div>
      ),
    }),

    assignmentColumnHelper.accessor((row) => row.createdAt, {
      id: "createdAt",
      header: "Creation Date",
      cell: ({ getValue }) => formatDate(getValue()),
    }),

    assignmentColumnHelper.accessor((row) => row.dueDate, {
      id: "dueDate",
      header: "Due Date",
      cell: ({ getValue }) => (getValue() ? formatDateWithTime(getValue()!) : "-"),
    }),

    assignmentColumnHelper.display({
      id: "actions",
      header: () => <div className="text-right">Actions</div>,
      cell: ({ row: { original: assignment } }) => (
        <div className="flex flex-wrap flex-row-reverse gap-2">
          {!deactivateInteraction && (
            <Button variant="ghost" size="icon" asChild>
              <Link
                to="/classrooms/$classroomId/assignments/$assignmentId"
                params={{ classroomId, assignmentId: assignment.id }}
              >
                <ArrowRight className="text-gray-600 dark:text-white h-6 w-6" />
              </Link>
            </Button>
          )}
        </div>
      ),
    }),
  ];
};

function AssignmentTable({
  assignments,
  classroomId,
  deactivateInteraction,
}: {
  assignments: Assignment[];
  classroomId: string;
  deactivateInteraction: boolean;
}) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);

  const assignmentColumns = useMemo(
    () => createAssignmentColumns(classroomId, deactivateInteraction),
    [classroomId, deactivateInteraction],
  );

  const table = useReactTable({
    data: assignments,
    columns: assignmentColumns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    state: {
      sorting,
      columnFilters,
    },
  });
  return (
    <>
      <Input
        placeholder="Filter assignments..."
        value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
        onChange={(event) => table.getColumn("name")?.setFilterValue(event.target.value)}
        className="max-w-sm"
      />
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <TableHead key={header.id}>
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                  </TableHead>
                );
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id} data-state={row.getIsSelected() && "selected"}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={assignmentColumns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </>
  );
}
