import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { ArrowRight, LogIn, SearchCode } from "lucide-react";
import { cn, formatDate, formatDateWithTime } from "@/lib/utils.ts";
import { Link } from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";
import { projectsQueryOptions } from "@/api/project";
import { ProjectResponse, UserClassroomResponse } from "@/swagger-client";
import { getStatusProps, Status } from "@/types/projects";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import { classroomQueryOptions } from "@/api/classroom.ts";
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
import { useMemo, useState } from "react";
import { Input } from "./ui/input";

/**
 * ProjectListSection is a React component that displays a list of projects in a classroom.
 * It includes a table of projects and a button to show more projects.
 *
 * @param {Object} props - The properties passed to the component.
 * @param {Array} props.projects - An array of Project objects representing the projects in the classroom.
 * @param {string} props.classroomId - The ID of the classroom.
 * @param {string} props.classroomName - The name of the classroom.
 * @param {boolean} props.deactivateInteraction - A boolean indicating whether the user can interact with the projects.
 * @returns {JSX.Element} A React component that displays a card with the list of projects in a classroom.
 * @constructor
 */
export function ProjectListSection({ classroomId }: { classroomId: string }): JSX.Element {
  const { data: projects } = useSuspenseQuery(projectsQueryOptions(classroomId));
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  return (
    <>
      <Card className="p-2">
        <CardHeader className="md:flex md:flex-row md:items-center justify-between space-y-0 pb-2 mb-4">
          <div className="mb-4 md:mb-0">
            <CardTitle className="mb-1">Assignments</CardTitle>
            <CardDescription>Your accepted or invited assignments for this classroom</CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <ProjectTable projects={projects} userClassroom={userClassroom} />
        </CardContent>
      </Card>
    </>
  );
}

const createProjectColumns = (user: UserClassroomResponse) => {
  const projectColumnHelper = createColumnHelper<ProjectResponse>();

  return [
    projectColumnHelper.accessor((row) => row.assignment.name, {
      id: "name",
      header: "Name",
      cell: ({ row: { original: project } }) => (
        <div className="cursor-default flex justify-between">
          <a href={project.webUrl} target="_blank" referrerPolicy="no-referrer">
            <div className="font-medium">{project.assignment.name}</div>
            <div className="text-sm text-muted-foreground md:inline">{project.assignment.description}</div>
          </a>
        </div>
      ),
    }),

    projectColumnHelper.display({
      id: "status",
      header: "Status",
      cell: ({ row: { original: project } }) => {
        const statusProps = getStatusProps(project.projectStatus);
        return project.assignment.dueDate && new Date(project.assignment.dueDate) < new Date() ? (
          <div className="flex pl-1 gap-3 items-center">
            <span className="relative flex h-3 w-3">
              <span className="relative inline-flex rounded-full h-3 w-3 bg-gray-400"></span>
            </span>
            Closed
          </div>
        ) : (
          <div className="flex pl-1 gap-3 items-center">
            <span className="relative flex h-3 w-3">
              <span
                className={cn(
                  "animate-ping absolute inline-flex h-full w-full rounded-full opacity-75",
                  statusProps.color.secondary,
                )}
              ></span>
              <span className={cn("relative inline-flex rounded-full h-3 w-3", statusProps.color.primary)}></span>
            </span>
            {statusProps.name}
          </div>
        );
      },
    }),

    projectColumnHelper.accessor((row) => row.createdAt, {
      id: "createdAt",
      header: "Creation Date",
      cell: ({ getValue }) => formatDate(getValue()),
    }),

    projectColumnHelper.accessor((row) => row.assignment.dueDate, {
      id: "dueDate",
      header: "Due Date",
      cell: ({ getValue }) => (getValue() ? formatDateWithTime(getValue()!) : "-"),
    }),

    projectColumnHelper.display({
      id: "actions",
      header: () => <div className="text-right">Actions</div>,
      cell: ({ row: { original: project } }) => (
        <div className="flex flex-wrap flex-row-reverse gap-2">
          {project.projectStatus === Status.Accepted ? (
            <>
              {user.classroom.studentsViewAllProjects && (
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button variant="ghost" size="icon" title="Go to assignment" asChild>
                      <Link
                        to="/classrooms/$classroomId/assignments/$assignmentId"
                        params={{ classroomId: user.classroom.id, assignmentId: project.assignment.id }}
                      >
                        <ArrowRight className="h-6 w-6 text-gray-600 dark:text-white" />
                      </Link>
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p>Go to assignment</p>
                  </TooltipContent>
                </Tooltip>
              )}
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="ghost" size="icon" title="Go to code" asChild>
                    <a href={project.webUrl} target="_blank" referrerPolicy="no-referrer">
                      <SearchCode className="h-6 w-6 text-gray-600 dark:text-white" />
                    </a>
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  <p>Go to code</p>
                </TooltipContent>
              </Tooltip>
            </>
          ) : project.projectStatus === Status.Pending || project.projectStatus === Status.Failed ? (
            <Tooltip delayDuration={0}>
              <TooltipTrigger asChild>
                <Button variant="ghost" size="icon" asChild>
                  <Link
                    to="/classrooms/$classroomId/projects/$projectId/accept"
                    params={{ classroomId: user.classroom.id, projectId: project.id }}
                  >
                    <LogIn className="text-gray-600 dark:text-white h-6 w-6" />
                  </Link>
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Accept assignment</p>
              </TooltipContent>
            </Tooltip>
          ) : (
            <Button variant="ghost" size="icon" asChild>
              <div>
                <SearchCode className="text-gray-600 dark:text-white h-6 w-6" />
              </div>
            </Button>
          )}
        </div>
      ),
    }),
  ];
};

function ProjectTable({
  projects,
  userClassroom,
}: {
  projects: ProjectResponse[];
  userClassroom: UserClassroomResponse;
}) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);

  const projectColumns = useMemo(() => createProjectColumns(userClassroom), [userClassroom]);

  const table = useReactTable({
    data: projects,
    columns: projectColumns,
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
        placeholder="Filter projects..."
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
              <TableCell colSpan={projectColumns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </>
  );
}
