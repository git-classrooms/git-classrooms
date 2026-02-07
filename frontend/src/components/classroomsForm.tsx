import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { createFormSchema, updateFormSchema } from "@/types/classroom";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, Loader2, Pencil } from "lucide-react";
import { getUUIDFromLocation, unwrapApiError } from "@/lib/utils.ts";
import { Switch } from "@/components/ui/switch";
import { toast } from "sonner";
import { useCreateClassroom, useUpdateClassroom } from "@/api/classroom";
import { UserClassroomResponse } from "@/swagger-client";
import { Header } from "./header";
import { useUnsavedChanges } from "@/hooks/useUnsavedChanges";
import { Card, CardContent } from "@/components/ui/card";

export const ClassroomCreateForm = () => {
  const navigate = useNavigate();
  const { mutateAsync, error, isPending } = useCreateClassroom();

  const createClassroomError = unwrapApiError(error);

  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    mode: "onBlur",
    reValidateMode: "onChange",
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

  useUnsavedChanges({ isDirty: form.formState.isDirty });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    const location = await mutateAsync(values);
    const classroomId = getUUIDFromLocation(location);
    await navigate({ to: "/classrooms/$classroomId", search: { tab: "assignments" }, params: { classroomId } });
  }

  return (
    <div className="p-2">
      <Header title="Create a classroom" subtitle="Add the details you need and submit when you're done." />
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
                <FormDescription>The description of your classroom</FormDescription>
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

          {createClassroomError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{createClassroomError.message}</AlertDescription>
            </Alert>
          )}
        </form>
      </Form>
    </div>
  );
};

export const ClassroomEditForm = ({ userClassroom }: { userClassroom: UserClassroomResponse }) => {
  const { mutateAsync, isError, isPending } = useUpdateClassroom(userClassroom.classroom.id);

  const form = useForm<z.infer<typeof updateFormSchema>>({
    resolver: zodResolver(updateFormSchema),
    mode: "onBlur",
    reValidateMode: "onChange",
    defaultValues: {
      name: userClassroom.classroom.name,
      description: userClassroom.classroom.description,
    },
  });

  useUnsavedChanges({ isDirty: form.formState.isDirty });

  async function onSubmit(values: z.infer<typeof updateFormSchema>) {
    await mutateAsync(values);
    toast.success("Classroom updated!");
  }

  return (
    <div className="w-full">
      <div className="flex items-center gap-3 mb-6">
        <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-[hsl(142,71%,45%)]/20 to-[hsl(142,71%,45%)]/5 border border-[hsl(142,71%,45%)]/20 flex items-center justify-center">
          <Pencil className="w-5 h-5 text-[hsl(142,71%,45%)]" />
        </div>
        <div>
          <h2 className="text-lg font-semibold font-mono">Edit Classroom</h2>
          <p className="text-sm text-muted-foreground">Update name and description</p>
        </div>
      </div>

      <Card className="border-border/50">
        <CardContent className="p-4">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs text-muted-foreground uppercase tracking-wide">
                      Name
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Programming classroom"
                        {...field}
                        className="bg-background"
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      The display name of this classroom
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs text-muted-foreground uppercase tracking-wide">
                      Description
                    </FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="This is my awesome classroom for..."
                        className="resize-none bg-background min-h-[100px]"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      A brief description of the classroom's purpose
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="flex items-center justify-between pt-4 border-t border-border/50">
                <p className="text-xs text-muted-foreground">
                  Changes are saved immediately
                </p>
                <Button
                  type="submit"
                  variant="glow"
                  size="sm"
                  disabled={isPending || !form.formState.isDirty}
                >
                  {isPending ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : null}
                  Save Changes
                </Button>
              </div>

              {isError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>Error</AlertTitle>
                  <AlertDescription>The classroom could not be updated!</AlertDescription>
                </Alert>
              )}
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
};
