import { createGradingApi } from "@/lib/utils";
import { useCsrf } from "@/provider/csrfProvider";
import { UpdateAssignmentRubricsRequest, UpdateAssignmentTestRequest, UpdateGradingRequest } from "@/swagger-client";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { authCsrfQueryOptions } from "./auth";
import { assignmentQueryOptions } from "./assignment";

const apiClient = createGradingApi();

export const classroomGradingRubricsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "grading"],
    queryFn: async () => {
      const res = await apiClient.getGradingRubrics(classroomId);
      return res.data;
    },
  });

export const assignmentGradingRubricsQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "assignments", assignmentId, "grading"],
    queryFn: async () => {
      const res = await apiClient.getAssignmentGradingRubrics(classroomId, assignmentId);
      return res.data;
    },
  });

export const assignmentTestsQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "assignments", assignmentId, "tests"],
    queryFn: async () => {
      const res = await apiClient.getClassroomAssignmentTests(classroomId, assignmentId);
      return res.data;
    },
  });

export const projectGradingResultsQueryOptions = (classroomId: string, projectId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "projects", projectId, "grading"],
    queryFn: async () => {
      const res = await apiClient.apiV2ClassroomsClassroomIdProjectsProjectIdGradingGet(classroomId, projectId);
      return res.data;
    },
  });

export const useUpdateClassroomRubrics = (classroomId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (data: UpdateGradingRequest) => {
      return apiClient.updateGradingRubrics(data, csrfToken, classroomId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(classroomGradingRubricsQueryOptions(classroomId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useUpdateAssignmentRubrics = (classroomId: string, assignmentId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (data: UpdateAssignmentRubricsRequest) => {
      return apiClient.updateAssignmentGradingRubrics(data, csrfToken, classroomId, assignmentId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(assignmentGradingRubricsQueryOptions(classroomId, assignmentId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useUpdateAssignmentTests = (classroomId: string, assignmentId: string) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (data: UpdateAssignmentTestRequest) => {
      return apiClient.updateAssignmentTests(data, csrfToken, classroomId, assignmentId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(assignmentTestsQueryOptions(classroomId, assignmentId));
      queryClient.invalidateQueries(assignmentQueryOptions(classroomId, assignmentId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};
