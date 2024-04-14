import { joinedClassroomsQueryOptions, ownedClassroomsQueryOptions } from "@/api/classrooms";
import { Header } from "@/components/header";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Link, Outlet, createFileRoute } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { Code } from "lucide-react";
import { UserClassroom, OwnedClassroom } from "@/types/classroom.ts";

export const Route = createFileRoute("/_auth/classrooms")({
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
      <Header title="Own Classrooms">
        <Button asChild variant="default">
          <Link to="/classrooms/create/modal">Create</Link>
        </Button>
      </Header>
      <OwnedClassroomTable classrooms={ownClassrooms} />
      <Header title="Joined Classrooms" />
      <JoinedClassroomTable classrooms={joinedClassrooms} />
      <Outlet />
    </div>
  );
}

function OwnedClassroomTable({ classrooms }: { classrooms: OwnedClassroom[] }) {
  return (
    <Table>
      <TableCaption>Own Classrooms</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Gitlab-Link</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {classrooms.map((c) => (
          <TableRow key={c.id}>
            <TableCell>{c.name}</TableCell>
            <TableCell>
              <a href={c.gitlabUrl} target="_blank" rel="noreferrer">
                <Code />
              </a>
            </TableCell>
            <TableCell className="text-right">
              <Button asChild variant="outline">
                <Link to="/classrooms/owned/$classroomId" params={{ classroomId: c.id }}>
                  Show classroom
                </Link>
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function JoinedClassroomTable({ classrooms }: { classrooms: UserClassroom[] }) {
  return (
    <Table>
      <TableCaption>Joined Classrooms</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Owner</TableHead>
          <TableHead>Gitlab-Link</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
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
              <Button variant="outline">TBD</Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
