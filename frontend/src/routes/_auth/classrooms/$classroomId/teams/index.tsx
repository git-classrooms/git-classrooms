import { Button } from "@/components/ui/button";
import { Outlet, createFileRoute } from "@tanstack/react-router";
import { teamsQueryOptions } from "@/api/team.ts";
import { Loader } from "@/components/loader.tsx";
import { Role } from "@/types/classroom.ts";
import { classroomQueryOptions } from "@/api/classroom.ts";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Header } from "@/components/header.tsx";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { useState } from "react";
import { CreateTeamForm } from "@/components/createTeamForm";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teams/")({
  loader: async ({ context: { queryClient }, params }) => {
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    const classroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    return { teams, classroom };
  },
  pendingComponent: Loader,
  component: TeamsIndex,
});

function TeamsIndex() {
  const { classroomId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));

  const [open, setOpen] = useState(false);

  return (
    <div className="pt-2">
      <Header title="Teams">
        {classroom.role !== Role.Student && (
          <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
              <Button variant="default">Create a team</Button>
            </DialogTrigger>
            <DialogContent>
              <CreateTeamForm onSuccess={() => setOpen(false)} classroomId={classroomId} />
            </DialogContent>
          </Dialog>
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
