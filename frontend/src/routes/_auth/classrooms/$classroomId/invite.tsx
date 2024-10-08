import { zodResolver } from "@hookform/resolvers/zod";
import { createFileRoute, Link, redirect, useRouter } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { getStatus, InviteForm, inviteFormSchema, Status } from "@/types/classroom";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, Clipboard, Loader2 } from "lucide-react";
import { Loader } from "@/components/loader.tsx";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { Header } from "@/components/header.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { formatDate, isStudent } from "@/lib/utils.ts";
import { ClassroomInvitation, UserClassroomResponse } from "@/swagger-client";
import { classroomInvitationsQueryOptions, classroomQueryOptions, useInviteClassroomMembers } from "@/api/classroom";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { toast } from "sonner";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/invite")({
  loader: async ({ context: { queryClient }, params }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    if (isStudent(userClassroom)) {
      throw redirect({
        to: "/classrooms/$classroomId",
        search: { tab: "assignments" },
        params,
        replace: true,
      });
    }
  },
  pendingComponent: Loader,
  component: ClassroomInviteForm,
});

function ClassroomInviteForm() {
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: invitations } = useSuspenseQuery(classroomInvitationsQueryOptions(classroomId));

  const { mutateAsync, isError, isPending } = useInviteClassroomMembers(classroomId);

  const form = useForm<z.infer<typeof inviteFormSchema>>({
    resolver: zodResolver(inviteFormSchema),
    defaultValues: {
      memberEmails: "",
    },
  });

  async function onSubmit(values: InviteForm) {
    await mutateAsync(values);
    form.reset();
  }

  return (
    <>
      <Breadcrumb className="mb-5">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId }}>
                {userClassroom.classroom.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Invitations</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <Header title="Invitations" />
      <InvitationsTable userClassroom={userClassroom} invitations={invitations} />

      <div className="flex flex-row justify-between">
        <h1 className="text-xl font-bold">Send new invitations</h1>
      </div>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <FormField
            control={form.control}
            name="memberEmails"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Textarea placeholder="toni@test.com" className="resize-none" {...field} />
                </FormControl>
                <FormDescription>E-Mails to invite into your Classroom</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit" disabled={isPending}>
            {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Submit"}
          </Button>

          {isError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>The classroom could not be created!</AlertDescription>
            </Alert>
          )}
        </form>
      </Form>
    </>
  );
}

function InvitationsTable({
  userClassroom,
  invitations,
}: {
  userClassroom: UserClassroomResponse;
  invitations: ClassroomInvitation[];
}) {
  const router = useRouter();
  return (
    <Table>
      <TableCaption>Classrooms</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>E-Mail</TableHead>
          <TableHead>Created At</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {invitations.map((i) => {
          const path = router.buildLocation({
            to: "/classrooms/$classroomId/invitations/$invitationId",
            params: { classroomId: userClassroom.classroom.id, invitationId: i.id },
          });
          return (
            <TableRow key={i.email}>
              <TableCell>{i.email}</TableCell>
              <TableCell>{formatDate(i.createdAt)}</TableCell>
              <TableCell>{getStatus(i.status)}</TableCell>
              <TableCell className="text-right">
                {i.status !== Status.Accepted && i.status !== Status.Revoked && (
                  <Button
                    variant="outline"
                    onClick={() => {
                      navigator.clipboard.writeText(`${location.origin}${path.pathname}`);
                      toast.success("Link copied to clipboard");
                    }}
                  >
                    <Clipboard className="mr-2 h-4 w-4" /> Get Link
                  </Button>
                )}
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
