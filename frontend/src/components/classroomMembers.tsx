import { getRole } from "@/types/classroom.ts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Link } from "@tanstack/react-router";
import { Clipboard, Gitlab } from "lucide-react";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Avatar } from "@/components/avatar.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { UserClassroomResponse } from "@/swagger-client";
import List from "@/components/ui/list.tsx";
import ListItem from "@/components/ui/listItem.tsx";
import { ClassroomTeamModal } from "./classroomTeam";
import { isModerator, isOwner, isStudent } from "@/lib/utils";

/**
 * MemberListCard is a React component that displays a list of members in a classroom.
 * It includes a table of members and a button to invite more members, if the user has the appropriate role.
 *
 * @param {Object} props - The properties passed to the component.
 * @param {Array} props.classroomMembers - An array of UserClassroom objects representing the members of the classroom.
 * @param {string} props.classroomId - The ID of the classroom.
 * @param {Role} props.userRole - The role of the current user in the classroom. This determines whether the invite button and view assignments-button is displayed.
 * @param {boolean} props.showTeams - A boolean indicating whether to show the teams of the members.
 * @param {boolean} props.deactivateInteraction - A boolean indicating whether the user can interact with the members.
 * @returns {JSX.Element} A React component that displays a card with the list of members in a classroom.
 */
export function MemberListCard({
  classroomMembers,
  classroomId,
  userClassroom,
  showTeams,
  deactivateInteraction,
}: {
  classroomMembers: UserClassroomResponse[];
  classroomId: string;
  userClassroom: UserClassroomResponse;
  showTeams: boolean;
  deactivateInteraction: boolean;
}): JSX.Element {
  return (
    <Card className="p-2">
      <CardHeader className="md:flex md:flex-row md:items-center justify-between space-y-0 pb-2 mb-4">
        <div className="mb-4 md:mb-0">
          <CardTitle className="mb-1">Members</CardTitle>
          <CardDescription>Members in this this classroom</CardDescription>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
          {!deactivateInteraction && isOwner(userClassroom) && (
            <Button variant="outline" asChild>
              <Link to="/classrooms/$classroomId/members" params={{ classroomId }}>
                Change roles
              </Link>
            </Button>
          )}

          {!deactivateInteraction && isModerator(userClassroom) && (
            <Button variant="outline" asChild>
              <Link to="/classrooms/$classroomId/invite" params={{ classroomId }}>
                Invite members
              </Link>
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <MemberTable
          members={classroomMembers}
          classroomId={classroomId}
          userClassroom={userClassroom}
          showTeams={showTeams}
        />
      </CardContent>
    </Card>
  );
}

function MemberTable({
  members,
  classroomId,
  userClassroom,
  showTeams,
}: {
  members: UserClassroomResponse[];
  classroomId: string;
  userClassroom: UserClassroomResponse;
  showTeams: boolean;
}) {
  return (
    <List
      items={members}
      renderItem={(m) => (
        <ListItem
          leftContent={<MemberListElement member={m} showTeams={showTeams} />}
          rightContent={
            <>
              <Button variant="ghost" size="icon" asChild>
                <a href={m.webUrl} target="_blank" rel="noreferrer">
                  <Gitlab className="h-6 w-6 text-gray-600" />
                </a>
              </Button>
              {(!isStudent(userClassroom) || userClassroom.classroom.studentsViewAllProjects) && m.team ? (
                <ClassroomTeamModal classroomId={classroomId} teamId={m.team.id} />
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
