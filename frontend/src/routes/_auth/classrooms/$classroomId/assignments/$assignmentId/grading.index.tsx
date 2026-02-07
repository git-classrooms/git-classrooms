import { assignmentQueryOptions } from "@/api/assignment";
import { classroomQueryOptions } from "@/api/classroom";
import { assignmentGradingRubricsQueryOptions, useGradeProject, useStartAutoGrading } from "@/api/grading";
import { assignmentProjectsQueryOptions } from "@/api/project";
import { assignmentReportQueryOptions } from "@/api/report";
import { Header } from "@/components/header";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Button } from "@/components/ui/button";
import { cn, isModerator } from "@/lib/utils";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { AlertCircle, Bot, Download, Loader2, SearchCheck, SearchCode } from "lucide-react";
import { useMemo, useRef } from "react";
import { toast } from "sonner";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import {
  Assignment,
  ManualGradingRubric,
  ProjectResponse,
  ReportApiAxiosParamCreator,
  UtilsReportDataItem,
} from "@/swagger-client";
import { Status } from "@/types/projects";
import { useFieldArray, useForm } from "react-hook-form";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AutosizeTextarea } from "@/components/ui/autosize-textarea";
import { DataTable, DataTableColumnHeader, ColumnDef } from "@/components/ui/data-table";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/$assignmentId/grading/")({
  beforeLoad: async ({ context: { queryClient }, params: { classroomId, assignmentId } }) => {
    const userClassroom = await queryClient.ensureQueryData(classroomQueryOptions(classroomId));
    if (!isModerator(userClassroom)) {
      throw redirect({
        to: "/classrooms/$classroomId/assignments/$assignmentId",
        params: { classroomId, assignmentId },
        replace: true,
      });
    }
  },
  loader: async ({ context: { queryClient }, params: { classroomId, assignmentId } }) => {
    const assignment = await queryClient.ensureQueryData(assignmentQueryOptions(classroomId, assignmentId));
    const rubrics = await queryClient.ensureQueryData(assignmentGradingRubricsQueryOptions(classroomId, assignmentId));
    const report = await queryClient.ensureQueryData(assignmentReportQueryOptions(classroomId, assignmentId));
    const projects = await queryClient.ensureQueryData(assignmentProjectsQueryOptions(classroomId, assignmentId));
    const { url: reportDownloadUrl } = await ReportApiAxiosParamCreator().getClassroomAssignmentReport(
      classroomId,
      assignmentId,
    );
    return { reportDownloadUrl, assignment, report, projects, rubrics };
  },
  component: GradingIndex,
});

function GradingIndex() {
  const { classroomId, assignmentId } = Route.useParams();
  const { reportDownloadUrl } = Route.useLoaderData();
  const { data: classroom } = useSuspenseQuery(classroomQueryOptions(classroomId));
  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));

  const { mutateAsync, isPending } = useStartAutoGrading({
    classroomId,
    assignmentId,
    onError: (error) => toast.error(error.message),
  });

  const onClick = async () => {
    await mutateAsync();
  };

  return (
    <>
      <Breadcrumb className="mb-5">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms/$classroomId" search={{ tab: "assignments" }} params={{ classroomId }}>
                {classroom.classroom.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms/$classroomId/assignments/$assignmentId" params={{ classroomId, assignmentId }}>
                {assignment.name}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Grading</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <div className="md:flex justify-between gap-1 mb-4">
        <Header className="grow" title="Grading" subtitle={`Overview of the current grading of ${assignment.name}`} />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          <Tooltip delayDuration={0}>
            <TooltipTrigger asChild>
              <Button variant="secondary" asChild size="sm" title="Download report">
                <a href={reportDownloadUrl} target="_blank" referrerPolicy="no-referrer">
                  <Download className="mr-2 h-4 w-4" /> Download report
                </a>
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p>Download a grading report as CSV-File for this assignment.</p>
            </TooltipContent>
          </Tooltip>
          <Tooltip delayDuration={0}>
            <TooltipTrigger asChild>
              <Button
                variant="secondary"
                onClick={onClick}
                size="sm"
                title="Test-driven grading"
                disabled={!assignment.gradingJUnitAutoGradingActive || isPending}
              >
                {isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <>
                    <Bot className="mr-2 h-4 w-4" /> Refresh test-driven grading
                  </>
                )}
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p>Triggers the start of a call for the results of test-driven grading.</p>
            </TooltipContent>
          </Tooltip>
        </div>
      </div>
      <GradingOverview classroomId={classroomId} assignmentId={assignmentId} />
    </>
  );
}

function GradingOverview({ assignmentId, classroomId }: { classroomId: string; assignmentId: string }) {
  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  const { data: gradingResults } = useSuspenseQuery(assignmentReportQueryOptions(classroomId, assignmentId));
  const { data: projects } = useSuspenseQuery(assignmentProjectsQueryOptions(classroomId, assignmentId));
  const { data: rubrics } = useSuspenseQuery(assignmentGradingRubricsQueryOptions(classroomId, assignmentId));

  const zippedProjects = useMemo(
    () =>
      projects
        .filter((p) => p.projectStatus === Status.Accepted)
        .map((project) => ({
          ...project,
          gradingResult: gradingResults.find((result) => result.projectId === project.id),
        })),
    [projects, gradingResults],
  );

  return (
    <div>
      <AssignmentProjectTable assignment={assignment} zippedProjects={zippedProjects} rubrics={rubrics} />
    </div>
  );
}

type ZippedProject = ProjectResponse & { gradingResult?: UtilsReportDataItem };

function AssignmentProjectTable({
  assignment,
  zippedProjects,
  rubrics,
}: {
  assignment: Assignment;
  zippedProjects: ZippedProject[];
  rubrics: ManualGradingRubric[];
}) {
  const maxManualScore = rubrics.reduce((acc, e) => acc + e.maxScore, 0);

  const columns = useMemo<ColumnDef<ZippedProject, unknown>[]>(() => {
    const baseColumns: ColumnDef<ZippedProject, unknown>[] = [
      {
        accessorKey: "team.name",
        header: ({ column }) => <DataTableColumnHeader column={column} title="Name" />,
        cell: ({ row }) => <span className="font-medium">{row.original.team.name}</span>,
      },
      {
        id: "status",
        accessorFn: (row) => {
          const alreadyGraded =
            Object.keys(row.gradingResult?.rubricResults ?? {}).length === rubrics.length &&
            (row.gradingResult?.autogradingMaxScore === 0 || row.gradingResult?.autogradingScore !== 0);
          return alreadyGraded ? "Graded" : "Not graded";
        },
        header: ({ column }) => <DataTableColumnHeader column={column} title="Status" />,
        cell: ({ row }) => {
          const alreadyGraded =
            Object.keys(row.original.gradingResult?.rubricResults ?? {}).length === rubrics.length &&
            (row.original.gradingResult?.autogradingMaxScore === 0 ||
              row.original.gradingResult?.autogradingScore !== 0);
          return (
            <div className="flex pl-1 gap-3 items-center">
              <span className="relative flex h-3 w-3">
                <span
                  className={cn(
                    "animate-ping absolute inline-flex h-full w-full rounded-full opacity-75",
                    alreadyGraded ? "bg-emerald-400" : "bg-gray-400",
                  )}
                ></span>
                <span
                  className={cn(
                    "relative inline-flex rounded-full h-3 w-3",
                    alreadyGraded ? "bg-emerald-500" : "bg-gray-500",
                  )}
                ></span>
              </span>
              {alreadyGraded ? "Graded" : "Not graded"}
            </div>
          );
        },
        filterFn: (row, id, value) => {
          return value.includes(row.getValue(id));
        },
      },
      {
        id: "manualScore",
        accessorFn: (row) => row.gradingManualResults?.reduce((acc, e) => acc + (e.score || 0), 0) ?? 0,
        header: ({ column }) => <DataTableColumnHeader column={column} title="Manual grading" />,
        cell: ({ row }) => (
          <span>
            {row.original.gradingManualResults?.reduce((acc, e) => acc + (e.score || 0), 0)}/{maxManualScore}
          </span>
        ),
      },
    ];

    if (assignment.gradingJUnitAutoGradingActive) {
      baseColumns.push({
        id: "autoScore",
        accessorFn: (row) => row.gradingResult?.autogradingScore ?? 0,
        header: ({ column }) => <DataTableColumnHeader column={column} title="Test-driven grading" />,
        cell: ({ row }) => (
          <Tooltip delayDuration={0}>
            <TooltipTrigger>
              <a
                className="flex items-center"
                href={row.original.reportWebUrl}
                target="_blank"
                referrerPolicy="no-referrer"
              >
                {row.original.gradingResult?.autogradingScore ?? 0}/
                {row.original.gradingResult?.autogradingMaxScore ?? 0} <SearchCheck className="ml-1 h-3.5 w-3.5" />
              </a>
            </TooltipTrigger>
            <TooltipContent>
              <p>Open details for test-driven grading</p>
            </TooltipContent>
          </Tooltip>
        ),
      });
    }

    baseColumns.push(
      {
        id: "totalScore",
        accessorFn: (row) => row.gradingResult?.score ?? 0,
        header: ({ column }) => <DataTableColumnHeader column={column} title="Score" />,
        cell: ({ row }) => (
          <span>
            {row.original.gradingResult?.score ?? 0}/{row.original.gradingResult?.maxScore}
          </span>
        ),
      },
      {
        id: "actions",
        cell: ({ row }) => (
          <div className="text-right">
            <DrawerForm zippedProject={row.original} rubrics={rubrics} />
          </div>
        ),
      },
    );

    return baseColumns;
  }, [assignment.gradingJUnitAutoGradingActive, rubrics, maxManualScore]);

  return (
    <DataTable
      columns={columns}
      data={zippedProjects}
      searchKey="team.name"
      searchPlaceholder="Search by team name..."
      showPagination={zippedProjects.length > 10}
      emptyMessage="No projects found."
    />
  );
}

const rowSchema = z
  .object({
    score: z.number().min(0),
    feedback: z.string().optional(),
    rubricId: z.string().uuid(),
  })
  .transform((data) => ({
    ...data,
    id: data.rubricId,
  }));

const formSchema = z.object({
  gradingManualRubrics: z.array(rowSchema),
});

const DrawerForm = ({
  zippedProject,
  rubrics,
}: {
  zippedProject: ProjectResponse & { gradingResult?: UtilsReportDataItem };
  rubrics: ManualGradingRubric[];
}) => {
  const { classroomId, assignmentId } = Route.useParams();

  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      gradingManualRubrics: rubrics.map((rubric) => ({
        rubricId: rubric.id,
        score: zippedProject.gradingResult?.rubricResults?.[rubric.name]?.score ?? 0,
        feedback: zippedProject.gradingResult?.rubricResults?.[rubric.name]?.feedback ?? "",
      })),
    },
  });

  const { mutateAsync, isPending, error } = useGradeProject(classroomId, assignmentId, zippedProject.id);

  const { fields } = useFieldArray({
    control: form.control,
    name: "gradingManualRubrics",
  });

  const closeModalButtonRef = useRef<HTMLButtonElement>(null);

  const handleSubmit = async (data: z.infer<typeof formSchema>) => {
    await mutateAsync(data);
    closeModalButtonRef.current?.click();
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Grade</Button>
      </DialogTrigger>
      <DialogContent className="lg:min-w-[600px] overflow-y-auto max-h-screen">
        <DialogHeader>
          <DialogTitle>Grade project</DialogTitle>
          <DialogDescription>Overview for grading of {zippedProject.team.name}.</DialogDescription>
        </DialogHeader>
        <Button className="mb-8" variant="default" asChild size="sm" title="Grading">
          <a href={zippedProject.webUrl} target="_blank" referrerPolicy="no-referrer">
            <SearchCode className="mr-2 h-4 w-4" />
            Go to code
          </a>
        </Button>
        {assignment.gradingJUnitAutoGradingActive && (
          <div className="md:flex justify-between gap-1 mb-2">
            <div className="grow mb-2">
              <h3 className="text-l font-bold">Test-based grading</h3>
              <span className="text-s font-light text-muted-foreground">The current result of test-based grading</span>
              <p className="mt-2 mb-2">
                {zippedProject.gradingResult?.autogradingScore ?? 0}/
                {zippedProject.gradingResult?.autogradingMaxScore ?? 0} points.
              </p>
            </div>
            <div className="flex-none">
              <Button className="w-full" variant="secondary" asChild size="sm" title="Details of test-based grading">
                <Link
                  to="/classrooms/$classroomId/assignments/$assignmentId/grading"
                  params={{ classroomId, assignmentId }}
                >
                  <Bot className="mr-2 h-4 w-4" />
                  Details
                </Link>
              </Button>
            </div>
          </div>
        )}
        <div className="gap-1">
          <h3 className="text-l font-bold">Manual grading</h3>
          <span className="text-s font-light text-muted-foreground">
            View and adjust the manual grades for this project
          </span>
        </div>
        {fields.length === 0 ? (
          <div className="text-center text-muted-foreground">No manual grading rubrics available.</div>
        ) : (
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className="w-full flex-col">
              {fields.map((field, index) => {
                const rubric = rubrics.find((e) => e.id === field.rubricId)!;
                return (
                  <div
                    key={field.id}
                    className="w-full grid grid-cols-1 md:grid-cols-[1fr_4fr] gap-2 rounded-md border p-4 mb-4"
                  >
                    <div className="gap-1 md:col-span-2">
                      <div className="flex justify-between items-center">
                        <h4 className="text-sm font-bold">{rubric.name}</h4>
                        <span className="text-sm text-foreground">Max Score: {rubric.maxScore}</span>
                      </div>
                      <span className="text-sm font-light text-muted-foreground">{rubric.description}</span>
                    </div>

                    <FormField
                      control={form.control}
                      name={`gradingManualRubrics.${index}.rubricId`}
                      render={({ field }) => <input value={field.value} readOnly hidden />}
                    />

                    <FormField
                      control={form.control}
                      name={`gradingManualRubrics.${index}.score`}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Points</FormLabel>
                          <FormControl>
                            <Input
                              type="number"
                              min={0}
                              step={1}
                              disabled={isPending}
                              {...field}
                              onChange={(e) => {
                                const value = e.target.value;
                                const numberValue = value ? Number(value) : "";
                                field.onChange(numberValue);
                              }}
                              className="text-base rounded-r"
                            />
                          </FormControl>
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name={`gradingManualRubrics.${index}.feedback`}
                      render={({ field }) => (
                        <FormItem className="grow">
                          <FormLabel>Feedback</FormLabel>
                          <FormControl>
                            <AutosizeTextarea
                              minHeight={1}
                              disabled={isPending}
                              {...field}
                              className={"rounded-r mt-5"}
                            />
                          </FormControl>
                        </FormItem>
                      )}
                    />
                  </div>
                );
              })}

              <Button>Save</Button>
              {error && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>Error</AlertTitle>
                  <AlertDescription>{error.message}</AlertDescription>
                </Alert>
              )}
            </form>
          </Form>
        )}
        <DialogClose ref={closeModalButtonRef} className="hidden">
          Close
        </DialogClose>
      </DialogContent>
    </Dialog>
  );
};
