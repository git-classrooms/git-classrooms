import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Loader2 } from "lucide-react";
import { Header } from "@/components/header.tsx";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { useAcceptAssignment } from "@/api/assignments.ts";

export const Route = createFileRoute('/_auth/classrooms/$classroomId/assignments/$assignmentId/accept/')({
  component: AcceptAssignment
})

function AcceptAssignment(){
  const navigate = useNavigate({
    from: "/_auth/classrooms/$classroomId/assignments/$assignmentId/accept/",
  });
  const { classroomId, assignmentId } = Route.useParams();
  const { mutateAsync, isError, isPending } = useAcceptAssignment(classroomId, assignmentId);
  const onClick = async ()=>{
    await mutateAsync()
    await navigate({ to: "/classrooms" });
  }


  return(
    <div className="p-2 space-y-6">
      <Header title="Accept Assignment">
        <Button onClick={onClick} disabled={isPending}>
          {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Accept Assignment"}
        </Button>
      </Header>
      {isError && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>Can't accept assignment!</AlertDescription>
        </Alert>
      )}
    </div>
  );

}
