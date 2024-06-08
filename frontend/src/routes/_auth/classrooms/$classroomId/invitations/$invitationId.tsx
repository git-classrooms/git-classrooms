import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Loader2 } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Separator } from "@/components/ui/separator";
import { getUUIDFromLocation } from "@/lib/utils";
import { Action } from "@/swagger-client";
import { useSuspenseQuery } from "@tanstack/react-query";
import { classroomInvitationQueryOptions, useJoinClassroom } from "@/api/classroom";

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
  const { data } = useSuspenseQuery(classroomInvitationQueryOptions(classroomId, invitationId));
  const { mutateAsync, isError, isPending } = useJoinClassroom(classroomId, invitationId);

  const onAccept = async () => {
    const location = await mutateAsync(Action.Accept);
    const classroomId = getUUIDFromLocation(location);
    await navigate({ to: "/classrooms/$classroomId", params: { classroomId } });
  };

  const onReject = async () => {
    await mutateAsync(Action.Reject);
    await navigate({ to: "/classrooms" });
  };

  return (
    <div className="rounded-lg border flex flex-col gap-5 max-w-5xl mx-auto">
      <h1 className="text-5xl font-bold text-center mb-5">Join Classroom</h1>
      <Separator />
      <p className="text-slate-500">
        You have been invited to join the classroom <span className="font-bold">{data.classroom.name}</span> by{" "}
        <span className="font-bold">{data.classroom.owner.name}</span>
      </p>
      <Separator />
      <p className="text-slate-500">{data.classroom.description}</p>
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
            <AlertDescription>Can't join classroom!</AlertDescription>
          </Alert>
        </>
      )}
    </div>
  );
}
