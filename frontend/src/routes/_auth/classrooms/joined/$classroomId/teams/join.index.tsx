import { joinedClassroomQueryOptions } from "@/api/classrooms";
import { joinedClassroomTeamsQueryOptions, useJoinTeam } from "@/api/teams";
import { CreateJoinedTeamForm } from "@/components/createJoinedTeamForm";
import { Header } from "@/components/header";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTrigger } from "@/components/ui/dialog";
import { TableCaption, TableHeader, TableRow, TableHead, TableBody, TableCell, Table } from "@/components/ui/table";
import { GetJoinedClassroomTeamResponse } from "@/swagger-client";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { Code } from "lucide-react";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/teams/join/")({
  loader: async ({ context, params }) => {
    const teams = await context.queryClient.ensureQueryData(joinedClassroomTeamsQueryOptions(params.classroomId));
    const joinedClassroom = await context.queryClient.ensureQueryData(joinedClassroomQueryOptions(params.classroomId));

    if (joinedClassroom.team) {
      throw redirect({
        to: "/classrooms/joined/$classroomId",
        params,
        replace: true,
      });
    }

    return { teams, joinedClassroom };
  },
  component: JoinTeam,
});

function JoinTeam() {
  const navigate = useNavigate();
  const { classroomId } = Route.useParams();
  const { data: joinedClassroom } = useSuspenseQuery(joinedClassroomQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(joinedClassroomTeamsQueryOptions(classroomId));

  const { mutateAsync, isPending } = useJoinTeam(classroomId);

  const joinTeam = async (teamId: string) => {
    await mutateAsync(teamId);
    await navigate({
      to: "/classrooms/joined/$classroomId",
      params: { classroomId },
    });
  };

  return (
    <div className="p-2">
      <Header title={`Join a team of ${joinedClassroom.classroom.name}`}>
        {joinedClassroom.classroom.createTeams && teams.length < joinedClassroom.classroom.maxTeams && (
          <Dialog>
            <DialogTrigger asChild>
              <Button variant="default">Create</Button>
            </DialogTrigger>
            <DialogHeader>Create a new Team</DialogHeader>
            <DialogContent>
              <CreateJoinedTeamForm classroomId={classroomId} />
            </DialogContent>
          </Dialog>
        )}
      </Header>
      <TeamsTable
        teams={teams}
        isPending={isPending}
        joinTeam={joinTeam}
        maxTeamSize={joinedClassroom.classroom.maxTeamSize}
      />
    </div>
  );
}

interface TeamsTableProps {
  teams: GetJoinedClassroomTeamResponse[];
  isPending: boolean;
  joinTeam: (teamId: string) => Promise<void>;
  maxTeamSize: number;
}

function TeamsTable({ teams, isPending, joinTeam, maxTeamSize }: TeamsTableProps) {
  return (
    <Table>
      <TableCaption>Teams to join</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Member Count</TableHead>
          <TableHead>Gitlab-Link</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {teams.map((t) => (
          <TableRow key={t.id}>
            <TableCell>{t.name}</TableCell>
            <TableCell>{t.members!.length}</TableCell>
            <TableCell>
              <a href={t.gitlabUrl} target="_blank" rel="noreferrer">
                <Code />
              </a>
            </TableCell>
            <TableCell className="text-right">
              <Button
                disabled={isPending || t.members!.length >= maxTeamSize}
                onClick={() => joinTeam(t.id!)}
                variant="outline"
              >
                Join Team
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
