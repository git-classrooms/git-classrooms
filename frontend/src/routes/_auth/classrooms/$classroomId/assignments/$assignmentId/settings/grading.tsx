import { assignmentQueryOptions } from "@/api/assignment";
import {
  assignmentGradingRubricsQueryOptions,
  assignmentTestsQueryOptions,
  classroomGradingRubricsQueryOptions,
  useUpdateAssignmentRubrics,
  useUpdateAssignmentTests,
} from "@/api/grading";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import {
  AlertCircle,
  Beaker,
  BookOpen,
  BookOpenCheck,
  CheckCircle2,
  Code2,
  Loader2,
} from "lucide-react";
import { Suspense, useMemo, useState } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { StatusBadge } from "@/components/ui/status-badge";
import { Markdown } from "@/components/ui/markdown";
import { useTranslation } from "react-i18next";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/assignments/$assignmentId/settings/grading")({
  loader: async ({ params: { classroomId, assignmentId }, context: { queryClient } }) => {
    const assignment = await queryClient.ensureQueryData(assignmentQueryOptions(classroomId, assignmentId));
    const rubrics = await queryClient.ensureQueryData(classroomGradingRubricsQueryOptions(classroomId));
    const assignmentRubrics = await queryClient.ensureQueryData(
      assignmentGradingRubricsQueryOptions(classroomId, assignmentId),
    );
    return { assignment, rubrics, assignmentRubrics };
  },
  component: Grading,
});

function Grading() {
  const { classroomId, assignmentId } = Route.useParams();

  return (
    <div className="space-y-8">
      <Suspense fallback={<TestsFormSkeleton />}>
        <TestsForm classroomId={classroomId} assignmentId={assignmentId} />
      </Suspense>
      <RubricForm classroomId={classroomId} assignmentId={assignmentId} />
    </div>
  );
}

function TestsFormSkeleton() {
  return (
    <section>
      <div className="flex items-center gap-3 mb-4">
        <Skeleton className="w-10 h-10 rounded-lg" />
        <div className="space-y-2">
          <Skeleton className="h-5 w-40" />
          <Skeleton className="h-4 w-60" />
        </div>
      </div>
      <Card className="border-border/50">
        <CardContent className="p-4">
          <Skeleton className="h-32 w-full" />
        </CardContent>
      </Card>
    </section>
  );
}

const testsFormSchema = z.object({
  junitAutoGradingActive: z.boolean(),
  assignmentTests: z.array(
    z.object({
      name: z.string(),
      score: z.number(),
      active: z.boolean(),
    }),
  ),
});

const TestsForm = (props: { classroomId: string; assignmentId: string }) => {
  const { t } = useTranslation("assignment");
  const { t: tc } = useTranslation("common");
  const { classroomId, assignmentId } = props;

  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  const { data: tests } = useSuspenseQuery(assignmentTestsQueryOptions(classroomId, assignmentId));

  const assignmentTests = useMemo<z.infer<typeof testsFormSchema>["assignmentTests"]>(
    () =>
      tests.report.map((test) => {
        const selected = tests.selectedTests.find((selectedTest) => selectedTest.name === test.name);
        return {
          name: test.name,
          score: selected?.score ?? 1,
          active: !!selected,
        };
      }),
    [tests],
  );

  const { mutateAsync, isPending, error } = useUpdateAssignmentTests(classroomId, assignmentId);

  const form = useForm<z.infer<typeof testsFormSchema>>({
    resolver: zodResolver(testsFormSchema),
    mode: "onBlur",
    reValidateMode: "onChange",
    defaultValues: { junitAutoGradingActive: assignment.gradingJUnitAutoGradingActive, assignmentTests },
  });

  const { fields } = useFieldArray({
    control: form.control,
    name: "assignmentTests",
  });

  const onSubmit = async (data: z.infer<typeof testsFormSchema>) => {
    await mutateAsync({
      junitAutoGradingActive: data.junitAutoGradingActive,
      assignmentTests: data.assignmentTests.filter((test) => test.active),
    });
    toast.success(t("settings.tests.success"));
  };

  const [showTests, setShowTests] = useState(assignment.gradingJUnitAutoGradingActive);

  const activeTestsCount = fields.filter((_, index) => form.watch(`assignmentTests.${index}.active`)).length;
  const totalScore = fields.reduce((sum, _, index) => {
    if (form.watch(`assignmentTests.${index}.active`)) {
      return sum + (form.watch(`assignmentTests.${index}.score`) || 0);
    }
    return sum;
  }, 0);

  return (
    <section>
      {/* Header */}
      <div className="flex items-start justify-between gap-4 mb-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center">
            <Beaker className="w-5 h-5 text-primary" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h2 className="text-lg font-semibold font-mono">{t("settings.tests.title")}</h2>
              <StatusBadge variant={showTests ? "success" : "neutral"} size="sm">
                {showTests ? t("settings.tests.enabled") : t("settings.tests.disabled")}
              </StatusBadge>
            </div>
            <p className="text-sm text-muted-foreground">
              {t("settings.tests.subtitle")}
            </p>
          </div>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <Card className="border-border/50">
            <CardContent className="p-4">
              {/* Enable Toggle */}
              <FormField
                control={form.control}
                name="junitAutoGradingActive"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between p-3 rounded-lg bg-muted/30 border border-border/50 mb-4">
                    <div className="flex items-center gap-3">
                      <FormControl>
                        <Checkbox
                          disabled={!tests.activatible || isPending}
                          checked={field.value}
                          onCheckedChange={(state) => {
                            field.onChange(state);
                            setShowTests(state !== "indeterminate" && state);
                          }}
                          className="data-[state=checked]:bg-primary data-[state=checked]:border-primary"
                        />
                      </FormControl>
                      <FormLabel className="font-medium cursor-pointer">
                        {t("settings.tests.enableAutoGrading")}
                      </FormLabel>
                    </div>
                    {showTests && tests.report.length > 0 && (
                      <div className="flex items-center gap-4 text-sm text-muted-foreground">
                        <span>{t("settings.tests.testsSelected", { count: activeTestsCount })}</span>
                        <span className="font-mono font-semibold text-foreground">{totalScore} pts</span>
                      </div>
                    )}
                  </FormItem>
                )}
              />

              {/* No Tests Available */}
              {(!tests.activatible || tests.report.length === 0) && (
                <div className="space-y-4">
                  <div className="p-4 rounded-lg border border-dashed border-border/50 bg-muted/20 text-center">
                    <Code2 className="w-8 h-8 mx-auto text-muted-foreground mb-2" />
                    <p className="text-sm text-muted-foreground mb-2">
                      {t("settings.tests.noTestsFound")}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {t("settings.tests.noTestsHelp")} <code className="font-mono bg-muted px-1 rounded">.gitlab-ci.yml</code> {t("settings.tests.noTestsHelp2")}
                    </p>
                  </div>

                  <Card className="border-border/50 bg-muted/10">
                    <CardContent className="p-4">
                      <div className="flex items-center gap-2 mb-3">
                        <BookOpenCheck className="w-4 h-4 text-muted-foreground" />
                        <span className="text-sm font-medium">{t("settings.tests.exampleTitle")}</span>
                      </div>
                      <pre className="text-xs font-mono p-3 rounded-lg bg-background border border-border/50 overflow-x-auto">
                        {tests.example.example}
                      </pre>
                    </CardContent>
                  </Card>
                </div>
              )}

              {/* Test List */}
              {showTests && tests.report.length > 0 && (
                <div className="space-y-2">
                  {fields.map((field, index) => {
                    const test = tests.report.find((t) => t.name === field.name)!;
                    const isActive = form.watch(`assignmentTests.${index}.active`);

                    return (
                      <div
                        key={field.id}
                        className={cn(
                          "flex items-center gap-4 p-3 rounded-lg border transition-all duration-200",
                          isActive
                            ? "border-primary/30 bg-primary/5"
                            : "border-border/50 bg-muted/10 opacity-60"
                        )}
                      >
                        <FormField
                          control={form.control}
                          name={`assignmentTests.${index}.active`}
                          render={({ field }) => (
                            <FormItem className="flex items-center space-y-0">
                              <FormControl>
                                <Checkbox
                                  disabled={isPending}
                                  checked={field.value}
                                  onCheckedChange={field.onChange}
                                  className="data-[state=checked]:bg-primary data-[state=checked]:border-primary"
                                />
                              </FormControl>
                            </FormItem>
                          )}
                        />

                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2">
                            {isActive && <CheckCircle2 className="w-3.5 h-3.5 text-success shrink-0" />}
                            <span className={cn("font-medium text-sm truncate", !isActive && "text-muted-foreground")}>
                              {test.testName}
                            </span>
                          </div>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <p className="text-xs text-muted-foreground truncate cursor-help">
                                {test.testSuite}
                              </p>
                            </TooltipTrigger>
                            <TooltipContent>{t("settings.tests.testSuite")}: {test.testSuite}</TooltipContent>
                          </Tooltip>
                        </div>

                        <FormField
                          control={form.control}
                          name={`assignmentTests.${index}.name`}
                          render={({ field }) => <input value={field.value} readOnly hidden />}
                        />

                        <FormField
                          control={form.control}
                          name={`assignmentTests.${index}.score`}
                          render={({ field }) => (
                            <FormItem className="flex items-center gap-2 space-y-0">
                              <FormLabel className="text-xs text-muted-foreground shrink-0">{t("settings.tests.points")}</FormLabel>
                              <FormControl>
                                <Input
                                  type="number"
                                  min={0}
                                  step={1}
                                  disabled={isPending || !isActive}
                                  {...field}
                                  onChange={(e) => {
                                    const value = e.target.value;
                                    field.onChange(value ? Number(value) : 0);
                                  }}
                                  className="w-20 h-8 text-center bg-background"
                                />
                              </FormControl>
                            </FormItem>
                          )}
                        />
                      </div>
                    );
                  })}
                </div>
              )}

              {/* Submit Button */}
              {tests.activatible && tests.report.length > 0 && (
                <div className="flex justify-end mt-4 pt-4 border-t border-border/50">
                  <Button type="submit" variant="glow" size="sm" disabled={isPending}>
                    {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                    {t("settings.tests.saveButton")}
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {error && (
            <Alert variant="destructive" className="mt-4">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>{tc("status.error")}</AlertTitle>
              <AlertDescription>{error.message}</AlertDescription>
            </Alert>
          )}
        </form>
      </Form>
    </section>
  );
};

const RubricForm = (props: { classroomId: string; assignmentId: string }) => {
  const { t } = useTranslation("assignment");
  const { t: tc } = useTranslation("common");
  const { classroomId, assignmentId } = props;

  const { data: rubrics } = useSuspenseQuery(classroomGradingRubricsQueryOptions(classroomId));
  const { data: assignmentRubrics } = useSuspenseQuery(assignmentGradingRubricsQueryOptions(classroomId, assignmentId));
  const { mutateAsync, isPending, error } = useUpdateAssignmentRubrics(classroomId, assignmentId);

  const rubricList = useMemo(() => {
    return rubrics.map((rubric) => ({
      ...rubric,
      active: assignmentRubrics.findIndex((r) => r.id === rubric.id) !== -1,
    }));
  }, [rubrics, assignmentRubrics]);

  const form = useForm<Record<string, boolean>>({
    defaultValues: Object.fromEntries(rubricList.map((rubric) => [rubric.id, rubric.active])),
  });

  const onSubmit = async (data: Record<string, boolean>) => {
    await mutateAsync({
      rubricIds: Object.entries(data)
        .filter((arg) => arg[1])
        .map(([key]) => key),
    });
    toast.success(t("settings.rubrics.success"));
  };

  const selectedCount = rubricList.filter((r) => form.watch(r.id)).length;
  const totalMaxScore = rubricList
    .filter((r) => form.watch(r.id))
    .reduce((sum, r) => sum + r.maxScore, 0);

  return (
    <section>
      {/* Header */}
      <div className="flex items-start justify-between gap-4 mb-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-[hsl(38,92%,55%)]/20 to-[hsl(38,92%,55%)]/5 border border-[hsl(38,92%,55%)]/20 flex items-center justify-center">
            <BookOpen className="w-5 h-5 text-[hsl(38,92%,55%)]" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h2 className="text-lg font-semibold font-mono">{t("settings.rubrics.title")}</h2>
              {selectedCount > 0 && (
                <StatusBadge variant="neutral" size="sm">
                  {t("settings.rubrics.selected", { count: selectedCount })}
                </StatusBadge>
              )}
            </div>
            <p className="text-sm text-muted-foreground">
              {t("settings.rubrics.subtitle")}
            </p>
          </div>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <Card className="border-border/50">
            <CardContent className="p-4">
              {rubricList.length === 0 ? (
                <div className="p-6 text-center">
                  <BookOpen className="w-8 h-8 mx-auto text-muted-foreground mb-2" />
                  <p className="text-sm text-muted-foreground mb-1">{t("settings.rubrics.noRubrics")}</p>
                  <p className="text-xs text-muted-foreground">
                    {t("settings.rubrics.noRubricsHelp")}
                  </p>
                </div>
              ) : (
                <>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                    {rubricList.map((rubric) => {
                      const isSelected = form.watch(rubric.id);
                      return (
                        <FormField
                          key={rubric.id}
                          control={form.control}
                          name={rubric.id}
                          render={({ field }) => (
                            <FormItem
                              className={cn(
                                "flex items-start gap-3 p-3 rounded-lg border transition-all duration-200",
                                isSelected
                                  ? "border-primary/30 bg-primary/5"
                                  : "border-border/50 hover:border-border"
                              )}
                            >
                              <FormControl>
                                <Checkbox
                                  id={`rubric-${rubric.id}`}
                                  disabled={isPending}
                                  checked={field.value}
                                  onCheckedChange={field.onChange}
                                  className="mt-0.5 data-[state=checked]:bg-primary data-[state=checked]:border-primary"
                                />
                              </FormControl>
                              <label
                                htmlFor={`rubric-${rubric.id}`}
                                className="flex-1 min-w-0 cursor-pointer"
                              >
                                <div className="flex items-center justify-between gap-2">
                                  <span className="font-medium">{rubric.name}</span>
                                  <span className="text-xs font-mono text-muted-foreground">
                                    {rubric.maxScore} pts
                                  </span>
                                </div>
                                {rubric.description && (
                                  <div className="text-xs text-muted-foreground mt-0.5 line-clamp-2">
                                    <Markdown compact inline>{rubric.description}</Markdown>
                                  </div>
                                )}
                              </label>
                            </FormItem>
                          )}
                        />
                      );
                    })}
                  </div>

                  {/* Summary & Submit */}
                  <div className="flex items-center justify-between mt-4 pt-4 border-t border-border/50">
                    <div className="text-sm text-muted-foreground">
                      {t("settings.rubrics.totalMaxScore")}: <span className="font-mono font-semibold text-foreground">{totalMaxScore}</span>
                    </div>
                    <Button type="submit" variant="glow" size="sm" disabled={isPending}>
                      {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                      {t("settings.rubrics.saveButton")}
                    </Button>
                  </div>
                </>
              )}
            </CardContent>
          </Card>

          {error && (
            <Alert variant="destructive" className="mt-4">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>{tc("status.error")}</AlertTitle>
              <AlertDescription>{error.message}</AlertDescription>
            </Alert>
          )}
        </form>
      </Form>
    </section>
  );
};
