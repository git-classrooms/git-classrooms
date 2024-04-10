import { ownedClassroomMemberQueryOptions, ownedClassroomQueryOptions } from "@/api/classrooms";
import { Header } from "@/components/header";
import { Loader } from "@/components/loader";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Assignment } from "@/types/assignments";
import { User } from "@/types/user";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { ownedAssignmentsQueryOptions } from "@/api/assignments.ts";
import { formatDate } from "@/lib/utils.ts";
import { useMemo } from "react";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/")({
  component: ClassroomDetail,
  loader: ({ context, params }) => {
    const classroom = context.queryClient.ensureQueryData(ownedClassroomQueryOptions(params.classroomId));
    const assignments = context.queryClient.ensureQueryData(ownedAssignmentsQueryOptions(params.classroomId));
    const members = context.queryClient.ensureQueryData(ownedClassroomMemberQueryOptions(params.classroomId));

    return { classroom, assignments, members };
  },
  pendingComponent: Loader,
});

function ClassroomDetail() {
  const { classroomId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(ownedClassroomQueryOptions(classroomId));
  const { data: assignments } = useSuspenseQuery(ownedAssignmentsQueryOptions(classroomId));
  const { data: members } = useSuspenseQuery(ownedClassroomMemberQueryOptions(classroomId));

  const users = useMemo(() => members.map((m) => m.user), [members]);

  return (
    <div className="p-2 space-y-6">
      <div className="flex flex-row justify-between">
        <h1 className="text-xl font-bold">Classroom Details </h1>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{classroom.name}</CardTitle>
          <CardDescription>{classroom.description}</CardDescription>
        </CardHeader>
      </Card>

      <Header title="Assignments">
        <Button variant="default" asChild>
          <Link to="/classrooms/$classroomId/assignments/create" params={{ classroomId }}>
            Create assignment
          </Link>
        </Button>
      </Header>
      <AssignmentTable assignments={assignments} classroomId={classroomId} />

      <Header title="Members">
        <Button variant="default" asChild>
          <Link to="/classrooms/$classroomId/invite" params={{ classroomId }}>
            Invite members
          </Link>
        </Button>
      </Header>
      <MemberTable members={users} />
    </div>
  );
}

function AssignmentTable({ assignments, classroomId }: { assignments: Assignment[]; classroomId: string }) {
  return (
    <Table>
      <TableCaption>Assignments</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Due date</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {assignments.map((a) => (
          <TableRow key={a.id}>
            <TableCell>{a.name}</TableCell>
            <TableCell>{a.dueDate ? formatDate(a.dueDate) : "No Due Date"}</TableCell>
            <TableCell className="text-right">
              <Button asChild>
                <Link
                  to="/classrooms/$classroomId/assignments/$assignmentId"
                  params={{ classroomId, assignmentId: a.id }}
                >
                  Show Assignment
                </Link>
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function MemberTable({ members }: { members: User[] }) {
  return (
    <Table>
      <TableCaption>Members</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>E-Mail</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {members.map((m) => (
          <TableRow key={m.id}>
            <TableCell>{m.name}</TableCell>
            <TableCell>{m.gitlabEmail}</TableCell>
            <TableCell className="text-right"></TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
