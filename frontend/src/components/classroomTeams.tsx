import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table.tsx";
import { AlertCircle, Edit, Loader2, SearchCode, UserPlus } from "lucide-react";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { Team, TeamResponse, UserClassroomResponse } from "@/swagger-client";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { CreateTeamForm } from "./createTeamForm";
import { useState } from "react";
import { ClassroomTeamModal } from "./classroomTeam";
import { isModerator, isStudent } from "@/lib/utils";
import { useUpdateTeam } from "@/api/team";
import { useForm } from "react-hook-form";
import { createFormSchema } from "@/types/team";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage, Form } from "./ui/form";
import { Alert, AlertTitle, AlertDescription } from "./ui/alert";
import { Input } from "./ui/input";
import { toast } from "sonner";
/**
 * TeamListCard is a React component that displays a list of members in a classroom.
 * It includes a table of members and a button to invite more members, if the user has the appropriate role.
 *
 * @param {Object} props - The properties passed to the component.
 * @param {Array} props.teams - An array of DefaultControllerGetOwnedClassroomTeamResponse[] objects representing the teams of the classroom.
 * @param {string} props.classroomId - The ID of the classroom.
 * @param {Role} props.userRole - The role of the current user in the classroom. This determines whether the invite button and view assignments-button is displayed.
 * @returns {JSX.Element} A React component that displays a card with the list of members in a classroom.
 */

export function TeamListCard({
  teams,
  classroomId,
  userClassroom,
  maxTeamSize,
  numInvitedMembers,
  studentsCanCreateTeams,
  deactivateInteraction,
  teamsReportUrls,
}: {
  teams: TeamResponse[];
  classroomId: string;
  userClassroom: UserClassroomResponse;
  maxTeamSize: number;
  numInvitedMembers: number;
  studentsCanCreateTeams: boolean;
  deactivateInteraction: boolean;
  teamsReportUrls: Map<string, string>;
}): JSX.Element {
  const teamSlots = teams.length * maxTeamSize;

  const [open, setOpen] = useState(false);

  return (
    <Card className="p-2">
      <CardHeader className="md:flex md:flex-row md:items-center justify-between space-y-0 pb-2 mb-4">
        <div className="mb-4 md:mb-1">
          <CardTitle className="mb-1">Teams</CardTitle>
          <CardDescription>All teams of this classroom</CardDescription>
        </div>
        <div className="grid grid-cols-1 gap-2">
          {isModerator(userClassroom) && !deactivateInteraction && (
            <Dialog open={open} onOpenChange={setOpen}>
              <DialogTrigger asChild>
                <Button variant="outline">Create a team</Button>
              </DialogTrigger>
              <DialogContent>
                <CreateTeamForm onSuccess={() => setOpen(false)} classroomId={classroomId} />
              </DialogContent>
            </Dialog>
          )}
        </div>
      </CardHeader>
      <CardContent>
        {teamSlots < numInvitedMembers && isModerator(userClassroom) && (
          <div>
            <p className="text-sm text-muted-foreground text-red-600">
              Not enough team spots to accommodate all classroom members.
            </p>
            {!studentsCanCreateTeams && (
              <p className="text-sm text-muted-foreground text-red-600">
                You have to add more teams, because students can't create teams by their own.
              </p>
            )}
          </div>
        )}
        <TeamTable
          teams={teams}
          classroomId={classroomId}
          userClassroom={userClassroom}
          maxTeamSize={maxTeamSize}
          teamsReportUrls={teamsReportUrls}
          deactivateInteraction={deactivateInteraction}
        />
      </CardContent>
    </Card>
  );
}

export function TeamTable({
  teams,
  teamsReportUrls,
  classroomId,
  userClassroom,
  maxTeamSize,
  isPending,
  onTeamSelect,
  deactivateInteraction,
}: {
  teams: TeamResponse[];
  teamsReportUrls: Map<string, string>;
  classroomId: string;
  userClassroom: UserClassroomResponse;
  maxTeamSize: number;
  isPending?: boolean;
  onTeamSelect?: (teamId: string) => void;
  deactivateInteraction: boolean;
}) {
  return (
    <Table>
      <TableBody>
        {teams.map((t) => {
          const reportUrl = teamsReportUrls.get(t.id)!;

          return (
            <TableRow key={t.id}>
              <TableCell className="p-2">
                <TeamListElement team={t} maxTeamSize={maxTeamSize} />
              </TableCell>
              <TableCell className="p-2 flex justify-end align-middle">
                {isModerator(userClassroom) && <ChangeTeamDialog classroomId={classroomId} team={t} />}
                <Button variant="ghost" size="icon" asChild title="Go to team">
                  <a href={t.webUrl} target="_blank" rel="noreferrer">
                    <SearchCode className="h-6 w-6 text-gray-600 dark:text-white" />
                  </a>
                </Button>
                {!deactivateInteraction && (
                  <>
                    {(!isStudent(userClassroom) || userClassroom.classroom.studentsViewAllProjects) && (
                      <ClassroomTeamModal
                        userClassroom={userClassroom}
                        classroomId={classroomId}
                        teamId={t.id}
                        reportUrl={reportUrl}
                      />
                    )}
                    {onTeamSelect && (
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => onTeamSelect?.(t.id)}
                        disabled={isPending || t.members.length >= maxTeamSize}
                        title="Get details"
                      >
                        <UserPlus className="text-gray-600 dark:text-white" />
                      </Button>
                    )}
                  </>
                )}
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}

function TeamListElement({ team, maxTeamSize }: { team: TeamResponse; maxTeamSize: number }) {
  return (
    <HoverCard>
      <HoverCardTrigger className="cursor-default flex">
        <div className="cursor-default">
          <div className="font-medium">{team.name}</div>
          <div className="text-sm text-muted-foreground md:inline">
            {team.members.length}/{maxTeamSize} member
          </div>
        </div>
      </HoverCardTrigger>
      <HoverCardContent className="w-100">
        <p className="text-lg font-semibold">{team.name}</p>
        <p className="text-sm text-muted-foreground mt-[-0.3rem]">
          {team.members.length}/{maxTeamSize} member
        </p>
        {team.members.length >= 1 && (
          <>
            <Separator className="my-1" />
            <div className="text-muted-foreground">
              {team.members.map((m) => (
                <div key={m.user.id}>
                  {m.user.gitlabUsername} - {m.user.name}
                </div>
              ))}
            </div>
          </>
        )}
      </HoverCardContent>
    </HoverCard>
  );
}

const ChangeTeamDialog = ({ classroomId, team }: { classroomId: string; team: Team }) => {
  const { mutateAsync, isError, isPending } = useUpdateTeam(classroomId, team.id);

  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    defaultValues: {
      name: team.name,
    },
  });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    console.log(values);
    await mutateAsync(values);
    toast.success("Team updated successfully!");
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="ghost" size="icon" title="Edit team">
          <Edit className="h-6 w-6 text-gray-600 dark:text-white" />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Team</DialogTitle>
          <DialogDescription>Change the name of the team</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="team name" {...field} />
                  </FormControl>
                  <FormDescription>This is your team name.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogClose asChild>
              <Button type="submit" disabled={isPending}>
                {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Submit"}
              </Button>
            </DialogClose>

            {isError && (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>The team could not be updated!</AlertDescription>
              </Alert>
            )}
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};
