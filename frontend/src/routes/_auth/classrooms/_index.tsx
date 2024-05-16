import { joinedClassroomsQueryOptions, ownedClassroomsQueryOptions } from "@/api/classrooms";
import { Header } from "@/components/header";
import { Button } from "@/components/ui/button";
import { Avatar } from "@/components/avatar";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Link, Outlet, createFileRoute } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { Code } from "lucide-react";
import { GetJoinedClassroomResponse, GetOwnedClassroomResponse } from "@/swagger-client";
import { ArrowRight as ArrowRight } from "lucide-react";

export const Route = createFileRoute("/_auth/classrooms/_index")({
  component: Classrooms,
  loader: ({ context }) => {
    const ownClassrooms = context.queryClient.ensureQueryData(ownedClassroomsQueryOptions);
    const joinedClassrooms = context.queryClient.ensureQueryData(joinedClassroomsQueryOptions);
    return {
      ownClassrooms,
      joinedClassrooms,
    };
  },
  pendingComponent: Loader,
});

function Classrooms() {
  const { data: ownClassrooms } = useSuspenseQuery(ownedClassroomsQueryOptions);
  const { data: joinedClassrooms } = useSuspenseQuery(joinedClassroomsQueryOptions);
  return (
    <div className="p-10">
      <Header title="Dashboard" className="text-5xl" />
      <div className="grid grid-cols-1 lg:grid-cols-2 justify-between gap-20">
        <OwnedClassroomTable classrooms={ownClassrooms} />
        <JoinedClassroomTable classrooms={joinedClassrooms} />
        <ActiveAssignmentsTable classrooms={joinedClassrooms} />
        <Outlet />
      </div>
    </div>
  );
}

function OwnedClassroomTable({ classrooms }: { classrooms: GetOwnedClassroomResponse[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Owned Classrooms</CardTitle>
        <CardDescription>Classrooms which are managed by you</CardDescription>
      </CardHeader>
      <Table className="flex-auto flex-wrap justify-end">
        <TableBody>
          {classrooms.map((c) => (
            <TableRow key={c.id}>
              <TableCell className="flex flex-wrap content-center gap-4">
                {
                  <Avatar
                    avatarUrl={c.owner.gitlabAvatar.avatarURL}
                    fallbackUrl={c.owner.gitlabAvatar.fallbackAvatarURL}
                    name="classroom-avatar"
                  />
                }
                <div className="justify-center">
                  <div className="">{c.name}</div>
                  <div className="text-sm text-muted-foreground"> {c.description}</div>
                </div>
              </TableCell>
              <TableCell className="text-right">
                <Button asChild variant="outline">
                  <Link to="/classrooms/owned/$classroomId" params={{ classroomId: c.id }}>
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

function JoinedClassroomTable({ classrooms }: { classrooms: GetJoinedClassroomResponse[] }) {
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
                <a href={c.gitlabUrl} target="_blank" rel="noreferrer">
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

function ActiveAssignmentsTable({ classrooms }: { classrooms: GetJoinedClassroomResponse[] }) {
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
                <a href={c.gitlabUrl} target="_blank" rel="noreferrer">
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
