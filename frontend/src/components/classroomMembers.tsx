import { Role } from "@/types/classroom";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link } from "@tanstack/react-router";
import { ExternalLink, Mail, Settings, UserPlus, Users } from "lucide-react";
import { Avatar } from "@/components/avatar";
import { UserClassroomResponse } from "@/swagger-client";
import { ClassroomTeamModal } from "./classroomTeam";
import { isModerator, isStudent } from "@/lib/utils";
import { StatusBadge } from "@/components/ui/status-badge";

const roleConfig: Record<Role, { variant: "success" | "info" | "neutral"; label: string }> = {
  [Role.Owner]: { variant: "success", label: "Owner" },
  [Role.Moderator]: { variant: "info", label: "Moderator" },
  [Role.Student]: { variant: "neutral", label: "Student" },
};

export function MemberListCard({
  classroomMembers,
  teamsReportUrls,
  classroomId,
  userClassroom,
  showTeams,
  deactivateInteraction,
}: {
  classroomMembers: UserClassroomResponse[];
  teamsReportUrls: Map<string, string>;
  classroomId: string;
  userClassroom: UserClassroomResponse;
  showTeams: boolean;
  deactivateInteraction: boolean;
}) {
  const owners = classroomMembers.filter((m) => m.role === Role.Owner);
  const moderators = classroomMembers.filter((m) => m.role === Role.Moderator);
  const students = classroomMembers.filter((m) => m.role === Role.Student);

  return (
    <div className="space-y-6">
      {/* Header with Actions */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
            <Users className="w-5 h-5 text-primary" />
          </div>
          <div>
            <h3 className="font-semibold">Members</h3>
            <p className="text-sm text-muted-foreground">
              {classroomMembers.length} member{classroomMembers.length !== 1 ? "s" : ""} in this classroom
            </p>
          </div>
        </div>

        {!deactivateInteraction && isModerator(userClassroom) && (
          <div className="flex gap-2">
            <Button variant="outline" size="sm" asChild>
              <Link to="/classrooms/$classroomId/members" params={{ classroomId }}>
                <Settings className="w-4 h-4 mr-2" />
                Manage
              </Link>
            </Button>
            <Button variant="glow" size="sm" asChild>
              <Link to="/classrooms/$classroomId/invite" params={{ classroomId }}>
                <UserPlus className="w-4 h-4 mr-2" />
                Invite
              </Link>
            </Button>
          </div>
        )}
      </div>

      {/* Member Sections */}
      <div className="space-y-6">
        {owners.length > 0 && (
          <MemberSection
            title="Owners"
            members={owners}
            teamsReportUrls={teamsReportUrls}
            classroomId={classroomId}
            userClassroom={userClassroom}
            showTeams={showTeams}
          />
        )}

        {moderators.length > 0 && (
          <MemberSection
            title="Moderators"
            members={moderators}
            teamsReportUrls={teamsReportUrls}
            classroomId={classroomId}
            userClassroom={userClassroom}
            showTeams={showTeams}
          />
        )}

        {students.length > 0 && (
          <MemberSection
            title="Students"
            members={students}
            teamsReportUrls={teamsReportUrls}
            classroomId={classroomId}
            userClassroom={userClassroom}
            showTeams={showTeams}
          />
        )}
      </div>
    </div>
  );
}

function MemberSection({
  title,
  members,
  teamsReportUrls,
  classroomId,
  userClassroom,
  showTeams,
}: {
  title: string;
  members: UserClassroomResponse[];
  teamsReportUrls: Map<string, string>;
  classroomId: string;
  userClassroom: UserClassroomResponse;
  showTeams: boolean;
}) {
  return (
    <div>
      <h4 className="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-3">
        {title} ({members.length})
      </h4>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
        {members.map((member) => (
          <MemberCard
            key={member.user.id}
            member={member}
            teamsReportUrls={teamsReportUrls}
            classroomId={classroomId}
            userClassroom={userClassroom}
            showTeams={showTeams}
          />
        ))}
      </div>
    </div>
  );
}

function MemberCard({
  member,
  teamsReportUrls,
  classroomId,
  userClassroom,
  showTeams,
}: {
  member: UserClassroomResponse;
  teamsReportUrls: Map<string, string>;
  classroomId: string;
  userClassroom: UserClassroomResponse;
  showTeams: boolean;
}) {
  const reportUrl = teamsReportUrls.get(member.team?.id ?? "");
  const config = roleConfig[member.role as Role];

  return (
    <Card className="group transition-all duration-200 hover:border-primary/30 hover:shadow-lg hover:shadow-background/50">
      <CardContent className="p-4">
        <div className="flex items-start gap-3">
          <Avatar
            avatarUrl={member.user.avatarURL}
            fallbackUrl={member.user.fallbackAvatarURL}
            name={member.user.name!}
            className="w-10 h-10"
          />
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <span className="font-medium truncate">{member.user.name}</span>
              <StatusBadge variant={config.variant} size="sm">
                {config.label}
              </StatusBadge>
            </div>
            <p className="text-sm text-muted-foreground truncate">
              @{member.user.gitlabUsername}
            </p>
            {showTeams && member.team && (
              <p className="text-xs text-muted-foreground mt-1">
                Team: <span className="text-foreground">{member.team.name}</span>
              </p>
            )}
          </div>
        </div>

        {/* Action buttons */}
        <div className="flex items-center gap-1 mt-3 pt-3 border-t border-border opacity-0 group-hover:opacity-100 transition-opacity">
          <Button variant="ghost" size="sm" className="h-7 px-2" asChild>
            <a href={member.webUrl} target="_blank" rel="noopener noreferrer">
              <ExternalLink className="w-3 h-3 mr-1" />
              GitLab
            </a>
          </Button>
          {member.user.gitlabEmail && (
            <Button variant="ghost" size="sm" className="h-7 px-2" asChild>
              <a href={`mailto:${member.user.gitlabEmail}`}>
                <Mail className="w-3 h-3 mr-1" />
                Email
              </a>
            </Button>
          )}
          {(!isStudent(userClassroom) || userClassroom.classroom.studentsViewAllProjects) && member.team && reportUrl && (
            <ClassroomTeamModal
              userClassroom={userClassroom}
              classroomId={classroomId}
              teamId={member.team.id}
              reportUrl={reportUrl}
            />
          )}
        </div>
      </CardContent>
    </Card>
  );
}
