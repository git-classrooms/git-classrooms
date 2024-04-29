import { z } from "zod";
import { User } from "@/types/user.ts";

export type Team = {
  id: string;
  name: string;
  groupId: number;
  createdAt: string;
  updatedAt: string;
  members: User[];
  gitlabUrl: string;
};

export const createFormSchema = z.object({
  name: z.string().min(3),
});

export type TeamForm = z.infer<typeof createFormSchema>;
