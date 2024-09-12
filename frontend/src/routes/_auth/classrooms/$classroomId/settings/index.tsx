import { classroomQueryOptions } from "@/api/classroom";
import { ClassroomEditForm } from "@/components/classroomsForm";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { Info } from "lucide-react";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/settings/")({
  loader: async ({ params: { classroomId }, context: { queryClient } }) => {
    const userClassroom = await queryClient.fetchQuery(classroomQueryOptions(classroomId));
    return { userClassroom };
  },
  component: Index,
});

function Index() {
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  return (
    <div className="w-full">
      <div className="p-2 mb-8">
        <div className="flex mb-6">
          <div className="grow">
            <div className="flex items-center">
              <h2 className="text-xl font-bold mr-2.5">General information</h2>
            </div>
            <p className="text-sm text-muted-foreground">Overview of the fixed configuration of the classroom</p>
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-2">
          <div className="rounded-md border p-4">
            <span className="font-medium text-lg">Creator</span>
            <p>{userClassroom.classroom.owner.name}</p>
          </div>
          <div className="rounded-md border p-4">
            <span className="font-medium text-lg">Teams</span>
            <p>
              {userClassroom.classroom.maxTeamSize > 1 ? (
                <>
                  <p>{"Max. " + userClassroom.classroom.maxTeamSize + " members"}</p>
                  <p>
                    {userClassroom.classroom.maxTeams !== 0 && ` (with max. ${userClassroom.classroom.maxTeams} Teams)`}
                  </p>
                </>
              ) : (
                "Disabled"
              )}
            </p>
          </div>
          {userClassroom.classroom.maxTeamSize > 1 && (
            <div className="rounded-md border p-4">
              <span className="font-medium text-lg">Teams by students</span>
              <p>{userClassroom.classroom.createTeams ? "Enabled" : "Disabled"}</p>
            </div>
          )}
          <Tooltip>
            <TooltipTrigger asChild>
              <div className="rounded-md border p-4">
                <span className="font-medium text-lg flex items-center gap-2">
                  Mutual code view <Info className="w-3.5 h-3.5" />
                </span>
                <p>{userClassroom.classroom.studentsViewAllProjects ? "Enabled" : "Disabled"}</p>
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <p>Members can see the code of other members.</p>
            </TooltipContent>
          </Tooltip>
        </div>
      </div>
      <ClassroomEditForm userClassroom={userClassroom} />
    </div>
  );
}
