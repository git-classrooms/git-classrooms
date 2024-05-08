import { createClassroomApi, createMemberApi } from "@/lib/utils";
import { ClassroomForm, InviteForm } from "@/types/classroom";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { authCsrfQueryOptions } from "@/api/auth.ts";
import { useCsrf } from "@/provider/csrfProvider";
import { Action } from "@/swagger-client";

export const ownedClassroomsQueryOptions = queryOptions({
  queryKey: ["ownedClassrooms"],
  queryFn: async () => {
    const api = createClassroomApi();
    const res = await api.getOwnedClassrooms();
    return res.data;
  },
});

export const joinedClassroomsQueryOptions = queryOptions({
  queryKey: ["joinedClassrooms"],
  queryFn: async () => {
    const api = createClassroomApi();
    const res = await api.getJoinedClassrooms();
    return res.data;
  },
});

export const invitationInfoQueryOptions = (invitationId: string) =>
  queryOptions({
    queryKey: ["invitations", invitationId],
    queryFn: async () => {
      const api = createClassroomApi();
      const res = await api.getInvitationInfo(invitationId);
      return res.data;
    },
  });

export const joinedClassroomQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["joinedClassrooms", classroomId],
    queryFn: async () => {
      const api = createClassroomApi();
      const res = await api.getJoinedClassroom(classroomId);
      return res.data;
    },
  });

export const ownedClassroomQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`],
    queryFn: async () => {
      const api = createClassroomApi();
      const res = await api.getOwnedClassroom(classroomId);
      return res.data;
    },
  });

export const ownedClassroomMemberQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "members"],
    queryFn: async () => {
      const api = createMemberApi();
      const res = await api.getOwnedClassroomMembers(classroomId);
      return res.data;
    },
  });

export const ownedClassroomInvitationsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "invitations"],
    queryFn: async () => {
      const api = createClassroomApi();
      const res = await api.getOwnedClassroomInvitations(classroomId);
      return res.data;
    },
  });

export const useCreateClassroom = () => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: ClassroomForm) => {
      const api = createClassroomApi();
      const res = await api.createClassroom(values, csrfToken);
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
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: InviteForm) => {
      const data = { memberEmails: values.memberEmails.split("\n").filter(Boolean) };
      const api = createClassroomApi();
      const res = await api.inviteToClassroom(data, csrfToken, classroomId);
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
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (action: Action) => {
      const api = createClassroomApi();
      const res = await api.joinClassroom({ invitationId, action }, csrfToken);
      return res.headers.location as string;
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};
