import { createMemberApi } from "@/lib/utils";
import { useCsrf } from "@/provider/csrfProvider";
import { queryOptions, useMutation, useQueryClient } from "@tanstack/react-query";
import { authCsrfQueryOptions } from "./auth";
import { Role } from "@/types/classroom";
import { MemberForm } from "@/types/member";

const apiClient = createMemberApi();

// Queries

export const membersQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "members"],
    queryFn: async () => {
      const res = await apiClient.getClassroomMembers(classroomId);
      return res.data;
    },
  });

export const memberQueryOptions = (classroomId: string, memberId: number) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "members", memberId],
    queryFn: async () => {
      const res = await apiClient.getClassroomMember(classroomId, memberId);
      return res.data;
    },
  });

export const teamMembersQueryOptions = (classroomId: string, teamId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "teams", teamId, "members"],
    queryFn: async () => {
      const res = await apiClient.getClassroomTeamMembers(classroomId, teamId);
      return res.data;
    },
  });

export const teamMemberQueryOptions = (classroomId: string, teamId: string, memberId: number) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "teams", teamId, "members", memberId],
    queryFn: async () => {
      const res = await apiClient.getClassroomTeamMember(classroomId, teamId, memberId);
      return res.data;
    },
  });

// Mutations

export const useUpdateMemberRole = (classroomId: string, memberId: number) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (values: MemberForm) => {
      const res = await apiClient.updateMemberRole({ role: Role[values.role] }, csrfToken, classroomId, memberId);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(membersQueryOptions(classroomId));
      queryClient.invalidateQueries(memberQueryOptions(classroomId, memberId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};

export const useUpdateMemberTeam = (classroomId: string, memberId: number) => {
  const queryClient = useQueryClient();
  const { csrfToken } = useCsrf();
  return useMutation({
    mutationFn: async (teamId: string) => {
      const res = await apiClient.updateMemberTeam({ teamId }, csrfToken, classroomId, memberId);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(membersQueryOptions(classroomId));
      queryClient.invalidateQueries(memberQueryOptions(classroomId, memberId));
    },
    onSettled: () => {
      queryClient.invalidateQueries(authCsrfQueryOptions);
    },
  });
};
