import { Role } from "@/types/classroom.ts";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Link } from "@tanstack/react-router";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table.tsx";
import { Clipboard, Gitlab, UserPlus } from "lucide-react";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { TeamResponse } from "@/swagger-client";
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
  userRole,
  maxTeamSize,
  numInvitedMembers,
  deactivateInteraction,
}: {
  teams: TeamResponse[];
  classroomId: string;
  userRole: Role;
  maxTeamSize: number;
  numInvitedMembers: number;
  deactivateInteraction: boolean;
}): JSX.Element {
  const teamSlots = teams.length * maxTeamSize;
  return (
    <Card className="p-2">
      <CardHeader>
        <CardTitle>Teams</CardTitle>
        <CardDescription>Every team in this classroom</CardDescription>
        {teamSlots < numInvitedMembers && userRole != Role.Student && (
          <div>
            <p className="text-sm text-muted-foreground text-red-600">Not enough team spots to accommodate all classroom members.</p>
            <p className="text-sm text-muted-foreground text-red-600">Please add more teams!</p>
          </div>
        )}
      </CardHeader>
      <CardContent>
        <TeamTable 
          teams={teams} 
          classroomId={classroomId} 
          userRole={userRole} 
          maxTeamSize={maxTeamSize}
          deactivateInteraction={deactivateInteraction} 
        />
      </CardContent>
      {(userRole != Role.Student && !deactivateInteraction )&& (
        <CardFooter className="flex justify-end">
          <Button variant="default" asChild>
            <Link to="/classrooms/$classroomId/team/create/modal" replace params={{ classroomId }}>
              Create a team
            </Link>
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}

export function TeamTable({
  teams,
  classroomId,
  userRole,
  maxTeamSize,
  isPending,
  onTeamSelect,
  deactivateInteraction,
}: {
  teams: TeamResponse[];
  classroomId: string;
  userRole: Role;
  maxTeamSize: number;
  isPending?: boolean;
  onTeamSelect?: (teamId: string) => void;
  deactivateInteraction: boolean;
}) {
  return (
    <Table>
      <TableBody>
        {teams.map((t) => (
          <TableRow key={t.id}>
            <TableCell className="p-2">
              <TeamListElement team={t} maxTeamSize={maxTeamSize} />
            </TableCell>
            <TableCell className="p-2 flex justify-end align-middle">
              <Button variant="ghost" size="icon" asChild>
                <a href={t.webUrl} target="_blank" rel="noreferrer">
                  <Gitlab className="h-6 w-6 text-gray-600 dark:text-white" />
                </a>
              </Button>
              {!deactivateInteraction && (
                <>
                {userRole != Role.Student && (
                  <Button variant="ghost" size="icon" asChild>
                    <Link
                      to="/classrooms/$classroomId/teams/$teamId/modal"
                      params={{ classroomId: classroomId, teamId: t.id }}
                    >
                      <Clipboard className="h-6 w-6 text-gray-600 dark:text-white" />
                    </Link>
                  </Button>
                )}
                {onTeamSelect && (
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => onTeamSelect?.(t.id)}
                    disabled={isPending || t.members.length >= maxTeamSize}
                  >
                    <UserPlus className="text-gray-600 dark:text-white" />
                  </Button>
                )}
                </>
              )}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function TeamListElement({
  team,
  maxTeamSize,
}: {
  team: TeamResponse;
  maxTeamSize: number;
}) {
  return (
    <HoverCard>
      <HoverCardTrigger className="cursor-default flex">
        <div className="cursor-default">
          <div className="font-medium">{team.name}</div>
          <div className="text-sm text-muted-foreground md:inline">
            {team.members.length} / {maxTeamSize} member
          </div>
        </div>
      </HoverCardTrigger>
      <HoverCardContent className="w-100">
        <p className="text-lg font-semibold">{team.name}</p>
        <p className="text-sm text-muted-foreground mt-[-0.3rem]">
          {team.members.length} / {maxTeamSize} member
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
