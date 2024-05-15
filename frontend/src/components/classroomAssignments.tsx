import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";

import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table.tsx";
import { Edit, Gitlab } from "lucide-react";
import { formatDate } from "@/lib/utils.ts";
import { Link } from "@tanstack/react-router";
import { Assignment } from "@/swagger-client";

/**
 * AssignmentListCard is a React component that displays a list of assignments in a classroom.
 * It includes a table of assignments and a button to show more assignments.
 *
 * @param {Object} props - The properties passed to the component.
 * @param {Array} props.assignments - An array of Assignment objects representing the assignments in the classroom.
 * @param {string} props.classroomId - The ID of the classroom.
 * @param {string} props.classroomName - The name of the classroom.
 * @returns {JSX.Element} A React component that displays a card with the list of assignments in a classroom.
 * @constructor
 */
export function AssignmentListCard({
  assignments,
  classroomId,
  classroomName,
}: {
  assignments: Assignment[];
  classroomId: string;
  classroomName: string;
}): JSX.Element {
  return (
    <Card className="p-2">
      <CardHeader>
        <CardTitle>Created Assignments</CardTitle>
        <CardDescription>Assignments you have created in this classroom</CardDescription>
      </CardHeader>
      <CardContent>
        <AssignmentTable assignments={assignments} classroomId={classroomId} classroomName={classroomName} />
      </CardContent>
      <CardFooter className="flex justify-end">
        <Button variant="default" asChild>
          <Link to="" params={{}}>
            View all assignments
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}

function AssignmentTable({
  assignments,
  classroomId,
  classroomName,
}: {
  assignments: Assignment[];
  classroomId: string;
  classroomName: string;
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
                    <Button variant="ghost" size="icon" asChild>
                      <Link
                        to="/classrooms/owned/$classroomId/assignments/$assignmentId"
                        params={{ classroomId, assignmentId: a.id }}
                      >
                        <Edit className="h-6 w-6 text-gray-600" />
                      </Link>
                    </Button>
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
