import { apiClient } from "@/lib/utils";
import { useCsrf } from "@/provider/csrfProvider";
import { TeamApi } from "@/swagger-client";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { authCsrfQueryOptions } from "./auth";
import { TeamForm } from "@/types/team";
import { joinedClassroomQueryOptions } from "./classrooms";

export const ownedClassroomTeamsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["teams"],
    queryFn: async () => {
      const api = new TeamApi(undefined, "", apiClient);
      const res = await api.getOwnedClassroomTeams(classroomId);
      return res.data;
    },
  });

export const joinedClassroomTeamsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["joinedClassroom", classroomId, "teams"],
    queryFn: async () => {
      const api = new TeamApi(undefined, "", apiClient);
      const res = await api.getJoinedClassroom(classroomId);
      return res.data;
    },
  });

export const useJoinTeam = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (teamId: string) => {
      const api = new TeamApi(undefined, "", apiClient);
      const res = await api.classroomsJoinedClassroomIdTeamsTeamIdJoinPost(classroomId, teamId, csrfToken);
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
      const api = new TeamApi(undefined, "", apiClient);
      const res = await api.classroomsJoinedClassroomIdTeamsPost(values, csrfToken, classroomId);
      return res.data;
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
      queryClient.invalidateQueries(joinedClassroomTeamsQueryOptions(classroomId));
      queryClient.invalidateQueries(joinedClassroomQueryOptions(classroomId));
    },
  });
};
