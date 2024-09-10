import { classroomGradingRubricsQueryOptions, useUpdateClassroomRubrics } from "@/api/grading";
import { classroomAvailableRunnersQueryOptions } from "@/api/runners";
import { Loader } from "@/components/loader";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { AlertCircle, Edit2, FolderPlus, RefreshCcw, Trash } from "lucide-react";
import React, { useEffect } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/settings/grading")({
  loader: async ({ params: { classroomId }, context: { queryClient } }) => {
    const rubrics = await queryClient.ensureQueryData(classroomGradingRubricsQueryOptions(classroomId));
    return { rubrics };
  },
  component: Grading,
});

const rubricSchema = z.object({
  id: z.string().uuid().optional(),
  name: z.string().min(3),
  description: z.string(),
  maxScore: z.number().int().positive(),
});

const formSchmema = z.object({
  gradingManualRubrics: z.array(rubricSchema),
});

function Grading() {
  const { classroomId } = Route.useParams();

  const {
    data: isRunnerAvailable,
    refetch: runnerStatusRefetch,
    isFetching: isRunnerAvailableFetching,
  } = useQuery(classroomAvailableRunnersQueryOptions(classroomId));

  const { data } = useSuspenseQuery(classroomGradingRubricsQueryOptions(classroomId));
  const { mutateAsync, isPending, error } = useUpdateClassroomRubrics(classroomId);

  const [editing, setEditing] = React.useState(false);

  const form = useForm<z.infer<typeof formSchmema>>({
    resolver: zodResolver(formSchmema),
    defaultValues: {
      gradingManualRubrics: data,
    },
  });

  useEffect(() => {
    if (!editing) form.reset({ gradingManualRubrics: data });
  }, [data, editing, form]);

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "gradingManualRubrics",
  });

  const disabled = !editing || isPending;
  const onSubmit = async (values: z.infer<typeof formSchmema>) => {
    await mutateAsync(values);
    toast.success("Rubrics saved successfully.");
    setEditing(false);
  };

  const onCancel = () => {
    form.reset({ gradingManualRubrics: data });
    setEditing(false);
  };

  return (
    <div className="p-2 w-full">
      <div className="flex mb-6">
        <div className="grow">
          <div className="flex items-center">
            <h2 className="text-xl font-bold mr-2.5">Test-driven grading</h2>
            {isRunnerAvailableFetching ? (
              <Skeleton className="rounded-full h-3 w-3" />
            ) : (
              <Tooltip delayDuration={0}>
                <TooltipTrigger asChild>
                  <span className="relative flex h-3 w-3">
                    <span
                      className={cn(
                        "animate-ping absolute inline-flex h-full w-full rounded-full opacity-75",
                        isRunnerAvailable ? "bg-emerald-400" : "bg-red-500",
                      )}
                    ></span>
                    <span
                      className={cn(
                        "relative inline-flex rounded-full h-3 w-3",
                        isRunnerAvailable ? "bg-emerald-500" : "bg-red-600",
                      )}
                    ></span>
                  </span>
                </TooltipTrigger>
                <TooltipContent>
                  {isRunnerAvailable ? (
                    <div>Test-driven grading available.</div>
                  ) : (
                    <div>Test-driven grading not available.</div>
                  )}
                </TooltipContent>
              </Tooltip>
            )}
          </div>
          <p className="text-sm text-muted-foreground">
            Status of automatic test-driven grading using CI/CD test reports for this classroom.
          </p>
        </div>
        <Button
          className="flex-none items-center"
          disabled={isRunnerAvailableFetching}
          onClick={() => runnerStatusRefetch()}
          variant="outline"
        >
          <RefreshCcw className="mr-2 h-4 w-4" /> Refresh
        </Button>
      </div>

      <p className="mt-2">
        An automated grading can be carried out using test results that are generated as a result of executing a CI/CD
        pipeline in GitLab. The executed automated tests must generate a report artifact in JUnit XML report format.
      </p>

      <div className="mt-2">
        {isRunnerAvailableFetching ? (
          <Skeleton className="h-6 rounded-lg w-full" />
        ) : isRunnerAvailable ? (
          <p>
            <b>
              At least one runner is available for the current classroom. Automatic test-driven grading is available.
            </b>
          </p>
        ) : (
          <p>
            <b>The associated GitLab group of this classroom does not yet have a runner or no runner is available.</b>
          </p>
        )}
      </div>

      <Separator className="my-6" />

      <div className="flex mb-6">
        <div className="grow">
          <h2 className="text-xl font-bold">Manual grading</h2>
          <p className="text-sm text-muted-foreground">Configure the manual grading rubrics for this classroom.</p>
        </div>
        {!editing && (
          <Button className="flex-none items-center" onClick={() => setEditing(true)} variant="outline">
            <Edit2 className="mr-2 h-4 w-4" /> Edit
          </Button>
        )}
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="mt-4 w-full">
          <div className="grid gap-2 grid-cols-1 md:grid-cols-[2fr_4fr_1fr_auto] w-full">
            <FormLabel className="hidden md:block">Name</FormLabel>
            <FormLabel className="hidden md:block">Description</FormLabel>
            <FormLabel className="hidden md:block">Max. score</FormLabel>
            <div></div>
            {fields.map((field, index) => (
              <React.Fragment key={field.id}>
                <FormField
                  control={form.control}
                  name={`gradingManualRubrics.${index}.id`}
                  render={({ field }) => <input hidden readOnly value={field.value} />}
                />

                <FormField
                  control={form.control}
                  name={`gradingManualRubrics.${index}.name`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="md:hidden">Name</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Name of the rubric"
                          disabled={disabled}
                          {...field}
                          className={"text-base border-r-none rounded-r-none"}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name={`gradingManualRubrics.${index}.description`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="md:hidden">Description</FormLabel>
                      <FormControl>
                        <Input
                          className={"text-base rounded-none"}
                          placeholder="Description"
                          type="text"
                          disabled={disabled}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name={`gradingManualRubrics.${index}.maxScore`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="md:hidden">Max. score</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          placeholder="Max. score"
                          min={0}
                          step={1}
                          disabled={disabled}
                          {...field}
                          onChange={(e) => {
                            const value = e.target.value;
                            const numberValue = value ? Number(value) : "";
                            field.onChange(numberValue);
                          }}
                          className={"text-base border-r-none rounded-l-none"}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <Button
                  onClick={() => remove(index)}
                  disabled={disabled}
                  type="button"
                  variant="destructive"
                  size="icon"
                  className="mt-4 justify-self-end md:mt-2"
                >
                  <Trash />
                </Button>

                <Separator className="md:hidden my-6" />
              </React.Fragment>
            ))}
          </div>
          {editing && (
            <div className="flex justify-end mt-4 gap-4">
              <Button
                onClick={() => append({ description: "", name: "", maxScore: 0 })}
                disabled={disabled}
                variant="secondary"
                type="button"
              >
                <FolderPlus className="mr-2 h-4 w-4" /> Add rubric
              </Button>

              <Button disabled={disabled} type="submit">
                {isPending ? <Loader /> : "Save"}
              </Button>
              <Button onClick={onCancel} variant="destructive" disabled={disabled} type="button">
                Cancel
              </Button>
            </div>
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
    </div>
  );
}
