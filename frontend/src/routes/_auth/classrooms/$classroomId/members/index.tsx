import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { getRole, Role } from "@/types/classroom.ts";
import { createFormSchema } from "@/types/member.ts";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Header } from "@/components/header.tsx";
import { membersQueryOptions, useUpdateMemberRole } from "@/api/member.ts";
import { ReportApiAxiosParamCreator, UserClassroomResponse } from "@/swagger-client";
import List from "@/components/ui/list.tsx";
import ListItem from "@/components/ui/listItem.tsx";
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card.tsx";
import { Avatar } from "@/components/avatar.tsx";
import { Separator } from "@/components/ui/separator.tsx";
import { teamsQueryOptions } from "@/api/team.ts";
import { classroomQueryOptions } from "@/api/classroom.ts";
import { assignmentsQueryOptions } from "@/api/assignment.ts";
import { Loader } from "@/components/loader.tsx";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue, } from "@/components/ui/select"
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form } from "@/components/ui/form";
import { AlertCircle } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";

export const Route = createFileRoute('/_auth/classrooms/$classroomId/members/')({
  component: Members,
  loader: async ({ context: { queryClient }, params }) => {
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    if (userClassroom.role === Role.Student && !userClassroom.team) {
      throw redirect({
        to: "/classrooms/$classroomId/teams/join",
        params,
      });
    }
    const { url: reportDownloadUrl } = await ReportApiAxiosParamCreator().getClassroomReport(params.classroomId)
    const members = await queryClient.ensureQueryData(membersQueryOptions(params.classroomId));
    if(userClassroom.role !== Role.Student) {
      const assignments = await queryClient.ensureQueryData(assignmentsQueryOptions(params.classroomId));
      return { userClassroom, assignments, members, teams, reportDownloadUrl };
    }else{
      return { userClassroom, members, teams, reportDownloadUrl };
    }
  },
  pendingComponent: Loader,
});

function Members() {
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));

  const classroomMembersSorted = classroomMembers.sort((a, b) => {
    if (a.role === Role.Owner) {
      return -1;
    } else if (b.role === Role.Owner) {
      return 1;
    } else if (a.role === Role.Moderator) {
      return -1;
    } else if (b.role === Role.Moderator) {
      return 1;
    } else {
      return 0;
    }
  });

  return (
    <div>
      <Header title="Members" />
      <div className="justify-between gap-10" >
        <MemberTable
          members={classroomMembersSorted}
          classroomId={classroomId}
          userRole={userClassroom.role}
          showTeams={userClassroom.classroom.maxTeamSize > 1}
        />
        <Outlet />
      </div>
    </div>
  )
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
              {userRole == Role.Owner && m.role != Role.Owner ? (
                <RoleDropdown role={m.role} memberID={m.user.id} classroomID={classroomId}/>
              ) : (
                ""
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

function RoleDropdown({ role, memberID, classroomID }: { role: Role, memberID: number, classroomID: string }) {
  const { mutateAsync, isError } = useUpdateMemberRole(classroomID, memberID);

  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    defaultValues: {
      role: role,
    },
  });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    await mutateAsync(values);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <Select onValueChange={onSubmit}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder={role === 2 ? "Student" : role === 1 ? "Moderator" : "Owner"} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="2">Student</SelectItem>
            <SelectItem value="1">Moderator</SelectItem>
          </SelectContent>
        </Select>

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
