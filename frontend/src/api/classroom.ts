import { createClassroomApi } from "@/lib/utils";
import { ClassroomForm, Filter, InviteForm } from "@/types/classroom";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { authCsrfQueryOptions } from "@/api/auth.ts";
import { useCsrf } from "@/provider/csrfProvider";
import { Action, CreateClassroomRequest } from "@/swagger-client";

const apiClient = createClassroomApi();

// Queries

export const classroomsQueryOptions = (filter: Filter | undefined = undefined) =>
  queryOptions({
    queryKey: ["classrooms", filter],
    queryFn: async () => {
      const res = await apiClient.getClassrooms(filter);
      return res.data;
    },
  });

export const classroomQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId],
    queryFn: async () => {
      const res = await apiClient.getClassroom(classroomId);
      return res.data;
    },
  });

export const classroomInvitationsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "invitations"],
    queryFn: async () => {
      const res = await apiClient.getClassroomInvitations(classroomId);
      return res.data;
    },
  });

export const classroomInvitationQueryOptions = (classroomId: string, invitationId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "invitations", invitationId],
    queryFn: async () => {
      const res = await apiClient.getClassroomInvitation(classroomId, invitationId);
      return res.data;
    },
  });

export const classroomTemplatesQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "templates"],
    queryFn: async () => {
      const res = await apiClient.getClassroomTemplates(classroomId);
      return res.data;
    },
  });

// Mutations
export const useCreateClassroom = () => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: ClassroomForm) => {
      const body: CreateClassroomRequest = {
        name: values.name,
        description: values.description,
        studentsViewAllProjects: values.studentsViewAllProjects,
        createTeams: values.teamsEnabled? values.createTeams : false,
        maxTeamSize: values.teamsEnabled? values.maxTeamSize : 1,
        maxTeams: values.teamsEnabled? values.maxTeams : 0,
      };
      const res = await apiClient.createClassroomV2(body, csrfToken);
      return res.headers.location as string;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(classroomsQueryOptions());
      queryClient.invalidateQueries(classroomsQueryOptions(Filter.Owned));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useUpdateClassroom = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: ClassroomForm) => {
      const res = await apiClient.updateClassroomV2(values, csrfToken, classroomId);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(classroomsQueryOptions());
      queryClient.invalidateQueries(classroomsQueryOptions(Filter.Owned));
      queryClient.invalidateQueries(classroomQueryOptions(classroomId));
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
      const res = await apiClient.inviteToClassroomV2(data, csrfToken, classroomId);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(classroomInvitationsQueryOptions(classroomId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useJoinClassroom = (classroomId: string, invitationId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (action: Action) => {
      const res = await apiClient.joinClassroomV2({ invitationId, action }, csrfToken, classroomId);
      return res.headers.location as string;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(classroomsQueryOptions());
      queryClient.invalidateQueries(classroomsQueryOptions(Filter.Student));
      queryClient.invalidateQueries(classroomInvitationQueryOptions(classroomId, invitationId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};
