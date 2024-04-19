import { z } from "zod";
import { Team } from "@/types/team";

export type Assignment = {
  id: string;
  name: string;
  description: string;
  dueDate: string;
};

export type AssignmentProject = {
  id: string;
  assignmentId: string;
  team: Team;
  assignmentAccepted: boolean;
  projectId: number;
  projectPath: string;
};

export type TemplateProject = {
  name: string;
  id: number;
};

export const createAssignmentFormSchema = z.object({
  name: z.string().min(3),
  description: z.string().min(3),
  templateProjectId: z.number().min(1, "Please select a template project"),
  dueDate: z.coerce.date(),
});

export type CreateAssignmentForm = z.infer<typeof createAssignmentFormSchema>;
