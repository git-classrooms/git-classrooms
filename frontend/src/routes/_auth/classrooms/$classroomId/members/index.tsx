import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { getRole, Role } from "@/types/classroom.ts";
import { createFormSchema } from "@/types/member.ts";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Header } from "@/components/header.tsx";
import { membersQueryOptions, useRemoveTeamMember, useUpdateMemberRole, useUpdateMemberTeam } from "@/api/member.ts";
import { ReportApiAxiosParamCreator, Team, TeamResponse, UserClassroomResponse } from "@/swagger-client";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Avatar } from "@/components/avatar.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { teamsQueryOptions } from "@/api/team.ts";
import { classroomQueryOptions } from "@/api/classroom.ts";
import { assignmentsQueryOptions } from "@/api/assignment.ts";
import { Loader } from "@/components/loader.tsx";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form";
import { AlertCircle } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { useMemo } from "react";
import { isCreator, isStudent } from "@/lib/utils";
import { Table, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/members/")({
  component: Members,
  beforeLoad: async ({ context: { queryClient }, params: { classroomId } }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(classroomId));
    if (isStudent(userClassroom)) {
      throw redirect({ to: "/classrooms/$classroomId", params: { classroomId }, search: { tab: "assignments" } });
    }
  },
  loader: async ({ context: { queryClient }, params }) => {
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));

    const { url: reportDownloadUrl } = await ReportApiAxiosParamCreator().getClassroomReport(params.classroomId);
    const members = await queryClient.ensureQueryData(membersQueryOptions(params.classroomId));
    if (isStudent(userClassroom)) {
      const assignments = await queryClient.ensureQueryData(assignmentsQueryOptions(params.classroomId));
      return { userClassroom, assignments, members, teams, reportDownloadUrl };
    } else {
      return { userClassroom, members, teams, reportDownloadUrl };
    }
  },
  pendingComponent: Loader,
});

function Members() {
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));

  const classroomMembersSorted = useMemo(
    () =>
      [...classroomMembers].sort((a, b) => {
        if (a.role !== b.role) {
          return a.role - b.role;
        }

        if (isCreator(a)) return -1;
        if (isCreator(b)) return 1;

        return a.user.name.localeCompare(b.user.name);
      }),
    [classroomMembers],
  );

  return (
    <>
      <Breadcrumb className="mb-5">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId }}>
                {userClassroom.classroom.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Manage members</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <Header title="Manage members" subtitle="Change the roles and teams of members" />
      <div className="justify-between gap-10">
        <MemberTable
          userClassroom={userClassroom}
          members={classroomMembersSorted}
          teams={teams}
          classroomId={classroomId}
          userRole={userClassroom.role}
          showTeams={userClassroom.classroom.maxTeamSize > 1}
        />
      </div>
    </>
  );
}

function MemberTable({
  userClassroom,
  members,
  classroomId,
  userRole,
  showTeams,
  teams,
}: {
  userClassroom: UserClassroomResponse;
  members: UserClassroomResponse[];
  classroomId: string;
  userRole: Role;
  showTeams: boolean;
  teams: TeamResponse[];
}) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-full">Member</TableHead>
          {userClassroom.classroom.maxTeamSize > 1 && <TableHead>Team</TableHead>}
          <TableHead className="text-right">Role</TableHead>
        </TableRow>
      </TableHeader>
      {members.map((m) => (
        <TableRow>
          <TableCell className="w-full">
            <MemberListElement member={m} showTeams={showTeams} />
          </TableCell>
          {userClassroom.classroom.maxTeamSize > 1 && (
            <TableCell>
              <div className="flex justify-end">
                {isStudent(m) && (
                  <TeamDropdown team={m.team} memberID={m.user.id} classroomID={classroomId} teams={teams} />
                )}
              </div>
            </TableCell>
          )}
          <TableCell className="grid place-content-end">
            <div className="flex justify-end">
              {m.user.id !== userClassroom.user.id &&
                (userClassroom.classroom.ownerId === userClassroom.user.id ||
                  (userRole === Role.Owner && m.role !== Role.Owner)) && (
                  <RoleDropdown
                    role={m.role}
                    memberID={m.user.id}
                    classroomID={classroomId}
                    userClassroom={userClassroom}
                  />
                )}
            </div>
          </TableCell>
        </TableRow>
      ))}
    </Table>
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

function RoleDropdown({
  role,
  memberID,
  classroomID,
  userClassroom,
}: {
  role: Role;
  memberID: number;
  classroomID: string;
  userClassroom: UserClassroomResponse;
}) {
  const { mutateAsync, isError, isPending } = useUpdateMemberRole(classroomID, memberID);

  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    defaultValues: {
      role: getRole(role),
    },
  });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    await mutateAsync(values);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="role"
          render={({ field }) => (
            <FormItem>
              <Select
                disabled={isPending}
                onValueChange={(role: keyof typeof Role) => onSubmit({ role })}
                defaultValue={field.value}
              >
                <FormControl>
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="Change the role from the person" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value={getRole(Role.Student)}>{getRole(Role.Student)}</SelectItem>
                  <SelectItem value={getRole(Role.Moderator)}>{getRole(Role.Moderator)}</SelectItem>
                  {userClassroom.classroom.ownerId === userClassroom.user.id && (
                    <SelectItem value={getRole(Role.Owner)}>{getRole(Role.Owner)}</SelectItem>
                  )}
                </SelectContent>
              </Select>
            </FormItem>
          )}
        />

        {isError && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>The role could not be switched!</AlertDescription>
          </Alert>
        )}
      </form>
    </Form>
  );
}

export const updateTeamSchema = z.object({
  teamId: z.string().uuid(),
});

const REMOVE_TEAM = "remove-team";

function TeamDropdown({
  team,
  memberID,
  classroomID,
  teams,
}: {
  team?: Team;
  memberID: number;
  classroomID: string;
  teams: TeamResponse[];
}) {
  const {
    mutateAsync: updateTeam,
    error: updateTeamError,
    isPending: updateTeamIsPending,
  } = useUpdateMemberTeam(classroomID, memberID);
  const {
    mutateAsync: removeTeam,
    error: removeTeamError,
    isPending: removeTeamPending,
  } = useRemoveTeamMember(classroomID, memberID);

  const isPending = updateTeamIsPending || removeTeamPending;
  const error = updateTeamError || removeTeamError;

  const form = useForm<z.infer<typeof updateTeamSchema>>({
    resolver: zodResolver(updateTeamSchema),
    defaultValues: {
      teamId: team?.id ?? "",
    },
  });

  async function onSubmit(values: z.infer<typeof updateTeamSchema>) {
    await updateTeam(values.teamId);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="teamId"
          render={({ field }) => (
            <FormItem>
              <Select
                disabled={isPending}
                onValueChange={(teamId: string) =>
                  teamId === REMOVE_TEAM ? removeTeam(team?.id) : onSubmit({ teamId })
                }
                defaultValue={field.value}
              >
                <FormControl>
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="Select a team..." />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {teams.map((t) => (
                    <SelectItem key={t.id} value={t.id}>
                      {t.name}
                    </SelectItem>
                  ))}
                  {team && <SelectItem value={REMOVE_TEAM}>No team</SelectItem>}
                </SelectContent>
              </Select>
            </FormItem>
          )}
        />

        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{error.message}</AlertDescription>
          </Alert>
        )}
      </form>
    </Form>
  );
}
