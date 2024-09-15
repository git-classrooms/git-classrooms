import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table.tsx";
import { Button } from "@/components/ui/button.tsx";
import { ArrowRight, Clipboard, Download, Gitlab } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { formatDateWithTime, isOwner } from "@/lib/utils.ts";
import { Avatar } from "@/components/avatar.tsx";
import { ProjectResponse, UserClassroomResponse } from "@/swagger-client";
import { useQuery } from "@tanstack/react-query";
import { teamProjectsQueryOptions } from "@/api/project";
import { teamQueryOptions } from "@/api/team";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "./ui/dialog";
import { Skeleton } from "./ui/skeleton";
import { Separator } from "./ui/separator";

interface ClassroomTeamModalProps {
  userClassroom: UserClassroomResponse;
  classroomId: string;
  teamId: string;
  reportUrl: string;
}

const ClassroomModalContent = ({ classroomId, teamId, reportUrl, userClassroom }: ClassroomTeamModalProps) => {
  const { data: team, isLoading: teamIsLoading, error: teamError } = useQuery(teamQueryOptions(classroomId, teamId));
  const {
    data: projects,
    isLoading: projectsIsLoading,
    error: projectsError,
  } = useQuery(teamProjectsQueryOptions(classroomId, teamId));

  const isLoading = teamIsLoading || projectsIsLoading;
  const error = teamError || projectsError;

  if (error) throw error;

  return (
    <>
      <DialogHeader>
        <DialogTitle>
          {isLoading
            ? "Loading..."
            : (team!.members.length === 0 || team!.members[0].user.gitlabUsername != team!.name) && team!.name}
        </DialogTitle>
        <DialogDescription>
          {isLoading
            ? "Loading..."
            : (team!.members.length === 0 || team!.members[0].user.gitlabUsername != team!.name) && "Members"}
        </DialogDescription>
      </DialogHeader>
      {isLoading ? (
        <Skeleton className="max-w-[462px] max-h-[206px] w-full h-full" />
      ) : (
        <>
          <ClassroomTeamMemberTable members={team!.members} />
          <Separator className="my-1" />
          <h2 className="text-xl mt-4">Assignments</h2>
          <ClassroomTeamAssignmentTable classroomId={classroomId} projects={projects!} />
          {isOwner(userClassroom) && (
            <>
              <Separator className="my-1" />
              <Button asChild variant="outline">
                <a href={reportUrl} target="_blank" rel="noreferrer">
                  <Download className="h-4 m-4" />
                  Download grading report
                </a>
              </Button>
            </>
          )}
        </>
      )}
    </>
  );
};

export const ClassroomTeamModal = (props: ClassroomTeamModalProps) => (
  <Dialog>
    <DialogTrigger asChild>
      <Button variant="ghost" size="icon">
        <Clipboard className="h-6 w-6 text-gray-600 dark:text-white" />
      </Button>
    </DialogTrigger>
    <DialogContent>
      <ClassroomModalContent {...props} />
    </DialogContent>
  </Dialog>
);

function ClassroomTeamMemberTable({ members }: { members: UserClassroomResponse[] }) {
  return (
    <Table>
      <TableBody>
        {members.length > 0 ? (
          members.map((m) => (
            <TableRow key={m.user.id}>
              <TableCell className="p-2">
                <ClassroomTeamMemberListElement member={m} />
              </TableCell>
            </TableRow>
          ))
        ) : (
          <TableRow>
            <TableCell className="p-2">No member in this team</TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}

function ClassroomTeamMemberListElement({ member }: { member: UserClassroomResponse }) {
  return (
    <div className="flex">
      <div className="pr-2">
        <Avatar
          avatarUrl={member.user.gitlabAvatar?.avatarURL}
          fallbackUrl={member.user.gitlabAvatar?.fallbackAvatarURL}
          name={member.user.name}
        />
      </div>
      <div>
        <div className="font-medium">{member.user.name}</div>
        <div className="text-sm text-muted-foreground mt-[-0.3rem]">@{member.user.gitlabUsername}</div>
      </div>
    </div>
  );
}

export function ClassroomTeamAssignmentTable({
  classroomId,
  projects,
}: {
  classroomId: string;
  projects: ProjectResponse[];
}) {
  return (
    <Table>
      <TableBody>
        {projects.map((p) => (
          <TableRow key={p.id}>
            <TableCell className="p-2">
              <div className="cursor-default flex justify-between">
                <div>
                  <div className="font-medium">{p.assignment.name}</div>
                  <div className="text-sm text-muted-foreground md:inline">{p.projectStatus}</div>
                </div>
                <div className="flex items-end">
                  <div className="ml-auto">
                    <div className="font-medium text-right">Due date</div>
                    <div className="text-sm text-muted-foreground md:inline">
                      {p.assignment.dueDate ? formatDateWithTime(p.assignment.dueDate) : "No Due Date"}
                    </div>
                  </div>
                  <Button className="ml-2" variant="ghost" size="icon" title="Go to assignment" asChild>
                    <Link
                      to="/classrooms/$classroomId/assignments/$assignmentId"
                      params={{ classroomId: classroomId, assignmentId: p.assignment.id }}
                    >
                      <ArrowRight className="h-6 w-6 text-gray-600" />
                    </Link>
                  </Button>
                  <Button variant="ghost" size="icon" title="Go to project" asChild>
                    {p.projectStatus === "accepted" ? (
                      <a href={p.webUrl} target="_blank" rel="noreferrer">
                        <Gitlab className="h-6 w-6 text-gray-600" />
                      </a>
                    ) : (
                      <div>
                        <Gitlab className="h-6 w-6 text-gray-400" />
                      </div>
                    )}
                  </Button>
                </div>
              </div>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
