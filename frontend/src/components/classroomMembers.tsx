import { getRole, Role } from "@/types/classroom.ts";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Link } from "@tanstack/react-router";
import { Clipboard, Gitlab } from "lucide-react";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Avatar } from "@/components/avatar.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { UserClassroomResponse } from "@/swagger-client";
import List from "@/components/ui/list.tsx";
import ListItem from "@/components/ui/listItem.tsx";

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
export function MemberListCard({
  classroomMembers,
  classroomId,
  userRole,
  showTeams,
}: {
  classroomMembers: UserClassroomResponse[];
  classroomId: string;
  userRole: Role;
  showTeams: boolean;
}): JSX.Element {
  return (
    <Card className="p-2">
      <CardHeader>
        <CardTitle>Members</CardTitle>
        <CardDescription>Every person in this classroom</CardDescription>
      </CardHeader>
      <CardContent>
        <MemberTable members={classroomMembers} classroomId={classroomId} userRole={userRole} showTeams={showTeams} />
      </CardContent>
      {userRole != 2 && (
        <CardFooter className="flex justify-end">
          <Button variant="default" asChild className="mr-2">
            <Link to="/classrooms/$classroomId/members" params={{ classroomId }}>
              Show all members
            </Link>
          </Button>
          <Button variant="default" asChild>
            <Link to="/classrooms/$classroomId/invite" params={{ classroomId }}>
              Invite members
            </Link>
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}

function MemberTable({
  members,
  classroomId,
  userRole,
  showTeams,
}: {
  members: UserClassroomResponse[];
  classroomId: string;
  userRole: Role;
  showTeams: boolean;
}) {
  return (
    <List
      items={members}
      renderItem={(m) => (
        <ListItem
          leftContent={
            <MemberListElement member={m} showTeams={showTeams} />
          }
          rightContent={
            <>
              <Button variant="ghost" size="icon" asChild>
                <a href={m.webUrl} target="_blank" rel="noreferrer">
                  <Gitlab className="h-6 w-6 text-gray-600" />
                </a>
              </Button>
              {userRole != Role.Student && m.team ? (
                <Button variant="ghost" size="icon" asChild>
                  <Link
                    to="/classrooms/$classroomId/teams/$teamId/modal"
                    params={{ classroomId: classroomId, teamId: m.team.id }}
                  >
                    <Clipboard className="h-6 w-6 text-gray-600" />
                  </Link>
                </Button>
              ) : (
                <Button variant="ghost" size="icon" asChild>
                  <div>
                    <Clipboard className="h-6 w-6 text-gray-400" />
                  </div>
                </Button>
              )}
            </>
          }
        />
      )}
    />
  );
}

function MemberListElement({ member, showTeams }: { member: UserClassroomResponse; showTeams: boolean }) {
  return (
    <HoverCard>
      <HoverCardTrigger className="cursor-default flex">
        <div className="pr-2">
          <Avatar
            avatarUrl={member.user.gitlabAvatar?.avatarURL}
            fallbackUrl={member.user.gitlabAvatar?.fallbackAvatarURL}
            name={member.user.name!}
          />
        </div>
        <div>
          <div className="font-medium">{member.user.name}</div>
          <div className="text-sm text-muted-foreground md:inline">
            {getRole(member.role)} {showTeams && member.team ? `- ${member.team.name}` : ""}
          </div>
        </div>
      </HoverCardTrigger>
      <HoverCardContent className="w-100">
        <p className="text-lg font-semibold">{member.user.name}</p>
        <p className="text-sm text-muted-foreground mt-[-0.3rem]">@{member.user.gitlabUsername}</p>
        <Separator className="my-1" />
        <p className="text-muted-foreground">{member.user.gitlabEmail}</p>
        <Separator className="my-1" />
        <div className="text-muted-foreground">
          <span className="font-bold">{getRole(member.role)}</span> of this classroom{" "}
          {showTeams && member.team ? (
            <>
              {" "}
              in team <span className="font-bold">{member.team?.name ?? ""}</span>
            </>
          ) : (
            ""
          )}
        </div>
      </HoverCardContent>
    </HoverCard>
  );
}
