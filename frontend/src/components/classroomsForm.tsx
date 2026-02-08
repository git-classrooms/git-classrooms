import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { createFormSchema, updateFormSchema } from "@/types/classroom";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  AlertCircle,
  BookOpen,
  ChevronRight,
  FolderGit2,
  Loader2,
  Pencil,
  Settings2,
  Users,
} from "lucide-react";
import { cn, getUUIDFromLocation, unwrapApiError } from "@/lib/utils";
import { Switch } from "@/components/ui/switch";
import { toast } from "sonner";
import { useCreateClassroom, useUpdateClassroom } from "@/api/classroom";
import { UserClassroomResponse } from "@/swagger-client";
import { useUnsavedChanges } from "@/hooks/useUnsavedChanges";
import { Card, CardContent } from "@/components/ui/card";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";

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

  const teamsEnabled = form.watch("teamsEnabled");

  return (
    <div className="max-w-2xl mx-auto">
      {/* Breadcrumb */}
      <Breadcrumb className="mb-6">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms">Classrooms</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Create Classroom</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-4 mb-3">
          <div className="w-14 h-14 rounded-xl bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center">
            <FolderGit2 className="w-7 h-7 text-primary" />
          </div>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">Create Classroom</h1>
            <p className="text-muted-foreground">Set up a new classroom for your students</p>
          </div>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          {/* Basic Info Section */}
          <Card className="border-border/50 overflow-hidden">
            <div className="px-5 py-4 border-b border-border/50 bg-muted/30">
              <div className="flex items-center gap-2">
                <BookOpen className="w-4 h-4 text-muted-foreground" />
                <h2 className="font-semibold text-sm">Basic Information</h2>
              </div>
            </div>
            <CardContent className="p-5 space-y-5">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Classroom Name</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="e.g., Introduction to Programming WS24"
                        className="bg-background"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Choose a descriptive name for your classroom
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
                    <FormLabel>Description</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="Brief description of the classroom's purpose and content..."
                        className="resize-none bg-background min-h-[100px]"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Help students understand what this classroom is about
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Team Configuration Section */}
          <Card className="border-border/50 overflow-hidden">
            <div className="px-5 py-4 border-b border-border/50 bg-muted/30">
              <div className="flex items-center gap-2">
                <Users className="w-4 h-4 text-muted-foreground" />
                <h2 className="font-semibold text-sm">Team Configuration</h2>
              </div>
            </div>
            <CardContent className="p-5 space-y-5">
              <FormField
                control={form.control}
                name="teamsEnabled"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between p-4 rounded-lg border border-border/50 bg-muted/20">
                    <div className="space-y-0.5">
                      <FormLabel className="font-medium">Enable Teams</FormLabel>
                      <FormDescription className="text-xs">
                        Allow students to work together in teams
                      </FormDescription>
                    </div>
                    <FormControl>
                      <Switch checked={field.value} onCheckedChange={field.onChange} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <div
                className={cn(
                  "space-y-4 overflow-hidden transition-all duration-300",
                  teamsEnabled ? "opacity-100 max-h-[500px]" : "opacity-0 max-h-0 pointer-events-none"
                )}
              >
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <FormField
                    control={form.control}
                    name="maxTeams"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Max Teams</FormLabel>
                        <FormControl>
                          <Input
                            type="number"
                            min={0}
                            step={1}
                            className="bg-background"
                            {...field}
                          />
                        </FormControl>
                        <FormDescription className="text-xs">
                          0 = unlimited teams
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="maxTeamSize"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Max Team Size</FormLabel>
                        <FormControl>
                          <Input
                            type="number"
                            min={2}
                            step={1}
                            className="bg-background"
                            {...field}
                          />
                        </FormControl>
                        <FormDescription className="text-xs">
                          Members per team (min. 2)
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <FormField
                  control={form.control}
                  name="createTeams"
                  render={({ field }) => (
                    <FormItem className="flex items-center justify-between p-4 rounded-lg border border-border/50">
                      <div className="space-y-0.5">
                        <FormLabel className="font-medium">Student Team Creation</FormLabel>
                        <FormDescription className="text-xs">
                          Allow students to create their own teams
                        </FormDescription>
                      </div>
                      <FormControl>
                        <Switch checked={field.value} onCheckedChange={field.onChange} />
                      </FormControl>
                    </FormItem>
                  )}
                />
              </div>
            </CardContent>
          </Card>

          {/* Privacy Settings Section */}
          <Card className="border-border/50 overflow-hidden">
            <div className="px-5 py-4 border-b border-border/50 bg-muted/30">
              <div className="flex items-center gap-2">
                <Settings2 className="w-4 h-4 text-muted-foreground" />
                <h2 className="font-semibold text-sm">Privacy Settings</h2>
              </div>
            </div>
            <CardContent className="p-5">
              <FormField
                control={form.control}
                name="studentsViewAllProjects"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between p-4 rounded-lg border border-border/50">
                    <div className="space-y-0.5">
                      <FormLabel className="font-medium">Mutual Code Visibility</FormLabel>
                      <FormDescription className="text-xs">
                        Students can view other students' repositories
                      </FormDescription>
                    </div>
                    <FormControl>
                      <Switch checked={field.value} onCheckedChange={field.onChange} />
                    </FormControl>
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Error Alert */}
          {createClassroomError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{createClassroomError.message}</AlertDescription>
            </Alert>
          )}

          {/* Submit Button */}
          <div className="flex items-center justify-end gap-3 pt-4">
            <Button type="button" variant="outline" asChild>
              <Link to="/classrooms">Cancel</Link>
            </Button>
            <Button type="submit" variant="glow" disabled={isPending} className="min-w-[140px]">
              {isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <>
                  Create Classroom
                  <ChevronRight className="ml-2 h-4 w-4" />
                </>
              )}
            </Button>
          </div>
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
