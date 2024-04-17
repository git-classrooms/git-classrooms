import { ownedClassroomMemberQueryOptions, ownedClassroomQueryOptions } from "@/api/classrooms";
import { Header } from "@/components/header";
import { Loader } from "@/components/loader";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Assignment } from "@/types/assignments";
import { User } from "@/types/user";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { ownedAssignmentsQueryOptions } from "@/api/assignments.ts";
import { formatDate } from "@/lib/utils.ts";
import { useMemo } from "react";
import { ownedClassroomTeamsQueryOptions } from "@/api/teams";
import { DefaultControllerGetOwnedClassroomTeamResponse } from "@/swagger-client";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export const Route = createFileRoute("/_auth/classrooms/owned/$classroomId")({
  component: ClassroomDetail,
  loader: async ({ context, params }) => {
    const classroom = await context.queryClient.ensureQueryData(ownedClassroomQueryOptions(params.classroomId));
    const assignments = await context.queryClient.ensureQueryData(ownedAssignmentsQueryOptions(params.classroomId));
    const members = await context.queryClient.ensureQueryData(ownedClassroomMemberQueryOptions(params.classroomId));
    const teams = await context.queryClient.ensureQueryData(ownedClassroomTeamsQueryOptions(params.classroomId));

    return { classroom, assignments, members, teams };
  },
  pendingComponent: Loader,
});

function ClassroomDetail() {
  const { classroomId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(ownedClassroomQueryOptions(classroomId));
  const { data: assignments } = useSuspenseQuery(ownedAssignmentsQueryOptions(classroomId));
  const { data: members } = useSuspenseQuery(ownedClassroomMemberQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(ownedClassroomTeamsQueryOptions(classroomId));

  const users = useMemo(() => members.map((m) => m.user), [members]);

  return (
    <div className="p-2 space-y-6">
      <Outlet />
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
          <Link to="/classrooms/owned/$classroomId/assignments/create" params={{ classroomId }}>
            Create assignment
          </Link>
        </Button>
      </Header>
      <AssignmentTable assignments={assignments} classroomId={classroomId} />

      <Header title="Members">
        <Button variant="default" asChild>
          <Link to="/classrooms/owned/$classroomId/invite" params={{ classroomId }}>
            Invite members
          </Link>
        </Button>
      </Header>
      <MemberTable members={users} />

      {classroom.maxTeamSize > 1 && (
        <>
          <Header title="Teams">
            <Button variant="default" asChild>
              <Link to="/classrooms/owned/$classroomId/teams/create/modal" params={{ classroomId }}>
                Create Teams
              </Link>
            </Button>
          </Header>
          <TeamTable teams={teams} />
        </>
      )}
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
                  to="/classrooms/owned/$classroomId/assignments/$assignmentId"
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

function TeamTable({ teams }: { teams: DefaultControllerGetOwnedClassroomTeamResponse[] }) {
  return (
    <Table>
      <TableCaption>Teams</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {teams.map((t) => (
          <TableRow key={t.id}>
            <TableCell>{t.name}</TableCell>
            <TableCell className="text-right">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button>Actions</Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuItem>Test</DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
