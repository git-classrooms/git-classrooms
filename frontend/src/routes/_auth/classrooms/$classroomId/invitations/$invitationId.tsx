import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Loader2 } from "lucide-react";
import { useJoinClassroom } from "@/api/classrooms.ts";
import { Header } from "@/components/header.tsx";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/invitations/$invitationId")({
  component: JoinClassroom,
});

function JoinClassroom() {
  const navigate = useNavigate({
    from: "/_auth/classrooms/$classroomId/invitations/$invitationId/",
  });
  const { invitationId } = Route.useParams();
  const { mutateAsync, isError, isPending } = useJoinClassroom(invitationId);
  const onClick = async () => {
    await mutateAsync();
    await navigate({ to: "/classrooms" });
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
