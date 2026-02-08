import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { getRole, Role } from "@/types/classroom.ts";
import { createFormSchema } from "@/types/member.ts";
import { useSuspenseQuery } from "@tanstack/react-query";
import { membersQueryOptions, useRemoveTeamMember, useUpdateMemberRole, useUpdateMemberTeam } from "@/api/member.ts";
import { TeamResponse, UserClassroomResponse } from "@/swagger-client";
import { Avatar } from "@/components/avatar.tsx";
import { teamsQueryOptions } from "@/api/team.ts";
import { classroomQueryOptions } from "@/api/classroom.ts";
import { Loader } from "@/components/loader.tsx";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form";
import {
  AlertCircle,
  ArrowLeft,
  Crown,
  ExternalLink,
  Loader2,
  Mail,
  Search,
  Shield,
  UserPlus,
  Users,
  GraduationCap,
} from "lucide-react";
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
import { isCreator, isStudent, cn } from "@/lib/utils";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { StatusBadge } from "@/components/ui/status-badge";
import { useTranslation } from "react-i18next";

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
    const members = await queryClient.ensureQueryData(membersQueryOptions(params.classroomId));
    return { userClassroom, members, teams };
  },
  pendingComponent: Loader,
});

const roleConfig: Record<Role, { icon: typeof Crown; color: string; bgColor: string }> = {
  [Role.Owner]: {
    icon: Crown,
    color: "text-warning",
    bgColor: "bg-warning/15",
  },
  [Role.Moderator]: {
    icon: Shield,
    color: "text-info",
    bgColor: "bg-info/15",
  },
  [Role.Student]: {
    icon: GraduationCap,
    color: "text-muted-foreground",
    bgColor: "bg-muted",
  },
};

function Members() {
  const { t } = useTranslation("classroom");
  const { t: tCommon } = useTranslation("common");
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));
  const [searchQuery, setSearchQuery] = useState("");

  const stats = useMemo(() => {
    const owners = classroomMembers.filter((m) => m.role === Role.Owner).length;
    const moderators = classroomMembers.filter((m) => m.role === Role.Moderator).length;
    const students = classroomMembers.filter((m) => m.role === Role.Student).length;
    const withTeam = classroomMembers.filter((m) => m.team).length;
    return { owners, moderators, students, withTeam, total: classroomMembers.length };
  }, [classroomMembers]);

  const filteredMembers = useMemo(() => {
    const query = searchQuery.toLowerCase();
    return [...classroomMembers]
      .filter((m) => {
        if (!query) return true;
        return (
          m.user.name.toLowerCase().includes(query) ||
          m.user.gitlabUsername?.toLowerCase().includes(query) ||
          m.team?.name.toLowerCase().includes(query)
        );
      })
      .sort((a, b) => {
        if (a.role !== b.role) return a.role - b.role;
        if (isCreator(a)) return -1;
        if (isCreator(b)) return 1;
        return a.user.name.localeCompare(b.user.name);
      });
  }, [classroomMembers, searchQuery]);

  const showTeams = userClassroom.classroom.maxTeamSize > 1;

  return (
    <div className="space-y-8">
      {/* Breadcrumb */}
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms">{t("title")}</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms/$classroomId" search={{ tab: "members" }} params={{ classroomId }}>
                {userClassroom.classroom.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{t("members.manage.title")}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {/* Header */}
      <div className="flex flex-col gap-6">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" asChild>
              <Link to="/classrooms/$classroomId" search={{ tab: "members" }} params={{ classroomId }}>
                <ArrowLeft className="w-5 h-5" />
              </Link>
            </Button>
            <div>
              <h1 className="text-3xl font-bold tracking-tight">{t("members.manage.title")}</h1>
              <p className="text-muted-foreground mt-1">
                {t("members.manage.subtitle")}
              </p>
            </div>
          </div>
          <Button variant="glow" asChild>
            <Link to="/classrooms/$classroomId/invite" params={{ classroomId }}>
              <UserPlus className="w-4 h-4 mr-2" />
              {t("members.invite")}
            </Link>
          </Button>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <StatCard
            label={tCommon("total")}
            value={stats.total}
            icon={Users}
            color="text-primary"
            bgColor="bg-primary/15"
          />
          <StatCard
            label={t("members.role.owners")}
            value={stats.owners}
            icon={Crown}
            color="text-warning"
            bgColor="bg-warning/15"
          />
          <StatCard
            label={t("members.role.moderators")}
            value={stats.moderators}
            icon={Shield}
            color="text-info"
            bgColor="bg-info/15"
          />
          <StatCard
            label={t("members.role.students")}
            value={stats.students}
            icon={GraduationCap}
            color="text-muted-foreground"
            bgColor="bg-muted"
          />
        </div>
      </div>

      {/* Search & Filter Bar */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="relative flex-1 max-w-md">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            placeholder={t("members.manage.searchPlaceholder")}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
        <div className="text-sm text-muted-foreground flex items-center">
          {t("members.manage.ofMembers", { filtered: filteredMembers.length, total: classroomMembers.length })}
        </div>
      </div>

      {/* Members List */}
      <div className="space-y-3">
        {filteredMembers.length === 0 ? (
          <EmptySearchState query={searchQuery} />
        ) : (
          filteredMembers.map((member) => (
            <MemberRow
              key={member.user.id}
              member={member}
              userClassroom={userClassroom}
              classroomId={classroomId}
              teams={teams}
              showTeams={showTeams}
            />
          ))
        )}
      </div>
    </div>
  );
}

function StatCard({
  label,
  value,
  icon: Icon,
  color,
  bgColor,
}: {
  label: string;
  value: number;
  icon: typeof Users;
  color: string;
  bgColor: string;
}) {
  return (
    <Card className="border-border/50">
      <CardContent className="p-4">
        <div className="flex items-center gap-3">
          <div className={cn("w-10 h-10 rounded-lg flex items-center justify-center", bgColor)}>
            <Icon className={cn("w-5 h-5", color)} />
          </div>
          <div>
            <p className="text-2xl font-bold font-mono">{value}</p>
            <p className="text-xs text-muted-foreground uppercase tracking-wide">{label}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

function EmptySearchState({ query }: { query: string }) {
  const { t } = useTranslation("classroom");
  return (
    <div className="border border-dashed border-border rounded-lg p-12 text-center">
      <Search className="w-10 h-10 text-muted-foreground mx-auto mb-3" />
      <h3 className="font-medium text-foreground mb-1">{t("members.manage.noMembersFound")}</h3>
      <p className="text-sm text-muted-foreground">
        {t("members.manage.noMembersMatch", { query })}
      </p>
    </div>
  );
}

function MemberRow({
  member,
  userClassroom,
  classroomId,
  teams,
  showTeams,
}: {
  member: UserClassroomResponse;
  userClassroom: UserClassroomResponse;
  classroomId: string;
  teams: TeamResponse[];
  showTeams: boolean;
}) {
  const { t } = useTranslation("classroom");
  const config = roleConfig[member.role as Role];
  const Icon = config.icon;
  const isCurrentUser = member.user.id === userClassroom.user.id;
  const isClassroomCreator = isCreator(member);

  const canEditRole =
    !isCurrentUser &&
    (userClassroom.classroom.ownerId === userClassroom.user.id ||
      (userClassroom.role === Role.Owner && member.role !== Role.Owner));

  const canEditTeam = isStudent(member);

  return (
    <Card className="group transition-all duration-200 hover:border-primary/30">
      <CardContent className="p-4">
        <div className="flex flex-col sm:flex-row sm:items-center gap-4">
          {/* Member Info */}
          <div className="flex items-center gap-3 flex-1 min-w-0">
            <div className="relative">
              <Avatar
                avatarUrl={member.user.avatarURL}
                fallbackUrl={member.user.fallbackAvatarURL}
                name={member.user.name!}
                className="w-12 h-12"
              />
              <div className="absolute -bottom-1 -right-1 w-5 h-5 rounded-full bg-card border-2 border-card">
                <div
                  className={cn(
                    "w-full h-full rounded-full flex items-center justify-center",
                    config.bgColor
                  )}
                >
                  <Icon className={cn("w-2.5 h-2.5", config.color)} />
                </div>
              </div>
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <span className="font-semibold truncate">{member.user.name}</span>
                {isCurrentUser && (
                  <StatusBadge variant="info" size="sm">{t("members.manage.you")}</StatusBadge>
                )}
                {isClassroomCreator && (
                  <StatusBadge variant="warning" size="sm">{t("members.manage.creator")}</StatusBadge>
                )}
              </div>
              <p className="text-sm text-muted-foreground truncate">
                @{member.user.gitlabUsername}
              </p>
            </div>
          </div>

          {/* Team Assignment */}
          {showTeams && (
            <div className="sm:w-48">
              {canEditTeam ? (
                <TeamDropdown
                  team={member.team}
                  memberID={member.user.id}
                  classroomID={classroomId}
                  teams={teams}
                />
              ) : (
                <div className="text-sm text-muted-foreground">
                  {member.team?.name || "—"}
                </div>
              )}
            </div>
          )}

          {/* Role Assignment */}
          <div className="sm:w-44">
            {canEditRole ? (
              <RoleDropdown
                role={member.role}
                memberID={member.user.id}
                classroomID={classroomId}
                userClassroom={userClassroom}
              />
            ) : (
              <StatusBadge
                variant={
                  member.role === Role.Owner
                    ? "warning"
                    : member.role === Role.Moderator
                      ? "info"
                      : "neutral"
                }
              >
                {getRole(member.role)}
              </StatusBadge>
            )}
          </div>

          {/* Quick Actions */}
          <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <Button variant="ghost" size="icon" className="h-8 w-8" asChild>
              <a href={member.webUrl} target="_blank" rel="noopener noreferrer">
                <ExternalLink className="w-4 h-4" />
              </a>
            </Button>
            {member.user.gitlabEmail && (
              <Button variant="ghost" size="icon" className="h-8 w-8" asChild>
                <a href={`mailto:${member.user.gitlabEmail}`}>
                  <Mail className="w-4 h-4" />
                </a>
              </Button>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
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
  const { t } = useTranslation("classroom");
  const { t: tCommon } = useTranslation("common");
  const { mutateAsync, isError, isPending } = useUpdateMemberRole(classroomID, memberID);

  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    mode: "onBlur",
    reValidateMode: "onChange",
    defaultValues: {
      role: getRole(role),
    },
  });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    await mutateAsync(values);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
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
                  <SelectTrigger className="w-full h-9">
                    <div className="flex items-center gap-2">
                      {isPending && <Loader2 className="w-3 h-3 animate-spin" />}
                      <SelectValue placeholder={t("members.manage.selectRole")} />
                    </div>
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value={getRole(Role.Student)}>
                    <div className="flex items-center gap-2">
                      <GraduationCap className="w-4 h-4" />
                      {t("members.role.student")}
                    </div>
                  </SelectItem>
                  <SelectItem value={getRole(Role.Moderator)}>
                    <div className="flex items-center gap-2">
                      <Shield className="w-4 h-4" />
                      {t("members.role.moderator")}
                    </div>
                  </SelectItem>
                  {userClassroom.classroom.ownerId === userClassroom.user.id && (
                    <SelectItem value={getRole(Role.Owner)}>
                      <div className="flex items-center gap-2">
                        <Crown className="w-4 h-4" />
                        {t("members.role.owner")}
                      </div>
                    </SelectItem>
                  )}
                </SelectContent>
              </Select>
              {isError && (
                <Alert variant="destructive" className="mt-2">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>{tCommon("status.error")}</AlertTitle>
                  <AlertDescription>{t("members.manage.failedToUpdateRole")}</AlertDescription>
                </Alert>
              )}
            </FormItem>
          )}
        />
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
  team?: { id: string; name: string };
  memberID: number;
  classroomID: string;
  teams: TeamResponse[];
}) {
  const { t } = useTranslation("classroom");
  const { t: tCommon } = useTranslation("common");
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
    mode: "onBlur",
    reValidateMode: "onChange",
    defaultValues: {
      teamId: team?.id ?? "",
    },
  });

  async function onSubmit(values: z.infer<typeof updateTeamSchema>) {
    await updateTeam(values.teamId);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
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
                  <SelectTrigger className="w-full h-9">
                    <div className="flex items-center gap-2">
                      {isPending && <Loader2 className="w-3 h-3 animate-spin" />}
                      <SelectValue placeholder={t("members.manage.selectTeam")} />
                    </div>
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {teams.map((teamItem) => (
                    <SelectItem key={teamItem.id} value={teamItem.id}>
                      {teamItem.name}
                    </SelectItem>
                  ))}
                  {team && (
                    <SelectItem value={REMOVE_TEAM} className="text-destructive">
                      {t("members.manage.removeFromTeam")}
                    </SelectItem>
                  )}
                </SelectContent>
              </Select>
              {error && (
                <Alert variant="destructive" className="mt-2">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>{tCommon("status.error")}</AlertTitle>
                  <AlertDescription>{t("members.manage.failedToUpdateTeam")}</AlertDescription>
                </Alert>
              )}
            </FormItem>
          )}
        />
      </form>
    </Form>
  );
}
