import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { createAssignmentApi, createClassroomApi } from "@/lib/utils.ts";
import { authCsrfQueryOptions } from "@/api/auth.ts";
import { useCsrf } from "@/provider/csrfProvider";
import { CreateAssignmentRequest } from "@/swagger-client";

export const ownedAssignmentsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "assignments"],
    queryFn: async () => {
      const api = createAssignmentApi();
      const res = await api.getOwnedClassroomAssignments(classroomId);
      return res.data;
    },
  });

export const ownedAssignmentQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "assignments", `classroom-${assignmentId}`],
    queryFn: async () => {
      const api = createAssignmentApi();
      const res = await api.getOwnedClassroomAssignment(classroomId, assignmentId);
      return res.data;
    },
  });

export const ownedAssignmentProjectsQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: [
      "ownedClassrooms",
      `ownedClassroom-${classroomId}`,
      "assignments",
      `classroom-${assignmentId}`,
      "projects",
    ],
    queryFn: async () => {
      const api = createAssignmentApi();
      const res = await api.getOwnedClassroomAssignmentProjects(classroomId, assignmentId);
      return res.data;
    },
  });

export const ownedTemplateProjectQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "templateProjects"],
    queryFn: async () => {
      const api = createClassroomApi();
      const res = await api.getOwnedClassroomTemplates(classroomId);
      return res.data;
    },
  });

export const joinedClassroomAssignmentQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: ["joinedClassrooms", `joinedClassroom-${classroomId}`, "assignments", `classroom-${assignmentId}`],
    queryFn: async () => {
      const api = createAssignmentApi();
      const res = await api.getJoinedClassroomAssignment(classroomId, assignmentId);
      return res.data;
    },
  });

export const useCreateAssignment = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: CreateAssignmentRequest) => {
      const api = createAssignmentApi();
      const res = await api.createAssignment(values, csrfToken, classroomId);
      return res.headers.location as string;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(ownedAssignmentsQueryOptions(classroomId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useInviteAssignmentMembers = (classroomId: string, assignmentId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async () => {
      const api = createAssignmentApi();
      const res = await api.inviteToAssignment(classroomId, assignmentId, csrfToken);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(ownedAssignmentProjectsQueryOptions(classroomId, assignmentId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};
export const useAcceptAssignment = (classroomId: string, assignmentId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async () => {
      const api = createAssignmentApi();
      const res = await api.acceptAssignment(classroomId, assignmentId, csrfToken);
      return res.data;
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};
export const ownedClassroomTeamProjectsQueryOptions = (classroomId: string, teamId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "teams", `team-${teamId}`, "assignments"],
    queryFn: async () => {
      const api = createAssignmentApi();
      const res = await api.getOwnedClassroomTeamProjects(classroomId, teamId);
      return res.data;
    },
  });

