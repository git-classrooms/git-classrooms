import { createFileRoute } from "@tanstack/react-router";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover.tsx";
import { Button } from "@/components/ui/button.tsx";
import { cn, formatDateWithTime } from "@/lib/utils.ts";
import { AlertTriangle, Calendar as CalendarIcon, Lock, Loader2, Pencil, X } from "lucide-react";
import { Calendar } from "@/components/ui/calendar.tsx";
import { assignmentQueryOptions, assignmentsQueryOptions, useUpdateAssignment } from "@/api/assignment.ts";
import { UpdateAssignmentForm, updateAssignmentFormSchema } from "@/types/assignments.ts";
import { Loader } from "@/components/loader.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { assignmentProjectsQueryOptions } from "@/api/project.ts";
import { Assignment, DatabaseStatus, ProjectResponse } from "@/swagger-client";
import { toast } from "sonner";
import { TimePicker } from "@/components/ui/timer-picker";
import { addSeconds } from "date-fns";
import { Card, CardContent } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/$assignmentId/settings/")({
  component: Index,
  loader: async ({ context: { queryClient }, params }) => {
    const assignment = await queryClient.ensureQueryData(
      assignmentQueryOptions(params.classroomId, params.assignmentId),
    );
    const assignmentProjects = await queryClient.ensureQueryData(
      assignmentProjectsQueryOptions(params.classroomId, params.assignmentId),
    );
    const assignments = await queryClient.ensureQueryData(assignmentsQueryOptions(params.classroomId));
    return { assignment, assignmentProjects, assignments };
  },
  pendingComponent: Loader,
});

function hasAcceptedAssignment(projects: ProjectResponse[]) {
  return projects.some((project) => project.projectStatus === DatabaseStatus.Accepted);
}

function checkNewAssignmentNameValid(assignment: Assignment, assignments: Assignment[], name: string) {
  return name === assignment.name || !assignments.some((assignment) => assignment.name === name);
}

function Index() {
  const { classroomId, assignmentId } = Route.useParams();

  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  const { data: assignmentProjects } = useSuspenseQuery(assignmentProjectsQueryOptions(classroomId, assignmentId));
  const { data: assignments } = useSuspenseQuery(assignmentsQueryOptions(classroomId));

  const { mutateAsync, isError, isPending } = useUpdateAssignment(classroomId, assignmentId);

  const isAccepted = hasAcceptedAssignment(assignmentProjects);

  const form = useForm<UpdateAssignmentForm>({
    resolver: zodResolver(updateAssignmentFormSchema(isAccepted)),
    defaultValues: {
      name: assignment.name,
      description: assignment.description,
      dueDate: assignment.dueDate ? new Date(assignment.dueDate) : undefined,
    },
  });

  async function onSubmit(values: UpdateAssignmentForm) {
    await mutateAsync({
      name: values.name ? values.name : "",
      description: values.description ? values.description : "",
      dueDate: values.dueDate?.toISOString(),
    });
    toast.success("Assignment updated successfully");
  }

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-2 duration-300">
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-[hsl(142,71%,45%)]/20 to-[hsl(142,71%,45%)]/5 border border-[hsl(142,71%,45%)]/20 flex items-center justify-center">
          <Pencil className="w-5 h-5 text-[hsl(142,71%,45%)]" />
        </div>
        <div>
          <h2 className="text-lg font-semibold font-mono">Edit Assignment</h2>
          <p className="text-sm text-muted-foreground">Update assignment details</p>
        </div>
      </div>

      {/* Locked Fields Warning */}
      {isAccepted && (
        <Alert className="border-warning/50 bg-warning/5">
          <AlertTriangle className="h-4 w-4 text-warning" />
          <AlertDescription className="text-warning">
            Name and description are locked because at least one team has accepted this assignment.
          </AlertDescription>
        </Alert>
      )}

      {/* Form Card */}
      <Card className="border-border/50">
        <CardContent className="p-4">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              {/* Name Field */}
              <FormField
                control={form.control}
                name="name"
                disabled={isAccepted}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs text-muted-foreground uppercase tracking-wide flex items-center gap-2">
                      Name
                      {isAccepted && <Lock className="w-3 h-3" />}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Programming Assignment"
                        {...field}
                        onBlur={async (e) => {
                          field.onBlur();
                          if (checkNewAssignmentNameValid(assignment, assignments, e.target.value)) {
                            form.clearErrors("name");
                          } else {
                            form.setError("name", {
                              type: "manual",
                              message: "This name is already taken.",
                            });
                          }
                        }}
                        className={cn("bg-background", isAccepted && "opacity-60")}
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      The display name of this assignment
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Description Field */}
              <FormField
                control={form.control}
                name="description"
                disabled={isAccepted}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs text-muted-foreground uppercase tracking-wide flex items-center gap-2">
                      Description
                      {isAccepted && <Lock className="w-3 h-3" />}
                    </FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="Describe what students need to accomplish..."
                        className={cn("resize-none bg-background min-h-[100px]", isAccepted && "opacity-60")}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      A brief description of the assignment's objectives
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Due Date Field */}
              <FormField
                control={form.control}
                name="dueDate"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs text-muted-foreground uppercase tracking-wide">
                      Due Date
                    </FormLabel>
                    <FormControl>
                      <div className="flex gap-2">
                        <Popover>
                          <PopoverTrigger asChild>
                            <Button
                              variant="outline"
                              className={cn(
                                "w-full sm:w-[280px] justify-start text-left font-normal bg-background",
                                !field.value && "text-muted-foreground"
                              )}
                            >
                              <CalendarIcon className="mr-2 h-4 w-4" />
                              {field.value ? formatDateWithTime(field.value) : <span>Pick a date</span>}
                            </Button>
                          </PopoverTrigger>
                          <PopoverContent className="w-auto p-0" align="start">
                            <Calendar
                              ISOWeek
                              fromDate={new Date()}
                              mode="single"
                              selected={field.value}
                              onSelect={(value) =>
                                field.onChange(value ? addSeconds(value, 23 * 60 * 60 + 59 * 60 + 59) : undefined)
                              }
                              initialFocus
                              defaultMonth={field.value}
                            />
                            <div className="p-3 border-t border-border">
                              <TimePicker setDate={field.onChange} date={field.value} />
                            </div>
                          </PopoverContent>
                        </Popover>
                        {field.value && (
                          <Button
                            type="button"
                            onClick={() => field.onChange(undefined)}
                            variant="ghost"
                            size="icon"
                            className="shrink-0 text-muted-foreground hover:text-destructive"
                          >
                            <X className="h-4 w-4" />
                          </Button>
                        )}
                      </div>
                    </FormControl>
                    <FormDescription className="text-xs">
                      Optional deadline for this assignment
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Submit Button */}
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
                  {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Save Changes
                </Button>
              </div>

              {isError && (
                <Alert variant="destructive">
                  <AlertDescription>
                    An error occurred while saving. Please try again.
                  </AlertDescription>
                </Alert>
              )}
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
