import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertCircle, Edit, ExternalLink, Loader2, Plus, Users2 } from "lucide-react";
import { TeamResponse, UserClassroomResponse } from "@/swagger-client";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { CreateTeamForm } from "./createTeamForm";
import { useState } from "react";
import { ClassroomTeamModal } from "./classroomTeam";
import { isModerator, isStudent } from "@/lib/utils";
import { useUpdateTeam } from "@/api/team";
import { useForm } from "react-hook-form";
import { createFormSchema } from "@/types/team";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { FormControl, FormField, FormItem, FormLabel, FormMessage, Form } from "./ui/form";
import { Alert, AlertTitle, AlertDescription } from "./ui/alert";
import { Input } from "./ui/input";
import { toast } from "sonner";
import { StatusBadge } from "@/components/ui/status-badge";
import { Avatar } from "./avatar";
import { cn } from "@/lib/utils";

export function TeamListCard({
  teams,
  classroomId,
  userClassroom,
  maxTeamSize,
  numInvitedMembers,
  studentsCanCreateTeams,
  deactivateInteraction,
  teamsReportUrls,
}: {
  teams: TeamResponse[];
  classroomId: string;
  userClassroom: UserClassroomResponse;
  maxTeamSize: number;
  numInvitedMembers: number;
  studentsCanCreateTeams: boolean;
  deactivateInteraction: boolean;
  teamsReportUrls: Map<string, string>;
}) {
  const teamSlots = teams.length * maxTeamSize;
  const [open, setOpen] = useState(false);

  const hasWarning = teamSlots < numInvitedMembers && isModerator(userClassroom);

  return (
    <div className="space-y-6">
      {/* Header with Actions */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
            <Users2 className="w-5 h-5 text-primary" />
          </div>
          <div>
            <h3 className="font-semibold">Teams</h3>
            <p className="text-sm text-muted-foreground">
              {teams.length} team{teams.length !== 1 ? "s" : ""} · {maxTeamSize} max members each
            </p>
          </div>
        </div>

        {isModerator(userClassroom) && !deactivateInteraction && (
          <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
              <Button variant="glow" size="sm">
                <Plus className="w-4 h-4 mr-2" />
                Create Team
              </Button>
            </DialogTrigger>
            <DialogContent>
              <CreateTeamForm onSuccess={() => setOpen(false)} classroomId={classroomId} />
            </DialogContent>
          </Dialog>
        )}
      </div>

      {/* Warning */}
      {hasWarning && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Not enough team spots</AlertTitle>
          <AlertDescription>
            {teamSlots} spots available but {numInvitedMembers} students need teams.
            {!studentsCanCreateTeams && " Students cannot create teams themselves."}
          </AlertDescription>
        </Alert>
      )}

      {/* Teams Grid */}
      {teams.length === 0 ? (
        <EmptyTeamsState canCreate={isModerator(userClassroom) && !deactivateInteraction} />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {teams.map((team) => (
            <TeamCard
              key={team.id}
              team={team}
              classroomId={classroomId}
              userClassroom={userClassroom}
              maxTeamSize={maxTeamSize}
              teamsReportUrls={teamsReportUrls}
              deactivateInteraction={deactivateInteraction}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function EmptyTeamsState({ canCreate }: { canCreate: boolean }) {
  return (
    <div className="border border-dashed border-border rounded-lg p-8 text-center">
      <Users2 className="w-10 h-10 text-muted-foreground mx-auto mb-3" />
      <h3 className="font-medium text-foreground mb-1">No teams yet</h3>
      <p className="text-sm text-muted-foreground">
        {canCreate ? "Create a team to get started" : "No teams have been created yet"}
      </p>
    </div>
  );
}

function TeamCard({
  team,
  classroomId,
  userClassroom,
  maxTeamSize,
  teamsReportUrls,
  deactivateInteraction,
}: {
  team: TeamResponse;
  classroomId: string;
  userClassroom: UserClassroomResponse;
  maxTeamSize: number;
  teamsReportUrls: Map<string, string>;
  deactivateInteraction: boolean;
}) {
  const reportUrl = teamsReportUrls.get(team.id)!;
  const memberCount = team.members.length;
  const isFull = memberCount >= maxTeamSize;

  return (
    <Card className="group transition-all duration-200 hover:border-primary/30 hover:shadow-lg hover:shadow-background/50">
      <CardContent className="p-4">
        {/* Team Header */}
        <div className="flex items-start justify-between mb-4">
          <div className="flex items-center gap-3">
            <TeamAvatar name={team.name} />
            <div>
              <h4 className="font-semibold">{team.name}</h4>
              <div className="flex items-center gap-2 mt-1">
                <StatusBadge variant={isFull ? "warning" : "success"} size="sm" showDot>
                  {memberCount}/{maxTeamSize}
                </StatusBadge>
              </div>
            </div>
          </div>
          {isModerator(userClassroom) && (
            <ChangeTeamDialog classroomId={classroomId} team={team} />
          )}
        </div>

        {/* Members */}
        <div className="space-y-2">
          {team.members.length > 0 ? (
            team.members.slice(0, 3).map((m) => (
              <div key={m.user.id} className="flex items-center gap-2">
                <Avatar
                  avatarUrl={m.user.avatarURL}
                  fallbackUrl={m.user.fallbackAvatarURL}
                  name={m.user.name}
                  className="w-6 h-6"
                />
                <span className="text-sm truncate">{m.user.name}</span>
              </div>
            ))
          ) : (
            <p className="text-sm text-muted-foreground italic">No members yet</p>
          )}
          {team.members.length > 3 && (
            <p className="text-xs text-muted-foreground">
              +{team.members.length - 3} more
            </p>
          )}
        </div>

        {/* Actions */}
        <div className="flex items-center gap-2 mt-4 pt-3 border-t border-border">
          <Button variant="ghost" size="sm" className="h-8" asChild>
            <a href={team.webUrl} target="_blank" rel="noopener noreferrer">
              <ExternalLink className="w-3 h-3 mr-1" />
              GitLab
            </a>
          </Button>
          {!deactivateInteraction &&
            (!isStudent(userClassroom) || userClassroom.classroom.studentsViewAllProjects) && (
              <ClassroomTeamModal
                userClassroom={userClassroom}
                classroomId={classroomId}
                teamId={team.id}
                reportUrl={reportUrl}
              />
            )}
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
        "w-10 h-10 rounded-lg flex items-center justify-center bg-gradient-to-br",
        gradient
      )}
    >
      <span className="font-mono font-semibold text-primary-foreground">
        {name.charAt(0).toUpperCase()}
      </span>
    </div>
  );
}

function ChangeTeamDialog({ classroomId, team }: { classroomId: string; team: TeamResponse }) {
  const [open, setOpen] = useState(false);
  const { mutateAsync, isError, isPending } = useUpdateTeam(classroomId, team.id);

  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    mode: "onBlur",
    reValidateMode: "onChange",
    defaultValues: {
      name: team.name,
    },
  });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    try {
      await mutateAsync(values);
      toast.success("Team updated successfully!");
      setOpen(false);
      form.reset({ name: values.name });
    } catch {
      // Error is handled by isError state
    }
  }

  // Reset form when dialog opens
  const handleOpenChange = (isOpen: boolean) => {
    setOpen(isOpen);
    if (isOpen) {
      form.reset({ name: team.name });
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button variant="ghost" size="icon" className="h-8 w-8">
          <Edit className="h-4 w-4" />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Team</DialogTitle>
          <DialogDescription>Change the name of the team</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Team Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter team name" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div className="flex gap-2 justify-end">
              <Button type="button" variant="outline" onClick={() => setOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                Save Changes
              </Button>
            </div>

            {isError && (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>The team could not be updated!</AlertDescription>
              </Alert>
            )}
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

// Re-export TeamTable for backwards compatibility
export function TeamTable({
  teams,
  teamsReportUrls,
  classroomId,
  userClassroom,
  maxTeamSize,
  isPending,
  onTeamSelect,
  deactivateInteraction,
}: {
  teams: TeamResponse[];
  teamsReportUrls: Map<string, string>;
  classroomId: string;
  userClassroom: UserClassroomResponse;
  maxTeamSize: number;
  isPending?: boolean;
  onTeamSelect?: (teamId: string) => void;
  deactivateInteraction: boolean;
}) {
  // These props are kept for backwards compatibility but not used in current implementation
  void teamsReportUrls;
  void classroomId;
  void userClassroom;
  void deactivateInteraction;
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {teams.map((team) => (
        <TeamSelectCard
          key={team.id}
          team={team}
          maxTeamSize={maxTeamSize}
          isPending={isPending}
          onSelect={onTeamSelect}
        />
      ))}
    </div>
  );
}

function TeamSelectCard({
  team,
  maxTeamSize,
  isPending,
  onSelect,
}: {
  team: TeamResponse;
  maxTeamSize: number;
  isPending?: boolean;
  onSelect?: (teamId: string) => void;
}) {
  const isFull = team.members.length >= maxTeamSize;

  return (
    <Card
      className={cn(
        "transition-all duration-200 cursor-pointer",
        isFull
          ? "opacity-50 cursor-not-allowed"
          : "hover:border-primary/30 hover:shadow-lg"
      )}
      onClick={() => !isFull && !isPending && onSelect?.(team.id)}
    >
      <CardContent className="p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <TeamAvatar name={team.name} />
            <div>
              <h4 className="font-semibold">{team.name}</h4>
              <StatusBadge variant={isFull ? "warning" : "success"} size="sm">
                {team.members.length}/{maxTeamSize} members
              </StatusBadge>
            </div>
          </div>
          {!isFull && (
            <Button variant="ghost" size="sm" disabled={isPending}>
              {isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : "Join"}
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
