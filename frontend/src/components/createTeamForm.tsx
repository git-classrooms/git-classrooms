import { createFormSchema } from "@/types/team";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2, AlertCircle, Users, ChevronRight } from "lucide-react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Alert, AlertTitle, AlertDescription } from "./ui/alert";
import { Form, FormField, FormItem, FormLabel, FormControl, FormDescription, FormMessage } from "./ui/form";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { useCreateTeam } from "@/api/team";
import React from "react";

interface CreateTeamFormProps {
  classroomId: string;
  onSuccess: () => Promise<void> | void;
  onCancel?: () => void;
}

export const CreateTeamForm: React.FC<CreateTeamFormProps> = ({ classroomId, onSuccess, onCancel }) => {
  const { mutateAsync, isError, isPending } = useCreateTeam(classroomId);

  const form = useForm<z.infer<typeof createFormSchema>>({
    resolver: zodResolver(createFormSchema),
    mode: "onBlur",
    reValidateMode: "onChange",
    defaultValues: {
      name: "",
    },
  });

  async function onSubmit(values: z.infer<typeof createFormSchema>) {
    await mutateAsync(values);
    onSuccess();
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-primary/20 to-primary/5 border border-primary/20 flex items-center justify-center">
          <Users className="w-5 h-5 text-primary" />
        </div>
        <div>
          <h2 className="text-lg font-semibold">Create Team</h2>
          <p className="text-sm text-muted-foreground">Start a new team for this classroom</p>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Team Name</FormLabel>
                <FormControl>
                  <Input
                    placeholder="e.g., Team Alpha"
                    className="bg-background"
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  Choose a unique name for your team
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          {isError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>The team could not be created. Please try again.</AlertDescription>
            </Alert>
          )}

          <div className="flex items-center justify-end gap-3 pt-2">
            {onCancel && (
              <Button type="button" variant="outline" onClick={onCancel}>
                Cancel
              </Button>
            )}
            <Button type="submit" variant="glow" disabled={isPending} className="min-w-[120px]">
              {isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <>
                  Create Team
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
