import { classroomQueryOptions } from "@/api/classroom";
import { teamsQueryOptions, useJoinTeam } from "@/api/team";
import { CreateTeamForm } from "@/components/createTeamForm";
import { Header } from "@/components/header";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { Role } from "@/types/classroom.ts";
import { TeamTable } from "@/components/classroomTeams";
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card.tsx";
import { ReportApiAxiosParamCreator } from "@/swagger-client";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/teams/join/")({
  loader: async ({ context: { queryClient }, params }) => {
    const userClassroom = await queryClient.fetchQuery(classroomQueryOptions(params.classroomId));

    if (userClassroom.classroom.maxTeamSize === 1) {
      throw redirect({
        to: "/classrooms/$classroomId",
        search: { tab: "assignments" },
        params,
        replace: true,
      });
    }

    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));

    const teamsReportUrls = (
      await Promise.all(
        teams.map(async (team) => ({
          teamId: team.id,
          url: (await ReportApiAxiosParamCreator().getClassroomTeamReport(params.classroomId, team.id)).url,
        })),
      )
    ).reduce((acc, { url, teamId }) => acc.set(teamId, url), new Map<string, string>());

    if (userClassroom.team || userClassroom.role !== Role.Student) {
      throw redirect({
        to: "/classrooms/$classroomId",
        params,
        search: { tab: "assignments" },
        replace: true,
      });
    }

    return { teams, userClassroom, teamsReportUrls };
  },
  component: JoinTeam,
});

function JoinTeam() {
  const navigate = Route.useNavigate();
  const { classroomId } = Route.useParams();
  const { data: joinedClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));
  const { teamsReportUrls } = Route.useLoaderData();

  const { mutateAsync, isPending } = useJoinTeam(classroomId);

  const joinTeam = async (teamId: string) => {
    await mutateAsync(teamId);
    await navigate({
      to: "/classrooms/$classroomId",
      search: { tab: "assignments" },
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
            teamsReportUrls={teamsReportUrls}
            isPending={isPending}
            classroomId={classroomId}
            userClassroom={joinedClassroom}
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
                    onSuccess={() =>
                      navigate({
                        to: "/classrooms/$classroomId/",
                        search: { tab: "assignments" },
                        params: { classroomId },
                      })
                    }
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
