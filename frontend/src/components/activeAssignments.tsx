import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { ArrowRight as ArrowRight } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import List from "@/components/ui/list.tsx";
import ListItem from "@/components/ui/listItem.tsx";
import { ActiveAssignmentResponse } from "@/swagger-client";
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import { formatDate } from "@/lib/utils";

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
  activeAssignments: ActiveAssignmentResponse[];
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

function AssignmentTable({ assignments }: { assignments: ActiveAssignmentResponse[] }) {
  return assignments.length === 0 ? (
    <p className="text-muted-foreground text-center">No active assignments.</p>
  ) : (
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
                  {assignment.dueDate ? formatDate(assignment.dueDate)   : "No due date"}
                </div>
              </div>
              <Button variant="ghost" size="icon" asChild>
                <Link to="/classrooms/$classroomId/assignments/$assignmentId"
                      params={{ classroomId: assignment.classroomId, assignmentId: assignment.id }}>
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

function AssignmentListElement({ assignment }: { assignment: ActiveAssignmentResponse }) {
  return (
    <HoverCard>
      <HoverCardTrigger className="cursor-default flex">
        <div className="pr-2">
          <Avatar>
            <AvatarFallback className="bg-[#FC6D25] text-black text-lg">
              {assignment.name.charAt(0)}
            </AvatarFallback>
          </Avatar>
        </div>
        <div>
          <div className="font-medium">{assignment.name}</div>
          <div className="text-sm text-muted-foreground md:inline">
            {assignment.classroom.name}
          </div>
        </div>
      </HoverCardTrigger>
      <HoverCardContent className="w-100">
        <p className="text-lg font-semibold">{assignment.name}</p>
        <div className="text-sm text-muted-foreground md:inline">
          {assignment.classroom.name}
        </div>
        <p className="text-sm text-muted-foreground my-1">
          Created at: {formatDate(assignment.createdAt)}
        </p>
        <Separator className="my-1" />
        <p className="text-muted-foreground">{assignment.description}</p>
      </HoverCardContent>
    </HoverCard>
  );
}
