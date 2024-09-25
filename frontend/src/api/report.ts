import { createReportApi } from "@/lib/utils";
import { queryOptions } from "@tanstack/react-query";

const apiClient = createReportApi();

export const classroomReportQueryOptions = (classroomId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "report"],
    queryFn: async () => {
      const res = await apiClient.getClassroomReport(classroomId);
      return res.data;
    },
  });

export const teamReportQueryOptions = (classroomId: string, teamId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "teams", teamId, "report"],
    queryFn: async () => {
      const res = await apiClient.getClassroomTeamReport(classroomId, teamId);
      return res.data;
    },
  });

export const assignmentReportQueryOptions = (classroomId: string, assignmentId: string) =>
  queryOptions({
    queryKey: ["classrooms", classroomId, "assignments", assignmentId, "report"],
    queryFn: async () => {
      const res = await apiClient.getClassroomAssignmentReport(classroomId, assignmentId);
      return res.data;
    },
  });
