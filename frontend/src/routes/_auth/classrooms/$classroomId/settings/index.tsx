import { classroomQueryOptions, useArchiveClassroom } from "@/api/classroom";
import { ClassroomEditForm } from "@/components/classroomsForm";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { AlertTriangle, Archive, Eye, EyeOff, Info, Lock, Users, Users2 } from "lucide-react";
import { StatusBadge } from "@/components/ui/status-badge";
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
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

export const Route = createFileRoute("/_auth/classrooms/$classroomId/settings/")({
  loader: async ({ params: { classroomId }, context: { queryClient } }) => {
    const userClassroom = await queryClient.fetchQuery(classroomQueryOptions(classroomId));
    return { userClassroom };
  },
  component: Index,
});

function Index() {
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
          <h2 className="text-lg font-semibold font-mono">Fixed Configuration</h2>
          <Tooltip>
            <TooltipTrigger>
              <Info className="w-3.5 h-3.5 text-muted-foreground" />
            </TooltipTrigger>
            <TooltipContent>
              <p>These settings were set during classroom creation and cannot be modified.</p>
            </TooltipContent>
          </Tooltip>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {/* Creator Card */}
          <ConfigCard
            icon={<Users className="w-4 h-4" />}
            label="Creator"
            value={classroom.owner.name}
            description="Classroom owner"
          />

          {/* Teams Configuration Card */}
          <ConfigCard
            icon={<Users2 className="w-4 h-4" />}
            label="Teams"
            value={
              teamsEnabled ? (
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <StatusBadge variant="success" size="sm">Enabled</StatusBadge>
                    <span className="text-sm text-muted-foreground">
                      max {classroom.maxTeamSize} members
                    </span>
                  </div>
                  {classroom.maxTeams > 0 && (
                    <p className="text-xs text-muted-foreground">
                      Limited to {classroom.maxTeams} teams
                    </p>
                  )}
                </div>
              ) : (
                <StatusBadge variant="neutral" size="sm">Disabled</StatusBadge>
              )
            }
            description="Team structure"
          />

          {/* Student Team Creation - only show if teams enabled */}
          {teamsEnabled && (
            <ConfigCard
              icon={<Users2 className="w-4 h-4" />}
              label="Student Team Creation"
              value={
                <StatusBadge variant={classroom.createTeams ? "success" : "neutral"} size="sm">
                  {classroom.createTeams ? "Allowed" : "Not Allowed"}
                </StatusBadge>
              }
              description="Can students create their own teams"
            />
          )}

          {/* Mutual Code View Card */}
          <Tooltip>
            <TooltipTrigger asChild>
              <div>
                <ConfigCard
                  icon={classroom.studentsViewAllProjects ? <Eye className="w-4 h-4" /> : <EyeOff className="w-4 h-4" />}
                  label="Mutual Code View"
                  value={
                    <StatusBadge
                      variant={classroom.studentsViewAllProjects ? "info" : "neutral"}
                      size="sm"
                    >
                      {classroom.studentsViewAllProjects ? "Enabled" : "Disabled"}
                    </StatusBadge>
                  }
                  description="Cross-team code visibility"
                  interactive
                />
              </div>
            </TooltipTrigger>
            <TooltipContent side="bottom">
              <p>
                {classroom.studentsViewAllProjects
                  ? "Students can view code from other teams"
                  : "Students can only view their own team's code"}
              </p>
            </TooltipContent>
          </Tooltip>
        </div>
      </section>

      {/* Divider */}
      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-border/50" />
        </div>
        <div className="relative flex justify-center">
          <span className="bg-card/30 px-3 text-xs text-muted-foreground uppercase tracking-wider">
            Editable Settings
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
                Danger Zone
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
                      <h3 className="font-semibold">Archive Classroom</h3>
                      <p className="text-sm text-muted-foreground mt-1">
                        Once archived, no new assignments can be created and students cannot join.
                        This action cannot be undone.
                      </p>
                    </div>
                  </div>
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <Button variant="destructive" size="sm">
                        <Archive className="w-4 h-4 mr-2" />
                        Archive
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Archive this classroom?</AlertDialogTitle>
                        <AlertDialogDescription>
                          This action cannot be undone. The classroom will be marked as archived
                          and no new assignments can be created.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={() => archiveClassroom()} variant="destructive">
                          Archive Classroom
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
