import axios from "axios";
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { format } from "date-fns";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const getUUIDFromLocation = (location: string) => location.split("/").pop()!;

export const formatDate = (date: Parameters<typeof format>[0]) => format(date, "PPP");

export async function isAuthenticated() {
  try {
    await axios.get("/api/auth", { withCredentials: true });
    return true;
  } catch (e) {
    return false;
  }
}

export const apiClient = axios.create({ baseURL: "/api", withCredentials: true });
