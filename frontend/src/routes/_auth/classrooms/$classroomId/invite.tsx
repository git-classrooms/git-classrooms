import { zodResolver } from "@hookform/resolvers/zod";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { getStatus, InviteForm, inviteFormSchema, Role } from "@/types/classroom";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, Loader2 } from "lucide-react";
import { Loader } from "@/components/loader.tsx";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { Header } from "@/components/header.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { formatDate } from "@/lib/utils.ts";
import { ClassroomInvitation } from "@/swagger-client";
import { classroomInvitationsQueryOptions, classroomQueryOptions, useInviteClassroomMembers } from "@/api/classroom";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/invite")({
  loader: async ({ context: { queryClient }, params }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    if (userClassroom.role === Role.Student) {
      throw redirect({
        to: "/classrooms/$classroomId",
        search: { tab: "assignments" },
        params,
      });
    }
  },
  pendingComponent: Loader,
  component: ClassroomInviteForm,
});

function ClassroomInviteForm() {
  const { classroomId } = Route.useParams();
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
    <div className="max-w-5xl mx-auto">
      <Header title="Invitations" />
      <InvitationsTable invitations={invitations} />

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
    </div>
  );
}

function InvitationsTable({ invitations }: { invitations: ClassroomInvitation[] }) {
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
        {invitations.map((i) => (
          <TableRow key={i.email}>
            <TableCell>{i.email}</TableCell>
            <TableCell>{formatDate(i.createdAt)}</TableCell>
            <TableCell>{getStatus(i.status)}</TableCell>
            <TableCell className="text-right">
              <Button variant="outline">Refresh status</Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
