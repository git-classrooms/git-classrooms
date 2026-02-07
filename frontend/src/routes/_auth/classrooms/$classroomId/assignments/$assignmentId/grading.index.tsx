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
import { Card, CardContent } from "@/components/ui/card";
import { cn, isModerator } from "@/lib/utils";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import {
  AlertCircle,
  Award,
  Bot,
  CheckCircle2,
  ChevronRight,
  Clock,
  Download,
  ExternalLink,
  Loader2,
  Target,
  TrendingUp,
  Users,
} from "lucide-react";
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
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AutosizeTextarea } from "@/components/ui/autosize-textarea";
import { StatusBadge } from "@/components/ui/status-badge";
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";

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
    <div className="animate-in fade-in slide-in-from-bottom-2 duration-300">
      <Breadcrumb className="mb-6">
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

      <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-4 mb-8">
        <Header className="grow mb-0" title="Grading Dashboard" subtitle={assignment.name} />
        <div className="flex flex-wrap gap-3">
          <Button variant="outline" asChild size="sm">
            <a href={reportDownloadUrl} target="_blank" referrerPolicy="no-referrer">
              <Download className="mr-2 h-4 w-4" />
              Export CSV
            </a>
          </Button>
          {assignment.gradingJUnitAutoGradingActive && (
            <Button variant="glow" onClick={onClick} size="sm" disabled={isPending}>
              {isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <Bot className="mr-2 h-4 w-4" />
              )}
              Refresh Auto-Grading
            </Button>
          )}
        </div>
      </div>

      <GradingOverview classroomId={classroomId} assignmentId={assignmentId} />
    </div>
  );
}

function CircularProgress({ value, max, size = 56, strokeWidth = 4 }: { value: number; max: number; size?: number; strokeWidth?: number }) {
  const radius = (size - strokeWidth) / 2;
  const circumference = radius * 2 * Math.PI;
  const percentage = max > 0 ? (value / max) * 100 : 0;
  const offset = circumference - (percentage / 100) * circumference;

  return (
    <div className="relative" style={{ width: size, height: size }}>
      <svg className="transform -rotate-90" width={size} height={size}>
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth={strokeWidth}
          className="text-muted/30"
        />
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth={strokeWidth}
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          strokeLinecap="round"
          className="text-primary transition-all duration-500 ease-out"
        />
      </svg>
      <div className="absolute inset-0 flex items-center justify-center">
        <span className="text-sm font-mono font-bold">{Math.round(percentage)}%</span>
      </div>
    </div>
  );
}

function ScoreBar({ value, max, variant = "default" }: { value: number; max: number; variant?: "default" | "success" | "warning" }) {
  const percentage = max > 0 ? (value / max) * 100 : 0;
  const colorClass = variant === "success" ? "bg-success" : variant === "warning" ? "bg-warning" : "bg-primary";

  return (
    <div className="flex items-center gap-3 w-full">
      <div className="flex-1 h-2 bg-muted/50 rounded-full overflow-hidden">
        <div
          className={cn("h-full rounded-full transition-all duration-500 ease-out", colorClass)}
          style={{ width: `${percentage}%` }}
        />
      </div>
      <span className="text-sm font-mono text-muted-foreground w-16 text-right">
        {value}/{max}
      </span>
    </div>
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

  const stats = useMemo(() => {
    const totalProjects = zippedProjects.length;
    const gradedProjects = zippedProjects.filter((p) => {
      const rubricCount = Object.keys(p.gradingResult?.rubricResults ?? {}).length;
      const hasAutoGrading = p.gradingResult?.autogradingMaxScore === 0 || p.gradingResult?.autogradingScore !== 0;
      return rubricCount === rubrics.length && hasAutoGrading;
    }).length;

    const totalScore = zippedProjects.reduce((acc, p) => acc + (p.gradingResult?.score ?? 0), 0);
    const maxTotalScore = zippedProjects.reduce((acc, p) => acc + (p.gradingResult?.maxScore ?? 0), 0);
    const avgPercentage = maxTotalScore > 0 ? (totalScore / maxTotalScore) * 100 : 0;

    return {
      totalProjects,
      gradedProjects,
      pendingProjects: totalProjects - gradedProjects,
      avgPercentage,
      totalScore,
      maxTotalScore,
    };
  }, [zippedProjects, rubrics.length]);

  return (
    <div className="space-y-6">
      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="border-border/50 bg-gradient-to-br from-card to-card/50">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">Total Projects</p>
                <p className="text-3xl font-mono font-bold">{stats.totalProjects}</p>
              </div>
              <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center">
                <Users className="w-6 h-6 text-primary" />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/50 bg-gradient-to-br from-card to-card/50">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">Grading Progress</p>
                <p className="text-3xl font-mono font-bold">
                  {stats.gradedProjects}<span className="text-lg text-muted-foreground">/{stats.totalProjects}</span>
                </p>
              </div>
              <CircularProgress value={stats.gradedProjects} max={stats.totalProjects} />
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/50 bg-gradient-to-br from-card to-card/50">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">Pending Review</p>
                <p className="text-3xl font-mono font-bold">{stats.pendingProjects}</p>
              </div>
              <div className="w-12 h-12 rounded-lg bg-warning/10 flex items-center justify-center">
                <Clock className="w-6 h-6 text-warning" />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/50 bg-gradient-to-br from-card to-card/50">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">Class Average</p>
                <p className="text-3xl font-mono font-bold">{stats.avgPercentage.toFixed(1)}%</p>
              </div>
              <div className="w-12 h-12 rounded-lg bg-success/10 flex items-center justify-center">
                <TrendingUp className="w-6 h-6 text-success" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Projects List */}
      <Card className="border-border/50">
        <CardContent className="p-0">
          <div className="p-4 border-b border-border/50">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center">
                  <Award className="w-5 h-5 text-primary" />
                </div>
                <div>
                  <h2 className="text-lg font-semibold font-mono">Project Grades</h2>
                  <p className="text-sm text-muted-foreground">Review and grade individual submissions</p>
                </div>
              </div>
            </div>
          </div>

          {zippedProjects.length === 0 ? (
            <div className="p-12 text-center">
              <Target className="w-12 h-12 mx-auto text-muted-foreground/50 mb-4" />
              <p className="text-muted-foreground">No accepted projects to grade yet.</p>
            </div>
          ) : (
            <div className="divide-y divide-border/50">
              {zippedProjects.map((project, index) => (
                <ProjectRow
                  key={project.id}
                  project={project}
                  rubrics={rubrics}
                  assignment={assignment}
                  index={index}
                />
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

type ZippedProject = ProjectResponse & { gradingResult?: UtilsReportDataItem };

function ProjectRow({
  project,
  rubrics,
  assignment,
  index,
}: {
  project: ZippedProject;
  rubrics: ManualGradingRubric[];
  assignment: Assignment;
  index: number;
}) {
  const maxManualScore = rubrics.reduce((acc, e) => acc + e.maxScore, 0);
  const manualScore = project.gradingManualResults?.reduce((acc, e) => acc + (e.score || 0), 0) ?? 0;
  const autoScore = project.gradingResult?.autogradingScore ?? 0;
  const autoMaxScore = project.gradingResult?.autogradingMaxScore ?? 0;
  const totalScore = project.gradingResult?.score ?? 0;
  const totalMaxScore = project.gradingResult?.maxScore ?? 0;

  const alreadyGraded =
    Object.keys(project.gradingResult?.rubricResults ?? {}).length === rubrics.length &&
    (project.gradingResult?.autogradingMaxScore === 0 || project.gradingResult?.autogradingScore !== 0);

  return (
    <div
      className="p-4 hover:bg-muted/30 transition-colors animate-in fade-in slide-in-from-bottom-1"
      style={{ animationDelay: `${index * 50}ms` }}
    >
      <div className="flex flex-col lg:flex-row lg:items-center gap-4">
        {/* Team Info */}
        <div className="flex items-center gap-3 lg:w-1/4 min-w-0">
          <div className="w-10 h-10 rounded-lg bg-muted flex items-center justify-center shrink-0">
            <span className="text-sm font-mono font-bold text-muted-foreground">
              {project.team.name.substring(0, 2).toUpperCase()}
            </span>
          </div>
          <div className="min-w-0">
            <p className="font-medium truncate">{project.team.name}</p>
            <div className="flex items-center gap-2">
              <StatusBadge variant={alreadyGraded ? "success" : "warning"} size="sm" showDot>
                {alreadyGraded ? "Graded" : "Pending"}
              </StatusBadge>
            </div>
          </div>
        </div>

        {/* Scores */}
        <div className="flex-1 grid grid-cols-1 md:grid-cols-3 gap-4">
          {/* Manual Score */}
          <div>
            <p className="text-xs text-muted-foreground mb-1.5">Manual</p>
            <ScoreBar value={manualScore} max={maxManualScore} />
          </div>

          {/* Auto Score */}
          {assignment.gradingJUnitAutoGradingActive && (
            <div>
              <p className="text-xs text-muted-foreground mb-1.5">Auto-Grading</p>
              <ScoreBar value={autoScore} max={autoMaxScore} variant="success" />
            </div>
          )}

          {/* Total Score */}
          <div>
            <p className="text-xs text-muted-foreground mb-1.5">Total</p>
            <div className="flex items-center gap-3">
              <div className="flex-1 h-2 bg-muted/50 rounded-full overflow-hidden">
                <div
                  className="h-full rounded-full bg-gradient-to-r from-primary to-primary/70 transition-all duration-500"
                  style={{ width: `${totalMaxScore > 0 ? (totalScore / totalMaxScore) * 100 : 0}%` }}
                />
              </div>
              <span className="text-sm font-mono font-bold w-16 text-right">
                {totalScore}/{totalMaxScore}
              </span>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-2 lg:w-auto">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="ghost" size="icon" asChild className="h-9 w-9">
                <a href={project.webUrl} target="_blank" rel="noopener noreferrer">
                  <ExternalLink className="h-4 w-4" />
                </a>
              </Button>
            </TooltipTrigger>
            <TooltipContent>View code on GitLab</TooltipContent>
          </Tooltip>

          <GradingSheet project={project} rubrics={rubrics} assignment={assignment} />
        </div>
      </div>
    </div>
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

function GradingSheet({
  project,
  rubrics,
  assignment,
}: {
  project: ZippedProject;
  rubrics: ManualGradingRubric[];
  assignment: Assignment;
}) {
  const { classroomId, assignmentId } = Route.useParams();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      gradingManualRubrics: rubrics.map((rubric) => ({
        rubricId: rubric.id,
        score: project.gradingResult?.rubricResults?.[rubric.name]?.score ?? 0,
        feedback: project.gradingResult?.rubricResults?.[rubric.name]?.feedback ?? "",
      })),
    },
  });

  const { mutateAsync, isPending, error } = useGradeProject(classroomId, assignmentId, project.id);

  const { fields } = useFieldArray({
    control: form.control,
    name: "gradingManualRubrics",
  });

  const closeModalButtonRef = useRef<HTMLButtonElement>(null);

  const handleSubmit = async (data: z.infer<typeof formSchema>) => {
    await mutateAsync(data);
    toast.success("Grades saved successfully");
    closeModalButtonRef.current?.click();
  };

  const totalManualScore = fields.reduce((sum, _, index) => {
    return sum + (form.watch(`gradingManualRubrics.${index}.score`) || 0);
  }, 0);

  const maxManualScore = rubrics.reduce((acc, e) => acc + e.maxScore, 0);

  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button variant="default" size="sm" className="gap-2">
          Grade
          <ChevronRight className="h-4 w-4" />
        </Button>
      </SheetTrigger>
      <SheetContent className="w-full sm:max-w-xl overflow-y-auto">
        <SheetHeader className="mb-6">
          <SheetTitle className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
              <Award className="w-5 h-5 text-primary" />
            </div>
            <div>
              <span className="block">{project.team.name}</span>
              <span className="text-sm font-normal text-muted-foreground">Grading Panel</span>
            </div>
          </SheetTitle>
          <SheetDescription>
            Review the submission and assign scores for each rubric criterion.
          </SheetDescription>
        </SheetHeader>

        {/* Quick Actions */}
        <div className="flex gap-2 mb-6">
          <Button variant="outline" size="sm" asChild className="flex-1">
            <a href={project.webUrl} target="_blank" rel="noopener noreferrer">
              <ExternalLink className="mr-2 h-4 w-4" />
              View Code
            </a>
          </Button>
          {assignment.gradingJUnitAutoGradingActive && project.reportWebUrl && (
            <Button variant="outline" size="sm" asChild className="flex-1">
              <a href={project.reportWebUrl} target="_blank" rel="noopener noreferrer">
                <Bot className="mr-2 h-4 w-4" />
                Test Results
              </a>
            </Button>
          )}
        </div>

        {/* Auto-Grading Section */}
        {assignment.gradingJUnitAutoGradingActive && (
          <div className="mb-6 p-4 rounded-lg bg-muted/30 border border-border/50">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2">
                <Bot className="w-4 h-4 text-success" />
                <span className="font-medium text-sm">Automated Test Results</span>
              </div>
              <span className="font-mono text-sm font-bold">
                {project.gradingResult?.autogradingScore ?? 0}/{project.gradingResult?.autogradingMaxScore ?? 0}
              </span>
            </div>
            <div className="h-2 bg-muted/50 rounded-full overflow-hidden">
              <div
                className="h-full rounded-full bg-success transition-all duration-500"
                style={{
                  width: `${
                    (project.gradingResult?.autogradingMaxScore ?? 0) > 0
                      ? ((project.gradingResult?.autogradingScore ?? 0) / (project.gradingResult?.autogradingMaxScore ?? 1)) * 100
                      : 0
                  }%`,
                }}
              />
            </div>
          </div>
        )}

        {/* Manual Grading Form */}
        <div className="mb-4">
          <div className="flex items-center justify-between mb-3">
            <h3 className="font-semibold flex items-center gap-2">
              <Target className="w-4 h-4 text-primary" />
              Manual Grading
            </h3>
            <span className="text-sm font-mono text-muted-foreground">
              {totalManualScore}/{maxManualScore} pts
            </span>
          </div>
        </div>

        {fields.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            <Target className="w-8 h-8 mx-auto mb-2 opacity-50" />
            <p>No manual grading rubrics configured.</p>
          </div>
        ) : (
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
              {fields.map((field, index) => {
                const rubric = rubrics.find((e) => e.id === field.rubricId)!;
                const currentScore = form.watch(`gradingManualRubrics.${index}.score`) || 0;
                const scorePercentage = rubric.maxScore > 0 ? (currentScore / rubric.maxScore) * 100 : 0;

                return (
                  <div
                    key={field.id}
                    className="p-4 rounded-lg border border-border/50 bg-card/50 space-y-3"
                  >
                    <div className="flex items-start justify-between gap-2">
                      <div className="flex-1 min-w-0">
                        <h4 className="font-medium text-sm">{rubric.name}</h4>
                        {rubric.description && (
                          <p className="text-xs text-muted-foreground mt-0.5">{rubric.description}</p>
                        )}
                      </div>
                      <span className="text-xs font-mono text-muted-foreground shrink-0">
                        max {rubric.maxScore} pts
                      </span>
                    </div>

                    {/* Score Progress */}
                    <div className="h-1.5 bg-muted/50 rounded-full overflow-hidden">
                      <div
                        className={cn(
                          "h-full rounded-full transition-all duration-300",
                          scorePercentage >= 80
                            ? "bg-success"
                            : scorePercentage >= 50
                              ? "bg-warning"
                              : scorePercentage > 0
                                ? "bg-destructive"
                                : "bg-muted"
                        )}
                        style={{ width: `${scorePercentage}%` }}
                      />
                    </div>

                    <FormField
                      control={form.control}
                      name={`gradingManualRubrics.${index}.rubricId`}
                      render={({ field }) => <input value={field.value} readOnly hidden />}
                    />

                    <div className="grid grid-cols-[100px_1fr] gap-3">
                      <FormField
                        control={form.control}
                        name={`gradingManualRubrics.${index}.score`}
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel className="text-xs">Score</FormLabel>
                            <FormControl>
                              <Input
                                type="number"
                                min={0}
                                max={rubric.maxScore}
                                step={1}
                                disabled={isPending}
                                {...field}
                                onChange={(e) => {
                                  const value = e.target.value;
                                  const num = value ? Number(value) : 0;
                                  field.onChange(Math.min(num, rubric.maxScore));
                                }}
                                className="h-9 text-center font-mono"
                              />
                            </FormControl>
                          </FormItem>
                        )}
                      />

                      <FormField
                        control={form.control}
                        name={`gradingManualRubrics.${index}.feedback`}
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel className="text-xs">Feedback</FormLabel>
                            <FormControl>
                              <AutosizeTextarea
                                minHeight={36}
                                maxHeight={120}
                                disabled={isPending}
                                placeholder="Optional feedback..."
                                {...field}
                                className="text-sm resize-none"
                              />
                            </FormControl>
                          </FormItem>
                        )}
                      />
                    </div>
                  </div>
                );
              })}

              {error && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>Error</AlertTitle>
                  <AlertDescription>{error.message}</AlertDescription>
                </Alert>
              )}

              {/* Summary & Submit */}
              <div className="pt-4 border-t border-border/50 space-y-4">
                <div className="flex items-center justify-between p-3 rounded-lg bg-muted/30">
                  <span className="text-sm font-medium">Total Score</span>
                  <div className="flex items-center gap-2">
                    <span className="text-2xl font-mono font-bold">
                      {totalManualScore + (project.gradingResult?.autogradingScore ?? 0)}
                    </span>
                    <span className="text-sm text-muted-foreground">
                      / {maxManualScore + (project.gradingResult?.autogradingMaxScore ?? 0)}
                    </span>
                  </div>
                </div>

                <div className="flex gap-3">
                  <SheetClose ref={closeModalButtonRef} asChild>
                    <Button type="button" variant="outline" className="flex-1">
                      Cancel
                    </Button>
                  </SheetClose>
                  <Button type="submit" variant="glow" className="flex-1" disabled={isPending}>
                    {isPending ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Saving...
                      </>
                    ) : (
                      <>
                        <CheckCircle2 className="mr-2 h-4 w-4" />
                        Save Grades
                      </>
                    )}
                  </Button>
                </div>
              </div>
            </form>
          </Form>
        )}
      </SheetContent>
    </Sheet>
  );
}
