import { z } from "zod";

export const createAssignmentFormSchema = z.object({
  name: z.string()
    .min(3)
    .regex(/^[\p{L}\p{N}\p{Emoji}_.+\-\s]+$/u, "Invalid characters in name"),
  description: z.string().min(3),
  templateProjectId: z.number().min(1, "Please select a template project"),
  dueDate: z.coerce.date(),
});


export type CreateAssignmentForm = z.infer<typeof createAssignmentFormSchema>;
