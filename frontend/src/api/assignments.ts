import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/utils.ts";
import { Assignment, AssignmentProject, CreateAssignmentForm, TemplateProject } from "@/types/assignments.ts";

export const assignmentsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", `classroom-${classroomId}`, "assignments"],
    queryFn: async () => {
      const res = await apiClient.get<Assignment[]>(`/api/classrooms/${classroomId}/assignments`);
      return res.data;
    },
  });

export const assignmentQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: ["classrooms", `classroom-${classroomId}`, "assignments", `classroom-${assignmentId}`],
    queryFn: async () => {
      const res = await apiClient.get<Assignment>(`/api/classrooms/${classroomId}/assignments/${assignmentId}`);
      return res.data;
    },
  });

export const assignmentProjectsQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: ["classrooms", `classroom-${classroomId}`, "assignments", `classroom-${assignmentId}`, "projects"],
    queryFn: async () => {
      const res = await apiClient.get<AssignmentProject[]>(
        `/api/classrooms/${classroomId}/assignments/${assignmentId}/projects`,
      );
      return res.data;
    },
  });

export const templateProjectQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", `classroom-${classroomId}`, "templateProjects"],
    queryFn: async () => {
      const res = await apiClient.get<TemplateProject[]>(`/api/me/classrooms/${classroomId}/templateProjects`);
      return res.data;
    },
  });

export const useCreateAssignment = (classroomId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (values: CreateAssignmentForm) => {
      const res = await apiClient.post<Assignment>(`/api/classrooms/${classroomId}/assignments`, values);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(assignmentsQueryOptions(classroomId));
    },
  });
};

export const useInviteAssignmentMembers = (classroomId: string, assignmentId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      await apiClient.post(`/api/classrooms/${classroomId}/assignments/${assignmentId}/projects`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(assignmentProjectsQueryOptions(classroomId, assignmentId));
    },
  });
};
