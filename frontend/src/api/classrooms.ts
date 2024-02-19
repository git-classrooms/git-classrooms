import { apiClient } from "@/lib/utils";
import { ClassroomForm, ClassroomInvitation, InviteForm, JoinedClassroom, OwnedClassroom } from "@/types/classroom";
import { User } from "@/types/user";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";

export const ownedClassroomsQueryOptions = queryOptions({
  queryKey: ["ownedClassrooms"],
  queryFn: async () => {
    const res = await apiClient.get<OwnedClassroom[]>("/classrooms/owned");
    return res.data;
  },
});

export const joinedClassroomsQueryOptions = queryOptions({
  queryKey: ["joinedClassrooms"],
  queryFn: async () => {
    const res = await apiClient.get<JoinedClassroom[]>("/classrooms/joined");
    return res.data;
  },
});

export const ownedClassroomQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`],
    queryFn: async () => {
      const res = await apiClient.get<OwnedClassroom>(`/classrooms/owned/${classroomId}`);
      return res.data;
    },
  });

export const ownedClassroomMemberQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "members"],
    queryFn: async () => {
      const res = await apiClient.get<User[]>(`/classrooms/owned/${classroomId}/members`);
      return res.data;
    },
  });

export const ownedClassroomInvitationsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "invitations"],
    queryFn: async () => {
      const res = await apiClient.get<ClassroomInvitation[]>(`/classrooms/owned/${classroomId}/invitations`);
      return res.data;
    },
  });

export const useCreateClassroom = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (values: ClassroomForm) => {
      const res = await apiClient.post<void>("/classrooms/owned", values);
      return res.headers.location as string;
    },
    onSuccess: () => queryClient.invalidateQueries(ownedClassroomsQueryOptions),
  });
};

export const useInviteClassroomMembers = (classroomId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (values: InviteForm) => {
      const res = await apiClient.post<void>(`/classrooms/owned/${classroomId}/invitations`, {
        memberEmails: values.memberEmails.split("\n").filter(Boolean),
      });
      return res.data;
    },
    onSuccess: () => queryClient.invalidateQueries(ownedClassroomInvitationsQueryOptions(classroomId)),
  });
};

export const useJoinClassroom = (invitationId: string) => {
  return useMutation({
    mutationFn: async () => {
      const res = await apiClient.post<void>("/classrooms/joined", { invitationId });
      return res.data;
    },
  });
};
