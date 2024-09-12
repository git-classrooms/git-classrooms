import { createRunnersApi } from "@/lib/utils";
import { queryOptions } from "@tanstack/react-query";

const apiClient = createRunnersApi();

export const classroomAvailableRunnersQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "runners"],
    queryFn: async () => {
      const res = await apiClient.getClassroomRunnersAreAvailable(classroomId);
      return res.data;
    },
  });
