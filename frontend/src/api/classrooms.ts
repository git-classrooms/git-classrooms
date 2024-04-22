import { apiClient } from "@/lib/utils";
import {
  ClassroomForm,
  ClassroomInvitation,
  InviteForm,
  UserClassroom,
  OwnedClassroom
} from "@/types/classroom";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { authCsrfQueryOptions } from "@/api/auth.ts";
import { useCsrf } from "@/provider/csrfProvider";

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
    const res = await apiClient.get<UserClassroom[]>("/classrooms/joined");
    return res.data;
  },
});

export const joinedClassroomQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["joinedClassrooms", classroomId],
    queryFn: async () => {
      const res = await apiClient.get<UserClassroom>(`/classrooms/joined/${classroomId}`);
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
      const res = await apiClient.get<UserClassroom[]>(`/classrooms/owned/${classroomId}/members`);
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
  const { apiClient } = useCsrf();
  return useMutation({
    mutationFn: async (values: ClassroomForm) => {
      const res = await apiClient.post<void>("/classrooms/owned", values);
      return res.headers.location as string;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(ownedClassroomsQueryOptions);
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useInviteClassroomMembers = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { apiClient } = useCsrf();
  return useMutation({
    mutationFn: async (values: InviteForm) => {
      const res = await apiClient.post<void>(`/classrooms/owned/${classroomId}/invitations`, {
        memberEmails: values.memberEmails.split("\n").filter(Boolean),
      });
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(ownedClassroomInvitationsQueryOptions(classroomId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useJoinClassroom = (invitationId: string) => {
  const queryClient = useQueryClient();
  const { apiClient } = useCsrf();
  return useMutation({
    mutationFn: async () => {
      const res = await apiClient.post<void>("/classrooms/joined", { invitationId });
      return res.headers.location as string;
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};
