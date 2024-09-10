import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Loader2 } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { Link } from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Separator } from "@/components/ui/separator";
import { projectQueryOptions, useAcceptAssignment } from "@/api/project";
import { classroomQueryOptions } from "@/api/classroom";
import { Role } from "@/types/classroom.ts";
import GitlabLogo from "@/assets/gitlab_logo.svg";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/projects/$projectId/accept")({
  loader: async ({ context: { queryClient }, params }) => {
    const project = await queryClient.ensureQueryData(projectQueryOptions(params.classroomId, params.projectId));
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    if (userClassroom.role !== Role.Student || userClassroom.team?.id !== project.team.id) {
      throw redirect({
        to: "/classrooms/$classroomId",
        search: { tab: "assignments" },
        params,
      });
    }
    return { project };
  },
  component: AcceptAssignment,
});

function AcceptAssignment() {
  const navigate = useNavigate({
    from: "/_auth/classrooms/$classroomId/projects/$projectId/accept/",
  });
  const { classroomId, projectId } = Route.useParams();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: project } = useSuspenseQuery(projectQueryOptions(classroomId, projectId));
  const { mutateAsync, isError, isPending } = useAcceptAssignment(classroomId, projectId);

  const onClick = async () => {
    await mutateAsync();
    await navigate({ to: "/dashboard" });
  };

  return (
    <div className="m-auto max-w-lg ">
      <div className="flex justify-center">
        <img src={GitlabLogo} className="max-w-xs" alt={"Logo"} />
      </div>

      <div className="p-6 rounded-lg border flex flex-col gap-5">
        <h1 className="text-5xl font-bold text-center mb-5">Accept Assignment</h1>
        <p className="text-slate-500 text-lg">
          You need to accept the assignment{" "}
          <span className="text-slate-900  dark:text-slate-300 font-bold">{project.assignment.name}</span> <br />
          in classroom <span className="text-slate-900  dark:text-slate-300 font-bold">{classroom.classroom.name}</span>
          .
        </p>
        <Separator />
        <p className="text-slate-500">Once you have accepted the assignment, you will get access to the repository.</p>
        <Separator />
        <div className="flex justify-between">
          <Button variant="destructive" asChild>
            <Link
              to="/classrooms/$classroomId"
              search={{ tab: "assignments" }}
              params={{ classroomId }}
              property="stylesheet"
            >
              Reject
            </Link>
          </Button>
          <Button onClick={onClick} disabled={isPending}>
            {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Accept"}
          </Button>
        </div>
        {isError && (
          <>
            <Separator />
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>Can't accept assignment!</AlertDescription>
            </Alert>
          </>
        )}
      </div>
    </div>
  );
}
