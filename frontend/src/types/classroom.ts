import { z } from "zod";
import { reversed } from "@/types/utils";

export const Role = {
  Owner: 0,
  Moderator: 1,
  Student: 2,
} as const;

export type Role = (typeof Role)[keyof typeof Role];

const ReverseRole = reversed(Role);
export const getRole = (role: Role) => ReverseRole[role];

export const createFormSchema = z.object({
  name: z.string().min(3),
  description: z.string().min(3),
  createTeams: z.boolean(),
  maxTeamSize: z.coerce.number().int().min(1),
  maxTeams: z.coerce.number().int().min(0),
});
export type ClassroomForm = z.infer<typeof createFormSchema>;

export const inviteFormSchema = z.object({
  memberEmails: z
    .string()
    .min(3)
    .refine((emails) =>
      emails
        .split("\n")
        .filter(Boolean)
        .every(
          (email) => {
            const result = z.string().email().safeParse(email);
            return result.success;
          },
          { message: "One or more Emails are not valid" },
        ),
    ),
});
export type InviteForm = z.infer<typeof inviteFormSchema>;

export const Status = {
  Pending: 0,
  Accepted: 1,
  Rejected: 2,
  Revoked: 3,
} as const;

const GetStatus = reversed(Status);
export const getStatus = (status: Status) => GetStatus[status];

export type Status = (typeof Status)[keyof typeof Status];
