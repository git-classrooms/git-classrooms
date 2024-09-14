import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { Gitlab, LogIn } from "lucide-react";
import { cn, formatDate } from "@/lib/utils.ts";
import { Link } from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";
import { projectsQueryOptions } from "@/api/project";
import { ProjectResponse } from "@/swagger-client";
import { getStatusProps, Status } from "@/types/projects";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";

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
          <ProjectTable projects={projects} classroomId={classroomId} />
        </CardContent>
      </Card>
    </>
  );
}

function ProjectTable({ projects, classroomId }: { projects: ProjectResponse[]; classroomId: string }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="hidden md:table-cell">Creation date</TableHead>
          <TableHead className="hidden md:table-cell">Due date</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {projects.map((p) => {
          const statusProps = getStatusProps(p.projectStatus);

          return (
            <TableRow key={p.id}>
              <TableCell>
                <div className="cursor-default flex justify-between">
                  <a href={p.webUrl} target="_blank" referrerPolicy="no-referrer">
                    <div className="font-medium">{p.assignment.name}</div>
                    <div className="text-sm text-muted-foreground md:inline">{p.assignment.description}</div>
                  </a>
                </div>
              </TableCell>
              <TableCell>
                {p.assignment.dueDate && new Date(p.assignment.dueDate) < new Date() ? (
                  <div className="flex pl-1 gap-3 items-center">
                  <span className="relative flex h-3 w-3">
                    <span className="relative inline-flex rounded-full h-3 w-3 bg-gray-400"></span>
                  </span>
                    Closed
                  </div>) : (<div className="flex pl-1 gap-3 items-center">
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
                </div>)}
              </TableCell>
              <TableCell className="hidden md:table-cell min-w-[30%]">{formatDate(p.createdAt)}</TableCell>
              <TableCell className="hidden md:table-cell">
                {p.assignment.dueDate ? formatDate(p.assignment.dueDate) : "-"}
              </TableCell>
              <TableCell className="flex flex-wrap flex-row-reverse gap-2">
                {p.projectStatus === Status.Accepted ? (
                  <Button disabled={p.projectStatus !== Status.Accepted} variant="outline" size="icon" asChild>
                    <a href={p.webUrl} target="_blank" referrerPolicy="no-referrer">
                      <Gitlab className="h-4 w-4" />
                    </a>
                  </Button>
                ) : p.projectStatus === Status.Pending || p.projectStatus === Status.Failed ? (
                  <Tooltip delayDuration={0}>
                    <TooltipTrigger asChild>
                      <Button variant="outline" size="icon" asChild>
                        <Link
                          to="/classrooms/$classroomId/projects/$projectId/accept"
                          params={{ classroomId, projectId: p.id }}
                        >
                          <LogIn className="h-4 w-4" />
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
                      <Gitlab className="h-4 w-4 text-gray-400" />
                    </div>
                  </Button>
                )}
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
