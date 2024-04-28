import { Role } from "@/types/classroom.ts";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Link } from "@tanstack/react-router";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table.tsx";
import { Clipboard, Gitlab } from "lucide-react";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { DefaultControllerGetOwnedClassroomTeamResponse } from "@/swagger-client";
type Team = DefaultControllerGetOwnedClassroomTeamResponse
/**
 * MemberListCard is a React component that displays a list of members in a classroom.
 * It includes a table of members and a button to invite more members, if the user has the appropriate role.
 *
 * @param {Object} props - The properties passed to the component.
 * @param {Array} props.teams - An array of DefaultControllerGetOwnedClassroomTeamResponse[] objects representing the teams of the classroom.
 * @param {string} props.classroomId - The ID of the classroom.
 * @param {Role} props.userRole - The role of the current user in the classroom. This determines whether the invite button and view assignments-button is displayed.
 * @returns {JSX.Element} A React component that displays a card with the list of members in a classroom.
 */
export function TeamListCard({ teams, classroomId, userRole }: {
  teams: Team[];
  classroomId: string,
  userRole: Role
}): JSX.Element {
  return (
    <Card className="p-2">
      <CardHeader>
        <CardTitle>Teams</CardTitle>
        <CardDescription>Every team in this classroom</CardDescription>
      </CardHeader>
      <CardContent>
        <TeamTable teams={teams} userRole={userRole} />
      </CardContent>
      {userRole != 2 &&
        <CardFooter className="flex justify-end">
          <Button variant="default" asChild>
            <Link to="/classrooms/owned/$classroomId/team/create/modal" replace params={{ classroomId }}>
              Create a team
            </Link>
          </Button>
        </CardFooter>}
    </Card>
  );
}

function TeamTable({ teams, userRole }: {
  teams: Team[];
  userRole: Role,
}) {
  return (
    <Table>
      <TableBody>
        {teams.map((t) => (
          <TableRow key={t.id}>
            <TableCell className="p-2">
              <TeamListElement team={t} />
            </TableCell>
            <TableCell className="p-2 flex justify-end align-middle">
              <Button variant="ghost" size="icon" asChild>
                <a href={t.gitlabUrl} target="_blank" rel="noreferrer">
                  <Gitlab className="h-6 w-6 text-gray-600" />
                </a>
              </Button>
              {userRole != 2 &&
                <Button variant="ghost"
                        size="icon"> {/* Should open a popup listing all assignments from that specific (team) */}
                  <Clipboard className="h-6 w-6 text-gray-600" />
                </Button>}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function TeamListElement({ team }: { team: Team }) {
  if(team.name == null || team.members == null) return (<>Error loading this team, team data is faulty</>);
  return (
    <HoverCard>
      <HoverCardTrigger className="cursor-default flex">
        <div className="cursor-default">
          <div className="font-medium">{team.name}</div>
          <div className="text-sm text-muted-foreground md:inline">
            {team.members.length} member{team.members.length != 1 ? "s" : ""}
          </div>
        </div>
      </HoverCardTrigger>
      <HoverCardContent className="w-100">
        <p className="text-lg font-semibold">{team.name}</p>
        <p
          className="text-sm text-muted-foreground mt-[-0.3rem]">{team.members.length} member{team.members.length != 1 ? "s" : ""}</p>
        {team.members.length >= 1 &&
          <>
            <Separator className="my-1" />
            <div className="text-muted-foreground">
              {team.members.map((m) => (
                <div key={m.id}>{m.gitlabUsername} - {m.name}</div>
              ))}
            </div>
          </>}
      </HoverCardContent>
    </HoverCard>
  );
}
