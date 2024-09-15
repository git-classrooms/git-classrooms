import axios, { isAxiosError } from "axios";
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { format } from "date-fns";
import {
  AssignmentApi,
  AuthApi,
  ClassroomApi,
  InfoApi,
  MemberApi,
  ProjectApi,
  TeamApi,
  GradingApi,
  ReportApi,
  RunnersApi,
  UserClassroomResponse,
  HTTPError,
} from "@/swagger-client";
import { Role } from "@/types/classroom";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const getUUIDFromLocation = (location: string) => location.split("/").pop()!;

export const formatDate = (date: Parameters<typeof format>[0]) => format(date, "PPP");

export const formatDateWithTime = (date: Parameters<typeof format>[0]) => format(date, "PPP HH:mm:ss");

const apiClient = axios.create({ withCredentials: true });

export const createAuthApi = () =>
  new AuthApi({
    baseOptions: {
      withCredentials: true,
    },
  });

/* eslint-disable @typescript-eslint/no-explicit-any */
export const unwrapApiError = <T = any, D = any>(error: Error | null): Error | null => {
  if (isAxiosError<T, D>(error)) {
    return new Error((error.response?.data as HTTPError | undefined)?.error ?? error.message);
  }
  return error;
};

export const createClassroomApi = () => new ClassroomApi(undefined, "", apiClient);
export const createAssignmentApi = () => new AssignmentApi(undefined, "", apiClient);
export const createProjectApi = () => new ProjectApi(undefined, "", apiClient);
export const createMemberApi = () => new MemberApi(undefined, "", apiClient);
export const createTeamApi = () => new TeamApi(undefined, "", apiClient);
export const createInfoApi = () => new InfoApi(undefined, "", apiClient);
export const createGradingApi = () => new GradingApi(undefined, "", apiClient);
export const createReportApi = () => new ReportApi(undefined, "", apiClient);
export const createRunnersApi = () => new RunnersApi(undefined, "", apiClient);

export const isCreator = (user: UserClassroomResponse) => user.user.id === user.classroom.ownerId;
export const isOwner = (user: UserClassroomResponse) => user.role === Role.Owner;
export const isModerator = (user: UserClassroomResponse) => user.role <= Role.Moderator;
export const isStudent = (user: UserClassroomResponse) => user.role === Role.Student;
