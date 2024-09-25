import { z } from "zod";
import { zodEnumFromObjectKeys } from "./utils";
import { Role } from "./classroom";

export const createFormSchema = z.object({
  role: zodEnumFromObjectKeys(Role),
});

export type MemberForm = z.infer<typeof createFormSchema>;
