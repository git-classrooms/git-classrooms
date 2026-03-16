import { zodResolver } from "@hookform/resolvers/zod";
import { createFileRoute, Link, redirect, useRouter } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { InviteForm, inviteFormSchema, Role, Status } from "@/types/classroom";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  AlertCircle,
  ArrowLeft,
  Check,
  Clock,
  Link as LinkIcon,
  Loader2,
  Mail,
  MailPlus,
  Send,
  UserPlus,
  X,
  XCircle,
} from "lucide-react";
import { Loader } from "@/components/loader.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { formatDate, formatRelativeTime, isStudent, cn } from "@/lib/utils.ts";
import { ClassroomInvitation, UserClassroomResponse } from "@/swagger-client";
import { classroomInvitationsQueryOptions, classroomQueryOptions, useInviteClassroomMembers } from "@/api/classroom";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { toast } from "sonner";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { StatusBadge } from "@/components/ui/status-badge";
import { useMemo } from "react";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useTranslation } from "react-i18next";
import { useRoleLabels, useInvitationStatusLabels } from "@/hooks/useClassroomLabels";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/invite")({
  loader: async ({ context: { queryClient }, params }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    if (isStudent(userClassroom)) {
      throw redirect({
        to: "/classrooms/$classroomId",
        search: { tab: "assignments" },
        params,
        replace: true,
      });
    }
  },
  pendingComponent: Loader,
  component: ClassroomInviteForm,
});

const statusConfigBase: Record<
  Status,
  { variant: "success" | "warning" | "destructive" | "neutral" | "info"; icon: typeof Check }
> = {
  [Status.Pending]: { variant: "warning", icon: Clock },
  [Status.Accepted]: { variant: "success", icon: Check },
  [Status.Rejected]: { variant: "destructive", icon: X },
  [Status.Revoked]: { variant: "neutral", icon: XCircle },
  [Status.Failed]: { variant: "destructive", icon: AlertCircle },
};

function ClassroomInviteForm() {
  const { t } = useTranslation("classroom");
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: invitations } = useSuspenseQuery(classroomInvitationsQueryOptions(classroomId));

  const stats = useMemo(() => {
    const pending = invitations.filter((i) => i.status === Status.Pending).length;
    const accepted = invitations.filter((i) => i.status === Status.Accepted).length;
    const rejected = invitations.filter((i) => i.status === Status.Rejected || i.status === Status.Revoked).length;
    return { pending, accepted, rejected, total: invitations.length };
  }, [invitations]);

  const sortedInvitations = useMemo(
    () =>
      [...invitations].sort((a, b) => {
        // Pending first, then by date
        if (a.status === Status.Pending && b.status !== Status.Pending) return -1;
        if (b.status === Status.Pending && a.status !== Status.Pending) return 1;
        return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
      }),
    [invitations]
  );

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
            <BreadcrumbPage>{t("invite.title")}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link to="/classrooms/$classroomId" search={{ tab: "members" }} params={{ classroomId }}>
              <ArrowLeft className="w-5 h-5" />
            </Link>
          </Button>
          <div>
            <h1 className="text-3xl font-bold tracking-tight">{t("invite.title")}</h1>
            <p className="text-muted-foreground mt-1">{t("invite.subtitle")}</p>
          </div>
        </div>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
        <StatCard label={t("invite.stats.totalSent")} value={stats.total} icon={Mail} color="text-primary" bgColor="bg-primary/15" />
        <StatCard label={t("invite.stats.pending")} value={stats.pending} icon={Clock} color="text-warning" bgColor="bg-warning/15" />
        <StatCard label={t("invite.stats.accepted")} value={stats.accepted} icon={Check} color="text-success" bgColor="bg-success/15" />
        <StatCard
          label={t("invite.stats.declined")}
          value={stats.rejected}
          icon={XCircle}
          color="text-muted-foreground"
          bgColor="bg-muted"
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
        {/* Invite Form Card */}
        <Card className="lg:col-span-2 h-fit">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-primary/15 flex items-center justify-center">
                <MailPlus className="w-5 h-5 text-primary" />
              </div>
              <div>
                <CardTitle className="text-lg">{t("invite.sendTitle")}</CardTitle>
                <CardDescription>{t("invite.sendSubtitle")}</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <InviteFormSection classroomId={classroomId} />
          </CardContent>
        </Card>

        {/* Invitations List Card */}
        <Card className="lg:col-span-3">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-info/15 flex items-center justify-center">
                  <UserPlus className="w-5 h-5 text-info" />
                </div>
                <div>
                  <CardTitle className="text-lg">{t("invite.sentTitle")}</CardTitle>
                  <CardDescription>{t("invite.sentSubtitle", { count: invitations.length })}</CardDescription>
                </div>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {sortedInvitations.length === 0 ? (
              <EmptyInvitationsState />
            ) : (
              <div className="space-y-2">
                {sortedInvitations.map((invitation) => (
                  <InvitationRow
                    key={invitation.id}
                    invitation={invitation}
                    userClassroom={userClassroom}
                  />
                ))}
              </div>
            )}
          </CardContent>
        </Card>
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
  icon: typeof Mail;
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

function InviteFormSection({ classroomId }: { classroomId: string }) {
  const { t } = useTranslation("classroom");
  const { t: tc } = useTranslation("common");
  const roleLabels = useRoleLabels();
  const { mutateAsync, isError, isPending } = useInviteClassroomMembers(classroomId);

  const form = useForm<z.infer<typeof inviteFormSchema>>({
    resolver: zodResolver(inviteFormSchema),
    mode: "onBlur",
    reValidateMode: "onChange",
    defaultValues: {
      memberEmails: "",
      role: Role.Student,
    },
  });

  async function onSubmit(values: InviteForm) {
    try {
      await mutateAsync(values);
      form.reset();
      toast.success(t("invite.success"));
    } catch {
      // Error handled by isError state
    }
  }

  const emailCount = form.watch("memberEmails")?.split("\n").filter(Boolean).length || 0;

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="memberEmails"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("invite.emailLabel")}</FormLabel>
              <FormControl>
                <Textarea
                  placeholder={`${t("invite.emailPlaceholder")}\nstudent2@example.com\nstudent3@example.com`}
                  className="resize-none min-h-[160px] font-mono text-sm"
                  {...field}
                />
              </FormControl>
              <FormDescription className="flex items-center justify-between">
                <span>{t("invite.emailDescription")}</span>
                {emailCount > 0 && (
                  <span className="text-xs text-muted-foreground">
                    {t("invite.emailCount", { count: emailCount })}
                  </span>
                )}
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="role"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("invite.roleLabel")}</FormLabel>
              <FormControl>
                <Select
                  value={String(field.value)}
                  onValueChange={(val) => field.onChange(Number(val))}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value={String(Role.Owner)}>{roleLabels[Role.Owner]}</SelectItem>
                    <SelectItem value={String(Role.Moderator)}>{roleLabels[Role.Moderator]}</SelectItem>
                    <SelectItem value={String(Role.Student)}>{roleLabels[Role.Student]}</SelectItem>
                  </SelectContent>
                </Select>
              </FormControl>
              <FormDescription>{t("invite.roleDescription")}</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit" disabled={isPending} className="w-full">
          {isPending ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              {t("invite.sending")}
            </>
          ) : (
            <>
              <Send className="mr-2 h-4 w-4" />
              {t("invite.sendButton")}
            </>
          )}
        </Button>

        {isError && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>{tc("status.error")}</AlertTitle>
            <AlertDescription>{t("invite.error")}</AlertDescription>
          </Alert>
        )}
      </form>
    </Form>
  );
}

function EmptyInvitationsState() {
  const { t } = useTranslation("classroom");

  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center mb-4">
        <Mail className="w-6 h-6 text-muted-foreground" />
      </div>
      <h3 className="font-medium text-foreground mb-1">{t("invite.empty.title")}</h3>
      <p className="text-sm text-muted-foreground max-w-[240px]">
        {t("invite.empty.description")}
      </p>
    </div>
  );
}

function InvitationRow({
  invitation,
  userClassroom,
}: {
  invitation: ClassroomInvitation;
  userClassroom: UserClassroomResponse;
}) {
  const { t } = useTranslation("classroom");
  const router = useRouter();
  const config = statusConfigBase[invitation.status as Status];
  const StatusIcon = config.icon;
  const statusLabels = useInvitationStatusLabels();
  const statusLabel = statusLabels[invitation.status as Status];
  const roleLabels = useRoleLabels();
  const roleLabel = roleLabels[invitation.role as Role];
  const isPending = invitation.status === Status.Pending;
  const canCopyLink = invitation.status !== Status.Accepted && invitation.status !== Status.Revoked;

  const copyLink = () => {
    const path = router.buildLocation({
      to: "/classrooms/$classroomId/invitations/$invitationId",
      params: { classroomId: userClassroom.classroom.id, invitationId: invitation.id },
      search: { groupLink: false },
    });
    navigator.clipboard.writeText(`${location.origin}${path.href}`);
    toast.success(t("invite.linkCopied"));
  };

  return (
    <div
      className={cn(
        "group flex items-center gap-4 p-3 rounded-lg border transition-all duration-200",
        isPending
          ? "border-warning/30 bg-warning/5 hover:border-warning/50"
          : "border-border/50 hover:border-border"
      )}
    >
      {/* Status Icon */}
      <div
        className={cn(
          "w-8 h-8 rounded-full flex items-center justify-center shrink-0",
          config.variant === "success" && "bg-success/15 text-success",
          config.variant === "warning" && "bg-warning/15 text-warning",
          config.variant === "destructive" && "bg-destructive/15 text-destructive",
          config.variant === "neutral" && "bg-muted text-muted-foreground"
        )}
      >
        <StatusIcon className="w-4 h-4" />
      </div>

      {/* Email & Info */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="font-medium truncate">{invitation.email}</span>
          {invitation.status !== Status.Accepted && (
            <StatusBadge variant="info" size="sm">
              {roleLabel}
            </StatusBadge>
          )}
          <StatusBadge variant={config.variant} size="sm">
            {statusLabel}
          </StatusBadge>
        </div>
        <p className="text-xs text-muted-foreground mt-0.5">
          {t("invite.sent")} {formatRelativeTime(invitation.createdAt)} · {formatDate(invitation.createdAt)}
        </p>
      </div>

      {/* Actions */}
      {canCopyLink && (
        <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8" onClick={copyLink}>
                <LinkIcon className="w-4 h-4" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>{t("invite.copyLinkTooltip")}</TooltipContent>
          </Tooltip>
        </div>
      )}
    </div>
  );
}
