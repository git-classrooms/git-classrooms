import { createProjectApi } from "@/lib/utils";
import { useCsrf } from "@/provider/csrfProvider";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { authCsrfQueryOptions } from "./auth";

const apiClient = createProjectApi();

// Queries

export const assignmentProjectsQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "assignments", assignmentId, "projects"],
    queryFn: async () => {
      const res = await apiClient.getClassroomAssignmentProjects(classroomId, assignmentId);
      return res.data;
    },
  });

export const assignmentProjectQueryOptions = (classroomId: string, assignmentId: string, projectId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "assignments", assignmentId, "projects", projectId],
    queryFn: async () => {
      const res = await apiClient.getClassroomAssignmentProject(classroomId, assignmentId, projectId);
      return res.data;
    },
  });

export const projectsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "projects"],
    queryFn: async () => {
      const res = await apiClient.getClassroomProjects(classroomId);
      return res.data;
    },
  });

export const projectQueryOptions = (classroomId: string, projectId: string, refetchInterval?: number) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "projects", projectId],
    queryFn: async () => {
      const res = await apiClient.getClassroomProject(classroomId, projectId);
      return res.data;
    },
    refetchInterval,
  });

export const teamProjectsQueryOptions = (classroomId: string, teamId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "teams", teamId, "projects"],
    queryFn: async () => {
      const res = await apiClient.getClassroomTeamProjects(classroomId, teamId);
      return res.data;
    },
  });

export const teamProjectQueryOptions = (classroomId: string, teamId: string, projectId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "teams", teamId, "projects", projectId],
    queryFn: async () => {
      const res = await apiClient.getClassroomTeamProject(classroomId, teamId, projectId);
      return res.data;
    },
  });

// Mutations

export const useInviteToAssignment = (classroomId: string, assignmentId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async () => {
      const res = await apiClient.inviteToAssignment(classroomId, assignmentId, csrfToken);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(assignmentProjectsQueryOptions(classroomId, assignmentId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useAcceptAssignment = (classroomId: string, projectId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async () => {
      const res = await apiClient.acceptAssignment(classroomId, projectId, csrfToken);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(projectsQueryOptions(classroomId));
      queryClient.invalidateQueries(projectQueryOptions(classroomId, projectId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};
