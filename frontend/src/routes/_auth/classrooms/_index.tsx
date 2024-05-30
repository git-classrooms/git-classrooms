import { Button } from "@/components/ui/button";
import { Avatar } from "@/components/avatar";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Link, Outlet, createFileRoute } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { Code } from "lucide-react";
import { ArrowRight as ArrowRight } from "lucide-react";
import { Header } from "@/components/header";
import { classroomsQueryOptions } from "@/api/classroom";
import { Filter } from "@/types/classroom";
import { useMemo } from "react";
import { UserClassroomResponse } from "@/swagger-client";

export const Route = createFileRoute("/_auth/classrooms/_index")({
  component: Classrooms,
  loader: async ({ context: { queryClient } }) => {
    const ownedClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Owned));
    const moderatorClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Moderator));
    const studentClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Student));
    return {
      ownedClassrooms,
      moderatorClassrooms,
      studentClassrooms,
    };
  },
  pendingComponent: Loader,
});

function Classrooms() {
  const { data: ownedClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Owned));
  const { data: moderatorClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Moderator));
  const { data: studentClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Student));

  const joinedClassrooms = useMemo(
    () => [...moderatorClassrooms, ...studentClassrooms],
    [moderatorClassrooms, studentClassrooms],
  );

  return (
    <div>
      <Header title="Dashboard" />
      <div className="grid grid-cols-1 lg:grid-cols-2 justify-between gap-10">
        <OwnedClassroomTable classrooms={ownedClassrooms} />
        <JoinedClassroomTable classrooms={joinedClassrooms} />
        <ActiveAssignmentsTable classrooms={joinedClassrooms} />
        <Outlet />
      </div>
    </div>
  );
}

function OwnedClassroomTable({ classrooms }: { classrooms: UserClassroomResponse[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Owned Classrooms</CardTitle>
        <CardDescription>Classrooms which are managed by you</CardDescription>
      </CardHeader>
      <Table className="flex-auto flex-wrap justify-end">
        <TableBody>
          {classrooms.map((c) => (
            <TableRow key={c.classroom.id}>
              <TableCell className="flex flex-wrap content-center gap-4">
                {
                  <Avatar
                    avatarUrl={c.classroom.owner.gitlabAvatar.avatarURL}
                    fallbackUrl={c.classroom.owner.gitlabAvatar.fallbackAvatarURL}
                    name="classroom-avatar"
                  />
                }
                <div className="justify-center">
                  <div className="">{c.classroom.name}</div>
                  <div className="text-sm text-muted-foreground"> {c.classroom.description}</div>
                </div>
              </TableCell>
              <TableCell className="text-right">
                <Button asChild variant="outline">
                  <Link to="/classrooms/owned/$classroomId" params={{ classroomId: c.classroom.id }}>
                    <ArrowRight />
                  </Link>
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <CardFooter className="flex justify-end gap-2">
        <Button asChild variant="default">
          <Link to="/classrooms/create/modal" replace>
            Create a new Classroom
          </Link>
        </Button>
        <Button asChild variant="default">
          <Link to="/classrooms/create/modal" replace>
            View all your Classrooms
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}

function JoinedClassroomTable({ classrooms }: { classrooms: UserClassroomResponse[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Joined Classrooms</CardTitle>
        <CardDescription>Classrooms which you have joined</CardDescription>
      </CardHeader>
      <Table className="flex-auto">
        <TableBody>
          {classrooms.map((c) => (
            <TableRow key={c.classroom.id}>
              <TableCell>{c.classroom.name}</TableCell>
              <TableCell>{c.classroom.owner.name}</TableCell>
              <TableCell>
                <a href={c.webUrl} target="_blank" rel="noreferrer">
                  <Code />
                </a>
              </TableCell>
              <TableCell className="text-right">
                <Button variant="outline">
                  <Link to="/classrooms/joined/$classroomId" params={{ classroomId: c.classroom.id }}>
                    <ArrowRight />
                  </Link>
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Card>
  );
}

function ActiveAssignmentsTable({ classrooms }: { classrooms: UserClassroomResponse[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Active Assignments</CardTitle>
        <CardDescription>Your assignments thaht are not overdue</CardDescription>
      </CardHeader>
      <Table className="flex-auto">
        <TableBody>
          {classrooms.map((c) => (
            <TableRow key={c.classroom.id}>
              <TableCell>{c.classroom.name}</TableCell>
              <TableCell>{c.classroom.owner.name}</TableCell>
              <TableCell>
                <a href={c.webUrl} target="_blank" rel="noreferrer">
                  <Code />
                </a>
              </TableCell>
              <TableCell className="text-right">
                <Button variant="outline">
                  <Link to="/classrooms/joined/$classroomId" params={{ classroomId: c.classroom.id }}>
                    <ArrowRight />
                  </Link>
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Card>
  );
}
