import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, ArrowLeft, Check, GraduationCap, Loader2, Sparkles, User } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { getUUIDFromLocation } from "@/lib/utils";
import { Action } from "@/swagger-client";
import { useSuspenseQuery } from "@tanstack/react-query";
import { classroomInvitationQueryOptions, useJoinClassroom } from "@/api/classroom";
import { AxiosError } from "axios";
import { z } from "zod";
import { Card, CardContent } from "@/components/ui/card";

const seachSchema = z.object({
  groupLink: z.boolean().catch(false),
});

export const Route = createFileRoute("/_auth/classrooms/$classroomId/invitations/$invitationId")({
  validateSearch: seachSchema,
  loaderDeps: ({ search }) => ({ search }),
  loader: async ({ context: { queryClient }, params, deps: { search: { groupLink } } }) => {
    const invitationInfo = await queryClient.ensureQueryData(
      classroomInvitationQueryOptions(params.classroomId, params.invitationId, groupLink)
    );
    return { invitationInfo };
  },
  component: JoinClassroom,
});

function JoinClassroom() {
  const navigate = useNavigate();
  const { classroomId, invitationId } = Route.useParams();
  const { groupLink } = Route.useSearch();
  const { data: invitation } = useSuspenseQuery(
    classroomInvitationQueryOptions(classroomId, invitationId, groupLink)
  );
  const { mutateAsync, isError, isPending, error } = useJoinClassroom(classroomId, invitationId, groupLink);

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
    <div className="space-y-8 animate-stagger-1">
      {/* Main Content */}
      <div className="flex justify-center">
        <Card className="w-full max-w-xl border-border/50">
          <CardContent className="p-8">
            {/* Header Icon */}
            <div className="flex justify-center mb-6">
              <div className="relative">
                <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center">
                  <GraduationCap className="w-10 h-10 text-primary" />
                </div>
                <div className="absolute -bottom-1 -right-1 w-7 h-7 rounded-full bg-primary flex items-center justify-center border-4 border-card">
                  <Check className="w-3.5 h-3.5 text-primary-foreground" />
                </div>
              </div>
            </div>

            {/* Title */}
            <div className="text-center mb-8">
              <h1 className="text-2xl font-bold tracking-tight mb-2">Join Classroom</h1>
              <p className="text-muted-foreground">Ready to join this classroom?</p>
            </div>

            {/* Classroom Info Card */}
            <div className="p-4 rounded-lg bg-muted/30 border border-border/50 mb-6">
              <div className="flex items-start gap-4">
                <div className="w-12 h-12 rounded-lg bg-primary/15 flex items-center justify-center shrink-0">
                  <GraduationCap className="w-6 h-6 text-primary" />
                </div>
                <div className="flex-1 min-w-0">
                  <h2 className="font-semibold text-lg">{invitation.classroom.name}</h2>
                  {invitation.classroom.description && (
                    <p className="text-sm text-muted-foreground mt-1">{invitation.classroom.description}</p>
                  )}
                  <div className="flex items-center gap-4 mt-3 text-xs text-muted-foreground">
                    <div className="flex items-center gap-1.5">
                      <User className="w-3.5 h-3.5" />
                      <span>{invitation.classroom.owner.name}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* What happens next */}
            <div className="p-4 rounded-lg bg-primary/5 border border-primary/20 mb-8">
              <div className="flex gap-3">
                <Sparkles className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <div>
                  <h3 className="font-medium text-sm mb-1">What happens next?</h3>
                  <p className="text-sm text-muted-foreground">
                    You'll become a member of this classroom and gain access to all assignments and resources.
                  </p>
                </div>
              </div>
            </div>

            {/* Actions */}
            <div className="flex gap-3">
              <Button
                onClick={onReject}
                variant="outline"
                className="flex-1"
                disabled={isPending}
              >
                {isPending ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <>
                    <ArrowLeft className="w-4 h-4 mr-2" />
                    Decline
                  </>
                )}
              </Button>
              <Button
                onClick={onAccept}
                variant="glow"
                className="flex-[2]"
                disabled={isPending}
              >
                {isPending ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <>
                    <Check className="w-4 h-4 mr-2" />
                    Accept & Join
                  </>
                )}
              </Button>
            </div>

            {/* Error State */}
            {isError && (
              <Alert variant="destructive" className="mt-6">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Unable to join</AlertTitle>
                <AlertDescription>
                  {error instanceof AxiosError
                    ? error.response?.data.error || "Can't join classroom!"
                    : "Can't join classroom!"}
                </AlertDescription>
              </Alert>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Helper text */}
      <p className="text-center text-sm text-muted-foreground">
        You can always access your classrooms from the dashboard
      </p>
    </div>
  );
}
