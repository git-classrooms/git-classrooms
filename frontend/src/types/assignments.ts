import { z } from "zod";

export const createAssignmentFormSchema = z.object({
  name: z.string()
    .min(3)
    .regex(/^[\p{L}\p{N}\p{Emoji}_.+\-\s]+$/u, "Invalid characters in name"),
  description: z.string().min(3),
  templateProjectId: z.number().min(1, "Please select a template project"),
  dueDate: z.coerce.date().optional(),
});


export type CreateAssignmentForm = z.infer<typeof createAssignmentFormSchema>;

export const updateAssignmentFormSchema = (isAccepted: boolean) =>
  z.object({
    name: isAccepted
      ? z.undefined()
      : z.string()
        .min(3)
        .regex(/^[\p{L}\p{N}\p{Emoji}_.+\-\s]+$/u, "Invalid characters in name"),
    description: isAccepted
      ? z.undefined()
      : z.string().min(3),
    dueDate: z.coerce.date().nullable(),
  });

export interface UpdateAssignmentForm {
  name?: string;
  description?: string;
  dueDate: Date | null;
}
