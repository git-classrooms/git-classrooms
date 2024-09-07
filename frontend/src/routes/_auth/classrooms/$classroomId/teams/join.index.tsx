import { classroomQueryOptions } from "@/api/classroom";
import { teamsQueryOptions, useJoinTeam } from "@/api/team";
import { CreateTeamForm } from "@/components/createTeamForm";
import { Header } from "@/components/header";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { Role } from "@/types/classroom.ts";
import { TeamTable } from "@/components/classroomTeams";
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card.tsx";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teams/join/")({
  loader: async ({ context: { queryClient }, params }) => {
    const userClassroom = await queryClient.fetchQuery(classroomQueryOptions(params.classroomId));
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));

    if (userClassroom.team || userClassroom.role !== Role.Student) {
      throw redirect({
        to: "/classrooms/$classroomId",
        params,
        replace: true,
      });
    }

    return { teams, userClassroom };
  },
  component: JoinTeam,
});

function JoinTeam() {
  const navigate = Route.useNavigate();
  const { classroomId } = Route.useParams();
  const { data: joinedClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));

  const { mutateAsync, isPending } = useJoinTeam(classroomId);

  const joinTeam = async (teamId: string) => {
    await mutateAsync(teamId);
    await navigate({
      to: "/classrooms/$classroomId",
      params: { classroomId },
    });
  };
  const freeTeamSlot = (): boolean => {
    return teams.some((team) => team.members.length < joinedClassroom.classroom.maxTeamSize);
  };

  return (
    <>
      <Header
        title={`Join a team of ${joinedClassroom.classroom.name}`}
        subtitle={joinedClassroom.classroom.description}
      />
      <Card className="p-2">
        <CardHeader>
          {joinedClassroom.classroom.createTeams
            ? "Choose a team you want to join or create a new team."
            : "Please select a team. "}
          {!joinedClassroom.classroom.createTeams && !freeTeamSlot() && (
            <div>
              <p className="text-sm text-muted-foreground text-red-600">There currently are no teams you can join.</p>
              <p className="text-sm text-muted-foreground text-red-600">
                Please contact the owner of this classroom to add more teams or raise the team-size
              </p>
            </div>
          )}
        </CardHeader>
        <CardContent>
          <TeamTable
            teams={teams}
            isPending={isPending}
            classroomId={classroomId}
            userRole={Role.Student}
            maxTeamSize={joinedClassroom.classroom.maxTeamSize}
            onTeamSelect={joinTeam}
            deactivateInteraction={false}
          />
        </CardContent>
        <CardFooter className="flex justify-end">
          {joinedClassroom.classroom.createTeams &&
            (joinedClassroom.classroom.maxTeams === 0 || teams.length < joinedClassroom.classroom.maxTeams) && (
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="default">Create new Team</Button>
                </DialogTrigger>
                <DialogContent>
                  <CreateTeamForm
                    onSuccess={() => navigate({ to: "/classrooms/$classroomId/", params: { classroomId } })}
                    classroomId={classroomId}
                  />
                </DialogContent>
              </Dialog>
            )}
        </CardFooter>
      </Card>
    </>
  );
}
