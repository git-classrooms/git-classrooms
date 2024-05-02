import { createTeamApi } from "@/lib/utils";
import { useCsrf } from "@/provider/csrfProvider";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { authCsrfQueryOptions } from "./auth";
import { TeamForm } from "@/types/team";
import { joinedClassroomQueryOptions, ownedClassroomQueryOptions } from "./classrooms";

export const ownedClassroomTeamsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", classroomId, "teams"],
    queryFn: async () => {
      const api = createTeamApi();
      const res = await api.getOwnedClassroomTeams(classroomId);
      return res.data;
    },
  });

export const ownedClassroomTeamQueryOptions = (classroomId: string, teamId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", classroomId, "teams", teamId],
    queryFn: async () => {
      const api = createTeamApi();
      const res = await api.getOwnedClassroomTeam(classroomId, teamId);
      return res.data;
    },
  });

export const joinedClassroomTeamsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["joinedClassrooms", classroomId, "teams"],
    queryFn: async () => {
      const api = createTeamApi();
      const res = await api.getJoinedClassroomTeams(classroomId);
      return res.data;
    },
  });

export const joinedClassroomTeamQueryOptions = (classroomId: string, teamId: string) =>
  queryOptions({
    queryKey: ["joinedClassrooms", classroomId, "teams", teamId],
    queryFn: async () => {
      const api = createTeamApi();
      const res = await api.getJoinedClassroomTeam(classroomId, teamId);
      return res.data;
    },
  });

export const useJoinTeam = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (teamId: string) => {
      const api = createTeamApi();
      const res = await api.joinJoinedClassroomTeam(classroomId, teamId, csrfToken);
      return res.data;
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
      queryClient.invalidateQueries(joinedClassroomTeamsQueryOptions(classroomId));
      queryClient.invalidateQueries(joinedClassroomQueryOptions(classroomId));
    },
  });
};

export const useCreateTeamJoinedClassroom = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: TeamForm) => {
      const api = createTeamApi();
      const res = await api.createJoinedClassroomTeam(values, csrfToken, classroomId);
      return res.data;
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
      queryClient.invalidateQueries(joinedClassroomTeamsQueryOptions(classroomId));
      queryClient.invalidateQueries(joinedClassroomQueryOptions(classroomId));
    },
  });
};

export const useCreateTeamOwnedClassroom = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: TeamForm) => {
      const api = createTeamApi();
      const res = await api.createOwnedClassroomTeam(values, csrfToken, classroomId);
      return res.data;
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
      queryClient.invalidateQueries(ownedClassroomTeamsQueryOptions(classroomId));
      queryClient.invalidateQueries(ownedClassroomQueryOptions(classroomId));
    },
  });
};
