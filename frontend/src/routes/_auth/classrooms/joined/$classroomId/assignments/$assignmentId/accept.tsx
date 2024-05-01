import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Loader2 } from "lucide-react";
import { Header } from "@/components/header.tsx";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { useAcceptAssignment } from "@/api/assignments.ts";

export const Route = createFileRoute('/_auth/classrooms/joined/$classroomId/assignments/$assignmentId/accept')({
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
    <div className="p-6 rounded-lg border">
      <h1 className="text-5xl font-bold text-center">Accept Assignment</h1>
      <div className="divide-y divide-solid">
        <div className="py-6">
          <p className="text-slate-500 text-lg">
            You need to accept the assignment 01-monads-and-other-awesome-stuff in the classroom 2024-SS-FFWE.
          </p>
        </div>
        <div className="py-6">
          <p className="text-slate-500 text-lg">
            Once you have accepted the assignment, you will get access to the repository monads-endboss. in the
            2024-SS-FFWE group.
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
      </div>

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
