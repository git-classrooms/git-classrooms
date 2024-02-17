import { z } from "zod";
import { User } from "@/types/user.ts";

export type Assignment = {
  id: string;
  name: string;
  description: string;
  dueDate: string;
};

export type AssignmentProject = {
  assignmentId: string;
  user: User;
  assignmentAccepted: boolean;
  projectId: number;
  projectPath: string;
};

export type TemplateProject = {
  name: string;
  id: number;
  visibility: number;
  webUrl: string;
  description: string;
};

export const createAssignmentFormSchema = z.object({
  name: z.string().min(3),
  description: z.string().min(3),
  templateProjectId: z.number().min(1, "Please select a template project"),
  dueDate: z.coerce.date(),
});

export type CreateAssignmentForm = z.infer<typeof createAssignmentFormSchema>;
