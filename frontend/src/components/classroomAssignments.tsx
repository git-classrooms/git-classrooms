import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { ArrowRight, Edit, Gitlab, Loader2 } from "lucide-react";
import { formatDate, formatDateWithTime } from "@/lib/utils.ts";
import { Link, Navigate, useNavigate } from "@tanstack/react-router";
import { Assignment } from "@/swagger-client";
import { useState } from "react";

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
  assignments,
  classroomId,
  classroomName,
  deactivateInteraction,
}: {
  assignments: Assignment[];
  classroomId: string;
  classroomName: string;
  deactivateInteraction: boolean;
}): JSX.Element {
  const [isLoading, setIsLoading] = useState(false);
  return (
    <>
      <Card className="p-2">
        <CardHeader className="md:flex flex-row items-center justify-between space-y-0 pb-2 mb-4">
          <div>
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
            classroomName={classroomName}
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
  classroomName,
  deactivateInteraction,
}: {
  assignments: Assignment[];
  classroomId: string;
  classroomName: string;
  deactivateInteraction: boolean;
}) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead className="hidden md:table-cell">Creation date</TableHead>
          <TableHead className="hidden md:table-cell">Due date</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {assignments.map((a) => (
          <TableRow key={a.id}>
            <TableCell>
              <div className="cursor-default flex justify-between">
                <Link
                  to="/classrooms/$classroomId/assignments/$assignmentId"
                  params={{ classroomId, assignmentId: a.id }}
                >
                  <div className="font-medium">{a.name}</div>
                  <div className="text-sm text-muted-foreground md:inline">{a.description}</div>
                </Link>
              </div>
            </TableCell>
            <TableCell className="hidden md:table-cell min-w-[30%]">{formatDate(a.createdAt)}</TableCell>
            <TableCell className="hidden md:table-cell">{formatDate(a.dueDate ?? "-")}</TableCell>
            <TableCell className="flex flex-wrap flex-row-reverse gap-2">
              {!deactivateInteraction && (
                <Button variant="outline" size="icon" asChild>
                  <Link
                    to="/classrooms/$classroomId/assignments/$assignmentId"
                    params={{ classroomId, assignmentId: a.id }}
                  >
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
              )}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
