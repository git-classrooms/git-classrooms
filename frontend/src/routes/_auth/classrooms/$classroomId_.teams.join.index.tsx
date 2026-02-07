import { classroomQueryOptions } from "@/api/classroom";
import { teamsQueryOptions, useJoinTeam } from "@/api/team";
import { CreateTeamForm } from "@/components/createTeamForm";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { Role } from "@/types/classroom.ts";
import { TeamTable } from "@/components/classroomTeams";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { ReportApiAxiosParamCreator } from "@/swagger-client";
import { AlertCircle, ArrowLeft, Plus, Users2 } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { useState } from "react";

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
        }))
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
  const [dialogOpen, setDialogOpen] = useState(false);

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

  const canCreateTeam =
    joinedClassroom.classroom.createTeams &&
    (joinedClassroom.classroom.maxTeams === 0 || teams.length < joinedClassroom.classroom.maxTeams);

  const noTeamsAvailable = !joinedClassroom.classroom.createTeams && !freeTeamSlot();

  return (
    <div className="space-y-8 animate-stagger-1">
      {/* Breadcrumb */}
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms">Classrooms</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId }}>
                {joinedClassroom.classroom.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Join Team</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId }}>
              <ArrowLeft className="w-5 h-5" />
            </Link>
          </Button>
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Join a Team</h1>
            <p className="text-muted-foreground mt-1">{joinedClassroom.classroom.name}</p>
          </div>
        </div>
        {canCreateTeam && (
          <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
            <DialogTrigger asChild>
              <Button variant="glow">
                <Plus className="w-4 h-4 mr-2" />
                Create Team
              </Button>
            </DialogTrigger>
            <DialogContent>
              <CreateTeamForm
                onSuccess={() => {
                  setDialogOpen(false);
                  navigate({
                    to: "/classrooms/$classroomId/",
                    search: { tab: "assignments" },
                    params: { classroomId },
                  });
                }}
                classroomId={classroomId}
              />
            </DialogContent>
          </Dialog>
        )}
      </div>

      {/* Warning if no teams available */}
      {noTeamsAvailable && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>No teams available</AlertTitle>
          <AlertDescription>
            There are currently no teams you can join. Please contact the classroom owner to add more teams or increase
            the team size.
          </AlertDescription>
        </Alert>
      )}

      {/* Teams Card */}
      <Card>
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-primary/15 flex items-center justify-center">
              <Users2 className="w-5 h-5 text-primary" />
            </div>
            <div>
              <CardTitle>Available Teams</CardTitle>
              <CardDescription>
                {joinedClassroom.classroom.createTeams
                  ? "Choose an existing team or create your own"
                  : "Select a team to join"}
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {teams.length === 0 ? (
            <div className="border border-dashed border-border rounded-lg p-12 text-center">
              <Users2 className="w-10 h-10 text-muted-foreground mx-auto mb-3" />
              <h3 className="font-medium text-foreground mb-1">No teams yet</h3>
              <p className="text-sm text-muted-foreground mb-4">
                {canCreateTeam ? "Be the first to create a team!" : "No teams have been created yet"}
              </p>
              {canCreateTeam && (
                <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
                  <DialogTrigger asChild>
                    <Button variant="outline">
                      <Plus className="w-4 h-4 mr-2" />
                      Create Team
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <CreateTeamForm
                      onSuccess={() => {
                        setDialogOpen(false);
                        navigate({
                          to: "/classrooms/$classroomId/",
                          search: { tab: "assignments" },
                          params: { classroomId },
                        });
                      }}
                      classroomId={classroomId}
                    />
                  </DialogContent>
                </Dialog>
              )}
            </div>
          ) : (
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
          )}
        </CardContent>
      </Card>
    </div>
  );
}
