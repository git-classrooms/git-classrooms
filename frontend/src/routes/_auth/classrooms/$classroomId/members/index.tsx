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
import { useMemo, useState } from "react";
import { isCreator, isOwner, isStudent } from "@/lib/utils";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import {
  ColumnFiltersState,
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  SortingState,
  useReactTable,
} from "@tanstack/react-table";
import { Input } from "@/components/ui/input";

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

    return { userClassroom, members, teams, reportDownloadUrl };
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

const createMemberColumns = (user: UserClassroomResponse, teams: TeamResponse[]) => {
  const memberColumnHelper = createColumnHelper<UserClassroomResponse>();
  const teamsEnabled = user.classroom.maxTeamSize > 1;

  return [
    memberColumnHelper.accessor((row) => row.user.name, {
      id: "user",
      header: "User",
      cell: ({ row: { original: member } }) => <MemberListElement member={member} showTeams={teamsEnabled} />,
      enableSorting: true,
      enableColumnFilter: true,
    }),

    teamsEnabled
      ? memberColumnHelper.accessor((row) => row.team?.name, {
          id: "team",
          header: "Team",
          cell: ({ row: { original: member } }) =>
            isStudent(member) && (
              <TeamDropdown
                team={member.team}
                memberID={member.user.id}
                classroomID={member.classroom.id}
                teams={teams}
              />
            ),
          enableSorting: true,
          enableColumnFilter: true,
        })
      : undefined!,

    isOwner(user)
      ? memberColumnHelper.display({
          id: "role",
          header: () => <div className="text-right">Role</div>,
          cell: ({ row: { original: member } }) =>
            member.user.id !== user.user.id &&
            (user.classroom.ownerId === user.user.id || (user.role === Role.Owner && member.role !== Role.Owner)) && (
              <RoleDropdown
                role={member.role}
                memberID={member.user.id}
                classroomID={user.classroom.id}
                userClassroom={user}
              />
            ),
        })
      : undefined!,
  ].filter(Boolean);
};

function MemberTable({
  userClassroom,
  members,
  teams,
}: {
  userClassroom: UserClassroomResponse;
  members: UserClassroomResponse[];
  classroomId: string;
  userRole: Role;
  showTeams: boolean;
  teams: TeamResponse[];
}) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);

  const memberColumns = useMemo(() => createMemberColumns(userClassroom, teams), [userClassroom, teams]);

  const table = useReactTable({
    data: members,
    columns: memberColumns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    state: {
      sorting,
      columnFilters,
    },
  });

  return (
    <>
      <Input
        placeholder="Filter users..."
        value={(table.getColumn("user")?.getFilterValue() as string) ?? ""}
        onChange={(event) => table.getColumn("user")?.setFilterValue(event.target.value)}
        className="max-w-sm"
      />
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header, i) => {
                return (
                  <TableHead className={i === 0 ? "w-full" : undefined} key={header.id}>
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                  </TableHead>
                );
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id} data-state={row.getIsSelected() && "selected"}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={memberColumns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </>
  );
}

function MemberListElement({ member, showTeams }: { member: UserClassroomResponse; showTeams: boolean }) {
  return (
    <HoverCard>
      <HoverCardTrigger className="cursor-default flex">
        <div className="pr-2">
          <Avatar
            avatarUrl={member.user.avatarURL}
            fallbackUrl={member.user.fallbackAvatarURL}
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
