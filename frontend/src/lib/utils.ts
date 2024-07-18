import axios from "axios";
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { format } from "date-fns";
import { AssignmentApi, AuthApi, ClassroomApi, InfoApi, MemberApi, ProjectApi, TeamApi } from "@/swagger-client";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const getUUIDFromLocation = (location: string) => location.split("/").pop()!;

export const formatDate = (date: Parameters<typeof format>[0]) => format(date, "PPP");

const apiClient = axios.create({ withCredentials: true });

export const createAuthApi = () =>
  new AuthApi({
    baseOptions: {
      withCredentials: true,
    },
  });

export const createClassroomApi = () => new ClassroomApi(undefined, "", apiClient);
export const createAssignmentApi = () => new AssignmentApi(undefined, "", apiClient);
export const createProjectApi = () => new ProjectApi(undefined, "", apiClient);
export const createMemberApi = () => new MemberApi(undefined, "", apiClient);
export const createTeamApi = () => new TeamApi(undefined, "", apiClient);
export const createInfoApi = () => new InfoApi(undefined, "", apiClient);
