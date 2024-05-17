import { createFileRoute } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import {
  ownedAssignmentProjectsQueryOptions,
  ownedAssignmentQueryOptions,
  useInviteAssignmentMembers,
} from "@/api/assignments.ts";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { Header } from "@/components/header.tsx";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Code, Loader2 } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { formatDate } from "@/lib/utils.ts";
import { GetOwnedClassroomAssignmentProjectResponse } from "@/swagger-client";

export const Route = createFileRoute("/_auth/classrooms/owned/$classroomId/assignments/$assignmentId/")({
  loader: async ({ context, params }) => {
    const assignment = await context.queryClient.ensureQueryData(
      ownedAssignmentQueryOptions(params.classroomId, params.assignmentId),
    );
    const assignmentProjects = await context.queryClient.ensureQueryData(
      ownedAssignmentProjectsQueryOptions(params.classroomId, params.assignmentId),
    );
    return { assignment, assignmentProjects };
  },
  component: AssignmentDetail,
  pendingComponent: Loader,
});

function AssignmentDetail() {
  const { classroomId, assignmentId } = Route.useParams();
  const { data: assignment } = useSuspenseQuery(ownedAssignmentQueryOptions(classroomId, assignmentId));
  const { data: assignmentProjects } = useSuspenseQuery(ownedAssignmentProjectsQueryOptions(classroomId, assignmentId));

  const { mutateAsync, isError, isPending } = useInviteAssignmentMembers(classroomId, assignmentId);

  return (
    <div className="max-w-5xl mx-auto">
      <Header title="Assignment Details" />

      <Card>
        <CardHeader>
          <CardTitle>{assignment.name}</CardTitle>
          <CardDescription>{assignment.description}</CardDescription>
          <CardFooter>Due date: {assignment.dueDate ? formatDate(assignment.dueDate) : "No Due Date"}</CardFooter>
        </CardHeader>
      </Card>

      <Header title="Member Assignments" className="mt-6">
        <Button onClick={() => mutateAsync()} disabled={isPending}>
          {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Send Invites"}
        </Button>
      </Header>
      {isError && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>The Invitation could not be send!</AlertDescription>
        </Alert>
      )}
      <AssignmentProjectTable assignmentProjects={assignmentProjects} />
    </div>
  );
}

function AssignmentProjectTable({
  assignmentProjects,
}: {
  assignmentProjects: GetOwnedClassroomAssignmentProjectResponse[];
}) {
  return (
    <Table>
      <TableCaption>AssignmentProjects</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>GitLab-URL</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {assignmentProjects.map((a) => (
          <TableRow key={`${a.assignmentId}-${a.team.id}`}>
            <TableHead>{a.team.name}</TableHead>
            <TableCell>{a.assignmentAccepted ? "Accepted" : "Pending"}</TableCell>
            <TableCell>
              {a.assignmentAccepted ? (
                <a href={a.projectPath} target="_blank" rel="noreferrer">
                  <Code />
                </a>
              ) : null}
            </TableCell>
            <TableCell className="text-right"></TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
