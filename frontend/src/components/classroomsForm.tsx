import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { createFormSchema } from "@/types/classroom";
import { useCreateClassroom } from "@/api/classrooms";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, Loader2 } from "lucide-react";
import { getUUIDFromLocation } from "@/lib/utils.ts";
import { Switch } from "@/components/ui/switch"
import { useEffect, useState } from "react";

export const ClassroomsForm = () => {
  const navigate = useNavigate();
  const [areTeamsEnabled, setAreTeamsEnabled] = useState(true);
  const [prevMaxTeams, setPrevMaxTeams] = useState(0);
  const [prevMaxTeamSize, setPrevMaxTeamSize] = useState(1);
  const [prevCanStudentsCreateTeams, setCanStudentsCreateTeams] = useState(true);
  const { mutateAsync, isError, isPending } = useCreateClassroom();

  useEffect(() => {
    if (areTeamsEnabled) {
      form.setValue("maxTeams", prevMaxTeams);
      form.setValue("maxTeamSize", prevMaxTeamSize);
      form.setValue("createTeams", prevCanStudentsCreateTeams);
    } else {
      setPrevMaxTeams(form.getValues().maxTeams);
      setPrevMaxTeamSize(form.getValues().maxTeamSize);
      setCanStudentsCreateTeams(form.getValues().createTeams);
      form.setValue("maxTeams", 0);
      form.setValue("maxTeamSize", 1);
      form.setValue("createTeams", false);
    }
  }, [areTeamsEnabled]);
  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    defaultValues: {
      name: "",
      description: "",
      maxTeams: 0,
      maxTeamSize: 1,
      createTeams: true,
      studentsViewAllProjects: false,
    },
  });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    const location = await mutateAsync(values);
    const classroomId = getUUIDFromLocation(location);
    await navigate({ to: "/classrooms/owned/$classroomId", params: { classroomId } });
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
          <div className="flex flex-row items-center space-x-3 space-y-0">
            <Switch checked={areTeamsEnabled}  onCheckedChange={setAreTeamsEnabled} />
            <FormLabel>Enable Teams</FormLabel>
          </div>
          <div className ="border-l px-4" hidden={!areTeamsEnabled}>
          <FormField
            control={form.control}
            name="maxTeams"
            render={({ field }) => (
              <FormItem className="space-y-1  my-2">
                <FormLabel>Max Teams</FormLabel>
                <FormControl>
                  <Input type="number" step={1} {...field} />
                </FormControl>
                <FormDescription>The maximum amount of teams</FormDescription>
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
                  <Input type="number" step={1} {...field} />
                </FormControl>
                <FormDescription>The maximum amount of members per team. Must be at least 2. For one-person teams deactivate teams.</FormDescription>
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
