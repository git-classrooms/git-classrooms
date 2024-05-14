import { GetOwnedClassroomTeamProjectResponse, GetOwnedClassroomTeamResponse } from "@/swagger-client";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Edit, Gitlab } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { formatDate } from "@/lib/utils.ts";
import { User } from "@/swagger-client/models/user.ts";
import { Avatar } from "@/components/avatar.tsx";
import { Separator } from "@/components/ui/separator.tsx";

export function ClassroomTeamModal({ classroomId, team, projects }: {
  classroomId: string,
  team: GetOwnedClassroomTeamResponse,
  projects: GetOwnedClassroomTeamProjectResponse[]
}) {
  return <div>
    <h1 className="text-2xl">{team.name}</h1>
    <Separator className="my-1" />
    <h2 className="text-xl mt-4">Members</h2>
    <ClassroomTeamMemberTable members={team.members} />
    <Separator className="my-1" />
    <h2 className="text-xl mt-4">Assignments</h2>
    <ClassroomTeamAssignmentTable classroomId={classroomId} projects={projects} />
  </div>;
}

function ClassroomTeamMemberTable({ members }: { members: User[]; }) {
  return (
    <Table>
      <TableBody>
        {members != null && members.map((m) => (
          <TableRow key={m.id}>
            <TableCell className="p-2">
              <ClassroomTeamMemberListElement member={m} />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function ClassroomTeamMemberListElement({ member }: { member: User }) {
  return (
  <div className="flex">
        <div className="pr-2">
          <Avatar
            avatarUrl={member.gitlabAvatar?.avatarURL}
            fallbackUrl={member.gitlabAvatar?.fallbackAvatarURL}
            name={member.name!}
          />
        </div>
        <div>
          <div className="font-medium">{member.name}</div>
          <div className="text-sm text-muted-foreground mt-[-0.3rem]">
            @{member.gitlabUsername}
          </div>
        </div>
  </div>
  );
}

export function ClassroomTeamAssignmentTable({ classroomId, projects }: {
  classroomId: string,
  projects: GetOwnedClassroomTeamProjectResponse[]
}) {
  return <Table>
    <TableBody>
      {projects.map(p => (
        <TableRow key={p.id}>
          <TableCell className="p-2">
            <div className="cursor-default flex justify-between">
              <div>
                <div className="font-medium">{p.assignment.name}</div>
                <div
                  className="text-sm text-muted-foreground md:inline">{p.assignmentAccepted ? "Accepted" : "Pending"}</div>
              </div>
              <div className="flex items-end">
                <div className="ml-auto">
                  <div className="font-medium text-right">Due date</div>
                  <div className="text-sm text-muted-foreground md:inline">
                    {p.assignment.dueDate ? formatDate(p.assignment.dueDate) : "No Due Date"}
                  </div>
                </div>
                <Button variant="ghost" size="icon" asChild>
                  <Link
                    to="/classrooms/owned/$classroomId/assignments/$assignmentId"
                    params={{ classroomId: classroomId, assignmentId: p.assignment.id }}
                  >
                    <Edit className="h-6 w-6 text-gray-600" />
                  </Link>
                </Button>
                <Button variant="ghost" size="icon" asChild>
                  {p.assignmentAccepted ?
                    <a href={p.projectPath} target="_blank" rel="noreferrer">
                      <Gitlab className="h-6 w-6 text-gray-600" />
                    </a> : <div><Gitlab className="h-6 w-6 text-gray-400" /></div>
                  }
                </Button>
              </div>
            </div>
          </TableCell>
        </TableRow>
      ))}
    </TableBody>
  </Table>;
}
