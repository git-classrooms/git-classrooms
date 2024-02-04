import { classroomsQueryOptions } from "@/api/classrooms";
import { Header } from "@/components/header";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Classroom } from "@/types/classroom";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Link, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/")({
  component: Classrooms,
  loader: ({ context }) => context.queryClient.ensureQueryData(classroomsQueryOptions),
});

function Classrooms() {
  const { data } = useSuspenseQuery(classroomsQueryOptions);
  return (
    <div className="p-2">
      <Header title="Classrooms">
        <Button asChild variant="default">
          <Link to="/classrooms/create">Create</Link>
        </Button>
      </Header>
      <ClassroomTable classrooms={data.ownClassrooms} />
    </div>
  );
}

function ClassroomTable({ classrooms }: { classrooms: Classroom[] }) {
  return (
    <Table>
      <TableCaption>Classrooms</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Owner</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {classrooms.map((c) => (
          <TableRow key={c.classroom.id}>
            <TableCell>{c.classroom.name}</TableCell>
            <TableCell>{c.classroom.ownerId}</TableCell>
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
