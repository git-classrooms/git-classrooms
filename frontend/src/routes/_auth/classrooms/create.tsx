import { zodResolver } from "@hookform/resolvers/zod";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { createFormSchema } from "@/types/classroom";
import { useCreateClassroom } from "@/api/classrooms";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, Loader2 } from "lucide-react";
import { getUUIDFromLocation } from "@/lib/utils.ts";

export const Route = createFileRoute("/_auth/classrooms/create")({
  component: ClassroomsForm,
});

function ClassroomsForm() {
  const navigate = useNavigate();
  const { mutateAsync, isError, isPending } = useCreateClassroom();
  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    defaultValues: {
      name: "",
      description: "",
    },
  });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    const location = await mutateAsync(values);
    const classroomId = getUUIDFromLocation(location);
    await navigate({ to: "/classrooms/$classroomId", params: { classroomId } });
  }

  return (
    <div className="p-2">
      <div className="flex flex-row justify-between">
        <h1 className="text-xl font-bold">Create a classroom</h1>
      </div>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name</FormLabel>
                <FormControl>
                  <Input placeholder="Programming classroom" {...field} />
                </FormControl>
                <FormDescription>This is your classroom name.</FormDescription>
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
                <FormDescription>This is the description of your classroom.</FormDescription>
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
    </div>
  );
}
