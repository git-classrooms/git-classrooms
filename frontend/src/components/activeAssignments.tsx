import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Gitlab } from "lucide-react";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import List from "@/components/ui/list.tsx";
import ListItem from "@/components/ui/listItem.tsx";
import { Assignment } from "@/swagger-client";

/**
 * AssignmentListCard is a React component that displays a list of assignments in a classroom.
 * It includes a table of assignments.
 *
 * @param {Object} props - The properties passed to the component.
 * @param {Array} props.assignments - An array of Assignment objects representing the assignments in the classroom.
 * @param {boolean} props.showActiveOnly - A boolean to determine if only active assignments should be shown.
 * @returns {JSX.Element} A React component that displays a card with the list of assignments in a classroom.
 */
export function AssignmentListCard({
  assignments,
  showActiveOnly,
}: {
  assignments: Assignment[];
  showActiveOnly: boolean;
}): JSX.Element {
  const filteredAssignments = showActiveOnly
    ? assignments.filter((assignment) => !assignment.dueDate || new Date(assignment.dueDate) > new Date())
    : assignments;

  return (
    <Card className="p-2">
      <CardHeader>
        <CardTitle>{showActiveOnly ? "Active Assignments" : "Assignments"}</CardTitle>
        <CardDescription>
          {showActiveOnly ? "Your assignements that are not yet overdue." : "Your assignements"}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <AssignmentTable assignments={filteredAssignments} />
      </CardContent>
    </Card>
  );
}

function AssignmentTable({ assignments }: { assignments: Assignment[] }) {
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
                <a href="#" target="_blank" rel="noreferrer">
                  <Gitlab className="h-6 w-6 text-gray-600" />
                </a>
              </Button>
            </div>
          }
        />
      )}
    />
  );
}

function AssignmentListElement({ assignment }: { assignment: Assignment }) {
  return (
    <HoverCard>
      <HoverCardTrigger className="cursor-default flex">
        <div>
          <div className="font-medium">{assignment.name}</div>
          <div className="text-sm text-muted-foreground">{assignment.description}</div>
        </div>
      </HoverCardTrigger>
      <HoverCardContent className="w-100">
        <p className="text-lg font-semibold">{assignment.name}</p>
        <p className="text-sm text-muted-foreground mt-[-0.3rem]">
          Created at: {new Date(assignment.createdAt).toLocaleDateString()}
        </p>
        <Separator className="my-1" />
        <p className="text-muted-foreground">{assignment.description}</p>
        <Separator className="my-1" />
        <div className="text-muted-foreground">
          <span className="font-bold">Template Project ID:</span> {assignment.templateProjectId}
        </div>
      </HoverCardContent>
    </HoverCard>
  );
}
