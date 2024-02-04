import { z } from "zod";

const role = {
  Owner: 0,
  Moderator: 1,
  Student: 2
} as const

export type Role = typeof role[keyof typeof role]

export type Classroom = {
  classroom:
  {
    id: string,
    name: string,
    ownerId: number,
    description: string,
    groupId: number
  },
  role: Role
}

export const createFormSchema = z.object({
  name: z.string().min(3),
  description: z.string().min(3),
})


export const inviteFormSchema = z.object({
  memberEmails: z.string()
})


export type ClassroomForm = z.infer<typeof createFormSchema>
