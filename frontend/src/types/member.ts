import { z } from "zod";

export const createFormSchema = z.object({
  role: z.number(),
});

export type MemberForm = z.infer<typeof createFormSchema>;
