import { Button } from "@/components/ui/button";
import { ArrowRight, ClipboardList, Download, ExternalLink, Users } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { formatRelativeTime, getDaysUntilDue, isOwner } from "@/lib/utils";
import { Avatar } from "@/components/avatar";
import { ProjectResponse, UserClassroomResponse } from "@/swagger-client";
import { useQuery } from "@tanstack/react-query";
import { teamProjectsQueryOptions } from "@/api/project";
import { teamQueryOptions } from "@/api/team";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "./ui/dialog";
import { Skeleton } from "./ui/skeleton";
import { StatusBadge } from "./ui/status-badge";
import { Card, CardContent } from "./ui/card";
import { cn } from "@/lib/utils";
import { useTranslation } from "react-i18next";

interface ClassroomTeamModalProps {
  userClassroom: UserClassroomResponse;
  classroomId: string;
  teamId: string;
  reportUrl: string;
}

export const ClassroomTeamModal = (props: ClassroomTeamModalProps) => {
  const { t } = useTranslation("team");
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="ghost" size="sm" className="h-7 px-2">
          <ClipboardList className="w-3 h-3 mr-1" />
          {t("modal.details")}
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-2xl">
        <ClassroomModalContent {...props} />
      </DialogContent>
    </Dialog>
  );
};

function ClassroomModalContent({ classroomId, teamId, reportUrl, userClassroom }: ClassroomTeamModalProps) {
  const { t } = useTranslation("team");
  const { data: team, isLoading: teamIsLoading, error: teamError } = useQuery(teamQueryOptions(classroomId, teamId));
  const {
    data: projects,
    isLoading: projectsIsLoading,
    error: projectsError,
  } = useQuery(teamProjectsQueryOptions(classroomId, teamId));

  const isLoading = teamIsLoading || projectsIsLoading;
  const error = teamError || projectsError;

  if (error) throw error;

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-24 w-full" />
        <Skeleton className="h-32 w-full" />
      </div>
    );
  }

  const isSoloTeam = team!.members.length === 1 && team!.members[0].user.gitlabUsername === team!.name;

  return (
    <>
      <DialogHeader>
        <DialogTitle className="flex items-center gap-3">
          <TeamAvatar name={team!.name} />
          <div>
            <span className="text-xl">{team!.name}</span>
            <p className="text-sm font-normal text-muted-foreground">
              {t("members.count", { count: team!.members.length })}
            </p>
          </div>
        </DialogTitle>
      </DialogHeader>

      <div className="space-y-6 mt-4">
        {/* Members Section */}
        {!isSoloTeam && (
          <section>
            <h3 className="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-3 flex items-center gap-2">
              <Users className="w-3 h-3" />
              {t("members.title")}
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
              {team!.members.map((member) => (
                <div
                  key={member.user.id}
                  className="flex items-center gap-3 p-3 rounded-lg bg-muted/50"
                >
                  <Avatar
                    avatarUrl={member.user.avatarURL}
                    fallbackUrl={member.user.fallbackAvatarURL}
                    name={member.user.name}
                    className="w-8 h-8"
                  />
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-sm truncate">{member.user.name}</p>
                    <p className="text-xs text-muted-foreground truncate">
                      @{member.user.gitlabUsername}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* Assignments Section */}
        <section>
          <h3 className="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-3 flex items-center gap-2">
            <ClipboardList className="w-3 h-3" />
            {t("modal.assignments", { count: projects!.length })}
          </h3>
          {projects!.length === 0 ? (
            <p className="text-sm text-muted-foreground text-center py-4">
              {t("modal.noAssignments")}
            </p>
          ) : (
            <div className="space-y-2">
              {projects!.map((project) => (
                <ProjectRow key={project.id} project={project} classroomId={classroomId} />
              ))}
            </div>
          )}
        </section>

        {/* Actions */}
        {isOwner(userClassroom) && (
          <div className="pt-4 border-t border-border">
            <Button variant="outline" className="w-full" asChild>
              <a href={reportUrl} target="_blank" rel="noopener noreferrer">
                <Download className="w-4 h-4 mr-2" />
                {t("modal.downloadReport")}
              </a>
            </Button>
          </div>
        )}
      </div>
    </>
  );
}

function ProjectRow({ project, classroomId }: { project: ProjectResponse; classroomId: string }) {
  const { t } = useTranslation("team");
  const daysUntil = getDaysUntilDue(project.assignment.dueDate);
  const isOverdue = daysUntil !== null && daysUntil < 0;
  const isAccepted = project.projectStatus === "accepted";

  const getStatusVariant = () => {
    if (project.projectStatus === "pending") return "warning";
    if (project.projectStatus === "accepted") return "success";
    return "neutral";
  };

  return (
    <Card className="transition-all duration-200 hover:border-primary/30">
      <CardContent className="p-3">
        <div className="flex items-center justify-between gap-4">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <span className="font-medium text-sm truncate">
                {project.assignment.name}
              </span>
              <StatusBadge variant={getStatusVariant()} size="sm">
                {project.projectStatus}
              </StatusBadge>
            </div>
            <p className="text-xs text-muted-foreground">
              {project.assignment.dueDate ? (
                <span className={cn(isOverdue && "text-destructive")}>
                  {t("modal.due")} {formatRelativeTime(project.assignment.dueDate)}
                </span>
              ) : (
                t("modal.noDueDate")
              )}
            </p>
          </div>

          <div className="flex items-center gap-1">
            <Button variant="ghost" size="icon" className="h-8 w-8" asChild>
              <Link
                to="/classrooms/$classroomId/assignments/$assignmentId"
                params={{ classroomId, assignmentId: project.assignment.id }}
              >
                <ArrowRight className="w-4 h-4" />
              </Link>
            </Button>
            {isAccepted ? (
              <Button variant="ghost" size="icon" className="h-8 w-8" asChild>
                <a href={project.webUrl} target="_blank" rel="noopener noreferrer">
                  <ExternalLink className="w-4 h-4" />
                </a>
              </Button>
            ) : (
              <Button variant="ghost" size="icon" className="h-8 w-8" disabled>
                <ExternalLink className="w-4 h-4 opacity-30" />
              </Button>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

function TeamAvatar({ name }: { name: string }) {
  const gradients = [
    "from-primary to-[hsl(280,100%,60%)]",
    "from-[hsl(142,71%,45%)] to-[hsl(185,100%,50%)]",
    "from-[hsl(38,92%,55%)] to-[hsl(0,72%,55%)]",
    "from-[hsl(280,65%,60%)] to-[hsl(210,100%,60%)]",
  ];
  const hash = name.split("").reduce((acc, char) => acc + char.charCodeAt(0), 0);
  const gradient = gradients[hash % gradients.length];

  return (
    <div
      className={cn(
        "w-12 h-12 rounded-xl flex items-center justify-center bg-gradient-to-br",
        gradient
      )}
    >
      <span className="font-mono font-bold text-lg text-primary-foreground">
        {name.charAt(0).toUpperCase()}
      </span>
    </div>
  );
}

// Re-export for backwards compatibility
export function ClassroomTeamAssignmentTable({
  classroomId,
  projects,
}: {
  classroomId: string;
  projects: ProjectResponse[];
}) {
  return (
    <div className="space-y-2">
      {projects.map((project) => (
        <ProjectRow key={project.id} project={project} classroomId={classroomId} />
      ))}
    </div>
  );
}
