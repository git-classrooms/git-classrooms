import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";

import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table.tsx";
import { Edit, Gitlab } from "lucide-react";
import { formatDate } from "@/lib/utils.ts";
import { Link } from "@tanstack/react-router";
import { Assignment } from "@/swagger-client";

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
  return (
    <>
      <div className="flex mt-16 mb-6">
        <div className="grow">
          <h2 className="text-xl font-bold">Assignments</h2>
          <p className="text-sm text-muted-foreground">Assignments managed by this classroom</p>
        </div>
      </div>
      <AssignmentTable
        assignments={assignments}
        classroomId={classroomId}
        classroomName={classroomName}
        deactivateInteraction={deactivateInteraction}
      />

      {!deactivateInteraction && (
        <Button variant="default" asChild>
          <Link to="/classrooms/$classroomId/assignments/create" params={{ classroomId }}>
            Create assignment
          </Link>
        </Button>
      )}
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
      <TableBody>
        {assignments.map((a) => (
          <TableRow key={a.id}>
            <TableCell className="p-2">
              <div className="cursor-default flex justify-between">
                <div>
                  <div className="font-medium">{a.name}</div>
                  <div className="text-sm text-muted-foreground md:inline">{classroomName}</div>
                </div>
                <div className="flex items-end">
                  <div className="ml-auto">
                    <div className="font-medium text-right">Due date</div>
                    <div className="text-sm text-muted-foreground md:inline">
                      {a.dueDate ? formatDate(a.dueDate) : "No Due Date"}
                    </div>
                  </div>
                  {!deactivateInteraction && (
                    <Button variant="ghost" size="icon" asChild>
                      <Link
                        to="/classrooms/$classroomId/assignments/$assignmentId"
                        params={{ classroomId, assignmentId: a.id }}
                      >
                        <Edit className="h-6 w-6 text-gray-600" />
                      </Link>
                    </Button>
                  )}
                  <Button variant="ghost" size="icon" asChild>
                    <Link to="" params={{}}>
                      <Gitlab className="h-6 w-6 text-gray-600" />
                    </Link>
                  </Button>
                </div>
              </div>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
