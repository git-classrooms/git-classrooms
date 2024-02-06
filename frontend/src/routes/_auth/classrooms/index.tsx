import { classroomsQueryOptions } from "@/api/classrooms";
import { Header } from "@/components/header";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Classroom } from "@/types/classroom";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Link, createFileRoute } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { Code } from "lucide-react";

export const Route = createFileRoute("/_auth/classrooms/")({
  component: Classrooms,
  loader: ({ context }) => context.queryClient.ensureQueryData(classroomsQueryOptions),
  pendingComponent: Loader,
});

function Classrooms() {
  const { data } = useSuspenseQuery(classroomsQueryOptions);
  return (
    <div className="p-2">
      <Header title="Own Classrooms">
        <Button asChild variant="default">
          <Link to="/classrooms/create">Create</Link>
        </Button>
      </Header>
      <ClassroomTable classrooms={data.ownClassrooms} />
      <Header title="Joined Classrooms" />
      <JoinedClassroomTable classrooms={data.joinedClassrooms} />
    </div>
  );
}

function ClassroomTable({ classrooms }: { classrooms: Classroom[] }) {
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
          <TableRow key={c.classroom.id}>
            <TableCell>{c.classroom.name}</TableCell>
            <TableCell>
              <a href={c.gitlabUrl} target="_blank" rel="noreferrer">
                <Code />
              </a>
            </TableCell>
            <TableCell className="text-right">
              <Button asChild variant="outline">
                <Link to="/classrooms/$classroomId" params={{ classroomId: c.classroom.id }}>
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

function JoinedClassroomTable({ classrooms }: { classrooms: Classroom[] }) {
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
