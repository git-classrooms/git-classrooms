import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Loader2 } from "lucide-react";
import { useJoinClassroom } from "@/api/classrooms.ts";
import { Header } from "@/components/header.tsx";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { getUUIDFromLocation } from "@/lib/utils";

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
    <div className="p-2 space-y-6">
      <Header title="Join Classroom">
        <Button onClick={onClick} disabled={isPending}>
          {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Join Classroom"}
        </Button>
      </Header>
      {isError && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>Can't join classroom!</AlertDescription>
        </Alert>
      )}
    </div>
  );
}
