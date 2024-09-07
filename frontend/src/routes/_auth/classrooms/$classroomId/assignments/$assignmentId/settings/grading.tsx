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
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { AlertCircle, BookOpenCheck, FolderGit2, Loader2 } from "lucide-react";
import test from "node:test";
import { Suspense, useMemo, useState } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

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
    <div className="p-2 w-full">
      <Suspense fallback={<Skeleton className="h-20" />}>
        <TestsForm classroomId={classroomId} assignmentId={assignmentId} />
      </Suspense>
      <Separator className="my-6" />
      <RubricForm classroomId={classroomId} assignmentId={assignmentId} />
    </div>
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
  const { classroomId, assignmentId } = props;

  const { data: assignment } = useSuspenseQuery(assignmentQueryOptions(classroomId, assignmentId));
  const { data: tests } = useSuspenseQuery(assignmentTestsQueryOptions(classroomId, assignmentId));

  const assignmentTests = useMemo<z.infer<typeof testsFormSchema>["assignmentTests"]>(
    () =>
      tests.report.map((test) => {
        const selected = tests.selectedTests.find((selectedTest) => selectedTest.name === test.name);
        return {
          name: test.name,
          score: selected?.score ?? 0,
          active: !!selected,
        };
      }),
    [tests],
  );

  const { mutateAsync, isPending, error } = useUpdateAssignmentTests(classroomId, assignmentId);

  const form = useForm<z.infer<typeof testsFormSchema>>({
    resolver: zodResolver(testsFormSchema),
    defaultValues: { junitAutoGradingActive: assignment.gradingJUnitAutoGradingActive, assignmentTests },
  });

  const { fields } = useFieldArray({
    control: form.control,
    name: "assignmentTests",
  });

  const onSubmit = async (data: z.infer<typeof testsFormSchema>) => {
    console.log(data);
    await mutateAsync({
      junitAutoGradingActive: data.junitAutoGradingActive,
      assignmentTests: data.assignmentTests.filter((test) => test.active),
    });
    toast.success("Tests updated");
  };

  const [showTests, setShowTests] = useState(assignment.gradingJUnitAutoGradingActive);

  return (
    <>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <div className="md:flex md:items-center mb-6">
            <div className="grow">
              <h2 className="text-xl font-bold mr-2.5">Test-driven grading</h2>
              <p className="text-sm text-muted-foreground">
                Configure which test are to be included in the automatic grading and how many points are awarded per
                successful test.
              </p>
            </div>

            <FormField
              control={form.control}
              name="junitAutoGradingActive"
              render={({ field }) => (
                <FormItem className="flex flex-row h-10 ml-0 mt-4 mb-10 md:mt-0 md:ml-4 md:mb-0 items-center space-x-3 space-y-0 rounded-md border p-4">
                  <FormControl>
                    <Checkbox
                      disabled={!tests.activatible || isPending}
                      checked={field.value}
                      onCheckedChange={(state) => {
                        field.onChange(state);
                        setShowTests(state !== "indeterminate" && state);
                      }}
                    />
                  </FormControl>
                  <FormLabel>Enable</FormLabel>
                </FormItem>
              )}
            />
          </div>

          {!tests.activatible || tests.report.length == 0 ? (
            <>
              <p>
                There are no tests found in your template project. To use this feature define some tests in the{" "}
                <i className="font-bold">.gitlab-ci.yml</i> of your template and make sure that the pipeline is running
                through to get a list with available tests.
              </p>
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Example for .gitlab-ci.yml</CardTitle>
                  <BookOpenCheck className="mr-2 h-4 w-4" />
                </CardHeader>
                <CardContent>
                  <pre className="text-xs">{tests.example.example}</pre>
                </CardContent>
              </Card>
            </>
          ) : (
            ""
          )}

          {fields.map((field, index) => {
            const test = tests.report.find((test) => test.name === field.name)!;
            return (
              <div className={cn(" border rounded-md", showTests ? "md:flex" : "hidden")} key={field.id}>
                <FormField
                  control={form.control}
                  name={`assignmentTests.${index}.active`}
                  render={({ field }) => (
                    <FormItem className="flex md:grow items-center space-x-3 space-y-0 p-4">
                      <FormControl>
                        <Checkbox disabled={isPending} checked={field.value} onCheckedChange={field.onChange} />
                      </FormControl>
                      <FormLabel>{test.testName}</FormLabel>
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <FormDescription>{test.testSuite}</FormDescription>
                        </TooltipTrigger>
                        <TooltipContent>Test originates within the »{test.testSuite}« test suite</TooltipContent>
                      </Tooltip>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name={`assignmentTests.${index}.name`}
                  render={({ field }) => <input value={field.value} readOnly hidden />}
                />
                <FormField
                  control={form.control}
                  name={`assignmentTests.${index}.score`}
                  render={({ field }) => (
                    <FormItem className="flex items-center mb-4 ml-4 mr-4 md:mt-2 md:mb-2 md:ml-0  md:w-1/4 space-x-3 space-y-0 ">
                      <FormLabel>Points</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          placeholder="Points"
                          min={0}
                          step={1}
                          disabled={isPending}
                          {...field}
                          onChange={(e) => {
                            const value = e.target.value;
                            const numberValue = value ? Number(value) : "";
                            field.onChange(numberValue);
                          }}
                          className={"text-base rounded-r"}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
              </div>
            );
          })}
          {tests.activatible && tests.report.length !== 0 && (
            <Button disabled={isPending} type="submit">
              {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Save"}
            </Button>
          )}
        </form>
        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{error.message}</AlertDescription>
          </Alert>
        )}
      </Form>
    </>
  );
};

const RubricForm = (props: { classroomId: string; assignmentId: string }) => {
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
        .filter(([_, value]) => value)
        .map(([key]) => key),
    });
    toast.success("Rubrics updated");
  };
  return (
    <>
      <div className="mb-6">
        <h2 className="text-xl font-bold mr-2.5">Manual grading</h2>
        <p className="text-sm text-muted-foreground">
          Configure which grading rubrics defined for this classroom are intended to use for manual grading in this
          assignment.
        </p>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 ">
          <div className="grid grid-cols-2 gap-6 items-center">
            {rubricList.map((rubric) => (
              <FormField
                key={rubric.id}
                control={form.control}
                name={rubric.id}
                render={({ field }) => (
                  <FormItem className="flex flex-row items-center space-x-3 space-y-0 rounded-md border p-4">
                    <FormControl>
                      <Checkbox disabled={isPending} checked={field.value} onCheckedChange={field.onChange} />
                    </FormControl>
                    <FormLabel>{rubric.name}</FormLabel>
                    <FormDescription>{rubric.description}</FormDescription>
                  </FormItem>
                )}
              />
            ))}
          </div>
          <Button disabled={isPending} type="submit">
            {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Save"}
          </Button>
        </form>
        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{error.message}</AlertDescription>
          </Alert>
        )}
      </Form>
    </>
  );
};
