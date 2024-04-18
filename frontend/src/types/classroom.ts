import { z } from "zod";
import { reversed } from "@/types/utils";
import { User } from "@/types/user";
import { Team } from "@/types/team";

export const Role = {
  Owner: 0,
  Moderator: 1,
  Student: 2,
} as const;

export type Role = (typeof Role)[keyof typeof Role];

export const GetRole = reversed(Role);

export type UserClassroom = {
  classroom: Omit<OwnedClassroom, "gitlabUrl">;
  user: User;
  role: Role;
  team?: Team;
  gitlabUrl: string;
};

export type OwnedClassroom = {
  id: string;
  name: string;
  ownerId: number;
  owner: User;
  description: string;
  groupId: number;
  gitlabUrl: string;
  maxTeamSize: number;
  maxTeams: number;
  createTeams: boolean;
};

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
  Declined: 2,
  Revoked: 3,
} as const;

export const GetStatus = reversed(Status);

export type Status = (typeof Status)[keyof typeof Status];

export type ClassroomInvitation = {
  id: string;
  status: Status;
  createdAt: string;
  email: string;
  expiryDate: string;
};
