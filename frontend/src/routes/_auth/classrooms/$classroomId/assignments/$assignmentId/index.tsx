import { createFileRoute, Link } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { Header } from "@/components/header.tsx";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Code, Edit, Loader2 } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { formatDate } from "@/lib/utils.ts";
import { assignmentQueryOptions } from "@/api/assignment";
import { assignmentProjectsQueryOptions, useInviteToAssignment } from "@/api/project";
import { ProjectResponse } from "@/swagger-client";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/$assignmentId/")({
  loader: async ({ context: { queryClient }, params }) => {
    const assignment = await queryClient.ensureQueryData(
      assignmentQueryOptions(params.classroomId, params.assignmentId),
    );
    const assignmentProjects = await queryClient.ensureQueryData(
      assignmentProjectsQueryOptions(params.classroomId, params.assignmentId),
    );
    return { assignment, assignmentProjects };
  },
  component: AssignmentDetail,
  pendingComponent: Loader,
});

function AssignmentDetail() {
  const { classroomId, assignmentId } = Route.useParams();
  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  const { data: assignmentProjects } = useSuspenseQuery(assignmentProjectsQueryOptions(classroomId, assignmentId));

  const { mutateAsync, isError, isPending } = useInviteToAssignment(classroomId, assignmentId);

  return (
    <div className="max-w-5xl mx-auto">
      <Header title="Assignment Details" />

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            {assignment.name}
            <Button className="mx-2" variant="ghost" size="icon" asChild>
              <Link
                to="/classrooms/$classroomId/assignments/$assignmentId/edit/modal"
                params={{ classroomId, assignmentId: assignmentId }}
              >
                <Edit className="h-6 w-6 text-gray-600" />
              </Link>
            </Button>
          </CardTitle>
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

function AssignmentProjectTable({ assignmentProjects }: { assignmentProjects: ProjectResponse[] }) {
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
          <TableRow key={`${a.assignment.id}-${a.team.id}`}>
            <TableHead>{a.team.name}</TableHead>
            <TableCell>{a.projectStatus}</TableCell>
            <TableCell>
              {a.projectStatus === "accepted" ? (
                <a href={a.webUrl} target="_blank" rel="noreferrer">
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
