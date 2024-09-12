import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { createAssignmentApi } from "@/lib/utils.ts";
import { authCsrfQueryOptions } from "@/api/auth.ts";
import { useCsrf } from "@/provider/csrfProvider";
import { CreateAssignmentRequest, UpdateAssignmentRequest } from "@/swagger-client";

const apiClient = createAssignmentApi();

// Queries

export const assignmentsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "assignments"],
    queryFn: async () => {
      const res = await apiClient.getClassroomAssignments(classroomId);
      return res.data;
    },
  });

export const assignmentQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "assignments", assignmentId],
    queryFn: async () => {
      const res = await apiClient.getClassroomAssignment(classroomId, assignmentId);
      return res.data;
    },
  });

export const activeAssignmentQueryOptions = () =>
  queryOptions({
    queryKey: ["activeAssignments"],
    queryFn: async () => {
      const res = await apiClient.getActiveAssignments();
      return res.data;
    },
  });

// Mutations

export const useCreateAssignment = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: CreateAssignmentRequest) => {
      const res = await apiClient.createAssignmentV2(values, csrfToken, classroomId);
      return res.headers.location as string;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(assignmentsQueryOptions(classroomId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useUpdateAssignment = (classroomId: string, assignmentId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: UpdateAssignmentRequest) => {
      const res = await apiClient.updateAssignment(values, csrfToken, classroomId, assignmentId);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(assignmentsQueryOptions(classroomId));
      queryClient.invalidateQueries(assignmentQueryOptions(classroomId, assignmentId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
}
