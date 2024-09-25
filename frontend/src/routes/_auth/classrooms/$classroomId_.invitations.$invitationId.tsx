import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Loader2 } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Separator } from "@/components/ui/separator";
import { getUUIDFromLocation } from "@/lib/utils";
import { Action } from "@/swagger-client";
import { useSuspenseQuery } from "@tanstack/react-query";
import { classroomInvitationQueryOptions, useJoinClassroom } from "@/api/classroom";
import GitlabLogo from "@/assets/gitlab_logo.svg";
import { AxiosError } from "axios";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/invitations/$invitationId")({
  loader: async ({ context: { queryClient }, params }) => {
    const invitationInfo = await queryClient.ensureQueryData(
      classroomInvitationQueryOptions(params.classroomId, params.invitationId),
    );
    return { invitationInfo };
  },
  component: JoinClassroom,
});

function JoinClassroom() {
  const navigate = useNavigate();
  const { classroomId, invitationId } = Route.useParams();
  const { data: invitation } = useSuspenseQuery(classroomInvitationQueryOptions(classroomId, invitationId));
  const { mutateAsync, isError, isPending, error } = useJoinClassroom(classroomId, invitationId);

  const onAccept = async () => {
    const location = await mutateAsync(Action.Accept);
    const classroomId = getUUIDFromLocation(location);
    await navigate({ to: "/classrooms/$classroomId", search: { tab: "assignments" }, params: { classroomId } });
  };

  const onReject = async () => {
    await mutateAsync(Action.Reject);
    await navigate({ to: "/dashboard" });
  };

  return (
    <div className="m-auto max-w-lg ">
      <div className="flex justify-center">
        <img src={GitlabLogo} className="max-w-xs" alt={"Logo"} />
      </div>

      <div className="p-6 rounded-lg border flex flex-col gap-5">
        <h1 className="text-5xl font-bold text-center mb-5">Join Classroom</h1>
        <p className="text-slate-500 text-lg">
          You have been invited to join the classroom <span className="font-bold">{invitation.classroom.name}</span> by{" "}
          <span className="font-bold">{invitation.classroom.owner.name}</span>
        </p>
        <Separator />
        <p className="text-slate-500">
          <p className="text-slate-500 ">Classroom Description:</p>
          <p className="text-slate-500 italic ml-5">{invitation.classroom.description}</p>
        </p>
        <Separator />
        <div className="flex justify-between">
          <Button onClick={onReject} variant="destructive" disabled={isPending}>
            {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Reject"}
          </Button>
          <Button onClick={onAccept} disabled={isPending}>
            {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Accept"}
          </Button>
        </div>
        {isError && (
          <>
            <Separator />
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>
                {error instanceof AxiosError
                  ? error.response?.data.error
                    ? error.response.data.error
                    : "Can't join classroom!"
                  : "Can't join classroom!"}
              </AlertDescription>
            </Alert>
          </>
        )}
      </div>
    </div>
  );
}
