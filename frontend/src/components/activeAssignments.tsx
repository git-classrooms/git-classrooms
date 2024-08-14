import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { ArrowRight as ArrowRight, Gitlab } from "lucide-react";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import List from "@/components/ui/list.tsx";
import ListItem from "@/components/ui/listItem.tsx";
import { AssignmentResponse } from "@/swagger-client";
import { Link } from "@tanstack/react-router";

/**
 * ActiveAssignmentListCard is a React component that displays a list of active assignments in a classroom.
 *
 * @param {Object} props - The properties passed to the component.
 * @param {Array} props.activeAssignments - An array of active Assignment objects representing the active assignments in the classroom.
 * @returns {JSX.Element} A React component that displays a card with the list of active assignments in a classroom.
 */
export function ActiveAssignmentListCard({
  activeAssignments,
}: {
  activeAssignments: AssignmentResponse[];
}): JSX.Element {
  return (
    <Card className="p-2">
      <CardHeader>
        <CardTitle>Active Assignments</CardTitle>
        <CardDescription>Your assignments that are not yet overdue.</CardDescription>
      </CardHeader>
      <CardContent>
        <AssignmentTable assignments={activeAssignments} />
      </CardContent>
    </Card>
  );
}

function AssignmentTable({ assignments }: { assignments: AssignmentResponse[] }) {
  return (
    <List
      items={assignments}
      renderItem={(assignment) => (
        <ListItem
          leftContent={<AssignmentListElement assignment={assignment} />}
          rightContent={
            <div className="flex text-end gap-2">
              <div>
                <div className="font-medium">Due Date</div>
                <div className="text-sm text-muted-foreground">
                  {assignment.dueDate ? new Date(assignment.dueDate).toLocaleDateString() : "No due date"}
                </div>
              </div>
              <Button variant="ghost" size="icon" asChild>
                <a href="#" target="_blank" rel="noreferrer"> {/* TODO: Replace with GitLab URL */}
                  <Gitlab className="h-6 w-6 text-gray-600" />
                </a>
              </Button>
              <Button variant="ghost" size="icon" asChild>
        <Link to="/classrooms/$classroomId" params={{ classroomId: assignment.classroomId }}>
          <ArrowRight className="text-slate-500 dark:text-white" />
        </Link>
      </Button>
            </div>
          }
        />
      )}
    />
  );
}

function AssignmentListElement({ assignment }: { assignment: AssignmentResponse }) {
  return (
    <HoverCard>
      <HoverCardTrigger className="cursor-default flex">
        <div>
          <div className="font-medium">{assignment.name}</div>
        </div>
      </HoverCardTrigger>
      <HoverCardContent className="w-100">
        <p className="text-lg font-semibold">{assignment.name}</p>
        <p className="text-sm text-muted-foreground my-1">
          Created at: {new Date(assignment.createdAt).toLocaleDateString()}
        </p>
        <Separator className="my-1" />
        <p className="text-muted-foreground">{assignment.description}</p>
      </HoverCardContent>
    </HoverCard>
  );
}