import { createFileRoute, Link } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { Header } from "@/components/header.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  Activity,
  AlertCircle,
  CalendarClock,
  Download,
  FolderGit2,
  Info,
  Loader2,
  MoreHorizontal,
  Scale,
  SearchCode,
  Settings,
  Text,
} from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { cn, formatDate, formatDateWithTime, isModerator, isOwner } from "@/lib/utils.ts";
import { assignmentQueryOptions } from "@/api/assignment";
import { assignmentProjectsQueryOptions, useInviteToAssignment } from "@/api/project";
import { ProjectResponse, ReportApiAxiosParamCreator, UserClassroomResponse } from "@/swagger-client";
import { classroomQueryOptions } from "@/api/classroom";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { formatDistanceToNow } from "date-fns";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { getStatusProps } from "@/types/projects";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/$assignmentId/")({
  loader: async ({ context: { queryClient }, params: { classroomId, assignmentId } }) => {
    const classroom = await queryClient.ensureQueryData(classroomQueryOptions(classroomId));
    const assignment = await queryClient.ensureQueryData(assignmentQueryOptions(classroomId, assignmentId));
    const assignmentProjects = await queryClient.ensureQueryData(
      assignmentProjectsQueryOptions(classroomId, assignmentId),
    );

    const { url: reportDownloadUrl } = await ReportApiAxiosParamCreator().getClassroomAssignmentReport(
      classroomId,
      assignmentId,
    );

    const urls = (
      await Promise.all(
        assignmentProjects.map(async (project) => ({
          url: (await ReportApiAxiosParamCreator().getClassroomTeamReport(classroomId, project.teamId)).url,
          projectId: project.id,
        })),
      )
    ).reduce((acc, { url, projectId }) => acc.set(projectId, url), new Map<string, string>());

    return { classroom, assignment, assignmentProjects, reportDownloadUrl, urls };
  },
  component: AssignmentDetail,
  pendingComponent: Loader,
});

function AssignmentDetail() {
  const { classroomId, assignmentId } = Route.useParams();
  const { reportDownloadUrl } = Route.useLoaderData();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  const { data: assignmentProjects } = useSuspenseQuery(assignmentProjectsQueryOptions(classroomId, assignmentId));

  const { mutateAsync, isError, isPending } = useInviteToAssignment(classroomId, assignmentId);

  return (
    <div>
      <Breadcrumb className="mb-5">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId }}>
                {classroom.classroom.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{assignment.name}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <div className="md:flex justify-between gap-1 mb-4">
        <Header title={assignment.name} subtitle="Assignment overview" />
        <div className="grid md:grid-cols-2 gap-3">
          {isModerator(classroom) && (
            <Button variant="secondary" asChild size="sm" title="Grading">
              <Link
                to="/classrooms/$classroomId/assignments/$assignmentId/grading"
                params={{ classroomId, assignmentId }}
              >
                <Scale className="mr-2 h-4 w-4" />
                Grading
              </Link>
            </Button>
          )}
          {/* <Button variant="secondary" asChild size="sm" title="Download Report">
            <a href={reportDownloadUrl} target="_blank" referrerPolicy="no-referrer">
              <Download className="mr-2 h-4 w-4" />
              Download Report
            </a>
          </Button> */}
          {isOwner(classroom) && (
            <Button variant="secondary" asChild size="sm" title="Settings">
              <Link
                to="/classrooms/$classroomId/assignments/$assignmentId/settings/"
                params={{ classroomId, assignmentId }}
              >
                <Settings className="mr-2 h-4 w-4" />
                Settings
              </Link>
            </Button>
          )}
        </div>
      </div>

      <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Due date</CardTitle>
            <CalendarClock className="mr-2 h-4 w-4" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {assignment.dueDate ? formatDate(new Date(assignment.dueDate)) : "No Expiry"}
            </div>
            {assignment.dueDate && (
              <p className="text-xs text-muted-foreground"> {formatDistanceToNow(new Date(assignment.dueDate))}</p>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Last changes</CardTitle>
            <Activity className="mr-2 h-4 w-4" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {assignment.updatedAt ? formatDate(new Date(assignment.updatedAt)) : "Unknown"}
            </div>
            {assignment.updatedAt && (
              <p className="text-xs text-muted-foreground">
                {" "}
                {assignment.updatedAt ? formatDistanceToNow(new Date(assignment.updatedAt)) + " ago" : ""}{" "}
              </p>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Projects</CardTitle>
            <FolderGit2 className="mr-2 h-4 w-4" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {assignmentProjects.filter((p) => p.projectStatus == "accepted").length}
            </div>
            <p className="text-xs text-muted-foreground">
              {assignmentProjects.length == 1 ? "accepted project" : "accepted projects"}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Status</CardTitle>
            <Info className="mr-2 h-4 w-4" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{assignment.closed === true ? "Closed" : "Open"}</div>
          </CardContent>
        </Card>

        <Card className="col-span-1 md:col-span-2 lg:col-span-4">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Description</CardTitle>
            <Text className="mr-2 h-4 w-4" />
          </CardHeader>
          <CardContent>
            <p>{assignment.description ?? <i>No description available</i>}</p>
          </CardContent>
        </Card>
      </div>

      <Card className="mt-16 mb-6 p-2">
        <CardHeader className="md:flex flex-row items-center justify-between space-y-0 pb-2 mb-4">
          <div>
            <CardTitle className="mb-1">Assignment projects</CardTitle>
            <CardDescription>
              {classroom.classroom.maxTeamSize === 1
                ? "All individual projects managed by the classroom"
                : "All projects per team managed by the classroom"}
            </CardDescription>
          </div>
          {isModerator(classroom) && (
            <Tooltip delayDuration={0}>
              <TooltipTrigger asChild>
                <Button variant="outline" onClick={() => mutateAsync()} disabled={isPending}>
                  {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Send Invites"}
                </Button>
              </TooltipTrigger>
              <TooltipContent>Sends invitations to members who have not yet accepted.</TooltipContent>
            </Tooltip>
          )}
        </CardHeader>
        <CardContent>
          <AssignmentProjectTable assignmentProjects={assignmentProjects} classroom={classroom} />
        </CardContent>
      </Card>

      {isError && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>The invitation could not be send!</AlertDescription>
        </Alert>
      )}
    </div>
  );
}

function AssignmentProjectTable({
  classroom,
  assignmentProjects,
}: {
  classroom: UserClassroomResponse;
  assignmentProjects: ProjectResponse[];
}) {
  const { urls } = Route.useLoaderData();
  return (
    <Table>
      <TableCaption>Projects</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Project name</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Invitet at</TableHead>
          <TableHead></TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {assignmentProjects.map((a) => {
          const statusProps = getStatusProps(a.projectStatus);
          const reportUrl = urls.get(a.id)!;
          return (
            <TableRow key={`${a.assignment.id}-${a.team.id}`}>
              <TableCell className="font-medium">{a.team.name}</TableCell>
              <TableCell>
                <div className="flex pl-1 gap-3 items-center">
                  <span className="relative flex h-3 w-3">
                    <span
                      className={cn(
                        "animate-ping absolute inline-flex h-full w-full rounded-full opacity-75",
                        statusProps.color.secondary,
                      )}
                    ></span>
                    <span className={cn("relative inline-flex rounded-full h-3 w-3", statusProps.color.primary)}></span>
                  </span>
                  {statusProps.name}
                </div>
              </TableCell>
              <TableCell>{formatDateWithTime(new Date(a.createdAt))}</TableCell>
              <TableCell className="text-right float-right">
                <>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">Open actions</span>
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuLabel>Actions</DropdownMenuLabel>
                      <DropdownMenuItem disabled={a.projectStatus !== "accepted"} asChild>
                        <a href={a.webUrl} target="_blank">
                          <SearchCode className="mr-2 h-4 w-4" />
                          Go to project
                        </a>
                      </DropdownMenuItem>
                      {isModerator(classroom) && (
                        <>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem disabled={a.projectStatus !== "accepted"} asChild>
                            <a href={reportUrl} target="_blank" referrerPolicy="no-referrer">
                              <Download className="mr-2 h-4 w-4" />
                              Download report
                            </a>
                          </DropdownMenuItem>
                        </>
                      )}
                    </DropdownMenuContent>
                  </DropdownMenu>
                </>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
