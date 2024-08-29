import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";
import { Loader } from "@/components/loader.tsx";
import { ArrowRight as ArrowRight, Code, Gitlab } from "lucide-react";
import { Header } from "@/components/header";
import { classroomsQueryOptions } from "@/api/classroom";
import { Filter } from "@/types/classroom";
import { useMemo } from "react";
import { UserClassroomResponse } from "@/swagger-client";
import List from "@/components/ui/list.tsx";
import ListItem from "@/components/ui/listItem.tsx";
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import { Tabs, TabsContent, TabsList, TabsTrigger, } from "@/components/ui/tabs"


export const Route = createFileRoute("/_auth/classrooms/_index")({
  component: Classrooms,
  loader: async ({ context: { queryClient } }) => {
    const ownedClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Owned));
    const moderatorClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Moderator));
    const studentClassrooms = await queryClient.ensureQueryData(classroomsQueryOptions(Filter.Student));
    
    return {
      ownedClassrooms,
      moderatorClassrooms,
      studentClassrooms,
    };
  },
  pendingComponent: Loader,
});

function Classrooms() {
  const { data: ownedClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Owned));
  const { data: moderatorClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Moderator));
  const { data: studentClassrooms } = useSuspenseQuery(classroomsQueryOptions(Filter.Student));

  const joinedClassrooms = useMemo(
    () => [...moderatorClassrooms, ...studentClassrooms],
    [moderatorClassrooms, studentClassrooms],
  );

  return (
    <div>
      <Header title="Classrooms" />
        <Tabs defaultValue="managed" className="w-[400]">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="managed">Managed</TabsTrigger>
            <TabsTrigger value="joined">Joined</TabsTrigger>
          </TabsList>
          <TabsContent value="managed">
            <OwnedClassroomTable classrooms={ownedClassrooms} />
          </TabsContent>
          <TabsContent value="joined">
            <JoinedClassroomTable classrooms={joinedClassrooms} />
          </TabsContent>
        </Tabs>
        <Outlet />
    </div>
  );
}

function OwnedClassroomTable({ classrooms }: { classrooms: UserClassroomResponse[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Managed Classrooms</CardTitle>
        <CardDescription>Classrooms which are managed by you</CardDescription>
      </CardHeader>

      <CardContent>
        <List
          items={classrooms}
          renderItem={(item) => (
            <ListItem
              leftContent={
                <ListLeftContent
                  classroomName={item.classroom.name}
                  assignmentsCount={item.assignmentsCount}
                />
              }
              rightContent={
                <ListRightContent gitlabUrl={item.webUrl} classroomId={item.classroom.id} />
              }
            />
          )}
        />
      </CardContent>


      <CardFooter className="flex justify-end gap-2">
        <Button asChild variant="default">
          <Link to="/classrooms/create/modal" replace>
            Create a new Classroom
          </Link>
        </Button>
        <Button asChild variant="default">
          <Link to="/classrooms/create/modal" replace>
            View all your Classrooms
          </Link>
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
                  <Link to="/classrooms/$classroomId" params={{ classroomId: c.classroom.id }}>
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

function ListLeftContent({ classroomName, assignmentsCount }: {
  classroomName: string,
  assignmentsCount: number
}) {
  const assignmentsText = assignmentsCount === 1
    ? `${assignmentsCount} Assignment`
    : `${assignmentsCount} Assignments`;
  return (
    <div className="cursor-default flex">
      <div className="pr-2">
        <Avatar>
          <AvatarFallback className="bg-[#FC6D25] text-black text-lg">
            {classroomName.charAt(0)}
          </AvatarFallback>
        </Avatar>
      </div>
      <div>
        <div className="font-medium">{classroomName}</div>
        <div className="text-sm text-muted-foreground md:inline">
          {assignmentsText}
        </div>
      </div>
    </div>
  );
}

function ListRightContent({ gitlabUrl, classroomId }: { gitlabUrl: string, classroomId: string }) {
  return (
    <>
      <Button variant="ghost" size="icon" asChild>
        <a href={gitlabUrl} target="_blank" rel="noreferrer">
          <Gitlab className="h-6 w-6 text-slate-500 dark:text-white" />
        </a>
      </Button>
      <Button variant="ghost" size="icon" asChild>
        <Link to="/classrooms/$classroomId" params={{ classroomId: classroomId }}>
          <ArrowRight className="text-slate-500 dark:text-white" />
        </Link>
      </Button>
    </>
  );
}
