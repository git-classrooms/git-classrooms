import { zodResolver } from "@hookform/resolvers/zod";
import { createFileRoute } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { ClassroomInvitation, GetStatus, InviteForm, inviteFormSchema } from "@/types/classroom";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, Loader2 } from "lucide-react";
import { classroomInvitationsQueryOptions, useInviteClassroomMembers } from "@/api/classrooms.ts";
import { Loader } from "@/components/loader.tsx";
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { Header } from "@/components/header.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { formatDate } from "@/lib/utils.ts";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/invite")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(classroomInvitationsQueryOptions(params.classroomId)),
  pendingComponent: Loader,
  component: ClassroomsForm,
});

function ClassroomsForm() {
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
    <div className="p-2">
      <div className="flex flex-row justify-between">
        <h1 className="text-xl font-bold">Current invitations</h1>
      </div>

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
            <TableCell>{GetStatus[i.status]}</TableCell>
            <TableCell className="text-right">
              <Button variant="outline">Refresh status</Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
