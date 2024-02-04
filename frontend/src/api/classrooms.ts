import { apiClient } from "@/lib/utils";
import { Classroom, ClassroomForm } from "@/types/classroom";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";

export const classRoomsQueryOptions = queryOptions({
  queryKey: ["classrooms"],
  queryFn: async () => {
    const res = await apiClient.get<{ ownClassrooms: Classroom[], joinedClassrooms: Classroom[] }>("/api/me/classrooms");
    return res.data;
  },
});


export const useCreateClassRoom = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (values: ClassroomForm) => {
      const res = await apiClient.post<void>("/api/classrooms", values);
      return res.headers.location as string;
    },
    onSuccess: () => queryClient.invalidateQueries(classRoomsQueryOptions),
  });
};

