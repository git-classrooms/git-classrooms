import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { createFormSchema, updateFormSchema } from "@/types/classroom";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, BookOpen, ChevronRight, FolderGit2, Loader2, Pencil, Settings2, Users } from "lucide-react";
import { cn, getUUIDFromLocation, unwrapApiError } from "@/lib/utils";
import { Switch } from "@/components/ui/switch";
import { toast } from "sonner";
import { useCreateClassroom, useUpdateClassroom } from "@/api/classroom";
import { UserClassroomResponse } from "@/swagger-client";
import { useUnsavedChanges } from "@/hooks/useUnsavedChanges";
import { Card, CardContent } from "@/components/ui/card";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { useTranslation } from "react-i18next";

export const ClassroomCreateForm = () => {
  const { t } = useTranslation("classroom");
  const { t: tc } = useTranslation("common");
  const navigate = useNavigate();
  const { mutateAsync, error, isPending } = useCreateClassroom();

  const createClassroomError = unwrapApiError(error);

  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    mode: "onBlur",
    reValidateMode: "onChange",
    defaultValues: {
      name: "",
      description: "",
      maxTeams: 0,
      maxTeamSize: 2,
      createTeams: true,
      studentsViewAllProjects: false,
      teamsEnabled: true,
      createTeachingMaterialGroup: true,
    },
  });

  useUnsavedChanges({ blocking: form.formState.isDirty && !form.formState.isSubmitting });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    const location = await mutateAsync(values);
    const classroomId = getUUIDFromLocation(location);
    await navigate({ to: "/classrooms/$classroomId", search: { tab: "assignments" }, params: { classroomId } });
  }

  const teamsEnabled = form.watch("teamsEnabled");

  return (
    <div className="max-w-2xl mx-auto">
      {/* Breadcrumb */}
      <Breadcrumb className="mb-6">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/classrooms">{t("title")}</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{t("create.title")}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-4 mb-3">
          <div className="w-14 h-14 rounded-xl bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center">
            <FolderGit2 className="w-7 h-7 text-primary" />
          </div>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">{t("create.title")}</h1>
            <p className="text-muted-foreground">{t("create.subtitle")}</p>
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
                <h2 className="font-semibold text-sm">{t("form.basicInfo")}</h2>
              </div>
            </div>
            <CardContent className="p-5 space-y-5">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("form.name")}</FormLabel>
                    <FormControl>
                      <Input placeholder={t("form.namePlaceholder")} className="bg-background" {...field} />
                    </FormControl>
                    <FormDescription>{t("form.nameDescription")}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("form.description")}</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder={t("form.descriptionPlaceholder")}
                        className="resize-none bg-background min-h-[100px]"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      {t("form.descriptionHelp")} · {tc("markdown.supported")}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Team Configuration Section */}
          <Card className="border-border/50 overflow-hidden">
            <div className="px-5 py-4 border-b border-border/50 bg-muted/30">
              <div className="flex items-center gap-2">
                <Users className="w-4 h-4 text-muted-foreground" />
                <h2 className="font-semibold text-sm">{t("teams.title")}</h2>
              </div>
            </div>
            <CardContent className="p-5 space-y-5">
              <FormField
                control={form.control}
                name="teamsEnabled"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between p-4 rounded-lg border border-border/50 bg-muted/20">
                    <div className="space-y-0.5">
                      <FormLabel className="font-medium">{t("teams.enabled")}</FormLabel>
                      <FormDescription className="text-xs">{t("teams.enabledDescription")}</FormDescription>
                    </div>
                    <FormControl>
                      <Switch checked={field.value} onCheckedChange={field.onChange} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <div
                className={cn(
                  "space-y-4 overflow-hidden transition-all duration-300",
                  teamsEnabled ? "opacity-100 max-h-[500px]" : "opacity-0 max-h-0 pointer-events-none",
                )}
              >
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <FormField
                    control={form.control}
                    name="maxTeams"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t("teams.maxTeams")}</FormLabel>
                        <FormControl>
                          <Input type="number" min={0} step={1} className="bg-background" {...field} />
                        </FormControl>
                        <FormDescription className="text-xs">{t("teams.maxTeamsDescription")}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="maxTeamSize"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t("teams.maxTeamSize")}</FormLabel>
                        <FormControl>
                          <Input type="number" min={2} step={1} className="bg-background" {...field} />
                        </FormControl>
                        <FormDescription className="text-xs">{t("teams.maxTeamSizeDescription")}</FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <FormField
                  control={form.control}
                  name="createTeams"
                  render={({ field }) => (
                    <FormItem className="flex items-center justify-between p-4 rounded-lg border border-border/50">
                      <div className="space-y-0.5">
                        <FormLabel className="font-medium">{t("teams.studentCreation")}</FormLabel>
                        <FormDescription className="text-xs">{t("teams.studentCreationDescription")}</FormDescription>
                      </div>
                      <FormControl>
                        <Switch checked={field.value} onCheckedChange={field.onChange} />
                      </FormControl>
                    </FormItem>
                  )}
                />
              </div>
            </CardContent>
          </Card>

          {/* Privacy Settings Section */}
          <Card className="border-border/50 overflow-hidden">
            <div className="px-5 py-4 border-b border-border/50 bg-muted/30">
              <div className="flex items-center gap-2">
                <Settings2 className="w-4 h-4 text-muted-foreground" />
                <h2 className="font-semibold text-sm">{t("privacy.title")}</h2>
              </div>
            </div>
            <CardContent className="p-5 space-y-5">
              <FormField
                control={form.control}
                name="studentsViewAllProjects"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between p-4 rounded-lg border border-border/50">
                    <div className="space-y-0.5">
                      <FormLabel className="font-medium">{t("privacy.mutualVisibility")}</FormLabel>
                      <FormDescription className="text-xs">{t("privacy.mutualVisibilityDescription")}</FormDescription>
                    </div>
                    <FormControl>
                      <Switch checked={field.value} onCheckedChange={field.onChange} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="createTeachingMaterialGroup"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between p-4 rounded-lg border border-border/50">
                    <div className="space-y-0.5">
                      <FormLabel className="font-medium">{t("privacy.teachingMaterial")}</FormLabel>
                      <FormDescription className="text-xs">{t("privacy.teachingMaterialDescription")}</FormDescription>
                    </div>
                    <FormControl>
                      <Switch checked={field.value} onCheckedChange={field.onChange} />
                    </FormControl>
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Error Alert */}
          {createClassroomError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>{tc("status.error")}</AlertTitle>
              <AlertDescription>{createClassroomError.message}</AlertDescription>
            </Alert>
          )}

          {/* Submit Button */}
          <div className="flex items-center justify-end gap-3 pt-4">
            <Button type="button" variant="outline" asChild>
              <Link to="/classrooms">{tc("actions.cancel")}</Link>
            </Button>
            <Button type="submit" variant="glow" disabled={isPending} className="min-w-[140px]">
              {isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <>
                  {t("create.button")}
                  <ChevronRight className="ml-2 h-4 w-4" />
                </>
              )}
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
};

export const ClassroomEditForm = ({ userClassroom }: { userClassroom: UserClassroomResponse }) => {
  const { t } = useTranslation("classroom");
  const { t: tc } = useTranslation("common");
  const { mutateAsync, isError, isPending } = useUpdateClassroom(userClassroom.classroom.id);

  const form = useForm<z.infer<typeof updateFormSchema>>({
    resolver: zodResolver(updateFormSchema),
    mode: "onBlur",
    reValidateMode: "onChange",
    defaultValues: {
      name: userClassroom.classroom.name,
      description: userClassroom.classroom.description,
    },
  });

  useUnsavedChanges({ blocking: form.formState.isDirty && !form.formState.isSubmitting });

  async function onSubmit(values: z.infer<typeof updateFormSchema>) {
    await mutateAsync(values);
    toast.success(t("edit.success"));
  }

  return (
    <div className="w-full">
      <div className="flex items-center gap-3 mb-6">
        <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-[hsl(142,71%,45%)]/20 to-[hsl(142,71%,45%)]/5 border border-[hsl(142,71%,45%)]/20 flex items-center justify-center">
          <Pencil className="w-5 h-5 text-[hsl(142,71%,45%)]" />
        </div>
        <div>
          <h2 className="text-lg font-semibold font-mono">{t("edit.title")}</h2>
          <p className="text-sm text-muted-foreground">{t("edit.subtitle")}</p>
        </div>
      </div>

      <Card className="border-border/50">
        <CardContent className="p-4">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs text-muted-foreground uppercase tracking-wide">
                      {t("form.nameLabel")}
                    </FormLabel>
                    <FormControl>
                      <Input placeholder={t("form.namePlaceholder")} {...field} className="bg-background" />
                    </FormControl>
                    <FormDescription className="text-xs">{t("form.descriptionLabel")}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs text-muted-foreground uppercase tracking-wide">
                      {t("form.description")}
                    </FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder={t("form.descriptionPlaceholder")}
                        className="resize-none bg-background min-h-[100px]"
                        {...field}
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      {t("form.descriptionPurpose")} · {tc("markdown.supported")}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="flex items-center justify-between pt-4 border-t border-border/50">
                <p className="text-xs text-muted-foreground">{t("edit.changesSaved")}</p>
                <Button type="submit" variant="glow" size="sm" disabled={isPending || !form.formState.isDirty}>
                  {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                  {tc("actions.saveChanges")}
                </Button>
              </div>

              {isError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>{tc("status.error")}</AlertTitle>
                  <AlertDescription>{t("errors.updateFailed")}</AlertDescription>
                </Alert>
              )}
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
};
