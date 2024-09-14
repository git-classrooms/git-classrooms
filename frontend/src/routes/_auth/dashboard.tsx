import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { ArrowRight as ArrowRight, Code, Gitlab } from "lucide-react";
import { Header } from "@/components/header";
import { classroomsQueryOptions } from "@/api/classroom";
import { Filter } from "@/types/classroom";
import { useMemo } from "react";
import { AssignmentResponse, UserClassroomResponse } from "@/swagger-client";
import List from "@/components/ui/list.tsx";
import ListItem from "@/components/ui/listItem.tsx";
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import { activeAssignmentsQueryOptions } from "@/api/assignment.ts";
import { formatDate } from "@/lib/utils.ts";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip.tsx";

export const Route = createFileRoute("/_auth/dashboard")({
  component: Classrooms,
  loader: async ({ context: { queryClient } }) => {
    const activeAssignments = await queryClient.ensureQueryData(activeAssignmentsQueryOptions());
    const ownedClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Owned));
    const moderatorClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Moderator));
    const studentClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Student));
    return {
      activeAssignments,
      ownedClassrooms,
      moderatorClassrooms,
      studentClassrooms,
    };
  },
  pendingComponent: Loader,
});

function Classrooms() {
  const { data: activeAssignments } = useSuspenseQuery(activeAssignmentsQueryOptions());
  const { data: ownedClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Owned));
  const { data: moderatorClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Moderator));
  const { data: studentClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Student));

  const joinedClassrooms = useMemo(
    () => [...moderatorClassrooms, ...studentClassrooms],
    [moderatorClassrooms, studentClassrooms],
  );

  const classroomIdsToNames = useMemo(() => {
    const classroomIdsToNames = new Map<string, string>();
    ownedClassrooms.forEach((c) => classroomIdsToNames.set(c.classroom.id, c.classroom.name));
    joinedClassrooms.forEach((c) => classroomIdsToNames.set(c.classroom.id, c.classroom.name));
    return classroomIdsToNames;
  }, [ownedClassrooms, joinedClassrooms]);

  const sortedAssignments = useMemo(() => {
    return [...activeAssignments].sort((a, b) => {
      if (a.dueDate === null) return 1;
      if (b.dueDate === null) return -1;
      return new Date(a.dueDate ?? 0).getTime() - new Date(b.dueDate ?? 0).getTime();
    });
  }, [activeAssignments]);

  return (
    <div>
      <div className="flex-1 space-y-4">
        <Header title="Dashboard" />
        <ActiveAssignmentsTable assignments={sortedAssignments} classroomIdsToNames={classroomIdsToNames} />
        <div className="grid grid-cols-1 lg:grid-cols-2 justify-between gap-4">
          <OwnedClassroomTable classrooms={ownedClassrooms} />
          <JoinedClassroomTable classrooms={joinedClassrooms} />
          <Outlet />
        </div>
      </div>
    </div>
  );
}

function OwnedClassroomTable({ classrooms }: { classrooms: UserClassroomResponse[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Owned Classrooms</CardTitle>
        <CardDescription>Classrooms which are managed by you</CardDescription>
      </CardHeader>

      <CardContent>
        <List
          items={classrooms}
          renderItem={(item) => (
            <ListItem
              leftContent={
                <ListLeftContent classroomName={item.classroom.name} assignmentsCount={item.assignmentsCount} />
              }
              rightContent={<ListRightContent gitlabUrl={item.webUrl} classroomId={item.classroom.id} />}
            />
          )}
        />
      </CardContent>

      <CardFooter className="flex justify-end gap-2">
        <Button asChild variant="default">
          <Link to="/classrooms/create">Create a new Classroom</Link>
        </Button>
        <Button asChild variant="default">
          <Link to="/classrooms">View all your Classrooms</Link>
        </Button>
      </CardFooter>
    </Card>
  );
}

function JoinedClassroomTable({ classrooms }: { classrooms: UserClassroomResponse[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Joined Classrooms</CardTitle>
        <CardDescription>Classrooms which you have joined</CardDescription>
      </CardHeader>
      <Table className="flex-auto">
        <TableBody>
          {classrooms.map((c) => (
            <TableRow key={c.classroom.id}>
              <TableCell>{c.classroom.name}</TableCell>
              <TableCell>{c.classroom.owner.name}</TableCell>
              <TableCell>
                <a href={c.webUrl} target="_blank" rel="noreferrer">
                  <Code />
                </a>
              </TableCell>
              <TableCell className="text-right">
                <Button variant="outline">
                  <Link
                    to="/classrooms/$classroomId"
                    search={{ tab: "assignments" }}
                    params={{ classroomId: c.classroom.id }}
                  >
                    <ArrowRight />
                  </Link>
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Card>
  );
}

function ActiveAssignmentsTable({ assignments, classroomIdsToNames }: { assignments: AssignmentResponse[], classroomIdsToNames: Map<string, string> }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Active Assignments</CardTitle>
        <CardDescription>Current assignments that are not yet closed</CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Classroom</TableHead>
              <TableHead>Due Date</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {assignments.map((a) => (
              <TableRow key={a.id}>
                <TableCell>{a.name}</TableCell>
                <TableCell>{classroomIdsToNames.get(a.classroomId)}</TableCell>
                <TableCell>{a.dueDate ? formatDate(a.dueDate) : "No expiry"}</TableCell>
                <TableCell className="text-right">
                  <Tooltip>
                    <TooltipTrigger>
                      <Button size="icon" variant="ghost" title="Go to classroom" asChild>
                        <Link
                          to="/classrooms/$classroomId"
                          search={{ tab: "assignments" }}
                          params={{ classroomId: a.classroomId }}
                        >
                          <ArrowRight className="h-6 w-6 text-gray-600 dark:text-white" />
                        </Link>
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>Go to classroom</p>
                    </TooltipContent>
                  </Tooltip>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}

function ListLeftContent({ classroomName, assignmentsCount }: { classroomName: string; assignmentsCount: number }) {
  const assignmentsText = assignmentsCount === 1 ? `${assignmentsCount} Assignment` : `${assignmentsCount} Assignments`;
  return (
    <div className="cursor-default flex">
      <div className="pr-2">
        <Avatar>
          <AvatarFallback className="bg-[#FC6D25] text-black text-lg">{classroomName.charAt(0)}</AvatarFallback>
        </Avatar>
      </div>
      <div>
        <div className="font-medium">{classroomName}</div>
        <div className="text-sm text-muted-foreground md:inline">{assignmentsText}</div>
      </div>
    </div>
  );
}

function ListRightContent({ gitlabUrl, classroomId }: { gitlabUrl: string; classroomId: string }) {
  return (
    <>
      <Button variant="ghost" size="icon" asChild>
        <a href={gitlabUrl} target="_blank" rel="noreferrer">
          <Gitlab className="h-6 w-6 text-slate-500 dark:text-white" />
        </a>
      </Button>
      <Button variant="ghost" size="icon" asChild>
        <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId: classroomId }}>
          <ArrowRight className="text-slate-500 dark:text-white" />
        </Link>
      </Button>
    </>
  );
}
