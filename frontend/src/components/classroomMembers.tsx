import { Role, UserClassroom } from "@/types/classroom.ts";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Link } from "@tanstack/react-router";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table.tsx";
import { Clipboard, Gitlab } from "lucide-react";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Avatar } from "@/components/avatar.tsx";
import { Separator } from "@/components/ui/separator.tsx";

/**
 * MemberListCard is a React component that displays a list of members in a classroom.
 * It includes a table of members and a button to invite more members, if the user has the appropriate role.
 *
 * @param {Object} props - The properties passed to the component.
 * @param {Array} props.classroomMembers - An array of UserClassroom objects representing the members of the classroom.
 * @param {string} props.classroomId - The ID of the classroom.
 * @param {Role} props.userRole - The role of the current user in the classroom. This determines whether the invite button and view assignments-button is displayed.
 * @returns {JSX.Element} A React component that displays a card with the list of members in a classroom.
 */
export function MemberListCard({ classroomMembers, classroomId, userRole, gitlabUserUrl }: {
  classroomMembers: UserClassroom[],
  classroomId: string,
  userRole: Role
  gitlabUserUrl: string
}): JSX.Element {
  return (
    <Card className="p-2">
      <CardHeader>
        <CardTitle>Members</CardTitle>
        <CardDescription>Every person in this classroom</CardDescription>
      </CardHeader>
      <CardContent>
        <MemberTable members={classroomMembers} userRole={userRole} gitlabUserUrl={gitlabUserUrl} />
      </CardContent>
      {userRole != 2 &&
        <CardFooter className="flex justify-end">
          <Button variant="default" asChild>
            <Link to="/classrooms/owned/$classroomId/invite" params={{ classroomId }}>
              Invite members
            </Link>
          </Button>
        </CardFooter>}
    </Card>
  );
}

function MemberTable({ members, userRole, gitlabUserUrl }: {
  members: UserClassroom[],
  userRole: Role,
  gitlabUserUrl: string
}) {
  return (
    <Table>
      <TableBody>
        {members.map((m) => (
          <TableRow key={m.user.id}>
            <TableCell className="p-2">
              <MemberListElement member={m} />
            </TableCell>
            <TableCell className="p-2 flex justify-end align-middle">
              <Button variant="ghost" size="icon" onClick={ () => location.href = (gitlabUserUrl.replace(":userId", m.user.id.toString()))}>
                <Gitlab color="#666" className="h-6 w-6" />
              </Button>
              {userRole != 2 &&
                <Button variant="ghost"
                        size="icon"> {/* Should open a popup listing all assignments from that specific user('s team) */}
                  <Clipboard color="#666" className="h-6 w-6" />
                </Button>}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function MemberListElement({ member }: { member: UserClassroom }) {
  return (
    <HoverCard>
      <HoverCardTrigger className="cursor-default flex">
        <div className="pr-2"><Avatar
          avatarUrl={member.user.gitlabAvatar?.avatarURL}
          fallbackUrl={member.user.gitlabAvatar?.fallbackAvatarURL}
          name={member.user.name!}
        /></div>
        <div>
          <div className="font-medium">{member.user.name}</div>
          <div className="hidden text-sm text-muted-foreground md:inline">
            {member.role === 0 ? "Owner" : member.role === 1 ? "Moderator" : "Student"} {member.team ? `- ${member.team.name}` : ""}
          </div>
        </div>
      </HoverCardTrigger>
      <HoverCardContent className="w-80">
        <div className="flex justify-between space-x-4">
          <div className="space-y-1">
            <p className="text-sm font-semibold">{member.user.name}</p>
            <p className="hidden text-sm text-muted-foreground md:inline">@{member.user.gitlabUsername}</p>
            <Separator className="my-4" />
            <p className="text-sm text-muted-foreground">{member.user.gitlabEmail}</p>
            <Separator className="my-4" />
            <div className="text-sm text-muted-foreground">
              <p className="font-bold md:inline">
                {member.role === 0 ? "Owner" : member.role === 1 ? "Moderator" : "Student"}</p> of
              this classroom {member.team ?
              <div className="md:inline">in team <p className="font-bold md:inline">${member.team.name}</p></div> : ""}
            </div>
          </div>
        </div>
      </HoverCardContent>
    </HoverCard>
  );
}
