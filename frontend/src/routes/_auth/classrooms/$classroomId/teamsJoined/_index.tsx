import { classroomQueryOptions } from "@/api/classroom";
import { teamsQueryOptions } from "@/api/team";
import { Header } from "@/components/header";
import { Loader } from "@/components/loader";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Role } from "@/types/classroom";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Link, Outlet, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teamsJoined/_index")({
  loader: async ({ context: { queryClient }, params }) => {
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    return { teams };
  },
  pendingComponent: Loader,
  component: Teams,
});

function Teams() {
  const { classroomId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));
  teams;
  return (
    <div className="pt-2">
      <Header title="Teams">
        {classroom.role !== Role.Student && (
          <Button variant="default" asChild>
            <Link to="/classrooms/owned/$classroomId/teams/create/modal" replace params={{ classroomId }}>
              Create Teams
            </Link>
          </Button>
        )}
      </Header>
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
      <Outlet />
    </div>
  );
}
