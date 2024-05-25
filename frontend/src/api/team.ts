import { createTeamApi } from "@/lib/utils";
import { useCsrf } from "@/provider/csrfProvider";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { authCsrfQueryOptions } from "./auth";
import { TeamForm } from "@/types/team";
import { classroomQueryOptions } from "./classroom";

const apiClient = createTeamApi();

// Queries

export const teamsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "teams"],
    queryFn: async () => {
      const res = await apiClient.getClassroomTeams(classroomId);
      return res.data;
    },
  });

export const teamQueryOptions = (classroomId: string, teamId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "teams", teamId],
    queryFn: async () => {
      const res = await apiClient.getClassroomTeam(classroomId, teamId);
      return res.data;
    },
  });

// Mutations

export const useCreateTeam = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: TeamForm) => {
      const res = await apiClient.createTeam(values, csrfToken, classroomId);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(teamsQueryOptions(classroomId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useUpdateTeam = (classroomId: string, teamId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: TeamForm) => {
      const res = await apiClient.updateTeam(values, csrfToken, classroomId, teamId);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(teamsQueryOptions(classroomId));
      queryClient.invalidateQueries(teamQueryOptions(classroomId, teamId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useJoinTeam = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (teamId: string) => {
      const res = await apiClient.joinTeam(classroomId, teamId, csrfToken);
      return res.headers.location as string;
    },
    onSuccess: (_, teamId) => {
      queryClient.invalidateQueries(classroomQueryOptions(classroomId));
      queryClient.invalidateQueries(teamsQueryOptions(classroomId));
      queryClient.invalidateQueries(teamQueryOptions(classroomId, teamId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};
