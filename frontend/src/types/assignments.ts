import { z } from "zod";

export const createAssignmentFormSchema = z.object({
  name: z.string().min(3),
  description: z.string().min(3),
  templateProjectId: z.number().min(1, "Please select a template project"),
  dueDate: z.coerce.date().optional(),
});

export type CreateAssignmentForm = z.infer<typeof createAssignmentFormSchema>;
