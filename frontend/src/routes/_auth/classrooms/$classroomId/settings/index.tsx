import { classroomQueryOptions, useArchiveClassroom, useCreateTeachingGroup } from "@/api/classroom";
import { ClassroomEditForm } from "@/components/classroomsForm";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { AlertTriangle, Archive, Eye, EyeOff, Info, Library, Loader2, Lock, Users, Users2 } from "lucide-react";
import { StatusBadge } from "@/components/ui/status-badge";
import { Card, CardContent } from "@/components/ui/card";
import { cn, isOwner, unwrapApiError } from "@/lib/utils";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { useTranslation } from "react-i18next";
import { useState } from "react";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/settings/")({
  loader: async ({ params: { classroomId }, context: { queryClient } }) => {
    const userClassroom = await queryClient.fetchQuery(classroomQueryOptions(classroomId));
    return { userClassroom };
  },
  component: Index,
});

function Index() {
  const { t } = useTranslation("classroom");
  const { t: tc } = useTranslation("common");
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const classroom = userClassroom.classroom;
  const { mutate: archiveClassroom } = useArchiveClassroom(classroomId);

  const teamsEnabled = classroom.maxTeamSize > 1;

  return (
    <div className="space-y-8">
      {/* Fixed Configuration Section */}
      <section>
        <div className="flex items-center gap-2 mb-4">
          <Lock className="w-4 h-4 text-muted-foreground" />
          <h2 className="text-lg font-semibold font-mono">{t("settings.fixedConfig")}</h2>
          <Tooltip>
            <TooltipTrigger>
              <Info className="w-3.5 h-3.5 text-muted-foreground" />
            </TooltipTrigger>
            <TooltipContent>
              <p>{t("settings.fixedConfigTooltip")}</p>
            </TooltipContent>
          </Tooltip>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {/* Creator Card */}
          <ConfigCard
            icon={<Users className="w-4 h-4" />}
            label={t("settings.creator")}
            value={classroom.owner.name}
            description={t("settings.creatorDescription")}
          />

          {/* Teams Configuration Card */}
          <ConfigCard
            icon={<Users2 className="w-4 h-4" />}
            label={t("settings.teams")}
            value={
              teamsEnabled ? (
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <StatusBadge variant="success" size="sm">{t("settings.teamsEnabled")}</StatusBadge>
                    <span className="text-sm text-muted-foreground">
                      {t("settings.teamsMax", { count: classroom.maxTeamSize })}
                    </span>
                  </div>
                  {classroom.maxTeams > 0 && (
                    <p className="text-xs text-muted-foreground">
                      {t("settings.teamsLimited", { count: classroom.maxTeams })}
                    </p>
                  )}
                </div>
              ) : (
                <StatusBadge variant="neutral" size="sm">{t("settings.teamsDisabled")}</StatusBadge>
              )
            }
            description={t("settings.teamsDescription")}
          />

          {/* Student Team Creation - only show if teams enabled */}
          {teamsEnabled && (
            <ConfigCard
              icon={<Users2 className="w-4 h-4" />}
              label={t("settings.studentTeamCreation")}
              value={
                <StatusBadge variant={classroom.createTeams ? "success" : "neutral"} size="sm">
                  {classroom.createTeams ? t("settings.allowed") : t("settings.notAllowed")}
                </StatusBadge>
              }
              description={t("settings.studentTeamCreationDescription")}
            />
          )}

          {/* Mutual Code View Card */}
          <Tooltip>
            <TooltipTrigger asChild>
              <div>
                <ConfigCard
                  icon={classroom.studentsViewAllProjects ? <Eye className="w-4 h-4" /> : <EyeOff className="w-4 h-4" />}
                  label={t("settings.mutualCodeView")}
                  value={
                    <StatusBadge
                      variant={classroom.studentsViewAllProjects ? "info" : "neutral"}
                      size="sm"
                    >
                      {classroom.studentsViewAllProjects ? t("settings.teamsEnabled") : t("settings.teamsDisabled")}
                    </StatusBadge>
                  }
                  description={t("settings.mutualCodeViewDescription")}
                  interactive
                />
              </div>
            </TooltipTrigger>
            <TooltipContent side="bottom">
              <p>
                {classroom.studentsViewAllProjects
                  ? t("settings.mutualCodeViewEnabled")
                  : t("settings.mutualCodeViewDisabled")}
              </p>
            </TooltipContent>
          </Tooltip>

          {/* Teaching Material Card */}
          <ConfigCard
            icon={<Library className="w-4 h-4" />}
            label={t("settings.teachingMaterial")}
            value={
              <StatusBadge
                variant={!!classroom.teachingGroupId ? "success" : "neutral"}
                size="sm"
              >
                {!!classroom.teachingGroupId ? t("settings.teamsEnabled") : t("settings.teamsDisabled")}
              </StatusBadge>
            }
            description={t("settings.teachingMaterialDescription")}
          />
        </div>
      </section>

      {/* Activate Teaching Material Section */}
      {!classroom.teachingGroupId && isOwner(userClassroom) && (
        <>
          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-primary/30" />
            </div>
            <div className="relative flex justify-center">
              <span className="bg-card/30 px-3 text-xs text-primary uppercase tracking-wider">
                {t("settings.teachingMaterialActivateSection")}
              </span>
            </div>
          </div>

          <section>
            <ActivateTeachingMaterialCard classroomId={classroomId} />
          </section>
        </>
      )}

      {/* Divider */}
      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-border/50" />
        </div>
        <div className="relative flex justify-center">
          <span className="bg-card/30 px-3 text-xs text-muted-foreground uppercase tracking-wider">
            {t("settings.editableSettings")}
          </span>
        </div>
      </div>

      {/* Editable Settings Section */}
      <section>
        <ClassroomEditForm userClassroom={userClassroom} />
      </section>

      {/* Danger Zone */}
      {!classroom.archived && (
        <>
          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-destructive/30" />
            </div>
            <div className="relative flex justify-center">
              <span className="bg-card/30 px-3 text-xs text-destructive uppercase tracking-wider">
                {t("settings.danger")}
              </span>
            </div>
          </div>

          <section>
            <Card className="border-destructive/30">
              <CardContent className="p-4">
                <div className="flex items-start justify-between gap-4">
                  <div className="flex items-start gap-3">
                    <div className="w-10 h-10 rounded-lg bg-destructive/10 flex items-center justify-center shrink-0">
                      <AlertTriangle className="w-5 h-5 text-destructive" />
                    </div>
                    <div>
                      <h3 className="font-semibold">{t("settings.archive.title")}</h3>
                      <p className="text-sm text-muted-foreground mt-1">
                        {t("settings.archive.description")}
                      </p>
                    </div>
                  </div>
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <Button variant="destructive" size="sm">
                        <Archive className="w-4 h-4 mr-2" />
                        {t("settings.archive.button")}
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>{t("settings.archive.confirm")}</AlertDialogTitle>
                        <AlertDialogDescription>
                          {t("settings.archive.confirmDescription")}
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>{tc("actions.cancel")}</AlertDialogCancel>
                        <AlertDialogAction
                          onClick={() =>
                            archiveClassroom(undefined, {
                              onSuccess: () => toast.success(t("settings.archive.success")),
                              onError: (error) => toast.error(unwrapApiError(error)?.message ?? t("settings.archive.error")),
                            })
                          }
                          variant="destructive"
                        >
                          {t("settings.archive.title")}
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </div>
              </CardContent>
            </Card>
          </section>
        </>
      )}
    </div>
  );
}

function ConfigCard({
  icon,
  label,
  value,
  description,
  interactive = false,
}: {
  icon: React.ReactNode;
  label: string;
  value: React.ReactNode;
  description: string;
  interactive?: boolean;
}) {
  return (
    <Card
      className={cn(
        "transition-all duration-200",
        interactive && "cursor-help hover:border-primary/30"
      )}
    >
      <CardContent className="p-4">
        <div className="flex items-start gap-3">
          <div className="w-8 h-8 rounded-lg bg-muted/50 flex items-center justify-center text-muted-foreground shrink-0">
            {icon}
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center justify-between gap-2 mb-1">
              <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                {label}
              </span>
            </div>
            <div className="font-medium">{value}</div>
            <p className="text-xs text-muted-foreground mt-1">{description}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

function ActivateTeachingMaterialCard({ classroomId }: { classroomId: string }) {
  const { t } = useTranslation("classroom");
  const { t: tc } = useTranslation("common");
  const [name, setName] = useState("Teaching Material");
  const { mutate, isPending } = useCreateTeachingGroup(classroomId);

  return (
    <Card className="border-primary/30">
      <CardContent className="p-4">
        <div className="flex items-start justify-between gap-4">
          <div className="flex items-start gap-3">
            <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center shrink-0">
              <Library className="w-5 h-5 text-primary" />
            </div>
            <div>
              <h3 className="font-semibold">{t("settings.teachingMaterial")}</h3>
              <p className="text-sm text-muted-foreground mt-1">
                {t("settings.teachingMaterialActivateDescription")}
              </p>
            </div>
          </div>
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button size="sm" disabled={isPending}>
                {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Library className="w-4 h-4 mr-2" />}
                {t("settings.teachingMaterialActivate")}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>{t("settings.teachingMaterialActivateConfirm")}</AlertDialogTitle>
                <AlertDialogDescription>
                  {t("settings.teachingMaterialActivateConfirmDescription")}
                </AlertDialogDescription>
              </AlertDialogHeader>
              <div className="py-2">
                <label className="text-sm font-medium mb-1.5 block">
                  {t("teachingMaterial.groupName")}
                </label>
                <Input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder={t("teachingMaterial.groupNamePlaceholder")}
                />
                <p className="text-xs text-muted-foreground mt-1.5">
                  {t("teachingMaterial.groupNameDescription")}
                </p>
              </div>
              <AlertDialogFooter>
                <AlertDialogCancel>{tc("actions.cancel")}</AlertDialogCancel>
                <AlertDialogAction
                  disabled={!name.trim()}
                  onClick={() =>
                    mutate(
                      { name: name.trim() },
                      {
                        onSuccess: () => toast.success(t("settings.teachingMaterialActivateSuccess")),
                        onError: (error) => toast.error(unwrapApiError(error)?.message ?? t("settings.teachingMaterialActivateError")),
                      },
                    )
                  }
                >
                  {t("settings.teachingMaterialActivate")}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </CardContent>
    </Card>
  );
}
