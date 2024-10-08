import { createFileRoute, Link, redirect, useNavigate } from "@tanstack/react-router";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import { Button } from "@/components/ui/button.tsx";
import { AlertCircle, Calendar as CalendarIcon, Check, ChevronsUpDown, Loader2 } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert.tsx";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { CreateAssignmentForm, createAssignmentFormSchema } from "@/types/assignments.ts";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover.tsx";
import { cn, formatDateWithTime, getUUIDFromLocation, isOwner } from "@/lib/utils.ts";
import { addSeconds } from "date-fns";
import { Calendar } from "@/components/ui/calendar.tsx";
import { useState } from "react";
import { Loader } from "@/components/loader.tsx";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem } from "@/components/ui/command.tsx";
import { Header } from "@/components/header";
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
  const [open, setOpen] = useState(false);

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
            <BreadcrumbPage>Create Assignment</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <Header title="Create Assignment" />
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input placeholder="Programming Assignment" {...field} />
                </FormControl>
                <FormDescription>This is your Assignment name.</FormDescription>
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
                  <Textarea placeholder="This is my awesome ..." className="resize-none" {...field} />
                </FormControl>
                <FormDescription>The description of your classroom</FormDescription>
                <FormMessage />
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
                      onClick={() => field.onChange(undefined, { shouldValidate: false })}
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

          <FormField
            control={form.control}
            name={"templateProjectId"}
            render={({ field }) => (
              <FormItem>
                <FormLabel>Template Project</FormLabel>
                <FormControl>
                  <div>
                    <Popover open={open} onOpenChange={setOpen}>
                      <PopoverTrigger asChild>
                        <Button
                          variant="outline"
                          role="combobox"
                          aria-expanded={open}
                          className="w-[200px] justify-between"
                        >
                          {field.value
                            ? templateProjects.find((template) => template.id === field.value)?.name
                            : "Select Template..."}
                          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="w-[200px] p-0">
                        <Command>
                          <CommandInput placeholder="Search project..." />
                          <CommandEmpty>No projects found.</CommandEmpty>
                          <CommandGroup>
                            {templateProjects.map((template) => (
                              <CommandItem
                                key={template.id}
                                value={template.name}
                                onSelect={() => {
                                  field.onChange(template.id);
                                  setOpen(false);
                                }}
                              >
                                <Check
                                  className={cn(
                                    "mr-2 h-4 w-4",
                                    field.value === template.id ? "opacity-100" : "opacity-0",
                                  )}
                                />
                                {template.name}
                              </CommandItem>
                            ))}
                          </CommandGroup>
                        </Command>
                      </PopoverContent>
                    </Popover>
                  </div>
                </FormControl>
                <FormDescription>
                  This is the Template Repository of your assignment.
                  <br /> It must be an public GitLab Project under your own account.
                </FormDescription>
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
