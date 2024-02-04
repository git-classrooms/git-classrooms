import { classRoomsQueryOptions } from "@/api/classrooms";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Classroom } from "@/types/classroom";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Link, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms")({
  component: Classrooms,
  loader: ({ context }) => context.queryClient.ensureQueryData(classRoomsQueryOptions)
});

function Classrooms() {
  const { data } = useSuspenseQuery(classRoomsQueryOptions)
  return (
    <div className="p-2">
      <div className="flex flex-row justify-between">
        <h1 className="text-xl font-bold">Classrooms</h1>
        <Button asChild variant="default">
          <Link to="/classrooms/create" >Create</Link>
        </Button>
      </div>
      <ClassroomTable title="Your classrooms" classRooms={data.ownClassrooms}>
        <Button variant="outline"> Show Classroom </Button>
      </ClassroomTable>
    </div>
  );
}


function ClassroomTable({ classRooms, title, children }: { children: React.ReactNode, title: string, classRooms: Classroom[] }) {
  return (
    <Table>
      <TableCaption>{title}</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Onwer</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {classRooms.map(c =>
          <TableRow key={c.classroom.id}>
            <TableCell>{c.classroom.name}</TableCell>
            <TableCell>{c.classroom.ownerId}</TableCell>
            <TableCell className="text-right">
              {children}
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  )
}
