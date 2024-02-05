import { apiClient } from "@/lib/utils";
import { Classroom, ClassroomForm, ClassroomInvitation, InviteForm } from "@/types/classroom";
import { User } from "@/types/user";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";

export const classroomsQueryOptions = queryOptions({
  queryKey: ["classrooms"],
  queryFn: async () => {
    const res = await apiClient.get<{
      ownClassrooms: Classroom[];
      joinedClassrooms: Classroom[];
    }>("/api/me/classrooms");
    return res.data;
  },
});

export const classroomQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", `classroom-${classroomId}`],
    queryFn: async () => {
      const res = await apiClient.get<Classroom>(`/api/me/classrooms/${classroomId}`);
      return res.data;
    },
  });

export const classroomMemberQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", `classroom-${classroomId}`, "members"],
    queryFn: async () => {
      const res = await apiClient.get<User[]>(`/api/me/classrooms/${classroomId}/members`);
      return res.data;
    },
  });

export const classroomInvitationsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", `classroom-${classroomId}`, "invitations"],
    queryFn: async () => {
      const res = await apiClient.get<ClassroomInvitation[]>(`/api/me/classrooms/${classroomId}/invitations`);
      return res.data;
    },
  });

export const useCreateClassroom = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (values: ClassroomForm) => {
      const res = await apiClient.post<void>("/api/classrooms", values);
      return res.headers.location as string;
    },
    onSuccess: () => queryClient.invalidateQueries(classroomsQueryOptions),
  });
};

export const useInviteClassroomMembers = (classroomId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (values: InviteForm) => {
      const res = await apiClient.post<void>(`/api/classrooms/${classroomId}/members`, {
        memberEmails: values.memberEmails.split("\n").filter(Boolean),
      });
      return res.data;
    },
    onSuccess: () => queryClient.invalidateQueries(classroomInvitationsQueryOptions(classroomId)),
  });
};

export const useJoinClassroom = (classroomId: string,invitationId: string) => {
  return useMutation({
    mutationFn: async () => {
      const res = await apiClient.post<void>(`/api/classrooms/${classroomId}/invitations/${invitationId}`);
      return res.data
    }
  })
}
