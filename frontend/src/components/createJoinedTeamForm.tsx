import { createFormSchema } from "@/types/team";
import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "@tanstack/react-router";
import { Loader2, AlertCircle } from "lucide-react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Alert, AlertTitle, AlertDescription } from "./ui/alert";
import { Form, FormField, FormItem, FormLabel, FormControl, FormDescription, FormMessage } from "./ui/form";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { useCreateTeam } from "@/api/team";

export const CreateJoinedTeamForm = ({ classroomId }: { classroomId: string }) => {
  const navigate = useNavigate();
  const { mutateAsync, isError, isPending } = useCreateTeam(classroomId);

  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    defaultValues: {
      name: "",
    },
  });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    await mutateAsync(values);
    navigate({ to: "/classrooms/joined/$classroomId", params: { classroomId } });
  }

  return (
    <div className="p-2">
      <div className="flex flex-row justify-between">
        <h1 className="text-xl font-bold">Create a new team</h1>
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
                  <Input placeholder="team name" {...field} />
                </FormControl>
                <FormDescription>This is your team name.</FormDescription>
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
              <AlertDescription>The team could not be created!</AlertDescription>
            </Alert>
          )}
        </form>
      </Form>
    </div>
  );
};
