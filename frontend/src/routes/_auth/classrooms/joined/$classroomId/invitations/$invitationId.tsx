import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Loader2 } from "lucide-react";
import { useJoinClassroom } from "@/api/classrooms.ts";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { getUUIDFromLocation } from "@/lib/utils";
import { OwnedClassroom } from "@/types/classroom.ts";

export const Route = createFileRoute("/_auth/classrooms/joined/$classroomId/invitations/$invitationId")({
  component: JoinClassroom,
});

function JoinClassroom() {
  const navigate = useNavigate();
  const { invitationId } = Route.useParams();
  const { mutateAsync, isError, isPending } = useJoinClassroom(invitationId);
  const onClick = async () => {
    const location = await mutateAsync();
    const classroomId = getUUIDFromLocation(location);
    await navigate({ to: "/classrooms/joined/$classroomId", params: { classroomId } });
  };

  return (
    <div className="p-6 rounded-lg  border">
      <h1 className="text-5xl font-bold text-slate-900 text-center">Join Classroom</h1>
      <div className="divide-y divide-solid">
        <div className="py-6">
          <p className="text-slate-500 text-lg">
            You have been invited to join the classroom NameOfClassroom by NameOwner
          </p>
        </div>
        <div className="py-6">
          <p className="text-slate-500 text-lg">
            This is the place for the classroom description. Where some details about the classroom are provided.
          </p>
        </div>
        <div className="pt-6 flex justify-between">
          <Button variant="destructive" disabled={isPending}>
            {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Reject"}
          </Button>
          <Button onClick={onClick} disabled={isPending}>
            {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Accept"}
          </Button>
        </div>
        {isError && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>Can't join classroom!</AlertDescription>
          </Alert>
        )}
      </div>
    </div>
  );
}
