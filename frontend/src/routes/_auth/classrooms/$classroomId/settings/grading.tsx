import { classroomGradingRubricsQueryOptions, useUpdateClassroomRubrics } from "@/api/grading";
import { classroomAvailableRunnersQueryOptions } from "@/api/runners";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import {
  AlertCircle,
  Beaker,
  BookOpen,
  Edit2,
  GripVertical,
  Loader2,
  Plus,
  RefreshCcw,
  Server,
  Trash2,
  X,
} from "lucide-react";
import React, { useEffect } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { StatusBadge } from "@/components/ui/status-badge";
import { Card, CardContent } from "@/components/ui/card";

export const Route = createFileRoute("/_auth/classrooms/$classroomId/settings/grading")({
  loader: async ({ params: { classroomId }, context: { queryClient } }) => {
    const rubrics = await queryClient.ensureQueryData(classroomGradingRubricsQueryOptions(classroomId));
    return { rubrics };
  },
  component: Grading,
});

const rubricSchema = z.object({
  id: z.string().uuid().optional(),
  name: z.string().min(3, "Name must be at least 3 characters"),
  description: z.string(),
  maxScore: z.number().int().positive("Score must be positive"),
});

const formSchema = z.object({
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

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    mode: "onBlur",
    reValidateMode: "onChange",
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

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    await mutateAsync(values);
    toast.success("Rubrics saved successfully.");
    setEditing(false);
  };

  const onCancel = () => {
    form.reset({ gradingManualRubrics: data });
    setEditing(false);
  };

  const totalMaxScore = fields.reduce((sum, field) => {
    const fieldValue = form.watch(`gradingManualRubrics.${fields.indexOf(field)}.maxScore`);
    return sum + (fieldValue || 0);
  }, 0);

  return (
    <div className="space-y-8">
      {/* Test-Driven Grading Section */}
      <section className="animate-in fade-in slide-in-from-bottom-2 duration-300">
        <div className="flex items-start justify-between gap-4 mb-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center">
              <Beaker className="w-5 h-5 text-primary" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h2 className="text-lg font-semibold font-mono">Test-Driven Grading</h2>
                {isRunnerAvailableFetching ? (
                  <Skeleton className="h-5 w-16 rounded-full" />
                ) : (
                  <RunnerStatusIndicator available={isRunnerAvailable ?? false} />
                )}
              </div>
              <p className="text-sm text-muted-foreground">
                Automated grading using CI/CD test reports
              </p>
            </div>
          </div>

          <Button
            variant="outline"
            size="sm"
            disabled={isRunnerAvailableFetching}
            onClick={() => runnerStatusRefetch()}
            className="shrink-0"
          >
            <RefreshCcw className={cn("w-4 h-4 mr-2", isRunnerAvailableFetching && "animate-spin")} />
            Refresh
          </Button>
        </div>

        {/* Runner Status Card */}
        <Card className="border-border/50">
          <CardContent className="p-4">
            <div className="flex items-start gap-4">
              <div
                className={cn(
                  "w-10 h-10 rounded-lg flex items-center justify-center shrink-0",
                  isRunnerAvailable ? "bg-success/10 text-success" : "bg-destructive/10 text-destructive"
                )}
              >
                <Server className="w-5 h-5" />
              </div>
              <div className="flex-1">
                {isRunnerAvailableFetching ? (
                  <div className="space-y-2">
                    <Skeleton className="h-5 w-48" />
                    <Skeleton className="h-4 w-full" />
                  </div>
                ) : isRunnerAvailable ? (
                  <>
                    <p className="font-medium text-success">Runner Available</p>
                    <p className="text-sm text-muted-foreground mt-1">
                      At least one GitLab runner is configured for this classroom. Automatic test-driven
                      grading using JUnit XML reports is available.
                    </p>
                  </>
                ) : (
                  <>
                    <p className="font-medium text-destructive">No Runner Available</p>
                    <p className="text-sm text-muted-foreground mt-1">
                      The associated GitLab group does not have a runner configured. Contact your GitLab
                      administrator to enable CI/CD runners for automated grading.
                    </p>
                  </>
                )}
              </div>
            </div>

            {/* Info box */}
            <div className="mt-4 p-3 rounded-lg bg-muted/30 border border-border/50">
              <p className="text-xs text-muted-foreground">
                <span className="font-medium text-foreground">How it works:</span> Automated tests in
                your CI/CD pipeline generate JUnit XML reports. These reports are parsed to calculate
                grades based on test results.
              </p>
            </div>
          </CardContent>
        </Card>
      </section>

      {/* Manual Grading Section */}
      <section className="animate-in fade-in slide-in-from-bottom-2 duration-300" style={{ animationDelay: "100ms" }}>
        <div className="flex items-start justify-between gap-4 mb-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-[hsl(38,92%,55%)]/20 to-[hsl(38,92%,55%)]/5 border border-[hsl(38,92%,55%)]/20 flex items-center justify-center">
              <BookOpen className="w-5 h-5 text-[hsl(38,92%,55%)]" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h2 className="text-lg font-semibold font-mono">Manual Grading</h2>
                {fields.length > 0 && (
                  <StatusBadge variant="neutral" size="sm">
                    {fields.length} rubric{fields.length !== 1 ? "s" : ""}
                  </StatusBadge>
                )}
              </div>
              <p className="text-sm text-muted-foreground">
                Configure rubrics for manual assessment
              </p>
            </div>
          </div>

          {!editing && (
            <Button variant="outline" size="sm" onClick={() => setEditing(true)} className="shrink-0">
              <Edit2 className="w-4 h-4 mr-2" />
              Edit Rubrics
            </Button>
          )}
        </div>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            {/* Rubrics List */}
            {fields.length === 0 ? (
              <Card className="border-dashed border-border/50">
                <CardContent className="p-6 text-center">
                  <BookOpen className="w-8 h-8 mx-auto text-muted-foreground mb-2" />
                  <h3 className="font-medium mb-1">No Rubrics Defined</h3>
                  <p className="text-sm text-muted-foreground mb-3">
                    Create rubrics to enable manual grading for assignments
                  </p>
                  {!editing && (
                    <Button variant="outline" size="sm" onClick={() => setEditing(true)}>
                      <Plus className="w-4 h-4 mr-2" />
                      Add Your First Rubric
                    </Button>
                  )}
                </CardContent>
              </Card>
            ) : (
              <div className="space-y-3">
                {fields.map((field, index) => (
                  <RubricCard
                    key={field.id}
                    index={index}
                    form={form}
                    disabled={disabled}
                    onRemove={() => remove(index)}
                    editing={editing}
                  />
                ))}
              </div>
            )}

            {/* Total Score Summary */}
            {fields.length > 0 && (
              <div className="mt-4 p-3 rounded-lg bg-muted/30 border border-border/50 flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Total Maximum Score</span>
                <span className="font-mono font-semibold text-lg">{totalMaxScore}</span>
              </div>
            )}

            {/* Action Buttons */}
            {editing && (
              <div className="flex flex-wrap items-center justify-between gap-4 mt-6 pt-4 border-t border-border/50">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => append({ description: "", name: "", maxScore: 0 })}
                  disabled={disabled}
                >
                  <Plus className="w-4 h-4 mr-2" />
                  Add Rubric
                </Button>

                <div className="flex items-center gap-2">
                  <Button type="button" variant="ghost" size="sm" onClick={onCancel} disabled={isPending}>
                    <X className="w-4 h-4 mr-2" />
                    Cancel
                  </Button>
                  <Button type="submit" variant="glow" size="sm" disabled={isPending}>
                    {isPending ? (
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    ) : null}
                    Save Changes
                  </Button>
                </div>
              </div>
            )}
          </form>

          {error && (
            <Alert variant="destructive" className="mt-4">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{error.message}</AlertDescription>
            </Alert>
          )}
        </Form>
      </section>
    </div>
  );
}

function RunnerStatusIndicator({ available }: { available: boolean }) {
  return (
    <Tooltip delayDuration={0}>
      <TooltipTrigger asChild>
        <div className="flex items-center gap-1.5">
          <span className="relative flex h-2.5 w-2.5">
            <span
              className={cn(
                "animate-ping absolute inline-flex h-full w-full rounded-full opacity-75",
                available ? "bg-success" : "bg-destructive"
              )}
            />
            <span
              className={cn(
                "relative inline-flex rounded-full h-2.5 w-2.5",
                available ? "bg-success" : "bg-destructive"
              )}
            />
          </span>
          <StatusBadge variant={available ? "success" : "destructive"} size="sm">
            {available ? "Online" : "Offline"}
          </StatusBadge>
        </div>
      </TooltipTrigger>
      <TooltipContent>
        {available ? "Test-driven grading is available" : "No CI/CD runner available"}
      </TooltipContent>
    </Tooltip>
  );
}

function RubricCard({
  index,
  form,
  disabled,
  onRemove,
  editing,
}: {
  index: number;
  form: ReturnType<typeof useForm<z.infer<typeof formSchema>>>;
  disabled: boolean;
  onRemove: () => void;
  editing: boolean;
}) {
  const rubricName = form.watch(`gradingManualRubrics.${index}.name`);
  const rubricScore = form.watch(`gradingManualRubrics.${index}.maxScore`);

  return (
    <Card
      className={cn(
        "transition-all duration-200",
        editing && "hover:border-primary/30"
      )}
    >
      <CardContent className="p-4">
        <div className="flex items-start gap-3">
          {/* Drag handle placeholder */}
          {editing && (
            <div className="pt-2.5 text-muted-foreground/50 cursor-grab">
              <GripVertical className="w-4 h-4" />
            </div>
          )}

          <div className="flex-1 min-w-0">
            {editing ? (
              <div className="grid grid-cols-1 md:grid-cols-[2fr_3fr_auto] gap-3">
                {/* Hidden ID field */}
                <FormField
                  control={form.control}
                  name={`gradingManualRubrics.${index}.id`}
                  render={({ field }) => <input hidden readOnly value={field.value ?? ""} />}
                />

                {/* Name Field */}
                <FormField
                  control={form.control}
                  name={`gradingManualRubrics.${index}.name`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-xs text-muted-foreground">Name</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Rubric name"
                          disabled={disabled}
                          {...field}
                          className="bg-background"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* Description Field */}
                <FormField
                  control={form.control}
                  name={`gradingManualRubrics.${index}.description`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-xs text-muted-foreground">Description</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Brief description"
                          disabled={disabled}
                          {...field}
                          className="bg-background"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {/* Max Score Field */}
                <FormField
                  control={form.control}
                  name={`gradingManualRubrics.${index}.maxScore`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-xs text-muted-foreground">Max Score</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          placeholder="0"
                          min={0}
                          step={1}
                          disabled={disabled}
                          {...field}
                          onChange={(e) => {
                            const value = e.target.value;
                            const numberValue = value ? Number(value) : 0;
                            field.onChange(numberValue);
                          }}
                          className="bg-background w-24"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            ) : (
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="font-medium">{rubricName || "Unnamed Rubric"}</h4>
                  <p className="text-sm text-muted-foreground">
                    {form.watch(`gradingManualRubrics.${index}.description`) || "No description"}
                  </p>
                </div>
                <div className="text-right">
                  <span className="font-mono text-lg font-semibold">{rubricScore}</span>
                  <p className="text-xs text-muted-foreground">points</p>
                </div>
              </div>
            )}
          </div>

          {/* Delete Button */}
          {editing && (
            <Button
              type="button"
              variant="ghost"
              size="icon"
              onClick={onRemove}
              disabled={disabled}
              className="shrink-0 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
            >
              <Trash2 className="w-4 h-4" />
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
