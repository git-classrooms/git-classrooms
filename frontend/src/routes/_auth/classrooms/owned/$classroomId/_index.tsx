import { ownedClassroomMemberQueryOptions, ownedClassroomQueryOptions } from "@/api/classrooms";
import { Header } from "@/components/header";
import { Loader } from "@/components/loader";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Assignment } from "@/types/assignments";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { ownedAssignmentsQueryOptions } from "@/api/assignments.ts";
import { formatDate } from "@/lib/utils.ts";
import { ownedClassroomTeamsQueryOptions } from "@/api/teams";
import { DefaultControllerGetOwnedClassroomTeamResponse } from "@/swagger-client";
import { Separator } from "@/components/ui/separator"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card";
import { UserClassroom } from "@/types/classroom.ts";
import { Gitlab, Clipboard } from 'lucide-react';
import { Avatar } from "@/components/avatar";
export const Route = createFileRoute("/_auth/classrooms/owned/$classroomId/_index")({
  component: ClassroomDetail,
  loader: async ({ context, params }) => {
    const classroom = await context.queryClient.ensureQueryData(ownedClassroomQueryOptions(params.classroomId));
    const assignments = await context.queryClient.ensureQueryData(ownedAssignmentsQueryOptions(params.classroomId));
    const members = await context.queryClient.ensureQueryData(ownedClassroomMemberQueryOptions(params.classroomId));
    const teams = await context.queryClient.ensureQueryData(ownedClassroomTeamsQueryOptions(params.classroomId));

    return { classroom, assignments, members, teams };
  },
  pendingComponent: Loader,
});

function ClassroomDetail() {
  const { classroomId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(ownedClassroomQueryOptions(classroomId));
  const { data: assignments } = useSuspenseQuery(ownedAssignmentsQueryOptions(classroomId));
  const { data: members } = useSuspenseQuery(ownedClassroomMemberQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(ownedClassroomTeamsQueryOptions(classroomId));

  return (
    <div className="p-2 space-y-6">
      <Outlet />
      <div className="flex flex-row justify-between">
        <h1 className="text-xl font-bold">Classroom Details </h1>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{classroom.name}</CardTitle>
          <CardDescription>{classroom.description}</CardDescription>
        </CardHeader>
      </Card>

      <Header title="Assignments">
        <Button variant="default" asChild>
          <Link to="/classrooms/owned/$classroomId/assignments/create" params={{ classroomId }}>
            Create assignment
          </Link>
        </Button>
      </Header>
      <AssignmentTable assignments={assignments} classroomId={classroomId} />

      <MemberListCard classroomMembers={members} classroomId={classroomId} />

      {classroom.maxTeamSize > 1 && (
        <>
          <Header title="Teams">
            <Button variant="default" asChild>
              <Link to="/classrooms/owned/$classroomId/team/create/modal" replace params={{ classroomId }}>
                Create Teams
              </Link>
            </Button>
          </Header>
          <TeamTable teams={teams} />
        </>
      )}
    </div>
  );
}

function AssignmentTable({ assignments, classroomId }: { assignments: Assignment[]; classroomId: string }) {
  return (
    <Table>
      <TableCaption>Assignments</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Due date</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {assignments.map((a) => (
          <TableRow key={a.id}>
            <TableCell>{a.name}</TableCell>
            <TableCell>{a.dueDate ? formatDate(a.dueDate) : "No Due Date"}</TableCell>
            <TableCell className="text-right">
              <Button asChild>
                <Link
                  to="/classrooms/owned/$classroomId/assignments/$assignmentId"
                  params={{ classroomId, assignmentId: a.id }}> Show Assignment
                </Link>
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function MemberListCard({ classroomMembers, classroomId }: { classroomMembers: UserClassroom[], classroomId: string }) {
  return (
    <Card className="p-2">
      <CardHeader>
        <CardTitle>Members</CardTitle>
        <CardDescription>Every person in this classroom</CardDescription>
      </CardHeader>
      <CardContent>
        <MemberTable members={classroomMembers} />
      </CardContent>
      <CardFooter className="flex justify-end">
        <Button variant="default" asChild>
          <Link to="/classrooms/owned/$classroomId/invite" params={{ classroomId }}>
            Invite members
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}

function MemberTable({ members }: { members: UserClassroom[] }) {
  return (
    <Table>
      <TableBody>
        {members.map((m) => (
          <TableRow key={m.user.id}>
            <TableCell className="p-2">
              <MemberListElement member={m} />
            </TableCell>
            <TableCell className="text-right p-2">
              <Button variant="ghost" size="icon" onClick={ () => location.href = `https://hs-flensburg.dev/${m.user.gitlabUsername}`}>
                <Gitlab color="#666" className="h-6 w-6" />
              </Button>
              <Button variant="ghost" size="icon"> {/* Should open a popup listing all assignments from that specific user('s team) */}
                <Clipboard color="#666" className="h-6 w-6" />
              </Button>
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
          {member.role === 0 ? "Owner" : member.role === 1 ? "Moderator" : "Member"} {member.team ? `- ${member.team.name}` : ""}
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
                {member.role === 0 ? "Owner" : member.role === 1 ? "Moderator" : "Member"}</p> of
              this classroom {member.team ?
              <div className="md:inline">in team <p className="font-bold md:inline">${member.team.name}</p></div> : ""}
            </div>
          </div>
        </div>
      </HoverCardContent>
    </HoverCard>
  );
}

function TeamTable({ teams }: { teams: DefaultControllerGetOwnedClassroomTeamResponse[] }) {
  return (
    <Table>
      <TableCaption>Teams</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {teams.map((t) => (
          <TableRow key={t.id}>
            <TableCell>{t.name}</TableCell>
            <TableCell className="text-right">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button>Actions</Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuItem>Test</DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
