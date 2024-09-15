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
import { Calendar as CalendarIcon, Loader2 } from "lucide-react";
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
    <div className="p-2 w-full">
      <div className="flex flex-row justify-between">
        <h1 className="text-xl font-bold">Edit assignment</h1>
      </div>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <FormField
            control={form.control}
            name="name"
            disabled={isAccepted}
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
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
                  />
                </FormControl>
                <FormDescription>This is your Assignment name.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="description"
            disabled={isAccepted}
            render={({ field }) => (
              <FormItem>
                <FormLabel>Description</FormLabel>
                <FormControl>
                  <Textarea placeholder="This is my awesome ..." className="resize-none" {...field} />
                </FormControl>
                <FormDescription>This is the description of your classroom.</FormDescription>
                <FormMessage />
                {isAccepted && (
                  <FormMessage>
                    Name and description cannot be changed once the Assignment has been accepted at least one team
                  </FormMessage>
                )}
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name={"dueDate"}
            render={({ field }) => (
              <FormItem>
                <FormLabel>Due Date</FormLabel>
                <FormControl>
                  <div className="flex gap-2">
                    <Popover>
                      <PopoverTrigger asChild>
                        <Button
                          variant={"outline"}
                          className={cn(
                            "w-[280px] justify-start text-left font-normal",
                            !field.value && "text-muted-foreground",
                          )}
                        >
                          <CalendarIcon className="mr-2 h-4 w-4" />
                          {field.value ? formatDateWithTime(field.value) : <span>Pick a date</span>}
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0">
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
                    <Button
                      type="button"
                      onClick={() => field.onChange(null, { shouldValidate: false })}
                      variant="outline"
                    >
                      Remove
                    </Button>
                  </div>
                </FormControl>
                <FormDescription>This is the due date of your assignment.</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit" disabled={isPending}>
            {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Save"}
          </Button>

          {isError && <div className="text-red-500">An error occurred. Please try again. </div>}
        </form>
      </Form>
    </div>
  );
}
