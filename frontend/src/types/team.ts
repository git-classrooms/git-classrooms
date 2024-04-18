import { z } from "zod";

export type Team = {
  id: string;
  name: string;
  groupId: number;
};

export const createFormSchema = z.object({
  name: z.string().min(3),
});

export type TeamForm = z.infer<typeof createFormSchema>;
