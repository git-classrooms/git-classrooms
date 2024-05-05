import { joinedClassroomsQueryOptions, ownedClassroomsQueryOptions } from "@/api/classrooms";
import { Header } from "@/components/header";
import { Button } from "@/components/ui/button";
import { Avatar } from "@/components/avatar";
import { Table, TableBody, TableCaption, TableCell, TableHeader, TableRow } from "@/components/ui/table";
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
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
    <div className="p-2">
      <Header title="Dashboard" size="5xl" margin="my-10" />

      <div className="flex flex-row justify-between w-full space-x-2.5">
        <div className="w-1/2">
          <OwnedClassroomTable classrooms={ownClassrooms} />
        </div>
        <div className="w-1/2">
          <Header title="Joined Classrooms" />
          <JoinedClassroomTable classrooms={joinedClassrooms} />
        </div>
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
              <TableCell className="flex flex-wrap">
                <div className="flex">
                  {
                    <Avatar
                      avatarUrl={c.owner.gitlabAvatar.avatarURL}
                      fallbackUrl={c.owner.gitlabAvatar.fallbackAvatarURL}
                      name="classroom-avatar"
                    />
                  }
                  {c.name}
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
          <div className="flex flex-wrap justify-end">
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
          </div>
        </TableBody>
      </Table>
    </Card>
  );
}

function JoinedClassroomTable({ classrooms }: { classrooms: GetJoinedClassroomResponse[] }) {
  return (
    <Table className="flex-auto">
      <TableCaption>Joined Classrooms</TableCaption>
      <TableHeader></TableHeader>
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
                  <ArrowRight />{" "}
                </Link>
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
