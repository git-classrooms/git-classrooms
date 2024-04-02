import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/utils.ts";
import { Assignment, AssignmentProject, CreateAssignmentForm, TemplateProject } from "@/types/assignments.ts";
import { authCsrfQueryOptions } from "@/api/auth.ts";
import { useCsrf } from "@/provider/csrfProvider";

export const ownedAssignmentsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "assignments"],
    queryFn: async () => {
      const res = await apiClient.get<Assignment[]>(`/classrooms/owned/${classroomId}/assignments`);
      return res.data;
    },
  });

export const ownedAssignmentQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "assignments", `classroom-${assignmentId}`],
    queryFn: async () => {
      const res = await apiClient.get<Assignment>(`/classrooms/owned/${classroomId}/assignments/${assignmentId}`);
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
      const res = await apiClient.get<AssignmentProject[]>(
        `/classrooms/owned/${classroomId}/assignments/${assignmentId}/projects`,
      );
      return res.data;
    },
  });

export const ownedTemplateProjectQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["ownedClassrooms", `ownedClassroom-${classroomId}`, "templateProjects"],
    queryFn: async () => {
      const res = await apiClient.get<TemplateProject[]>(`/classrooms/owned/${classroomId}/templateProjects`);
      return res.data;
    },
  });

export const useCreateAssignment = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { apiClient } = useCsrf();
  return useMutation({
    mutationFn: async (values: CreateAssignmentForm) => {
      const res = await apiClient.post<void>(`/classrooms/owned/${classroomId}/assignments`, values);
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
  const { apiClient } = useCsrf();
  return useMutation({
    mutationFn: async () => {
      await apiClient.post(`/classrooms/owned/${classroomId}/assignments/${assignmentId}/projects`);
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
  const { apiClient } = useCsrf();
  return useMutation({
    mutationFn: async () => {
      const res = await apiClient.post<void>(`/classrooms/joined/${classroomId}/assignments/${assignmentId}/accept`);
      return res.data;
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};
