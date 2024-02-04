import { apiClient } from "@/lib/utils";
import { Assignment } from "@/types/assignments";
import { Classroom, ClassroomForm } from "@/types/classroom";
import { User } from "@/types/user";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";

export const classroomsQueryOptions = queryOptions({
  queryKey: ["classrooms"],
  queryFn: async () => {
    const res = await apiClient.get<{ ownClassrooms: Classroom[], joinedClassrooms: Classroom[] }>("/api/me/classrooms");
    return res.data;
  },
});

export const classroomQueryOptions = (classroomId: string) => queryOptions({
  queryKey: ["classrooms", `classroom-${classroomId}`],
  queryFn: async () => {
    const res = await apiClient.get<Classroom>(`/api/me/classrooms/${classroomId}`);
    return res.data;
  },
});

export const classroomMemberQueryOptions = (classroomId: string) => queryOptions({
  queryKey: ["classrooms", `classroom-${classroomId}`, "members"],
  queryFn: async () => {
    const res = await apiClient.get<User[]>(`/api/me/classrooms/${classroomId}/members`);
    return res.data;
  },
});

export const assignmentsQueryOptions = (classroomId: string) => queryOptions({
  queryKey: ["classrooms", `classroom-${classroomId}`, "assignments"],
  queryFn: async () => {
    const res = await apiClient.get<Assignment[]>(`/api/classrooms/${classroomId}/assignments`);
    return res.data;
  },
});

export const useCreateClassroom = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (values: ClassroomForm) => {
      const res = await apiClient.post<void>("/api/classrooms", values);
      return res.headers.location as string;
    },
    onSuccess: () => queryClient.invalidateQueries(classroomsQueryOptions),
  });
};

