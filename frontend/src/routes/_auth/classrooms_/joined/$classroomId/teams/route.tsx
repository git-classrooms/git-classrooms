import { joinedClassroomQueryOptions } from "@/api/classrooms";
import { joinedClassroomTeamsQueryOptions } from "@/api/teams";
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
import { Link, Outlet, createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/teams")({
  loader: async ({ context, params }) => {
    const classroom = await context.queryClient.ensureQueryData(joinedClassroomQueryOptions(params.classroomId));
    const teams = await context.queryClient.ensureQueryData(joinedClassroomTeamsQueryOptions(params.classroomId));
    if (classroom.classroom.maxTeamSize === 1) {
      throw redirect({ to: "/classrooms/joined/$classroomId/", params });
    }
    return { teams };
  },
  pendingComponent: Loader,
  component: Teams,
});

function Teams() {
  const { classroomId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(joinedClassroomQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(joinedClassroomTeamsQueryOptions(classroomId));
  teams;
  return (
    <div className="pt-2">
      <Header title="Teams">
        {classroom.role === Role.Moderator && (
          <Button variant="default" asChild>
            <Link to="/classrooms/owned/$classroomId/teams/create/modal" params={{ classroomId }}>
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
