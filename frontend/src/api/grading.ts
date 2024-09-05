import { createGradingApi } from "@/lib/utils";
import { useCsrf } from "@/provider/csrfProvider";
import { UpdateGradingRequest } from "@/swagger-client";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { authCsrfQueryOptions } from "./auth";

const apiClient = createGradingApi();

export const classroomGradingRubricsQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "grading"],
    queryFn: async () => {
      const res = await apiClient.getGradingRubrics(classroomId);
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
