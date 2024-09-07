import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { createFormSchema, updateFormSchema } from "@/types/classroom";
import { classroomQueryOptions, useCreateClassroom, useUpdateClassroom } from "@/api/classroom";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, Loader2 } from "lucide-react";
import { getUUIDFromLocation } from "@/lib/utils.ts";
import { Switch } from "@/components/ui/switch";
import { useSuspenseQuery } from "@tanstack/react-query";
import { toast } from "sonner";

export const ClassroomCreateForm = () => {
  const navigate = useNavigate();
  const { mutateAsync, isError, isPending } = useCreateClassroom();

  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    defaultValues: {
      name: "",
      description: "",
      maxTeams: 0,
      maxTeamSize: 2,
      createTeams: true,
      studentsViewAllProjects: false,
      teamsEnabled: true,
    },
  });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    const location = await mutateAsync(values);
    const classroomId = getUUIDFromLocation(location);
    await navigate({ to: "/classrooms/$classroomId", params: { classroomId } });
  }

  return (
    <div className="p-2">
      <div>
        <h1 className="text-xl font-bold">Create a classroom</h1>
        <p className="text-sm text-muted-foreground">Add the details you need and submit when you're done.</p>
      </div>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem className="space-y-1 my-2">
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input placeholder="Programming classroom" {...field} />
                </FormControl>
                <FormDescription>The name of the new classroom</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="description"
            render={({ field }) => (
              <FormItem className="space-y-1  my-2">
                <FormLabel>Description</FormLabel>
                <FormControl>
                  <Textarea placeholder="This is my awesome ..." className="resize-none" {...field} />
                </FormControl>
                <FormDescription>This is the description of your classroom.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="teamsEnabled"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center space-x-3 space-y-0">
                <FormControl>
                  <Switch checked={field.value} onCheckedChange={field.onChange} />
                </FormControl>
                <FormLabel>Enable Teams</FormLabel>
                <FormMessage />
              </FormItem>
            )}
          />
          <div className="border-l px-4" hidden={!form.getValues("teamsEnabled")}>
            <FormField
              control={form.control}
              name="maxTeams"
              render={({ field }) => (
                <FormItem className="space-y-1  my-2">
                  <FormLabel>Max Teams</FormLabel>
                  <FormControl>
                    <Input type="number" min={0} step={1} {...field} />
                  </FormControl>
                  <FormDescription>The maximum amount of teams. Keep at 0 to have no limit.</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="maxTeamSize"
              render={({ field }) => (
                <FormItem className="space-y-1 my-2">
                  <FormLabel>Max Team Size</FormLabel>
                  <FormControl>
                    <Input type="number" min={2} step={1} {...field} />
                  </FormControl>
                  <FormDescription>
                    The maximum amount of members per team. Must be at least 2. For one-person teams deactivate teams.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="createTeams"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center space-x-3 space-y-0">
                  <FormControl>
                    <Switch checked={field.value} onCheckedChange={field.onChange} />
                  </FormControl>
                  <FormLabel>Students can create teams</FormLabel>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
          <FormField
            control={form.control}
            name="studentsViewAllProjects"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center space-x-3 space-y-0">
                <FormControl>
                  <Switch checked={field.value} onCheckedChange={field.onChange} />
                </FormControl>
                <FormLabel>Students can inspect other students' repositories</FormLabel>
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
};

export const ClassroomEditForm = ({ classroomId }: { classroomId: string }) => {
  const { data: userClaasroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { mutateAsync, isError, isPending } = useUpdateClassroom(classroomId);

  const form = useForm<z.infer<typeof updateFormSchema>>({
    resolver: zodResolver(updateFormSchema),
    defaultValues: {
      name: userClaasroom.classroom.name,
      description: userClaasroom.classroom.description,
    },
  });

  async function onSubmit(values: z.infer<typeof updateFormSchema>) {
    await mutateAsync(values);
    toast.success("Classroom updated!");
  }

  return (
    <div className="p-2 w-full">
      <div>
        <h2 className="text-xl font-bold">Edit the classroom</h2>
        <p className="text-sm text-muted-foreground">Change the details of this classroom.</p>
      </div>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem className="space-y-1 my-2">
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input placeholder="Programming classroom" {...field} />
                </FormControl>
                <FormDescription>The name of the new classroom</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="description"
            render={({ field }) => (
              <FormItem className="space-y-1  my-2">
                <FormLabel>Description</FormLabel>
                <FormControl>
                  <Textarea placeholder="This is my awesome ..." className="resize-none" {...field} />
                </FormControl>
                <FormDescription>This is the description of your classroom.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit" disabled={isPending}>
            {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Save"}
          </Button>

          {isError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>The classroom could not be updated!</AlertDescription>
            </Alert>
          )}
        </form>
      </Form>
    </div>
  );
};
