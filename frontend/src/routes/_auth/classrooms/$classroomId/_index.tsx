import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, redirect, Link } from "@tanstack/react-router";
import { MemberListCard } from "@/components/classroomMembers.tsx";
import { Role } from "@/types/classroom.ts";
import { TeamListCard } from "@/components/classroomTeams.tsx";
import { AssignmentListCard } from "@/components/classroomAssignments.tsx";
import { Header } from "@/components/header";
import { classroomQueryOptions } from "@/api/classroom";
import { assignmentsQueryOptions } from "@/api/assignment";
import { membersQueryOptions } from "@/api/member";
import { teamsQueryOptions } from "@/api/team";
import { ReportApiAxiosParamCreator, UserClassroomResponse } from "@/swagger-client";
import { Button } from "@/components/ui/button.tsx";
import { Pen, Archive } from "lucide-react";
import { useState, useCallback } from "react";
import { useArchiveClassroom } from "@/api/classroom";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/_index")({
  component: ClassroomDetail,
  loader: async ({ context: { queryClient }, params }) => {
    const teams = await queryClient.ensureQueryData(teamsQueryOptions(params.classroomId));
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    if (userClassroom.role === Role.Student && !userClassroom.team) {
      throw redirect({
        to: "/classrooms/$classroomId/teams/join",
        params,
      });
    }
    const { url: reportDownloadUrl } = await ReportApiAxiosParamCreator().getClassroomReport(params.classroomId)
    const members = await queryClient.ensureQueryData(membersQueryOptions(params.classroomId));
    if(userClassroom.role !== Role.Student) {
      const assignments = await queryClient.ensureQueryData(assignmentsQueryOptions(params.classroomId));
      return { userClassroom, assignments, members, teams, reportDownloadUrl };
    }else{
      return { userClassroom, members, teams, reportDownloadUrl };
    }
  },
  pendingComponent: Loader,
});

function ClassroomDetail() {
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  if (userClassroom.role === Role.Student) {
    return <ClassroomStudentView />;
  } else if (userClassroom.role === Role.Owner || userClassroom.role === Role.Moderator){
    return <ClassroomSupervisorView userClassroom={userClassroom} />;
  }
}

function ClassroomStudentView(){
  // Role.Student does not have access to assignments
  // This is a placeholder
  return (
    <div>
      <h1>Joined Classroom Info</h1>
    </div>
  );
}

function ClassroomSupervisorView( {userClassroom}: {userClassroom: UserClassroomResponse}){
  const { classroomId } = Route.useParams();
  // const { reportDownloadUrl } = Route.useLoaderData()
  const { data: classroomMembers } = useSuspenseQuery(membersQueryOptions(classroomId));
  const { data: teams } = useSuspenseQuery(teamsQueryOptions(classroomId));
  const { data: assignments } = useSuspenseQuery(assignmentsQueryOptions(classroomId));

  const [isDialogOpen, setDialogOpen] = useState(false);
  const archiveClassroom = useArchiveClassroom(classroomId);

  const handleArchiveClick = useCallback(() => {
    setDialogOpen(true);
  }, []);

  const handleConfirmArchive = useCallback(() => {
    archiveClassroom.mutate();
    setDialogOpen(false);
  }, [archiveClassroom]);

  const handleCancel = useCallback(() => {
    setDialogOpen(false);
  }, []);

  return (
    <div>
      <div className="grid grid-cols-[1fr,auto] justify-between gap-1">
        <Header title={`Classroom: ${userClassroom.classroom.name}`} subtitle={userClassroom.classroom.description} />
        <div className="grid grid-cols-2 gap-3">
          {!userClassroom.classroom.archived && (
            <Button className="col-start-1" variant="ghost" size="icon" onClick={handleArchiveClick} title="Archive classroom">
              <Archive className="text-slate-500 dark:text-white h-28 w-28" />
            </Button>
          )}
          <Button className="col-start-2" variant="ghost" size="icon" asChild title="Edit classroom">
            <Link to="/classrooms/$classroomId/edit/modal" params={{ classroomId: classroomId }} replace>
              <Pen className="text-slate-500 dark:text-white h-28 w-28" />
            </Link>
          </Button>
        </div>
        
      </div>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 justify-between gap-10">
        { /*<Button asChild><a href={reportDownloadUrl}>Download Report</a></Button>*/ }

        <AssignmentListCard
          assignments={assignments}
          classroomId={classroomId}
          classroomName={userClassroom.classroom.name}
        />
        <MemberListCard
          classroomMembers={classroomMembers}
          classroomId={classroomId}
          userRole={Role.Owner}
          showTeams={userClassroom.classroom.maxTeamSize > 1}
        />
        {/* uses Role.Owner, as you can only be the owner, making a check if GetMe.id == OwnedClassroom.ownerId unnecessary*/}
        {userClassroom.classroom.maxTeamSize > 1 && (
          <TeamListCard teams={teams} classroomId={classroomId} userRole={Role.Owner} maxTeamSize={userClassroom.classroom.maxTeamSize} numInvitedMembers={classroomMembers.length} />
        )}
        <Outlet />
      </div>

      <SimpleDialog
        isOpen={isDialogOpen}
        onConfirm={handleConfirmArchive}
        onCancel={handleCancel}
      />
    </div>
  );
}

type SimpleDialogProps = {
  isOpen: boolean;
  onConfirm: () => void;
  onCancel: () => void;
};

function SimpleDialog({ isOpen, onConfirm, onCancel }: SimpleDialogProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white p-6 rounded shadow-md">
        <h2 className="text-lg font-semibold">Archive classroom</h2>
        <p>Are you sure that you wanna archive this classroom? This action can not be undone!</p>
        <div className="mt-4 flex justify-end gap-2">
          <button onClick={onCancel} className="px-4 py-2 bg-gray-200 rounded">Cancel</button>
          <button onClick={onConfirm} className="px-4 py-2 bg-blue-500 text-white rounded">Confirm</button>
        </div>
      </div>
    </div>
  );
}