import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Link } from "@tanstack/react-router";
import { Clipboard, Gitlab } from "lucide-react";
import List from "@/components/ui/list.tsx";
import ListItem from "@/components/ui/listItem.tsx";
import { Assignment } from "@/swagger-client";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Separator } from "@/components/ui/separator.tsx";

/**
 * AssignmentListCard is a React component that displays a list of active assignments in a classroom.
 *
 * @param {Object} props - The properties passed to the component.
 * @param {Array} props.assignments - An array of Assignment objects representing the assignments of the classroom.
 * @returns {JSX.Element} A React component that displays a card with the list of assignments in a classroom.
 */
export function AssignmentListCard({ assignments }: { assignments: Assignment[] }): JSX.Element {
  return (
    <Card className="p-2">
      <CardHeader>
        <CardTitle>Assignments</CardTitle>
        <CardDescription>All active assignments for this classroom</CardDescription>
      </CardHeader>
      <CardContent>
        <AssignmentTable assignments={assignments} />
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
            <Button variant="ghost" size="icon" asChild>
              <a href={assignment.webUrl} target="_blank" rel="noreferrer">
                <Gitlab className="h-6 w-6 text-gray-600" />
              </a>
            </Button>
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
        <div className="pr-2">
          <div className="font-medium">{assignment.name}</div>
          <div className="text-sm text-muted-foreground">{assignment.description}</div>
        </div>
      </HoverCardTrigger>
      <HoverCardContent className="w-100">
        <p className="text-lg font-semibold">{assignment.name}</p>
        <p className="text-sm text-muted-foreground mt-[-0.3rem]">{assignment.description}</p>
        <Separator className="my-1" />
        <p className="text-muted-foreground">Due Date: {assignment.dueDate ? assignment.dueDate : "No due date"}</p>
      </HoverCardContent>
    </HoverCard>
  );
}
