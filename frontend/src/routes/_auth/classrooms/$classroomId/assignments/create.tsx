import { createFileRoute, Link, redirect, useNavigate } from "@tanstack/react-router";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import {
  AlertCircle,
  BookOpen,
  Calendar as CalendarIcon,
  ChevronRight,
  Clock,
  FileCode2,
  FolderGit2,
  Loader2,
} from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { CreateAssignmentForm, createAssignmentFormSchema } from "@/types/assignments";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { cn, formatDateWithTime, getUUIDFromLocation, isOwner } from "@/lib/utils";
import { addSeconds } from "date-fns";
import { Calendar } from "@/components/ui/calendar";
import { Loader } from "@/components/loader";
import { useSuspenseQuery } from "@tanstack/react-query";
import { classroomQueryOptions, classroomTemplatesQueryOptions } from "@/api/classroom";
import { useCreateAssignment } from "@/api/assignment";
import { TimePicker } from "@/components/ui/timer-picker";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Card, CardContent } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/create")({
  loader: async ({ context: { queryClient }, params }) => {
    const templateProjects = await queryClient.ensureQueryData(classroomTemplatesQueryOptions(params.classroomId));
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(params.classroomId));
    if (!isOwner(userClassroom)) {
      throw redirect({
        to: "/classrooms/$classroomId",
        search: { tab: "assignments" },
        params,
        replace: true,
      });
    }
    return { templateProjects };
  },
  component: CreateAssignment,
  pendingComponent: Loader,
});

function CreateAssignment() {
  const { classroomId } = Route.useParams();
  const { data: userClassroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const navigate = useNavigate();

  const { data: templateProjects } = useSuspenseQuery(classroomTemplatesQueryOptions(classroomId));

  const { mutateAsync, isError, isPending } = useCreateAssignment(classroomId);
  const form = useForm<CreateAssignmentForm>({
    resolver: zodResolver(createAssignmentFormSchema),
    defaultValues: {
      name: "",
      description: "",
      templateProjectId: 0,
    },
  });

  async function onSubmit(values: CreateAssignmentForm) {
    const location = await mutateAsync({ ...values, dueDate: values.dueDate?.toISOString() });
    const assignmentId = getUUIDFromLocation(location);
    await navigate({
      to: "/classrooms/$classroomId/assignments/$assignmentId",
      params: { classroomId, assignmentId },
    });
  }

  return (
    <div className="max-w-2xl mx-auto">
      {/* Breadcrumb */}
      <Breadcrumb className="mb-6">
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
            <BreadcrumbPage>Create Assignment</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-4 mb-3">
          <div className="w-14 h-14 rounded-xl bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center">
            <FileCode2 className="w-7 h-7 text-primary" />
          </div>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">Create Assignment</h1>
            <p className="text-muted-foreground">Set up a new assignment for your students</p>
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
                    <FormLabel>Assignment Name</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="e.g., Exercise 1: Hello World"
                        className="bg-background"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      A clear, descriptive name for the assignment
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
                        placeholder="Describe the assignment objectives and requirements..."
                        className="resize-none bg-background min-h-[100px]"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      Help students understand what they need to do
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Template Section */}
          <Card className="border-border/50 overflow-hidden">
            <div className="px-5 py-4 border-b border-border/50 bg-muted/30">
              <div className="flex items-center gap-2">
                <FolderGit2 className="w-4 h-4 text-muted-foreground" />
                <h2 className="font-semibold text-sm">Template Repository</h2>
              </div>
            </div>
            <CardContent className="p-5">
              <FormField
                control={form.control}
                name="templateProjectId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Select Template</FormLabel>
                    <Select
                      onValueChange={(value) => field.onChange(Number(value))}
                      value={field.value ? String(field.value) : undefined}
                    >
                      <FormControl>
                        <SelectTrigger className="bg-background">
                          <SelectValue placeholder="Select a template project..." />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {templateProjects.length === 0 ? (
                          <div className="py-6 text-center text-sm text-muted-foreground">
                            No templates found
                          </div>
                        ) : (
                          templateProjects.map((template) => (
                            <SelectItem key={template.id} value={String(template.id)}>
                              <span className="flex items-center gap-2">
                                <FolderGit2 className="w-4 h-4 text-muted-foreground" />
                                {template.name}
                              </span>
                            </SelectItem>
                          ))
                        )}
                      </SelectContent>
                    </Select>
                    <FormDescription>
                      The GitLab project that will be forked for each student/team
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Due Date Section */}
          <Card className="border-border/50 overflow-hidden">
            <div className="px-5 py-4 border-b border-border/50 bg-muted/30">
              <div className="flex items-center gap-2">
                <Clock className="w-4 h-4 text-muted-foreground" />
                <h2 className="font-semibold text-sm">Deadline</h2>
              </div>
            </div>
            <CardContent className="p-5">
              <FormField
                control={form.control}
                name="dueDate"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Due Date</FormLabel>
                    <FormControl>
                      <div className="flex gap-2">
                        <Popover>
                          <PopoverTrigger asChild>
                            <Button
                              variant="outline"
                              className={cn(
                                "flex-1 justify-start text-left font-normal bg-background",
                                !field.value && "text-muted-foreground"
                              )}
                            >
                              <CalendarIcon className="mr-2 h-4 w-4" />
                              {field.value ? formatDateWithTime(field.value) : "No deadline set"}
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
                            variant="outline"
                            size="icon"
                            className="shrink-0"
                          >
                            ✕
                          </Button>
                        )}
                      </div>
                    </FormControl>
                    <FormDescription>
                      Optional: Set a submission deadline for the assignment
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Error Alert */}
          {isError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>The assignment could not be created!</AlertDescription>
            </Alert>
          )}

          {/* Submit Button */}
          <div className="flex items-center justify-end gap-3 pt-4">
            <Button type="button" variant="outline" asChild>
              <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId }}>
                Cancel
              </Link>
            </Button>
            <Button type="submit" variant="glow" disabled={isPending} className="min-w-[160px]">
              {isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <>
                  Create Assignment
                  <ChevronRight className="ml-2 h-4 w-4" />
                </>
              )}
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
}
